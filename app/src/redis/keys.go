package redis

import "fmt"

const (
	BankStocksKey = "bank:stocks"
	WalletsSetKey = "wallets"
	AuditLogKey   = "audit:log"
)

func WalletStocksKey(walletID string) string {
	return fmt.Sprintf("wallet:%s:stocks", walletID)
}
