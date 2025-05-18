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

4. Install the application with make:

```bash
make install
```

This will install the application to `~/.dill_monitor/` (Linux/macOS) or `%USERPROFILE%\.dill_monitor` (Windows).

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

The server configuration is stored in `server_config.json`:

```json
{
    "metricsPort": 9090,
    "logLevel": "info",
    "host": "0.0.0.0"
}
```

## Usage

### Running the Application Directly

Run the application with default settings:

```bash
~/.dill_monitor/dill-monitor
```

Or specify custom configuration:

```bash
~/.dill_monitor/dill-monitor -config=/path/to/config.json -server-config=/path/to/server_config.json
```

### Running as a Service

#### Linux (systemd)

1. Create a systemd service file using tee:

```bash
sudo tee /etc/systemd/system/dill-monitor.service > /dev/null << EOT
[Unit]
Description=Dill Monitor Service
After=network.target

[Service]
Type=simple
User=$(whoami)
ExecStart=$(echo $HOME)/.dill_monitor/dill-monitor
Restart=on-failure
RestartSec=10
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
EOT
```

2. Enable and start the service:

```bash
sudo systemctl daemon-reload
sudo systemctl enable dill-monitor
sudo systemctl start dill-monitor
```

3. Check service status:

```bash
sudo systemctl status dill-monitor
```

#### Windows

1. Install NSSM (Non-Sucking Service Manager):
   Download from https://nssm.cc/download

2. Open Command Prompt as Administrator and navigate to NSSM directory:

```cmd
cd path\to\nssm\directory
```

3. Install the service:

```cmd
nssm install DillMonitor
```

4. In the GUI that appears:

    - Path: `%USERPROFILE%\.dill_monitor\dill-monitor.exe`
    - Startup directory: `%USERPROFILE%\.dill_monitor`
    - Service name: DillMonitor

5. Start the service:

```cmd
nssm start DillMonitor
```

### Using Docker

1. Build and run with Docker:

```bash
docker build -t dill-monitor .
docker run -d -p 9090:9090 -v $(pwd)/config:/app/config --name dill-monitor dill-monitor
```

2. Or use Docker Compose for a complete monitoring stack:

```bash
docker-compose up -d
```

This will start:

-   Dill Monitor on port 9090
-   Prometheus on port 9091
-   Grafana on port 3000

Access Grafana at http://localhost:3000 (default login: admin/admin)

### Command Line Arguments

-   `-config`: Path to the configuration file (default: "config/config.json")
-   `-server-config`: Path to the server configuration file (default: "config/server_config.json")

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
├── config/              # Configuration files
├── docker/              # Docker-related files
├── internal/
│   ├── api/             # API client implementations
│   ├── config/          # Configuration management
│   ├── models/          # Data models
│   ├── repository/      # Data storage interfaces and implementations
│   ├── service/         # Business logic
│   └── util/            # Utility functions
├── pkg/
│   └── metrics/         # Prometheus metrics
├── prometheus/          # Prometheus configuration
├── Dockerfile
├── docker-compose.yml
└── Makefile
```

### Adding New Features

1. Define new models in `internal/models/`
2. Add repository methods in `internal/repository/interface.go`
3. Implement business logic in `internal/service/`
4. Add metrics in `pkg/metrics/`
5. Update the main application in `cmd/server/main.go`
