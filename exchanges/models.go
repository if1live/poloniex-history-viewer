package exchanges

import (
	"time"

	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/thrasher-/gocryptotrader/exchanges/poloniex"
)

const (
	ExchangeSell = "sell"
	ExchangeBuy  = "buy"
)

func convertPoloniexDate(val string) time.Time {
	// date example : 2017-06-18 04:31:08
	t, _ := time.Parse("2006-01-02 15:04:05", val)
	return t
}

// PoloniexAuthenticatedTradeHistory
type Exchange struct {
	gorm.Model

	// from zenbot format
	// {exchange_slug}.{asset}-{currency}
	// poloniex.AMP-BTC
	Asset    string
	Currency string

	GlobalTradeID int64 `gorm:"unique"`
	TradeID       int64
	Date          time.Time
	Rate          float64
	Amount        float64
	Total         float64
	Fee           float64
	OrderNumber   int64
	Type          string
	Category      string
}

func NewExchange(asset, currency string, h poloniex.PoloniexAuthentictedTradeHistory) Exchange {
	return Exchange{
		Asset:    asset,
		Currency: currency,

		GlobalTradeID: h.GlobalTradeID,
		TradeID:       h.TradeID,
		Date:          convertPoloniexDate(h.Date),
		Rate:          h.Rate,
		Amount:        h.Amount,
		Total:         h.Total,
		Fee:           h.Fee,
		OrderNumber:   h.OrderNumber,
		Type:          h.Type,
		Category:      h.Category,
	}
}

func (r *Exchange) MakeHistory() poloniex.PoloniexAuthentictedTradeHistory {
	return poloniex.PoloniexAuthentictedTradeHistory{
		GlobalTradeID: r.GlobalTradeID,
		TradeID:       r.TradeID,
		Date:          r.Date.Format(time.RFC3339),
		Rate:          r.Rate,
		Amount:        r.Amount,
		Total:         r.Total,
		Fee:           r.Fee,
		OrderNumber:   r.OrderNumber,
		Type:          r.Type,
		Category:      r.Category,
	}
}

func (r *Exchange) FeeAmount() float64 {
	switch r.Type {
	case ExchangeBuy:
		return r.buyFeeAmount()
	case ExchangeSell:
		return r.sellFeeAmount()
	}
	return -1
}

// sell example
// rate: 0.00007900
// amount : 137.43455498
// fee : 0.00001629 BTC (0.15%)
// total in db : 0.01085732
// total in poloniex : 0.01084103 BTC
// 0.00001629 BTC = 0.01085732 BTC * (0.01) * (0.15)
// fee amount = (total in db) * fee
func (r *Exchange) sellFeeAmount() float64 {
	return r.Total * r.Fee
}

// buy example
// amount : 13.00373802
// fee : 0.03250935 SYS (0.25%)
// total : 0.00094368 BTC
// 0.03250935 SYS = 13.00373802 SYS * (0.01) * (0.25)
// fee amount = amount * fee
func (r *Exchange) buyFeeAmount() float64 {
	return r.Amount * r.Fee
}

func (r *Exchange) MyTotal() float64 {
	switch r.Type {
	case ExchangeBuy:
		return r.Total
	case ExchangeSell:
		return r.Total - r.FeeAmount()
	}
	return -1
}
func (r *Exchange) MyAmount() float64 {
	switch r.Type {
	case ExchangeBuy:
		return r.Amount - r.FeeAmount()
	case ExchangeSell:
		return r.Amount
	}
	return -1
}

func (r *Exchange) DateStr() string {
	return r.Date.Format("2006-01-02 15:04:05")
}

func (r *Exchange) CurrencyPair() string {
	return fmt.Sprintf("%s_%s", r.Currency, r.Asset)
}
