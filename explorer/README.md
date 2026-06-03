# ME Network Explorer

## Overview

The ME Network Explorer is a web-based interface for monitoring the ME Network blockchain. It provides real-time visibility into chain activity, including block production, transaction history, and network status.

## Features

- **Real-time Block Height Display**: The block height automatically updates every 5-6 seconds via polling mechanism
- **Network Status Dashboard**: View current block height, validator set, and network parameters
- **Transaction Explorer**: Search and view transaction details
- **Block Explorer**: Browse blocks and their contents

## Architecture

The explorer consists of:

1. **Frontend UI**: HTML/CSS/JavaScript interface that communicates with the backend API
2. **Backend API**: RESTful API server that proxies requests to the Cosmos SDK RPC endpoint
3. **Polling Mechanism**: The frontend polls the backend every 5 seconds for the latest block height

## Setup

### Prerequisites

- Node.js 18+ (for frontend development)
- Go 1.21+ (for backend development)
- Access to a running ME Network node with RPC enabled

### Running the Explorer

1. Start the backend API server:
   ```bash
   cd explorer/api
   go run main.go
   ```

2. Open the frontend:
   ```bash
   cd explorer/ui
   open index.html
   ```

## API Endpoints

### GET /api/block-height

Returns the latest block height from the connected node.

**Response:**
```json
{
  "height": 158432,
  "timestamp": "2024-01-15T10:30:00Z"
}
```

### GET /api/status

Returns the full node status including sync information.

**Response:**
```json
{
  "node_info": { ... },
  "sync_info": {
    "latest_block_height": "158432",
    "catching_up": false
  },
  "validator_info": { ... }
}
```

## Development

### Frontend

The frontend is a single-page application built with vanilla JavaScript. Key files:

- `ui/index.html` - Main HTML file
- `ui/js/app.js` - Application logic and polling mechanism
- `ui/css/style.css` - Styling

### Backend

The backend is a Go HTTP server that proxies requests to the Cosmos SDK RPC. Key files:

- `api/main.go` - Server entry point
- `api/handlers.go` - HTTP handlers
- `api/rpc.go` - RPC client for communicating with the node

## Troubleshooting

### Block height not updating

1. Check that the backend API server is running
2. Verify the node RPC endpoint is accessible
3. Check browser console for any JavaScript errors
4. Ensure the polling interval is not blocked by browser tab throttling

### Connection refused

1. Ensure the ME Network node is running
2. Verify the RPC port (default: 26657) is open
3. Check firewall settings

## License

This project is licensed under the Apache License 2.0.
