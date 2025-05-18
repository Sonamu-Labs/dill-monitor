package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	// ValidatorStatusActive 값은 활성 상태의 validator를 나타냅니다
	ValidatorStatusActive = 1.0
	// ValidatorStatusInactive 값은 비활성 상태의 validator를 나타냅니다
	ValidatorStatusInactive = 0.0
)

// PrometheusClient handles Prometheus metrics registration and updates
type PrometheusClient struct {
	// Balance metrics
	balanceGauge               *prometheus.GaugeVec
	stakingBalanceGauge        *prometheus.GaugeVec
	stakedAmountGauge          *prometheus.GaugeVec
	rewardGauge                *prometheus.GaugeVec
	dailyRewardGauge           *prometheus.GaugeVec
	latestIncomeGauge          *prometheus.GaugeVec
	lastEpochGauge             *prometheus.GaugeVec
	lastRewardTimeGauge        *prometheus.GaugeVec
	poolCreatedCountGauge      *prometheus.GaugeVec
	poolParticipatedCountGauge *prometheus.GaugeVec

	// Validator metrics
	validatorRewardGauge     *prometheus.GaugeVec
	validatorStatusGauge     *prometheus.GaugeVec
	validatorLastEpochGauge  *prometheus.GaugeVec
	validatorLastRewardGauge *prometheus.GaugeVec
	validatorBalanceGauge    *prometheus.GaugeVec
	validatorStatusInfoGauge *prometheus.GaugeVec

	// Summary metrics
	totalAddressCountGauge    prometheus.Gauge
	totalBalanceGauge         prometheus.Gauge
	totalRewardGauge          prometheus.Gauge
	totalStakedAmountGauge    prometheus.Gauge
	totalValidatorCountGauge  prometheus.Gauge
	activeValidatorCountGauge prometheus.Gauge
	validatorStatusCountGauge *prometheus.GaugeVec

	// API metrics
	requestCounter  *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
	requestErrors   *prometheus.CounterVec
}

// NewPrometheusClient creates a new Prometheus client with registered metrics
func NewPrometheusClient() *PrometheusClient {
	return &PrometheusClient{
		balanceGauge: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "account_balance",
				Help: "Current account balance in DILL",
			},
			[]string{"address", "label"},
		),
		stakingBalanceGauge: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "staking_balance",
				Help: "Current staking balance in DILL",
			},
			[]string{"address", "label"},
		),
		stakedAmountGauge: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "staked_amount",
				Help: "Total staked amount in DILL",
			},
			[]string{"address", "label"},
		),
		rewardGauge: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "reward_amount",
				Help: "Current reward amount in DILL",
			},
			[]string{"address", "label"},
		),
		dailyRewardGauge: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "daily_reward_amount",
				Help: "Daily reward amount in DILL",
			},
			[]string{"address", "label"},
		),
		latestIncomeGauge: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "latest_income_amount",
				Help: "Latest income amount in DILL",
			},
			[]string{"address", "label"},
		),
		lastEpochGauge: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "last_epoch",
				Help: "Last epoch number",
			},
			[]string{"address", "label"},
		),
		lastRewardTimeGauge: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "last_reward_time",
				Help: "Last reward time as unix timestamp",
			},
			[]string{"address", "label"},
		),
		poolCreatedCountGauge: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "pool_created_count",
				Help: "Number of pools created",
			},
			[]string{"address", "label"},
		),
		poolParticipatedCountGauge: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "pool_participated_count",
				Help: "Number of pools participated in",
			},
			[]string{"address", "label"},
		),
		validatorRewardGauge: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "validator_reward",
				Help: "Validator reward amount in DILL",
			},
			[]string{"validator_idx", "label"},
		),
		validatorStatusGauge: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "validator_status",
				Help: "Validator status (1 for active, 0 for inactive)",
			},
			[]string{"validator_idx", "label"},
		),
		validatorLastEpochGauge: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "validator_last_epoch",
				Help: "Validator's last processed epoch",
			},
			[]string{"validator_idx", "label"},
		),
		validatorLastRewardGauge: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "validator_last_reward_time",
				Help: "Validator's last reward time as unix timestamp",
			},
			[]string{"validator_idx", "label"},
		),
		validatorBalanceGauge: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "validator_balance",
				Help: "Validator's balance in DILL",
			},
			[]string{"validator_idx", "label"},
		),
		validatorStatusInfoGauge: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "validator_status_info",
				Help: "Validator status information",
			},
			[]string{"validator_idx", "label", "status"},
		),
		// Summary metrics
		totalAddressCountGauge: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "total_address_count",
				Help: "Total number of addresses being monitored",
			},
		),
		totalBalanceGauge: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "total_balance",
				Help: "Total balance across all addresses in DILL",
			},
		),
		totalRewardGauge: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "total_reward",
				Help: "Total rewards across all addresses in DILL",
			},
		),
		totalStakedAmountGauge: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "total_staked_amount",
				Help: "Total staked amount across all addresses in DILL",
			},
		),
		totalValidatorCountGauge: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "total_validator_count",
				Help: "Total number of validators",
			},
		),
		activeValidatorCountGauge: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "active_validator_count",
				Help: "Total number of active validators",
			},
		),
		validatorStatusCountGauge: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "validator_status_count",
				Help: "Number of validators in each status",
			},
			[]string{"status"},
		),
		requestCounter: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"endpoint", "method", "status"},
		),
		requestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "HTTP request duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"endpoint", "method"},
		),
		requestErrors: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_request_errors_total",
				Help: "Total number of HTTP request errors",
			},
			[]string{"endpoint", "method", "error_type"},
		),
	}
}

