package exchanges

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/deckarep/golang-set"
	"github.com/jinzhu/gorm"
	"github.com/thrasher-/gocryptotrader/exchanges/poloniex"
)

type Sync struct {
	db  *gorm.DB
	api *poloniex.Poloniex
}

func NewSync(db *gorm.DB, api *poloniex.Poloniex) *Sync {
	return &Sync{
		db:  db,
		api: api,
	}
}

// currencyPair example : "all", "BTC_DOGE"
func (sync *Sync) Sync(currencyPair string, start, end time.Time) (int, error) {
	startTime := strconv.FormatInt(start.Unix(), 10)
	endTime := strconv.FormatInt(end.Unix(), 10)
	retval, err := sync.api.GetAuthenticatedTradeHistory(currencyPair, startTime, endTime)

	if err != nil {
		return -1, err
	}

	if all, ok := retval.(poloniex.PoloniexAuthenticatedTradeHistoryAll); ok {
		rowcount := 0
		for key, histories := range all.Data {
			// key example : BTC_DOGE
			tokens := strings.Split(key, "_")
			asset := tokens[1]
			currency := tokens[0]
			rowcount += sync.insertHistories(asset, currency, histories)
		}
		return rowcount, nil
	}

	if resp, ok := retval.(poloniex.PoloniexAuthenticatedTradeHistoryResponse); ok {
		tokens := strings.Split(currencyPair, "_")
		rowcount := sync.insertHistories(tokens[1], tokens[0], resp.Data)
		return rowcount, nil
	}

	return -1, errors.New("unknown error : load from exchange")
}

func (sync *Sync) SyncAll() (int, error) {
	start := time.Unix(0, 0)
	end := time.Now()
	return sync.Sync("all", start, end)
}

func (sync *Sync) SyncRecent() (int, error) {
	start := sync.GetLastTime()
	end := time.Now()
	return sync.Sync("all", start, end)
}

func (sync *Sync) GetLastTime() time.Time {
	var last Exchange
	sync.db.Order("date desc").First(&last)
	if last.ID == 0 {
		return time.Unix(0, 0)
	}
	return last.Date
}

func (sync *Sync) insertHistories(asset, currency string, histories []poloniex.PoloniexAuthentictedTradeHistory) int {
	var existRows []Exchange
	sync.db.Where(&Exchange{Asset: asset, Currency: currency}).Select("trade_id").Find(&existRows)
	tradeIDSet := mapset.NewSet()
	for _, r := range existRows {
		tradeIDSet.Add(r.TradeID)
	}

	rows := []Exchange{}
	for _, history := range histories {
		if tradeIDSet.Contains(history.TradeID) {
			continue
		}
		row := NewExchange(asset, currency, history)
		rows = append(rows, row)
	}

	for _, row := range rows {
		sync.db.Create(&row)
	}
	return len(rows)
}
