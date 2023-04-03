package comparison

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gummy789j/telegram-quote-bot/internal/constant"
	domain "github.com/gummy789j/telegram-quote-bot/internal/domain/repo"
	"github.com/gummy789j/telegram-quote-bot/internal/transport"
)

type comparisonClient struct {
	cli      transport.HttpClient
	endpoint string
}

var _ domain.QuoteRepo = (*comparisonClient)(nil)

func NewComparisonClient(cli transport.HttpClient) domain.QuoteRepo {
	return &comparisonClient{
		cli:      cli,
		endpoint: "https://www.usdtwhere.com/wallet-api",
	}
}

var (
	pathComparison = "/v1/kgi/exchange-rates/comparison"
)

func (c *comparisonClient) GetQuotations(ctx context.Context, req domain.GetQuotationsRequest) (*domain.GetQuotationsResponse, error) {

	url := fmt.Sprintf("%s%s", c.endpoint, pathComparison)
	httpResp, err := c.cli.Send(ctx, &transport.HttpRequest{
		Method: http.MethodGet,
		URL:    url,
	})
	if err != nil {
		log.Println("get quotations failed", err.Error())
		return nil, err
	}

	respBody := &comparisonRespBody{}

	if err := json.Unmarshal(httpResp.Body, respBody); err != nil {
		log.Println("json marshal failed", err.Error())
		return nil, err
	}

	infos := make(map[constant.Exchange]domain.QuotationInfo)

	for _, v := range respBody.Data.Exchanges {
		updateTime := time.UnixMilli(v.UpdateTime)
		exchange := constant.Exchange(v.Name)
		if len(exchange) == 0 {
			continue // skip unknown exchange
		}

		infos[exchange] = domain.QuotationInfo{
			BuyPrice:   v.BuyRate,
			SellPrice:  v.SellRate,
			UpdateTime: updateTime,
		}
	}

	return &domain.GetQuotationsResponse{Infos: infos}, nil
}
