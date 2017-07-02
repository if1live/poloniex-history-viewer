package balances

import (
	"strconv"
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

func (sync *Sync) Sync(start, end time.Time) (int, error) {
	startTime := strconv.FormatInt(start.Unix(), 10)
	endTime := strconv.FormatInt(end.Unix(), 10)
	retval, err := sync.api.GetDepositsWithdrawals(startTime, endTime)
	if err != nil {
		return -1, err
	}

	var existRows []Transaction
	sync.db.Select("transaction_id").Find(&existRows)
	idSet := mapset.NewSet()
	for _, r := range existRows {
		idSet.Add(r.TransactionID)
	}

	rows := []Transaction{}
	retvals := NewTransactions(retval)
	for _, history := range retvals {
		if idSet.Contains(history.TransactionID) {
			continue
		}
		rows = append(rows, history)
	}
	for _, row := range rows {
		sync.db.Create(&row)
	}
	return len(rows), nil
}

func (sync *Sync) SyncAll() (int, error) {
	start := time.Unix(0, 0)
	end := time.Now()
	return sync.Sync(start, end)
}

func (sync *Sync) SyncRecent() (int, error) {
	start := sync.GetLastTime()
	end := time.Now()
	return sync.Sync(start, end)
}
func (sync *Sync) GetLastTime() time.Time {
	var last Transaction
	sync.db.Order("timestamp desc").First(&last)
	if last.ID == 0 {
		return time.Unix(0, 0)
	}
	return last.Timestamp
}
