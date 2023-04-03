package comparison

import (
	"github.com/shopspring/decimal"
)

type comparisonRespBody struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Exchanges []struct {
			Name       string          `json:"name"`
			BuyRate    decimal.Decimal `json:"buy_rate"`
			SellRate   decimal.Decimal `json:"sell_rate"`
			UpdateTime int64           `json:"update_time"`
		} `json:"exchanges"`
	} `json:"data"`
}