// UpdateBalanceMetrics updates all balance-related metrics
func (c *PrometheusClient) UpdateBalanceMetrics(
	address string,
	label string,
	balance,
	stakingBalance,
	stakedAmount,
	reward,
	dailyReward,
	latestIncome float64,
	lastEpoch float64,
	lastRewardTime float64,
	poolCreatedCount float64,
	poolParticipatedCount float64,
) {
	c.balanceGauge.WithLabelValues(address, label).Set(balance)
	c.stakingBalanceGauge.WithLabelValues(address, label).Set(stakingBalance)
	c.stakedAmountGauge.WithLabelValues(address, label).Set(stakedAmount)
	c.rewardGauge.WithLabelValues(address, label).Set(reward)
	c.dailyRewardGauge.WithLabelValues(address, label).Set(dailyReward)
	c.latestIncomeGauge.WithLabelValues(address, label).Set(latestIncome)
	c.lastEpochGauge.WithLabelValues(address, label).Set(lastEpoch)
	c.lastRewardTimeGauge.WithLabelValues(address, label).Set(lastRewardTime)
	c.poolCreatedCountGauge.WithLabelValues(address, label).Set(poolCreatedCount)
	c.poolParticipatedCountGauge.WithLabelValues(address, label).Set(poolParticipatedCount)
}

// UpdateBasicMetrics updates only basic metrics that should always be shown
func (c *PrometheusClient) UpdateBasicMetrics(
	address string,
	label string,
	balance float64,
	stakedAmount float64,
	reward float64,
	poolCreatedCount float64,
	poolParticipatedCount float64,
	lastRewardTime float64,
) {
	c.balanceGauge.WithLabelValues(address, label).Set(balance)
	c.stakedAmountGauge.WithLabelValues(address, label).Set(stakedAmount)
	c.rewardGauge.WithLabelValues(address, label).Set(reward)
	c.lastRewardTimeGauge.WithLabelValues(address, label).Set(lastRewardTime)
	c.poolCreatedCountGauge.WithLabelValues(address, label).Set(poolCreatedCount)
	c.poolParticipatedCountGauge.WithLabelValues(address, label).Set(poolParticipatedCount)
}

