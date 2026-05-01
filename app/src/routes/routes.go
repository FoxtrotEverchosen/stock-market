package routes

import "net/http"

func Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /wallets/{wallet_id}/stocks/{stock_name}", handleWalletStockOperation)
	mux.HandleFunc("GET /wallets/{wallet_id}/stocks/{stock_name}", handleGetWalletStock)
	mux.HandleFunc("GET /wallets/{wallet_id}", handleGetWallet)
	mux.HandleFunc("GET /stocks", handleGetBankStocks)
	mux.HandleFunc("POST /stocks", handleSetBankStocks)
	mux.HandleFunc("GET /log", handleGetLog)
	mux.HandleFunc("POST /chaos", handleChaos)
}

func handleWalletStockOperation(w http.ResponseWriter, r *http.Request) {}
func handleGetWalletStock(w http.ResponseWriter, r *http.Request)       {}
func handleGetWallet(w http.ResponseWriter, r *http.Request)            {}
func handleGetBankStocks(w http.ResponseWriter, r *http.Request)        {}
func handleSetBankStocks(w http.ResponseWriter, r *http.Request)        {}
func handleGetLog(w http.ResponseWriter, r *http.Request)               {}
func handleChaos(w http.ResponseWriter, r *http.Request)                {}
