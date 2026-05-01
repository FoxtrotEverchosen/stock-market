package services

import (
	"context"
	"fmt"

	rdb "github.com/bulanda/stock-market/src/redis"
)

type StockEntry struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}

func GetBankStocks(ctx context.Context) ([]StockEntry, error) {
	data, err := rdb.Client.HGetAll(ctx, rdb.BankStocksKey).Result()
	if err != nil {
		return nil, err
	}

	stocks := make([]StockEntry, 0, len(data))
	for name, qty := range data {
		q := 0
		if _, err := fmt.Sscanf(qty, "%d", &q); err != nil {
			return nil, fmt.Errorf("invalid quantity for stock %s: %w", name, err)
		}
		stocks = append(stocks, StockEntry{Name: name, Quantity: q})
	}
	return stocks, nil
}

func SetBankStocks(ctx context.Context, stocks []StockEntry) error {
	pipe := rdb.Client.Pipeline()
	pipe.Del(ctx, rdb.BankStocksKey)
	if len(stocks) > 0 {
		fields := make(map[string]any, len(stocks))
		for _, s := range stocks {
			fields[s.Name] = s.Quantity
		}
		pipe.HSet(ctx, rdb.BankStocksKey, fields)
	}
	_, err := pipe.Exec(ctx)
	return err
}
