package services

import (
	"context"
	"encoding/json"
	"fmt"

	rdb "github.com/bulanda/stock-market/src/redis"
)

type LogEntry struct {
	Type      string `json:"type"`
	WalletID  string `json:"wallet_id"`
	StockName string `json:"stock_name"`
}

func GetAuditLog(ctx context.Context) ([]LogEntry, error) {
	data, err := rdb.Client.LRange(ctx, rdb.AuditLogKey, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	entries := make([]LogEntry, 0, len(data))
	for _, raw := range data {
		var entry LogEntry
		if err := json.Unmarshal([]byte(raw), &entry); err != nil {
			return nil, fmt.Errorf("failed to parse log entry: %w", err)
		}
		entries = append(entries, entry)
	}
	return entries, nil
}
