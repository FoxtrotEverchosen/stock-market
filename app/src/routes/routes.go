package routes

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/bulanda/stock-market/src/services"
)

func Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /wallets/{wallet_id}/stocks/{stock_name}", handleWalletStockOperation)
	mux.HandleFunc("GET /wallets/{wallet_id}/stocks/{stock_name}", handleGetWalletStock)
	mux.HandleFunc("GET /wallets/{wallet_id}", handleGetWallet)
	mux.HandleFunc("GET /stocks", handleGetBankStocks)
	mux.HandleFunc("POST /stocks", handleSetBankStocks)
	mux.HandleFunc("GET /log", handleGetLog)
	mux.HandleFunc("POST /chaos", handleChaos)
}

func handleGetBankStocks(w http.ResponseWriter, r *http.Request) {
	stocks, err := services.GetBankStocks(r.Context())
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]any{"stocks": stocks}); err != nil {
		log.Printf("[handleGetBankStocks] failed to encode response: %v", err)
	}
}
func handleSetBankStocks(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Stocks []services.StockEntry `json:"stocks"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	if err := services.SetBankStocks(r.Context(), body.Stocks); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func handleGetWallet(w http.ResponseWriter, r *http.Request) {
	walletID := r.PathValue("wallet_id")

	wallet, err := services.GetWallet(r.Context(), walletID)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	if wallet == nil {
		http.Error(w, "wallet not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(wallet); err != nil {
		log.Printf("[handleGetWallet] failed to encode response: %v", err)
	}
}
func handleGetWalletStock(w http.ResponseWriter, r *http.Request) {
	walletID := r.PathValue("wallet_id")
	stockName := r.PathValue("stock_name")

	qty, walletExists, err := services.GetWalletStockQuantity(r.Context(), walletID, stockName)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	if !walletExists {
		http.Error(w, "wallet not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(qty); err != nil {
		log.Printf("[handleGetWalletStock] failed to encode response: %v", err)
	}
}

func handleWalletStockOperation(w http.ResponseWriter, r *http.Request) {
	walletID := r.PathValue("wallet_id")
	stockName := r.PathValue("stock_name")

	var body struct {
		Type string `json:"type"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	if body.Type != "buy" && body.Type != "sell" {
		http.Error(w, "type must be buy or sell", http.StatusBadRequest)
		return
	}

	err := services.ExecuteTrade(r.Context(), walletID, stockName, body.Type)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrStockNotFound):
			http.Error(w, "stock not found", http.StatusNotFound)
		case errors.Is(err, services.ErrInsufficientBankStock):
			http.Error(w, "insufficient stock in bank", http.StatusBadRequest)
		case errors.Is(err, services.ErrInsufficientWalletStock):
			http.Error(w, "insufficient stock in wallet", http.StatusBadRequest)
		default:
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

func handleGetLog(w http.ResponseWriter, r *http.Request) {
	entries, err := services.GetAuditLog(r.Context())
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]any{"log": entries}); err != nil {
		log.Printf("[handleGetLog] failed to encode response: %v", err)
	}
}

func handleChaos(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
	os.Exit(1)
}
