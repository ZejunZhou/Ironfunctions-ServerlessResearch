package main

import (
	"os"

	"github.com/ZejunZhou/Ironfunctions-ServerlessResearch/fn/app"
)

func main() {
	fn := app.NewFn()
	fn.Run(os.Args)
}
