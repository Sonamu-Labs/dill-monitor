package main

import (
	"context"
	"dill-monitor/internal/config"
	"dill-monitor/internal/models"
	"dill-monitor/internal/repository"
	"dill-monitor/internal/service"
	"dill-monitor/pkg/metrics"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	configPath         = flag.String("config", "", "path to config file")
	serverConfigPath   = flag.String("server-config", "", "path to server config file")
	defaultMetricsPort = 9090 // 기본값
)

func getDefaultConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Printf("Could not get home directory: %v, using current directory", err)
		return "config/config.json"
	}
	return filepath.Join(homeDir, ".dill_monitor", "config.json")
}

func getDefaultServerConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Printf("Could not get home directory: %v, using current directory", err)
		return "config/server_config.json"
	}
	return filepath.Join(homeDir, ".dill_monitor", "server_config.json")
}

func main() {
	flag.Parse()

	// Set default config paths if not specified
	if *configPath == "" {
		*configPath = getDefaultConfigPath()
	}
	if *serverConfigPath == "" {
		*serverConfigPath = getDefaultServerConfigPath()
	}

	log.Printf("Using config file: %s", *configPath)
	log.Printf("Using server config file: %s", *serverConfigPath)

	// config.json 파일이 실제로 존재하는지 확인
	if _, err := os.Stat(*configPath); os.IsNotExist(err) {
		log.Fatalf("Config file does not exist: %s", *configPath)
	}

	// 설정 파일의 내용을 로그에 출력
	configData, err := os.ReadFile(*configPath)
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}
	log.Printf("Config file contents: %s", string(configData))

	// Initialize configuration
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 초기화 후 cfg 내용 확인
	log.Printf("Loaded addresses: %d", len(cfg.ListAddresses()))

	// Initialize server configuration
	serverCfg, err := config.LoadServerConfig(*serverConfigPath)
	if err != nil {
		log.Printf("Failed to load server config: %v, using default values", err)
		serverCfg = &models.ServerConfig{
			MetricsPort: defaultMetricsPort,
			LogLevel:    "info",
			Host:        "0.0.0.0",
		}
	}

	// Initialize Prometheus metrics
	promClient := metrics.NewPrometheusClient()
	promRepo := repository.NewPrometheusRepository(promClient)

	// Initialize services
	balanceService := service.NewBalanceService(promRepo)

	// Create a new ServeMux for routing
	mux := http.NewServeMux()

	// Check if we're in test environment
	isTestEnv := os.Getenv("DILL_ENV") == "test"

	// Serve static files and web interface only in test environment
	if isTestEnv {
		log.Println("Running in test environment - web interface enabled")
		// Serve static files
		fs := http.FileServer(http.Dir("static"))
		mux.Handle("/static/", http.StripPrefix("/static/", fs))

		// Handle web interface
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			tmpl := template.Must(template.ParseFiles("templates/index.html"))
			addresses := cfg.ListAddresses()
			data := struct {
				Addresses []models.Address
			}{
				Addresses: addresses,
			}
			tmpl.Execute(w, data)
		})
	} else {
		log.Println("Running in production environment - web interface disabled")
		// Redirect root to metrics endpoint in production
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/metrics", http.StatusMovedPermanently)
		})
	}

	// Handle metrics endpoint
	mux.Handle("/metrics", promhttp.Handler())

	// Start the server
	addr := fmt.Sprintf("%s:%d", serverCfg.Host, serverCfg.MetricsPort)
	log.Printf("Starting server on %s", addr)
	go func() {
		if err := http.ListenAndServe(addr, mux); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Create context that can be cancelled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 종료 채널 및 WaitGroup 추가
	done := make(chan bool)
	var wg sync.WaitGroup

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		sig := <-sigChan
		log.Printf("Received shutdown signal: %v", sig)
		cancel() // 컨텍스트 취소

		// 모든 고루틴이 정리될 때까지 잠시 대기
		time.Sleep(2 * time.Second)

		close(done)
	}()

	// 메인 처리 루프 수정
	wg.Add(1)
	go func() {
		defer wg.Done()

		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		// 즉시 처리 시작
		log.Println("Processing addresses on startup...")
		processAddresses(ctx, cfg, balanceService)

		for {
			select {
			case <-ctx.Done():
				log.Println("Processing loop received cancel signal")
				return
			case <-ticker.C:
				log.Println("Processing addresses on schedule...")
				processAddresses(ctx, cfg, balanceService)
			}
		}
	}()

	// 모든 고루틴이 완료될 때까지 대기
	go func() {
		wg.Wait()
		done <- true
	}()

	// done 채널이 닫히거나 신호를 받을 때까지 대기
	<-done
	log.Println("Shutting down gracefully...")
}

func processAddresses(ctx context.Context, cfg *config.Config, balanceService *service.BalanceService) {
	addresses := cfg.ListAddresses()

	// Create a wait group to wait for all goroutines to finish
	var wg sync.WaitGroup
	wg.Add(len(addresses))

	// Create a channel for results
	resultCh := make(chan struct {
		addr    string
		balance *models.Balance
		err     error
	}, len(addresses))

	// Process each address in a separate goroutine
	for _, addr := range addresses {
		go func(addr models.Address) {
			defer wg.Done()

			balance, err := balanceService.ProcessAddress(ctx, addr)
			resultCh <- struct {
				addr    string
				balance *models.Balance
				err     error
			}{addr.Address, balance, err}
		}(addr)
	}

	// Close the result channel once all goroutines are done
	go func() {
		wg.Wait()
		close(resultCh)
	}()

	// Collect all processed balances for summary metrics
	var processedBalances []*models.Balance

	// Process results as they come in
	for result := range resultCh {
		if result.err != nil {
			log.Printf("Error processing address %s: %v", result.addr, result.err)
			continue
		}
		log.Printf("Processed address %s: balance=%s, staking=%s, reward=%s",
			result.addr, result.balance.Balance, result.balance.StakingBalance, result.balance.Reward)

		// Collect successful results for summary metrics
		processedBalances = append(processedBalances, result.balance)
	}

	// Update summary metrics with all processed balances
	if len(processedBalances) > 0 {
		if err := balanceService.UpdateSummaryMetrics(ctx, processedBalances); err != nil {
			log.Printf("Error updating summary metrics: %v", err)
		} else {
			log.Printf("Updated summary metrics with %d addresses", len(processedBalances))
		}

		// Process validator information to update validator-specific metrics
		if err := balanceService.ProcessValidators(ctx); err != nil {
			log.Printf("Error processing validators: %v", err)
		} else {
			log.Printf("Processed validator information")
		}
	}
}
