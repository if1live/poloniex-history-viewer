package yui

import (
	"testing"

	"fmt"

	"github.com/stretchr/testify/assert"
)

func Test_ToFix(t *testing.T) {
	cases := []struct {
		input float64
		prec  int
		fixed float64
		str   string
	}{
		// floor
		{0.1234567812, 8, 0.12345678, "0.12345678"},

		// celi
		{0.1234567891, 8, 0.12345679, "0.12345679"},
	}
	for _, c := range cases {
		v := ToFixed(c.input, c.prec)
		assert.Equal(t, c.fixed, v)
		assert.Equal(t, c.str, fmt.Sprint(c.fixed))
	}
}
