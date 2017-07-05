package histories

import (
	"fmt"
	"strconv"
	"time"

	"github.com/if1live/poloniex-history-viewer/balances"
	"github.com/if1live/poloniex-history-viewer/exchanges"
	"github.com/if1live/poloniex-history-viewer/lendings"
	"github.com/if1live/poloniex-history-viewer/yui"
	"github.com/jinzhu/gorm"
)

type API struct {
	db *gorm.DB
}

func NewAPI(db *gorm.DB) *API {
	return &API{
		db: db,
	}
}

const (
	tradeHistoryTypeSellsOnly    = 0
	tradeHistoryTypeBuysOnly     = 1
	tradeHistoryTypeBuysAndSells = 2
	tradeHistoryTypeLoadEarnings = 3
)

// https://poloniex.com/private.php
// command=returnPaginatedTradeHistory
// start=0
// end=1917895987
// page=1
// tradesPerPage=50
// type=0
func (s *API) PaginateTradeHistory(start, end time.Time, page, tradesPerPage int, apiType int) []TradeHistory {
	limit := tradesPerPage
	offset := 0

	if page > 1 {
		offset = (page-1)*(tradesPerPage) - 1
		limit = tradesPerPage + 1
	}

	if apiType == tradeHistoryTypeLoadEarnings {
		var rows []lendings.Lending
		q := s.db.Where("close between ? and ?", start, end)
		q = q.Order("close desc")
		q = q.Limit(limit).Offset(offset).Find(&rows)

		histories := make([]TradeHistory, len(rows))
		for i, r := range rows {
			histories[i] = NewTradeHistoryFromLending(&r)
		}
		return histories

	} else {
		var rows []exchanges.Exchange

		q := s.db.Where("date between ? and ?", start, end)
		switch apiType {
		case tradeHistoryTypeBuysOnly:
			q = q.Where("type = ?", exchanges.ExchangeBuy)
		case tradeHistoryTypeSellsOnly:
			q = q.Where("type = ?", exchanges.ExchangeSell)
		case tradeHistoryTypeBuysAndSells:
		}
		q = q.Order("date desc").Find(&rows)
		q = q.Limit(limit).Offset(offset).Find(&rows)

		histories := make([]TradeHistory, len(rows))
		for i, r := range rows {
			histories[i] = NewTradeHistoryFromExchange(&r)
		}
		return histories
	}
}

// https://poloniex.com/private.php
// command=returnNumberOfPagesInTradeHistory
// start=0
// end=1917895987
// tradesPerPage=50
// type=2
func (s *API) NumberOfPagesInTradeHistory(start, end time.Time, tradesPerPage int, apiType int) int {
	rowcount := 0
	if apiType == tradeHistoryTypeLoadEarnings {
		q := s.db.Model(&lendings.Lending{}).Where("close between ? and ?", start, end)
		q.Order("close desc").Count(&rowcount)

	} else {
		q := s.db.Model(&exchanges.Exchange{}).Where("date between ? and ?", start, end)
		switch apiType {
		case tradeHistoryTypeBuysOnly:
			q = q.Where("type = ?", exchanges.ExchangeBuy)
		case tradeHistoryTypeSellsOnly:
			q = q.Where("type = ?", exchanges.ExchangeSell)
		case tradeHistoryTypeBuysAndSells:
		}
		q.Order("date desc").Count(&rowcount)
	}

	if rowcount == 0 {
		return 0
	}
	page := ((rowcount - 1) / tradesPerPage) + 1
	return page
}

