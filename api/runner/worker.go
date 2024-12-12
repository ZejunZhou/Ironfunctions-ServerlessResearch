package runner

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/ZejunZhou/Ironfunctions-ServerlessResearch/api/runner/drivers"
	"github.com/ZejunZhou/Ironfunctions-ServerlessResearch/api/runner/protocol"
	"github.com/ZejunZhou/Ironfunctions-ServerlessResearch/api/runner/task"
)

// hot functions - theory of operation
//
// A function is converted into a hot function if its `Format` is either
// a streamable format/protocol. At the very first task request a hot
// container shall be started and run it. Each hot function has an internal
// clock that actually halts the container if it goes idle long enough. In the
// absence of workload, it just stops the whole clockwork.
//
// Internally, the hot function uses a modified Config whose Stdin and Stdout
// are bound to an internal pipe. This internal pipe is fed with incoming tasks
// Stdin and feeds incoming tasks with Stdout.
//
// Each execution is the alternation of feeding hot functions stdin with tasks
// stdin, and reading the answer back from containers stdout. For all `Format`s
// we send embedded into the message metadata to help the container to know when
// to stop reading from its stdin and Functions expect the container to do the
// same. Refer to api/runner/protocol.go for details of these communications.
//
// hot functions implementation relies in two moving parts (drawn below):
// htfnmgr and htfn. Refer to their respective comments for
// details.
//                             │
//                         Incoming
//                           Task
//                             │
//                             ▼
//                     ┌───────────────┐
//                     │ Task Request  │
//                     │   Main Loop   │
//                     └───────────────┘
//                             │
//                      ┌──────▼────────┐
//                     ┌┴──────────────┐│
//                     │  Per Function ││             non-streamable f()
//             ┌───────│   Container   │├──────┐───────────────┐
//             │       │    Manager    ├┘      │               │
//             │       └───────────────┘       │               │
//             │               │               │               │
//             ▼               ▼               ▼               ▼
//       ┌───────────┐   ┌───────────┐   ┌───────────┐   ┌───────────┐
//       │    Hot    │   │    Hot    │   │    Hot    │   │   Cold    │
//       │ Function  │   │ Function  │   │ Function  │   │ Function  │
//       └───────────┘   └───────────┘   └───────────┘   └───────────┘
//                                           Timeout
//                                           Terminate
//                                           (internal clock)

// RunTask helps sending a task.Request into the common concurrency stream.
// Refer to StartWorkers() to understand what this is about.
func RunTask(tasks chan task.Request, ctx context.Context, cfg *task.Config) (drivers.RunResult, error) {
	tresp := make(chan task.Response)
	treq := task.Request{Ctx: ctx, Config: cfg, Response: tresp}
	tasks <- treq
	resp := <-treq.Response
	return resp.Result, resp.Err
}

// StartWorkers operates the common concurrency stream, ie, it will process all
// IronFunctions tasks, either sync or async. In the process, it also dispatches
// the workload to either regular or hot functions.
func StartWorkers(ctx context.Context, rnr *Runner, tasks <-chan task.Request) {
	var wg sync.WaitGroup
	defer wg.Wait()
	var hfcmgr hotFnManager

	for {
		select {
		case <-ctx.Done():
			return
		case task := <-tasks:
			// p is the channel to send the task to the hot function
			p := hfcmgr.getTasks(ctx, rnr, task.Config)
			if p == nil {
				wg.Add(1)
				// normal and non-streamable function, use Stdin and Stdout
				go runTaskReq(rnr, &wg, task)
				continue
			}

			// enqueue the task
			rnr.Start()

			// runner will loop queueHandler to get the task and handleTask
			select {
			case <-ctx.Done():
				return
			case p <- task:
				rnr.Complete()
			}
		}
	}
}

// hotFnManager is the intermediate between the common concurrency stream and
// hot functions. All hot functions share a single task.Request stream per
// function (chn), but each function may have more than one hot function (hc).
type hotFnManager struct {
	chn      map[string]chan task.Request
	hotFnSvr map[string]*hotFnServer
}

