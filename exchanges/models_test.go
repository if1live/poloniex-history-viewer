package exchanges

import (
	"testing"

	"github.com/if1live/poloniex-history-viewer/yui"
	"github.com/stretchr/testify/assert"

	"strconv"
)

func TestRow_buyFeeAmount(t *testing.T) {
	// buy example
	// amount : 13.00373802
	// fee : 0.03250935 SYS (0.25%)
	// 0.03250935 SYS = 13.00373802 SYS * (0.01) * (0.25)
	cases := []struct {
		amountStr    string
		feeStr       string
		feeAmountStr string
	}{
		{"13.00373802", "0.25", "0.03250935"},
		{"20.75377718", "0.25", "0.05188444"},
	}
	for _, c := range cases {
		amount, _ := strconv.ParseFloat(c.amountStr, 64)
		fee, _ := strconv.ParseFloat(c.feeStr, 64)
		r := Exchange{
			Amount: amount,
			Fee:    fee * float64(0.01),
		}
		v := r.buyFeeAmount()
		assert.Equal(t, c.feeAmountStr, yui.ToFloatStr(v))
	}
}

func TestRow_sellFeeAmount(t *testing.T) {
	cases := []struct {
		totalStr     string
		feeStr       string
		feeAmountStr string
	}{
		{"0.01085732", "0.15", "0.00001629"},
		// trade id: 170322632
		{"0.01068353", "0.25", "0.00002671"},
	}
	for _, c := range cases {
		total, _ := strconv.ParseFloat(c.totalStr, 64)
		fee, _ := strconv.ParseFloat(c.feeStr, 64)
		r := Exchange{
			Total: total,
			Fee:   fee * float64(0.01),
		}
		v := r.sellFeeAmount()
		assert.Equal(t, c.feeAmountStr, yui.ToFloatStr(v))
	}
}

func TestRow_MyTotal(t *testing.T) {
	cases := []struct {
		typeStr    string
		totalStr   string
		feeStr     string
		myTotalStr string
	}{
		// trade id: 170322632
		{"sell", "0.01068353", "0.25", "0.01065682"},
	}
	for _, c := range cases {
		total, _ := strconv.ParseFloat(c.totalStr, 64)
		fee, _ := strconv.ParseFloat(c.feeStr, 64)
		r := Exchange{
			Type:  c.typeStr,
			Total: total,
			Fee:   fee * float64(0.01),
		}
		v := r.MyTotal()
		assert.Equal(t, c.myTotalStr, yui.ToFloatStr(v))
	}
}
