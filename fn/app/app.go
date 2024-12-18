package app

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	vers "github.com/ZejunZhou/Ironfunctions-ServerlessResearch/api/version"
	"github.com/ZejunZhou/Ironfunctions-ServerlessResearch/fn/commands"
	image_commands "github.com/ZejunZhou/Ironfunctions-ServerlessResearch/fn/commands/images"
	"github.com/ZejunZhou/Ironfunctions-ServerlessResearch/fn/common"
	"github.com/urfave/cli"
)

var aliases = map[string]cli.Command{
	"build":  image_commands.Build(),
	"bump":   image_commands.Bump(),
	"deploy": image_commands.Deploy(),
	"push":   image_commands.Push(),
	"run":    image_commands.Run(),
	"call":   commands.Call(),
}

func aliasesFn() []cli.Command {
	cmds := []cli.Command{}
	for alias, cmd := range aliases {
		cmd.Name = alias
		cmd.Hidden = true
		cmds = append(cmds, cmd)
	}
	return cmds
}

func NewFn() *cli.App {
	common.SetEnv()
	app := cli.NewApp()
	app.Name = "fn"
	app.Version = vers.Version
	app.Authors = []cli.Author{{Name: "iron.io"}}
	app.Description = "IronFunctions command line tools"
	app.UsageText = `Check the manual at https://github.com/iron-io/functions/blob/master/fn/README.md`

	cli.AppHelpTemplate = `{{.Name}} {{.Version}}{{if .Description}}

{{.Description}}{{end}}

USAGE:
   {{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}{{if .Commands}} command [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{end}}

ENVIRONMENT VARIABLES:
   API_URL - IronFunctions remote API address{{if .VisibleCommands}}

COMMANDS:{{range .VisibleCategories}}{{if .Name}}
   {{.Name}}:{{end}}{{range .VisibleCommands}}
     {{join .Names ", "}}{{"\t"}}{{.Usage}}{{end}}{{end}}{{end}}{{if .VisibleFlags}}

ALIASES:
     build    (images build)
     bump     (images bump)
     deploy   (images deploy)
     run      (images run)
     call     (routes call)
     push     (images push)

GLOBAL OPTIONS:
   {{range $index, $option := .VisibleFlags}}{{if $index}}
   {{end}}{{$option}}{{end}}{{end}}
`

	app.CommandNotFound = func(c *cli.Context, cmd string) {
		fmt.Fprintf(os.Stderr, "command not found: %v\n", cmd)
	}
	app.Commands = []cli.Command{
		commands.InitFn(),
		commands.Apps(),
		commands.Routes(),
		commands.Images(),
		commands.Lambda(),
		commands.Version(),
	}
	app.Commands = append(app.Commands, aliasesFn()...)

	prepareCmdArgsValidation(app.Commands)

	return app
}

func parseArgs(c *cli.Context) ([]string, []string) {
	args := strings.Split(c.Command.ArgsUsage, " ")
	var reqArgs []string
	var optArgs []string
	for _, arg := range args {
		if strings.HasPrefix(arg, "[") {
			optArgs = append(optArgs, arg)
		} else if strings.Trim(arg, " ") != "" {
			reqArgs = append(reqArgs, arg)
		}
	}
	return reqArgs, optArgs
}

func prepareCmdArgsValidation(cmds []cli.Command) {
	// TODO: refactor fn to use urfave/cli.v2
	// v1 doesn't let us validate args before the cmd.Action

	for i, cmd := range cmds {
		prepareCmdArgsValidation(cmd.Subcommands)
		if cmd.Action == nil {
			continue
		}
		action := cmd.Action
		cmd.Action = func(c *cli.Context) error {
			reqArgs, _ := parseArgs(c)
			if c.NArg() < len(reqArgs) {
				var help bytes.Buffer
				cli.HelpPrinter(&help, cli.CommandHelpTemplate, c.Command)
				return fmt.Errorf("ERROR: Missing required arguments: %s\n\n%s", strings.Join(reqArgs[c.NArg():], " "), help.String())
			}
			return cli.HandleAction(action, c)
		}
		cmds[i] = cmd
	}
}
