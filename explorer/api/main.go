package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	defaultPort     = "8080"
	defaultRPCURL   = "http://localhost:26657"
	requestTimeout  = 5 * time.Second
)

type BlockHeightResponse struct {
	Height    int64  `json:"height"`
	Timestamp string `json:"timestamp"`
}

type StatusResponse struct {
	NodeInfo    interface{} `json:"node_info"`
	SyncInfo    interface{} `json:"sync_info"`
	ValidatorInfo interface{} `json:"validator_info"`
}

var rpcClient *RPCClient

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	rpcURL := os.Getenv("RPC_URL")
	if rpcURL == "" {
		rpcURL = defaultRPCURL
	}

	rpcClient = NewRPCClient(rpcURL, requestTimeout)

	http.HandleFunc("/api/block-height", corsMiddleware(blockHeightHandler))
	http.HandleFunc("/api/status", corsMiddleware(statusHandler))
	http.HandleFunc("/health", corsMiddleware(healthHandler))

	log.Printf("Explorer API server starting on port %s, connecting to RPC at %s\n", port, rpcURL)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

func blockHeightHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	height, err := rpcClient.GetBlockHeight()
	if err != nil {
		log.Printf("Error getting block height: %v\n", err)
		http.Error(w, fmt.Sprintf("Failed to get block height: %v", err), http.StatusInternalServerError)
		return
	}

	resp := BlockHeightResponse{
		Height:    height,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	status, err := rpcClient.GetStatus()
	if err != nil {
		log.Printf("Error getting status: %v\n", err)
		http.Error(w, fmt.Sprintf("Failed to get status: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