// https://poloniex.com/private.php
// command=returnPersonalTradeHistory
// start=0
// end=1917895987
// retval: key=BTC, BTC_AMP, value=array
func (s *API) PersonalTradeHistory(start, end time.Time) map[string][]PersonalTradeHistory {
	var exchangerows []exchanges.Exchange
	s.db.Where("date between ? and ?", start, end).Order("date desc").Find(&exchangerows)
	exchangehistories := make([]PersonalTradeHistory, len(exchangerows))
	for i, r := range exchangerows {
		h := NewPersonalTradeHistoryFromExchange(&r)
		exchangehistories[i] = h
	}

	// tradeId - currencyPair
	tradeIDMap := map[string]string{}
	for _, r := range exchangerows {
		key := strconv.FormatInt(r.TradeID, 10)
		tradeIDMap[key] = r.CurrencyPair()
	}

	var lendingrows []lendings.Lending
	s.db.Where("close between ? and ?", start, end).Order("close desc").Find(&lendingrows)
	lendinghistories := make([]PersonalTradeHistory, len(lendingrows))
	for i, r := range lendingrows {
		h := NewPersonalTradeHistoryFromLending(&r)
		lendinghistories[i] = h
	}

	// lending id - currency
	lendingIDMap := map[string]string{}
	for _, r := range lendingrows {
		key := convertLendingIDtoTradeID(r.LendingID)
		lendingIDMap[key] = r.Currency
	}

	retval := map[string][]PersonalTradeHistory{}
	for _, h := range exchangehistories {
		currencyPair := tradeIDMap[h.TradeID]
		list, ok := retval[currencyPair]
		if ok {
			retval[currencyPair] = append(list, h)
		} else {
			retval[currencyPair] = []PersonalTradeHistory{h}
		}
	}

	for _, h := range lendinghistories {
		currency := lendingIDMap[h.TradeID]
		list, ok := retval[currency]
		if ok {
			retval[currency] = append(list, h)
		} else {
			retval[currency] = []PersonalTradeHistory{h}
		}
	}
	return retval
}

// https://poloniex.com/
// private
// command=returnDepositsAndWithdrawalsMobile
func (s *API) DepositsAndWithdrawals() {
}

// https://poloniex.com/
// private
// command=returnWithdrawalsDeposits
// limit=50
func (s *API) WithdrawalsDeposits(limit int) WithdrawalDepositHistory {
	var rows []balances.Transaction
	s.db.Order("timestamp desc").Find(&rows)

	deposits := []DepositHistory{}
	withdrawals := []WithdrawalHistory{}
	for _, r := range rows {
		if r.Type == balances.TypeDeposit {
			deposits = append(deposits, NewDepositHistory(&r))
		} else if r.Type == balances.TypeWithdrawal {
			withdrawals = append(withdrawals, NewWithdrawalHistory(&r))
		}
	}

	limitWithdraw := "2000.00000000"
	remaining := "1234.12345678"
	return WithdrawalDepositHistory{
		Deposits:    deposits,
		Withdrawals: withdrawals,
		Limit:       limitWithdraw,
		Remaining:   remaining,
	}
}

type TradeHistory struct {
	CurrencyPair string `json:"currencyPair"`
	Date         string `json:"date"`
	Type         string `json:"type"`
	Category     string `json:"category"`
	Rate         string `json:"rate"`
	Amount       string `json:"amount"`
	Total        string `json:"total"`
	Fee          string `json:"fee"`
}

func NewTradeHistoryFromExchange(r *exchanges.Exchange) TradeHistory {
	feeAmountStr := yui.ToFloatStr(r.FeeAmount())
	feeUnit := ""
	if r.Type == exchanges.ExchangeBuy {
		feeUnit = r.Asset
	} else if r.Type == exchanges.ExchangeSell {
		feeUnit = r.Currency
	}
	feePercent := r.Fee * 100
	fee := fmt.Sprintf("%s %s (%.2f%%)", feeAmountStr, feeUnit, feePercent)

	return TradeHistory{
		CurrencyPair: r.CurrencyPair(),
		Type:         r.Type,
		Category:     r.Category,
		Amount:       yui.ToFloatStr(r.Amount),
		Rate:         yui.ToFloatStr(r.Rate),
		Date:         r.DateStr(),
		Total:        yui.ToFloatStr(r.MyTotal()) + " " + r.Currency,
		Fee:          fee,
	}
}
func NewTradeHistoryFromLending(r *lendings.Lending) TradeHistory {
	return TradeHistory{
		CurrencyPair: r.Currency,
		Category:     "lendingEarning",
		Type:         "1",
		Date:         r.Close.Format("2006-01-02 15:04:05"),
		Amount:       yui.ToFloatStr(r.Amount),
		Fee:          fmt.Sprintf("%.0f%%", r.FeeRate()*100),
		Rate:         fmt.Sprintf("%.4f%%", r.Rate*100),
		Total:        yui.ToFloatStr(r.Interest),
	}
}

