package web

import (
	"fmt"
	"net/http"
	"path"
	"strconv"

	"github.com/if1live/poloniex-history-viewer/balances"
	"github.com/if1live/poloniex-history-viewer/exchanges"
	"github.com/if1live/poloniex-history-viewer/histories"
	"github.com/if1live/poloniex-history-viewer/lendings"
	"github.com/if1live/poloniex-history-viewer/yui"
	"github.com/thrasher-/gocryptotrader/exchanges/poloniex"
)

/*
reference
http://www.alexedwards.net/blog/golang-response-snippets
*/

type Server struct {
	addr string
	port int

	settings yui.Settings
	db       *histories.Database
	api      *poloniex.Poloniex
}

var svr *Server

func NewServer(addr string, port int, s yui.Settings) *Server {
	if svr != nil {
		panic("already server exists!")
	}

	api := s.MakePoloniex()

	// share single orm
	db, err := histories.NewDatabase(s.DatabaseFileName)
	if err != nil {
		panic(err)
	}

	svr = &Server{
		addr:     addr,
		port:     port,
		settings: s,
		db:       &db,
		api:      api,
	}

	return svr
}

func (s *Server) Run() {
	http.HandleFunc("/tradeHistory", handlerTradeHistory)
	http.HandleFunc("/depositHistory", handlerDepositHistory)
	http.HandleFunc("/balances", handlerBalances)

	http.HandleFunc("/", handlerIndex)
	http.HandleFunc("/js/", handlerJS)
	http.HandleFunc("/css/", handlerCSS)
	http.HandleFunc("/private.php", handlerPrivateAPI)
	http.HandleFunc("/private", handlerPrivateAPI)

	http.HandleFunc("/sync/balance", handlerSyncBalance)
	http.HandleFunc("/sync/exchange", handlerSyncExchange)
	http.HandleFunc("/sync/lending", handlerSyncLending)

	http.HandleFunc("/static/", handlerStatic)

	addr := s.addr + ":" + strconv.Itoa(s.port)
	fmt.Println("run server on", addr)
	http.ListenAndServe(addr, nil)
}

func (s *Server) Close() {
	s.db.Close()
}

func handlerStatic(w http.ResponseWriter, r *http.Request) {
	targetPath := r.URL.Path[len("/static/"):]
	renderStatic(w, r, targetPath)
}

func handlerJS(w http.ResponseWriter, r *http.Request) {
	targetPath := r.URL.Path[len("/js/"):]
	targetPath = path.Join("js", targetPath)
	renderPoloniexStatic(w, r, targetPath)
}
func handlerCSS(w http.ResponseWriter, r *http.Request) {
	targetPath := r.URL.Path[len("/css/"):]
	targetPath = path.Join("css", targetPath)
	renderPoloniexStatic(w, r, targetPath)
}

func handlerTradeHistory(w http.ResponseWriter, r *http.Request) {
	renderPoloniexStatic(w, r, "trade_history.html")
}

func handlerDepositHistory(w http.ResponseWriter, r *http.Request) {
	renderPoloniexStatic(w, r, "deposit_history.html")
}

func handlerBalances(w http.ResponseWriter, r *http.Request) {
	renderPoloniexStatic(w, r, "balances.html")
}

func handlerIndex(w http.ResponseWriter, r *http.Request) {
	type Context struct {
		ExchangeSync *exchanges.Sync
		LendingSync  *lendings.Sync
		BalanceSync  *balances.Sync
	}
	ctx := Context{
		ExchangeSync: svr.db.MakeExchangeSync(nil),
		LendingSync:  svr.db.MakeLendingSync(nil),
		BalanceSync:  svr.db.MakeBalanceSync(nil),
	}

	err := renderTemplate(w, "index.html", ctx)
	if err != nil {
		renderErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

}
