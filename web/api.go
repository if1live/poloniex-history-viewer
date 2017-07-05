package web

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/if1live/poloniex-history-viewer/histories"
)

func handler_returnDepositsAndWithdrawalsMobile(w http.ResponseWriter, r *http.Request) {
	//s := histories.NewAPI(svr.db.GetORM())
	data, _ := json.Marshal(10)
	w.Write(data)
}

func handler_returnWithdrawalsDeposits(w http.ResponseWriter, r *http.Request) {
	s := histories.NewAPI(svr.db.GetORM())
	limit, _ := strconv.Atoi(r.FormValue("limit"))

	result := s.WithdrawalsDeposits(limit)
	data, _ := json.Marshal(result)
	w.Write(data)
}

func handler_returnNumberOfPagesInTradeHistory(w http.ResponseWriter, r *http.Request) {
	s := histories.NewAPI(svr.db.GetORM())

	start, _ := strconv.ParseInt(r.FormValue("start"), 10, 64)
	end, _ := strconv.ParseInt(r.FormValue("end"), 10, 64)
	tradesPerPage, _ := strconv.Atoi(r.FormValue("tradesPerPage"))
	typeval, _ := strconv.Atoi(r.FormValue("type"))

	startTime := time.Unix(start, 0)
	endTime := time.Unix(end, 0)

	rowcount := s.NumberOfPagesInTradeHistory(startTime, endTime, tradesPerPage, typeval)

	data, _ := json.Marshal(rowcount)
	w.Write(data)
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

	// first : atLeastOne
	atLeastOne := ""
	if len(rows) > 0 {
		atLeastOne = "1"
	} else {
		atLeastOne = "0"
	}

	v := []interface{}{atLeastOne}
	for _, r := range rows {
		v = append(v, r)
	}
	data, _ := json.Marshal(v)
	w.Write(data)
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

	} else if cmd == "returnNumberOfPagesInTradeHistory" {
		handler_returnNumberOfPagesInTradeHistory(w, r)
		return

	} else if cmd == "returnWithdrawalsDeposits" {
		handler_returnWithdrawalsDeposits(w, r)
		return

	} else if cmd == "returnDepositsAndWithdrawalsMobile" {
		handler_returnDepositsAndWithdrawalsMobile(w, r)
		return
	}
}
