package utils_test

import (
	"testing"

	"github.com/izacgaldino23/binance-consult-trade-api/utils"
	"github.com/stretchr/testify/assert"
)

func TestSlice(t *testing.T) {

	t.Run("TestEach", func(t *testing.T) {
		var (
			temp  = make(utils.Slice[int], 0)
			total = 0

			value1 = 10
			value2 = 20
			value3 = 30
		)

		temp = append(temp, value1)
		temp = append(temp, value2)
		temp = append(temp, value3)

		temp.Each(func(i int, v *int) {
			total += *v
		})

		assert.Equal(t, total, value1+value2+value3)
	})
}
