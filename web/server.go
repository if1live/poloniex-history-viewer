package web

import (
	"encoding/json"
	"net/http"
	"path"
	"strconv"
	"time"

	"github.com/if1live/poloniex-history-viewer/histories"
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
	http.HandleFunc("/poloniex-trade-history", handlerTradeHistories)
	http.HandleFunc("/js/", handlerJS)
	http.HandleFunc("/css/", handlerCSS)
	http.HandleFunc("/private.php", handlerPrivateAPI)
}

func (s *Server) Close() {
	s.db.Close()
}

func handlerJS(w http.ResponseWriter, r *http.Request) {
	targetPath := r.URL.Path[len("/js/"):]
	targetPath = path.Join("js", targetPath)
	renderStatic(w, r, targetPath)
}
func handlerCSS(w http.ResponseWriter, r *http.Request) {
	targetPath := r.URL.Path[len("/css/"):]
	targetPath = path.Join("css", targetPath)
	renderStatic(w, r, targetPath)
}

func handlerTradeHistories(w http.ResponseWriter, r *http.Request) {
	renderHtml(w, r, "trade_history.html")
}

func handler_returnPaginatedTradeHistory(w http.ResponseWriter, r *http.Request) {
	s := histories.NewAPI(svr.db.GetORM())

	start, _ := strconv.ParseInt(r.FormValue("start"), 10, 64)
	end, _ := strconv.ParseInt(r.FormValue("end"), 10, 64)
	page, _ := strconv.Atoi(r.FormValue("page"))
	tradesPerPage, _ := strconv.Atoi(r.FormValue("tradesPerPage"))
	typeval, _ := strconv.Atoi(r.FormValue("type"))

	startTime := time.Unix(start, 0)
	endTime := time.Unix(end, 0)

	rows := s.PaginateTradeHistory(startTime, endTime, page, tradesPerPage, typeval)

	w.Write([]byte("["))

	// first : atLeastOne
	if len(rows) > 0 {
		w.Write([]byte("1"))
	} else {
		w.Write([]byte("0"))
	}

	if len(rows) > 0 {
		w.Write([]byte(","))
		for i, r := range rows {
			data, _ := json.Marshal(r)
			w.Write(data)

			if i < len(rows)-1 {
				w.Write([]byte(","))
			}
		}
	}
	w.Write([]byte("]"))
}

func handler_returnPersonalTradeHistory(w http.ResponseWriter, r *http.Request) {
	s := histories.NewAPI(svr.db.GetORM())

	start, _ := strconv.ParseInt(r.FormValue("start"), 10, 64)
	end, _ := strconv.ParseInt(r.FormValue("end"), 10, 64)

	startTime := time.Unix(start, 0)
	endTime := time.Unix(end, 0)

	retval := s.PersonalTradeHistory(startTime, endTime)
	data, _ := json.Marshal(retval)
	w.Write(data)
}

func handlerPrivateAPI(w http.ResponseWriter, r *http.Request) {
	cmd := r.FormValue("command")

	if cmd == "returnPaginatedTradeHistory" {
		handler_returnPaginatedTradeHistory(w, r)
		return

	} else if cmd == "returnPersonalTradeHistory" {
		handler_returnPersonalTradeHistory(w, r)
		return
	}
}