// return the stask when it is streamable, otherwise return nil
// it will also check if the hot function is already running
func (h *hotFnManager) getTasks(ctx context.Context, rnr *Runner, cfg *task.Config) chan task.Request {
	// returns true if the format is streamable: HTTP
	isStream, err := protocol.IsStreamable(cfg.Format)
	if err != nil {
		logrus.WithError(err).Info("could not detect container IO protocol")
		return nil
	} else if !isStream {
		return nil
	}

	// HTTP streamable format detected, let's start the hot function manager
	if h.chn == nil {
		h.chn = make(map[string]chan task.Request)
		h.hotFnSvr = make(map[string]*hotFnServer)
	}

	fn := fmt.Sprint(cfg.AppName, ",", cfg.Path, cfg.Image, cfg.Timeout, cfg.Memory, cfg.Format, cfg.MaxConcurrency)
	// logrus.WithField("fn", fn).Info("hot function detected")
	// e.g. hotel,/userhotel_user:0.0.11m0s 256http8
	tasks, ok := h.chn[fn] // tasks is a channel for specific hot function
	if !ok {
		// hot function not running, create a new channel for task requests
		h.chn[fn] = make(chan task.Request)
		tasks = h.chn[fn]

		// new and call server.pipe() to start the hot function
		// it will start looping and see if it needs to start more hot functions
		svr := newHotFnServer(ctx, cfg, rnr, tasks)
		// create a new hot function container
		if err := svr.launch(ctx); err != nil {
			logrus.WithError(err).Error("cannot start hot function supervisor")
			return nil
		}
		h.hotFnSvr[fn] = svr
	} else {
		// hot function already running, just return the channel
		// logrus.WithField("fn", fn).Info("hot function already running")
	}

	return tasks
}

// hotFnServer is part of htfnmgr, abstracted apart for simplicity, its only
// purpose is to test for hot functions saturation and try starting as many as
// needed. In case of absence of workload, it will stop trying to start new hot
// containers.
type hotFnServer struct {
	cfg      *task.Config
	rnr      *Runner
	tasksin  <-chan task.Request
	tasksout chan task.Request
	maxc     chan struct{}
}

func newHotFnServer(ctx context.Context, cfg *task.Config, rnr *Runner, tasks <-chan task.Request) *hotFnServer {
	svr := &hotFnServer{
		cfg:      cfg,
		rnr:      rnr,
		tasksin:  tasks,
		tasksout: make(chan task.Request, 1),
		maxc:     make(chan struct{}, cfg.MaxConcurrency),
	}

	// This pipe will take all incoming tasks and just forward them to the
	// started hot functions. The catch here is that it feeds a buffered
	// channel from an unbuffered one. And this buffered channel is
	// then used to determine the presence of running hot functions.
	// If no hot function is available, tasksout will fill up to its
	// capacity and pipe() will start them.
	go svr.pipe(ctx)
	return svr
}

func (svr *hotFnServer) pipe(ctx context.Context) {
	for {
		select {
		// tasksin is the channel that receives all incoming tasks
		case t := <-svr.tasksin:
			// tasksout is the channel that feeds the hot functions
			svr.tasksout <- t
			if len(svr.tasksout) > 0 {
				logrus.WithField("tasks", len(svr.tasksout)).Info("hot function saturation detected, starting more")
				if err := svr.launch(ctx); err != nil {
					logrus.WithError(err).Error("cannot start more hot functions")
				}
			}
		case <-ctx.Done():
			return
		}
	}
}

// launch starts a new hot function, pass input to container
func (svr *hotFnServer) launch(ctx context.Context) error {
	select {
	// if there is a slot available, start a new hot function
	case svr.maxc <- struct{}{}:
		hc, err := newHotFn(
			svr.cfg,
			protocol.Protocol(svr.cfg.Format),
			svr.tasksout,
			svr.rnr,
		)
		if err != nil {
			return err
		}
		go func() {
			// start the hot function, pass stdin/stdout to container
			hc.serve(ctx)
			// free up the slot
			<-svr.maxc
		}()
	default:
	}

	return nil
}

