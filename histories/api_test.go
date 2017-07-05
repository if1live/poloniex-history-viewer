package histories

import (
	"testing"
	"time"

	"github.com/if1live/poloniex-history-viewer/balances"
	"github.com/if1live/poloniex-history-viewer/exchanges"
	"github.com/if1live/poloniex-history-viewer/lendings"
	"github.com/stretchr/testify/assert"
)

func Test_NewTradeHistoryFromExchange_Sell(t *testing.T) {
	t.Parallel()

	expected := TradeHistory{
		CurrencyPair: "BTC_POT",
		Date:         "2017-06-26 15:18:39",
		Type:         "sell",
		Category:     "exchange",
		Rate:         "0.00004882",
		Amount:       "201.39169020",
		Total:        "0.00981719 BTC",
		Fee:          "0.00001475 BTC (0.15%)",
	}

	//  "2017-06-26 15:18:39+00:00"
	date, err := time.Parse("2006-01-02 15:04:05", "2017-06-26 15:18:39")
	if err != nil {
		t.Fatalf("time.Parse: %v", err)
	}

	r := exchanges.Exchange{
		Asset:         "POT",
		Currency:      "BTC",
		GlobalTradeID: 177732676,
		TradeID:       1519576,
		Date:          date,
		Rate:          4.882e-05,
		Amount:        201.3916902,
		Total:         0.00983194,
		Fee:           0.0015,
		OrderNumber:   12943626491,
		Type:          "sell",
		Category:      "exchange",
	}
	v := NewTradeHistoryFromExchange(&r)
	assert.Equal(t, expected, v)
}

func Test_NewTradeHistoryFromExchange_Buy(t *testing.T) {
	t.Parallel()

	expected := TradeHistory{
		CurrencyPair: "BTC_SC",
		Date:         "2017-06-23 16:22:30",
		Type:         "buy",
		Category:     "exchange",
		Rate:         "0.00000683",
		Amount:       "732.06442166",
		Total:        "0.00499999 BTC",
		Fee:          "1.83016105 SC (0.25%)",
	}

	date, err := time.Parse("2006-01-02 15:04:05", "2017-06-23 16:22:30")
	if err != nil {
		t.Fatalf("time.Parse: %v", err)
	}

	r := exchanges.Exchange{
		Asset:         "SC",
		Currency:      "BTC",
		GlobalTradeID: 175312662,
		TradeID:       2613513,
		Date:          date,
		Rate:          6.83e-06,
		Amount:        732.06442166,
		Total:         0.00499999,
		Fee:           0.0025,
		OrderNumber:   9175049916,
		Type:          "buy",
		Category:      "exchange",
	}
	v := NewTradeHistoryFromExchange(&r)
	assert.Equal(t, expected, v)
}

func Test_NewTradeHistoryFromLending(t *testing.T) {
	t.Parallel()

	expected := TradeHistory{
		Amount:       "0.10000000",
		Category:     "lendingEarning",
		CurrencyPair: "BTC",
		Date:         "2017-06-13 11:46:18",
		Fee:          "15%",
		Rate:         "0.1000%",
		Total:        "0.00000018",
		// use int 1
		Type: "1",
	}
	open, _ := time.Parse("2006-01-02 15:04:05", "2017-06-13 11:43:49")
	close, _ := time.Parse("2006-01-02 15:04:05", "2017-06-13 11:46:18")
	r := lendings.Lending{
		LendingID: 374173954,
		Currency:  "BTC",
		Rate:      0.001,
		Amount:    0.1,
		Duration:  0.0017,
		Interest:  1.8e-07,
		Fee:       -2e-08,
		Earned:    1.6e-07,
		Open:      open,
		Close:     close,
	}
	v := NewTradeHistoryFromLending(&r)
	assert.Equal(t, expected, v)
}

func Test_NewPersonTradeHistoryFromExchange(t *testing.T) {
	t.Parallel()

	expected := PersonalTradeHistory{
		Amount:        "1.86348482",
		Category:      "exchange",
		Date:          "2017-06-18 01:10:27",
		Fee:           "0.00150000",
		GlobalTradeID: "169934831",
		OrderNumber:   "20304929853",
		Rate:          "0.00536629",
		Total:         "0.00999999",
		TradeID:       "762087",
		Type:          "buy",
	}

	date, err := time.Parse("2006-01-02 15:04:05", "2017-06-18 01:10:27")
	if err != nil {
		t.Fatalf("time.Parse: %v", err)
	}

	r := exchanges.Exchange{
		Asset:         "XCP",
		Currency:      "BTC",
		GlobalTradeID: 169934831,
		TradeID:       762087,
		Date:          date,
		Rate:          0.00536629,
		Amount:        1.86348482,
		Total:         0.00999999,
		Fee:           0.0015,
		OrderNumber:   20304929853,
		Type:          "buy",
		Category:      "exchange",
	}
	v := NewPersonalTradeHistoryFromExchange(&r)
	assert.Equal(t, expected, v)
}

