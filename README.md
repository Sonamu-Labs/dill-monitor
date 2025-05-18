# Dill Monitor

A comprehensive monitoring application for tracking validator balances, rewards, and status in the Dill blockchain network.

## Features

-   Concurrent processing of wallet addresses using goroutines
-   Real-time tracking of validator status (active, inactive, etc.)
-   Detailed monitoring of wallet balances and staking information
-   Tracking of validator performance, rewards, and daily income
-   In-memory storage for reliable validator information retrieval
-   Aggregate statistics for total balances, rewards, and validator counts
-   User-defined labels for better validator identification
-   Exposes comprehensive metrics via Prometheus
-   Configurable address list with custom labels
-   Graceful shutdown handling

## Prerequisites

-   Go 1.21 or later
-   Prometheus (for metrics collection)
-   Grafana (optional, for visualization)

## Installation

1. Clone the repository:

```bash
git clone https://github.com/yourusername/dill-monitor.git
cd dill-monitor
```

2. Install dependencies:

```bash
go mod download
```

3. Build the application:

```bash
go build -o dill-monitor ./cmd/server
```

## Configuration

The application uses a JSON configuration file located at `config/config.json`. The file should contain a list of addresses to monitor with optional user-defined labels:

```json
{
    "addresses": [
        {
            "label": "MainValidator-1",
            "address": "0x...",
            "validator_address": "0x..."
        },
        {
            "label": "MainValidator-2",
            "address": "0x...",
            "validator_address": "0x..."
        }
    ]
}
```

## Usage

Run the application with default settings:

```bash
./dill-monitor
```

Or specify custom configuration:

```bash
./dill-monitor -config /path/to/config.json -metrics-port 9090
```

### Command Line Arguments

-   `-config`: Path to the configuration file (default: "config/config.json")
-   `-metrics-port`: Port for Prometheus metrics (default: 9090)

## Metrics

The application exposes the following Prometheus metrics:

### Balance Metrics

-   `dill_balance`: Current wallet balance with address and label labels
-   `dill_staking_balance`: Current staking balance for validators
-   `dill_staked_amount`: Amount staked by the address
-   `dill_reward`: Current reward amount for the address
-   `dill_daily_reward`: Estimated daily rewards for validators
-   `dill_latest_income`: Latest income amount for validators
-   `dill_pool_created_count`: Number of pools created by the address
-   `dill_pool_participated_count`: Number of pools participated in by the address
-   `dill_last_reward_time`: Unix timestamp of the last reward time

### Validator Metrics

-   `dill_validator_reward`: Validator reward amount
-   `dill_validator_balance`: Current validator balance
-   `dill_validator_active`: Validator active status (1 for active, 0 for inactive)
-   `dill_validator_last_epoch`: Last epoch number for the validator
-   `dill_validator_last_reward_time`: Unix timestamp of the last validator reward time
-   `dill_validator_status_info`: Status information for validators (with status label)

### Aggregate Metrics

-   `dill_address_count`: Total number of addresses being monitored
-   `dill_validator_count`: Total number of validators being monitored
-   `dill_active_validator_count`: Number of active validators
-   `dill_total_balance`: Sum of all wallet balances
-   `dill_total_reward`: Sum of all rewards
-   `dill_total_staked_amount`: Sum of all staked amounts
-   `dill_validator_status_count`: Count of validators by status

### API Metrics

-   `dill_api_requests_total`: Total number of API requests
-   `dill_api_request_duration_seconds`: API request duration
-   `dill_api_errors_total`: Total number of API errors

## Development

### Project Structure

```
.
├── cmd/
│   ├── server/          # Application entry point
│   └── test/            # Test utilities
├── internal/
│   ├── api/             # API client implementations
│   ├── config/          # Configuration management
│   ├── models/          # Data models
│   ├── repository/      # Data storage interfaces and implementations
│   ├── service/         # Business logic
│   └── util/            # Utility functions
└── pkg/
    └── metrics/         # Prometheus metrics
```

### Adding New Features

1. Define new models in `internal/models/`
2. Add repository methods in `internal/repository/interface.go`
3. Implement business logic in `internal/service/`
4. Add metrics in `pkg/metrics/`
5. Update the main application in `cmd/server/main.go`
