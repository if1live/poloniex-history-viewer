package commands

import (
	"flag"
	"fmt"
	"os"
	"reflect"

	"github.com/if1live/poloniex-history-viewer/histories"
	"github.com/if1live/poloniex-history-viewer/web"
	"github.com/if1live/poloniex-history-viewer/yui"
)

var command string

// for server
var port int
var host string

func init() {
	flag.StringVar(&command, "cmd", "", "command, svr or sync")
	flag.IntVar(&port, "port", 8000, "port to use")
	flag.StringVar(&host, "host", "127.0.0.1", "address to use")
}

func Dispatch(settings yui.Settings) {
	flag.Parse()

	var err error
	switch command {
	case "svr":
		err = executeSvr(settings)
	case "sync":
		err = executeSync(settings)
	default:
		fmt.Println("unknown command : ", command)
		os.Exit(-1)
	}
	if err != nil {
		panic(err)
	}
}

func executeSvr(s yui.Settings) error {
	svr := web.NewServer(host, port, s)
	svr.Run()
	svr.Close()
	return nil
}

func executeSync(s yui.Settings) error {
	db, err := histories.NewDatabase(s.DatabaseFileName)
	if err != nil {
		return err
	}
	defer db.Close()

	api := s.MakePoloniex()
	syncs := []histories.Synchronizer{
		db.MakeExchangeSync(api),
		db.MakeLendingSync(api),
		db.MakeBalanceSync(api),
	}
	for _, sync := range syncs {
		rowcount, err := sync.SyncRecent()
		if err != nil {
			return err
		}

		syncName := reflect.TypeOf(sync).String()
		// TODO use logger
		fmt.Printf("%s : %d exchange rows added\n", syncName, rowcount)
	}

	return nil

}