// HotFn actually interfaces an incoming task from the common concurrency
// stream into a long lived container. If idle long enough, it will stop. It
// uses route configuration to determine which protocol to use.
type HotFn struct {
	cfg   *task.Config
	proto protocol.ContainerIO
	tasks <-chan task.Request

	// Side of the pipe that takes information from outer world
	// and injects into the container.
	in  io.Writer
	out io.Reader

	// Receiving side of the container.
	containerIn  io.Reader
	containerOut io.Writer

	rnr *Runner
}

func newHotFn(cfg *task.Config, proto protocol.Protocol, tasks <-chan task.Request, rnr *Runner) (*HotFn, error) {
	stdinr, stdinw := io.Pipe()
	stdoutr, stdoutw := io.Pipe()

	p, err := protocol.New(proto, stdinw, stdoutr)
	if err != nil {
		return nil, err
	}

	hc := &HotFn{
		cfg:   cfg,
		proto: p,
		tasks: tasks,

		in:  stdinw,
		out: stdoutr,

		containerIn:  stdinr,
		containerOut: stdoutw,

		rnr: rnr,
	}

	return hc, nil
}

// serve is the main loop of a hot function. It will keep running until the
// context is canceled. It will also stop if it goes idle for too long.
func (hc *HotFn) serve(ctx context.Context) {
	lctx, cancel := context.WithCancel(ctx)
	var wg sync.WaitGroup
	cfg := *hc.cfg
	logger := logrus.WithFields(logrus.Fields{
		"app":             cfg.AppName,
		"route":           cfg.Path,
		"image":           cfg.Image,
		"memory":          cfg.Memory,
		"format":          cfg.Format,
		"max_concurrency": cfg.MaxConcurrency,
		"idle_timeout":    cfg.IdleTimeout,
	})

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			inactivity := time.After(cfg.IdleTimeout)

			select {
			case <-lctx.Done():
				return

			case <-inactivity:
				logger.Info("Canceling inactive hot function")
				cancel()

			// if there comes a task
			case t := <-hc.tasks:
				// send stdin to http request to channel
				if err := hc.proto.Dispatch(lctx, t); err != nil {
					logrus.WithField("ctx", lctx).Info("task failed")
					t.Response <- task.Response{
						Result: &runResult{StatusValue: "error", error: err},
						Err:    err,
					}
					continue
				}

				t.Response <- task.Response{
					Result: &runResult{StatusValue: "success"},
					Err:    nil,
				}
			}
		}
	}()

	cfg.Env["FN_FORMAT"] = cfg.Format
	cfg.Timeout = 0 // add a timeout to simulate ab.end. failure.
	cfg.Stdin = hc.containerIn
	cfg.Stdout = hc.containerOut

	errr, errw := io.Pipe()
	cfg.Stderr = errw
	wg.Add(1)
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(errr)
		for scanner.Scan() {
			logger.Info(scanner.Text())
		}
	}()

	// run the hot function
	result, err := hc.rnr.Run(lctx, &cfg)
	if err != nil {
		logrus.WithError(err).Error("hot function failure detected")
	}
	errw.Close()
	wg.Wait()
	logrus.WithField("result", result).Info("hot function terminated")
}

func runTaskReq(rnr *Runner, wg *sync.WaitGroup, t task.Request) {
	defer wg.Done()
	rnr.Start()
	defer rnr.Complete()
	result, err := rnr.Run(t.Ctx, t.Config)
	select {
	case t.Response <- task.Response{Result: result, Err: err}:
		close(t.Response)
	default:
	}
}

type runResult struct {
	error
	StatusValue string
}

func (r *runResult) Error() string {
	if r.error == nil {
		return ""
	}
	return r.error.Error()
}

func (r *runResult) Status() string { return r.StatusValue }