type DepositHistory struct {
	Currency      string `json:"currency"`
	Address       string `json:"address"`
	Amount        string `json:"amount"`
	Confirmations int    `json:"confirmations"`
	Txid          string `json:"txid"`
	Timestamp     int64  `json:"timestamp"`
	Status        string `json:"status"`
}

func NewDepositHistory(r *balances.Transaction) DepositHistory {
	return DepositHistory{
		Currency:      r.Currency,
		Address:       r.Address,
		Amount:        yui.ToFloatStr(r.Amount),
		Confirmations: r.Confirmations,
		Txid:          r.TransactionID,
		Timestamp:     r.Timestamp.Unix(),
		Status:        r.Status,
	}
}

type WithdrawalHistory struct {
	WithdrawalNumber int64  `json:"withdrawalNumber"`
	Currency         string `json:"currency"`
	Address          string `json:"address"`
	Amount           string `json:"amount"`
	Fee              string `json:"fee"`
	Timestamp        int64  `json:"timestamp"`
	Status           string `json:"status"`
	IPAddress        string `json:"ipAddress"`
}

func NewWithdrawalHistory(r *balances.Transaction) WithdrawalHistory {
	// currency-fee table
	feeTable := map[string]float64{
		"BTC": 0.0001,
	}
	fee, ok := feeTable[r.Currency]
	if !ok {
		fee = -1
	}

	return WithdrawalHistory{
		WithdrawalNumber: r.WithdrawalNumber,
		Currency:         r.Currency,
		Address:          r.Address,
		Amount:           yui.ToFloatStr(r.Amount),
		Fee:              yui.ToFloatStr(fee),
		Timestamp:        r.Timestamp.Unix(),
		Status:           r.Status,
		IPAddress:        r.IPAddress,
	}
}

type WithdrawalDepositHistory struct {
	Deposits    []DepositHistory    `json:"deposits"`
	Withdrawals []WithdrawalHistory `json:"withdrawals"`
	Limit       string              `json:"limit"`
	Remaining   string              `json:"remaining"`
}

type PersonalTradeHistory struct {
	Amount        string `json:"amount"`
	Category      string `json:"category"`
	Date          string `json:"date"`
	Fee           string `json:"fee"`
	GlobalTradeID string `json:"globalTradeID,int"`
	OrderNumber   string `json:"orderNumber"`
	Rate          string `json:"rate"`
	Total         string `json:"total"`
	TradeID       string `json:"tradeID"`
	Type          string `json:"type"`
}

func NewPersonalTradeHistoryFromExchange(r *exchanges.Exchange) PersonalTradeHistory {
	return PersonalTradeHistory{
		Amount:        yui.ToFloatStr(r.Amount),
		Category:      r.Category,
		Date:          r.DateStr(),
		Fee:           yui.ToFloatStr(r.Fee),
		GlobalTradeID: strconv.FormatInt(r.GlobalTradeID, 10),
		OrderNumber:   strconv.FormatInt(r.OrderNumber, 10),
		Rate:          yui.ToFloatStr(r.Rate),
		Total:         yui.ToFloatStr(r.Total),
		TradeID:       strconv.FormatInt(r.TradeID, 10),
		Type:          r.Type,
	}
}

func NewPersonalTradeHistoryFromLending(r *lendings.Lending) PersonalTradeHistory {
	return PersonalTradeHistory{
		Amount:        yui.ToFloatStr(r.Amount),
		Category:      "lendingEarning",
		Date:          r.Close.Format("2006-01-02 15:04:05"),
		Fee:           yui.ToFloatStr(r.FeeRate()),
		GlobalTradeID: "",
		OrderNumber:   "",
		Rate:          yui.ToFloatStr(r.Rate),
		Total:         yui.ToFloatStr(r.Interest),
		TradeID:       convertLendingIDtoTradeID(r.LendingID),
		Type:          "buy",
	}
}

func convertLendingIDtoTradeID(lendingID int64) string {
	return "s" + strconv.FormatInt(lendingID, 10)
}
