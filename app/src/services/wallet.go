package services

import (
	"context"
	"fmt"

	rdb "github.com/bulanda/stock-market/src/redis"
)

type Wallet struct {
	ID     string       `json:"id"`
	Stocks []StockEntry `json:"stocks"`
}

func WalletExists(ctx context.Context, walletID string) (bool, error) {
	result, err := rdb.Client.SIsMember(ctx, rdb.WalletsSetKey, walletID).Result()
	if err != nil {
		return false, err
	}
	return result, nil
}

func GetWallet(ctx context.Context, walletID string) (*Wallet, error) {
	exists, err := WalletExists(ctx, walletID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, nil
	}

	data, err := rdb.Client.HGetAll(ctx, rdb.WalletStocksKey(walletID)).Result()
	if err != nil {
		return nil, err
	}

	stocks := make([]StockEntry, 0, len(data))
	for name, qty := range data {
		q := 0
		if _, err := fmt.Sscanf(qty, "%d", &q); err != nil {
			return nil, fmt.Errorf("invalid quantity for stock %s: %w", name, err)
		}
		if q > 0 {
			stocks = append(stocks, StockEntry{Name: name, Quantity: q})
		}
	}

	return &Wallet{ID: walletID, Stocks: stocks}, nil
}

func GetWalletStockQuantity(ctx context.Context, walletID, stockName string) (int, bool, error) {
	exists, err := WalletExists(ctx, walletID)
	if err != nil {
		return 0, false, err
	}
	if !exists {
		return 0, false, nil
	}

	qty, err := rdb.Client.HGet(ctx, rdb.WalletStocksKey(walletID), stockName).Result()
	if err != nil {
		return 0, true, nil // wallet exists but stock not found, return 0
	}

	q := 0
	if _, err := fmt.Sscanf(qty, "%d", &q); err != nil {
		return 0, true, fmt.Errorf("invalid quantity: %w", err)
	}
	return q, true, nil
}