func Test_NewPersonalTradeHistoryFromLending(t *testing.T) {
	t.Parallel()

	expected := PersonalTradeHistory{
		Amount:   "0.10000000",
		Category: "lendingEarning",
		Date:     "2017-06-13 11:46:18",
		Fee:      "0.15000000",
		// null
		GlobalTradeID: "",
		OrderNumber:   "",
		Rate:          "0.00100000",
		Total:         "0.00000018",
		TradeID:       "s374173954",
		Type:          "buy",
	}

	open, _ := time.Parse("2006-01-02 15:04:05", "2017-06-13 11:43:49")
	close, _ := time.Parse("2006-01-02 15:04:05", "2017-06-13 11:46:18")
	r := lendings.Lending{
		LendingID: 374173954,
		Currency:  "BTC",
		Rate:      0.001,
		Amount:    0.1,
		Duration:  0.0017,
		Interest:  1.8e-07,
		Fee:       -2e-08,
		Earned:    1.6e-07,
		Open:      open,
		Close:     close,
	}
	v := NewPersonalTradeHistoryFromLending(&r)
	assert.Equal(t, expected, v)

}

func Test_NewWithdrawalHistory(t *testing.T) {
	t.Parallel()

	expected := WithdrawalHistory{
		WithdrawalNumber: 5485728,
		Currency:         "BTC",
		Address:          "1FfeRCpZNt8bskfqS7DCULE1JAZC32w9WN",
		Amount:           "0.01425148",
		Fee:              "0.00010000",
		Timestamp:        1499090630,
		Status:           "COMPLETE: 9010da9ee9eb67f0dd632fac6a768156aba1140e7d45e0dfee202e65b3651f92",
		IPAddress:        "1.2.3.4",
	}

	timestamp, err := time.Parse("2006-01-02 15:04:05-07:00", "2017-07-03 23:03:50+09:00")
	if err != nil {
		t.Errorf("time.Parse: %v", err)
	}

	r := balances.Transaction{
		Type:             "withdrawal",
		WithdrawalNumber: 5485728,
		Currency:         "BTC",
		Address:          "1FfeRCpZNt8bskfqS7DCULE1JAZC32w9WN",
		Amount:           0.01425148,
		Confirmations:    0,
		TransactionID:    "",
		Timestamp:        timestamp,
		Status:           "COMPLETE: 9010da9ee9eb67f0dd632fac6a768156aba1140e7d45e0dfee202e65b3651f92",
		IPAddress:        "1.2.3.4",
	}
	v := NewWithdrawalHistory(&r)

	assert.Equal(t, expected, v)
}

func Test_NewDepositHistory(t *testing.T) {
	t.Parallel()

	expected := DepositHistory{
		Currency:      "BTC",
		Address:       "13ie6ptXxdwpNkSTwowXzX3AsHUvKRic8v",
		Amount:        "0.30732863",
		Confirmations: 1,
		Txid:          "f2fd5472dc596c3e6f6b377f67a0cade5e65e491758907b937aa3c5c50055aed",
		Timestamp:     1497314214,
		Status:        "COMPLETE",
	}

	timestamp, err := time.Parse("2006-01-02 15:04:05-07:00", "2017-06-13 09:36:54+09:00")
	if err != nil {
		t.Errorf("time.Parse: %v", err)
	}

	r := balances.Transaction{
		Type:             "deposit",
		WithdrawalNumber: 0,
		Currency:         "BTC",
		Address:          "13ie6ptXxdwpNkSTwowXzX3AsHUvKRic8v",
		Amount:           0.30732863,
		Confirmations:    1,
		TransactionID:    "f2fd5472dc596c3e6f6b377f67a0cade5e65e491758907b937aa3c5c50055aed",
		Timestamp:        timestamp,
		Status:           "COMPLETE",
		IPAddress:        "",
	}
	v := NewDepositHistory(&r)

	assert.Equal(t, expected, v)
}