// UpdateValidatorRelatedMetrics updates metrics that should only be shown for accounts with validators
func (c *PrometheusClient) UpdateValidatorRelatedMetrics(
	address string,
	label string,
	stakingBalance float64,
	dailyReward float64,
	latestIncome float64,
	lastEpoch float64,
) {
	c.stakingBalanceGauge.WithLabelValues(address, label).Set(stakingBalance)
	c.dailyRewardGauge.WithLabelValues(address, label).Set(dailyReward)
	c.latestIncomeGauge.WithLabelValues(address, label).Set(latestIncome)
	c.lastEpochGauge.WithLabelValues(address, label).Set(lastEpoch)
}

// UpdateValidatorMetrics updates validator-related metrics
func (c *PrometheusClient) UpdateValidatorMetrics(
	validatorIdx string,
	label string,
	reward float64,
	balance float64,
	isActive bool,
	lastEpoch float64,
	lastRewardTime float64,
	statusString string,
) {
	c.validatorRewardGauge.WithLabelValues(validatorIdx, label).Set(reward)
	c.validatorBalanceGauge.WithLabelValues(validatorIdx, label).Set(balance)
	status := ValidatorStatusInactive
	if isActive {
		status = ValidatorStatusActive
	}
	c.validatorStatusGauge.WithLabelValues(validatorIdx, label).Set(status)
	c.validatorLastEpochGauge.WithLabelValues(validatorIdx, label).Set(lastEpoch)
	c.validatorLastRewardGauge.WithLabelValues(validatorIdx, label).Set(lastRewardTime)

	// 상태 정보 업데이트
	c.UpdateValidatorStatusInfo(validatorIdx, label, statusString)
}

// UpdateValidatorStatusInfo updates the validator status information
func (c *PrometheusClient) UpdateValidatorStatusInfo(validatorIdx string, label string, statusString string) {
	// status가 비어있는 경우 unknown으로 처리
	if statusString == "" {
		statusString = "unknown"
	}

	// 새 상태 정보 설정 (값은 1로 고정, 라벨에 상태 정보 포함)
	c.validatorStatusInfoGauge.WithLabelValues(validatorIdx, label, statusString).Set(1)
}

// RecordAPIMetrics records API-related metrics
func (c *PrometheusClient) RecordAPIMetrics(endpoint, method string, status int, duration float64) {
	c.requestCounter.WithLabelValues(endpoint, method, string(status)).Inc()
	c.requestDuration.WithLabelValues(endpoint, method).Observe(duration)
	if status >= 400 {
		c.requestErrors.WithLabelValues(endpoint, method, "http_error").Inc()
	}
}

// UpdateSummaryMetrics updates summary metrics with aggregated data
func (c *PrometheusClient) UpdateSummaryMetrics(
	addressCount int,
	validatorCount int,
	activeValidatorCount int,
	totalBalance float64,
	totalReward float64,
	totalStakedAmount float64,
) {
	c.totalAddressCountGauge.Set(float64(addressCount))
	c.totalValidatorCountGauge.Set(float64(validatorCount))
	c.activeValidatorCountGauge.Set(float64(activeValidatorCount))
	c.totalBalanceGauge.Set(totalBalance)
	c.totalRewardGauge.Set(totalReward)
	c.totalStakedAmountGauge.Set(totalStakedAmount)
}

// UpdateValidatorStatusMetrics updates the count of validators by status
func (c *PrometheusClient) UpdateValidatorStatusMetrics(statusCounts map[string]int) {
	// 모든 상태 카운터를 0으로 초기화 (기존 값 제거)
	c.validatorStatusCountGauge.Reset()

	// 각 상태별 카운트 설정
	for status, count := range statusCounts {
		c.validatorStatusCountGauge.WithLabelValues(status).Set(float64(count))
	}
}
