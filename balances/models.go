package balances

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/thrasher-/gocryptotrader/exchanges/poloniex"
)

const (
	TypeDeposit    = "deposit"
	TypeWithdrawal = "withdrawal"
)

// merge two struct
// - poloniex.PoloniexDepositsWithdrawals.Deposits
// - poloniex.PoloniexDepositsWithdrawals.Withdrawals
type Transaction struct {
	gorm.Model

	// deposit / withdrawal
	Type string

	WithdrawalNumber int64 // only withdrawals
	Currency         string
	Address          string
	Amount           float64
	Confirmations    int
	TransactionID    string `gorm:"unique"`
	Timestamp        time.Time
	Status           string
	IPAddress        string // only withdrawals
}

func NewTransactions(h poloniex.PoloniexDepositsWithdrawals) []Transaction {
	deposits := []Transaction{}
	for _, row := range h.Deposits {
		r := Transaction{
			Type:          TypeDeposit,
			Currency:      row.Currency,
			Address:       row.Address,
			Amount:        row.Amount,
			Confirmations: row.Confirmations,
			TransactionID: row.TransactionID,
			Timestamp:     time.Unix(row.Timestamp, 0),
			Status:        row.Status,
		}
		deposits = append(deposits, r)
	}

	withdrawals := []Transaction{}
	for _, row := range h.Withdrawals {
		r := Transaction{
			Type:             TypeWithdrawal,
			WithdrawalNumber: row.WithdrawalNumber,
			Currency:         row.Currency,
			Address:          row.Address,
			Amount:           row.Amount,
			Confirmations:    row.Confirmations,
			TransactionID:    row.TransactionID,
			Timestamp:        time.Unix(row.Timestamp, 0),
			Status:           row.Status,
			IPAddress:        row.IPAddress,
		}
		withdrawals = append(withdrawals, r)
	}

	rows := []Transaction{}
	rows = append(rows, deposits...)
	rows = append(rows, withdrawals...)
	return rows
}
