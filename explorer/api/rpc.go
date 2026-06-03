package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

type RPCClient struct {
	baseURL    string
	httpClient *http.Client
}

type RPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      int             `json:"id"`
	Result  json.RawMessage `json:"result"`
}

type StatusResult struct {
	NodeInfo    interface{} `json:"node_info"`
	SyncInfo    SyncInfo    `json:"sync_info"`
	ValidatorInfo interface{} `json:"validator_info"`
}

type SyncInfo struct {
	LatestBlockHeight string `json:"latest_block_height"`
	CatchingUp        bool   `json:"catching_up"`
}

func NewRPCClient(baseURL string, timeout time.Duration) *RPCClient {
	return &RPCClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *RPCClient) call(method string, params []interface{}) (*RPCResponse, error) {
	requestBody := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  method,
		"params":  params,
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.httpClient.Post(c.baseURL, "application/json", io.NopCloser(strings.NewReader(string(body))))
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var rpcResp RPCResponse
	if err := json.Unmarshal(respBody, &rpcResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &rpcResp, nil
}

func (c *RPCClient) GetBlockHeight() (int64, error) {
	resp, err := c.call("status", []interface{}{})
	if err != nil {
		return 0, fmt.Errorf("failed to get status: %w", err)
	}

	var status StatusResult
	if err := json.Unmarshal(resp.Result, &status); err != nil {
		return 0, fmt.Errorf("failed to unmarshal status: %w", err)
	}

	height, err := strconv.ParseInt(status.SyncInfo.LatestBlockHeight, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse block height: %w", err)
	}

	return height, nil
}

func (c *RPCClient) GetStatus() (*StatusResult, error) {
	resp, err := c.call("status", []interface{}{})
	if err != nil {
		return nil, fmt.Errorf("failed to get status: %w", err)
	}

	var status StatusResult
	if err := json.Unmarshal(resp.Result, &status); err != nil {
		return nil, fmt.Errorf("failed to unmarshal status: %w", err)
	}

	return &status, nil
}
