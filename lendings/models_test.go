package lendings

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLending_FeeRate(t *testing.T) {
	r := Lending{
		LendingID: 374173954,
		Currency:  "BTC",
		Rate:      0.001,
		Amount:    0.1,
		Duration:  0.0017,
		Interest:  1.8e-07,
		Fee:       -2e-08,
		Earned:    1.6e-07,
	}

	v := r.FeeRate()
	expected := 0.15
	assert.Equal(t, expected, v)
}
