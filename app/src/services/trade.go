package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	rdb "github.com/bulanda/stock-market/src/redis"
	"github.com/redis/go-redis/v9"
)

var ErrStockNotFound = errors.New("stock not found in bank")
var ErrInsufficientBankStock = errors.New("insufficient stock in bank")
var ErrInsufficientWalletStock = errors.New("insufficient stock in wallet")

var buyScript = redis.NewScript(`
local bank_key = KEYS[1]
local wallet_key = KEYS[2]
local audit_key = KEYS[3]
local stock_name = ARGV[1]
local wallet_id = ARGV[2]
local wallets_set_key = KEYS[4]

-- check stock exists in bank
local bank_qty = redis.call('HGET', bank_key, stock_name)
if bank_qty == false then
    return -1
end

-- check sufficient quantity in bank
if tonumber(bank_qty) < 1 then
    return -2
end

-- execute transfer
redis.call('HINCRBY', bank_key, stock_name, -1)
redis.call('HINCRBY', wallet_key, stock_name, 1)
redis.call('SADD', wallets_set_key, wallet_id)

-- audit log
local entry = cjson.encode({type='buy', wallet_id=wallet_id, stock_name=stock_name})
redis.call('RPUSH', audit_key, entry)

return 1
`)

var sellScript = redis.NewScript(`
local bank_key = KEYS[1]
local wallet_key = KEYS[2]
local audit_key = KEYS[3]
local stock_name = ARGV[1]
local wallet_id = ARGV[2]
local wallets_set_key = KEYS[4]

-- check stock exists in bank
local bank_qty = redis.call('HGET', bank_key, stock_name)
if bank_qty == false then
    return -1
end

-- check wallet has the stock
local wallet_qty = redis.call('HGET', wallet_key, stock_name)
if wallet_qty == false or tonumber(wallet_qty) < 1 then
    return -3
end

-- execute transfer
redis.call('HINCRBY', wallet_key, stock_name, -1)
redis.call('HINCRBY', bank_key, stock_name, 1)
redis.call('SADD', wallets_set_key, wallet_id)

-- audit log
local entry = cjson.encode({type='sell', wallet_id=wallet_id, stock_name=stock_name})
redis.call('RPUSH', audit_key, entry)

return 1
`)

func ExecuteTrade(ctx context.Context, walletID, stockName, tradeType string) error {
	keys := []string{
		rdb.BankStocksKey,
		rdb.WalletStocksKey(walletID),
		rdb.AuditLogKey,
		rdb.WalletsSetKey,
	}
	args := []any{stockName, walletID}

	var result interface{}
	var err error

	if tradeType == "buy" {
		result, err = buyScript.Run(ctx, rdb.Client, keys, args...).Result()
	} else {
		result, err = sellScript.Run(ctx, rdb.Client, keys, args...).Result()
	}

	if err != nil {
		return fmt.Errorf("script error: %w", err)
	}

	switch result.(int64) {
	case 1:
		return nil
	case -1:
		return ErrStockNotFound
	case -2:
		return ErrInsufficientBankStock
	case -3:
		return ErrInsufficientWalletStock
	default:
		return fmt.Errorf("unexpected result: %v", result)
	}
}

// keep compiler happy until audit log endpoint is wired
var _ = json.Marshal
