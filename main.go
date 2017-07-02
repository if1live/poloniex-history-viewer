package main

import (
	"path"
	"runtime"

	"github.com/if1live/poloniex-history-viewer/commands"
	"github.com/if1live/poloniex-history-viewer/yui"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	filename := "config.yaml"
	filepath := path.Join(yui.GetExecutablePath(), filename)
	s, err := yui.LoadSettings(filepath)
	if err != nil {
		yui.Check(err)
	}

	commands.Dispatch(s)
}
