package repository

import (
	"context"
	"dill-monitor/internal/models"
	"dill-monitor/pkg/metrics"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

// PrometheusRepository implements the Repository interface using Prometheus
type PrometheusRepository struct {
	client *metrics.PrometheusClient
	// 메모리 내 저장소: 최근 업데이트된 밸런스 정보를 저장
	balances      map[string]*models.Balance
	balancesMutex sync.RWMutex
}

// NewPrometheusRepository creates a new Prometheus repository
func NewPrometheusRepository(client *metrics.PrometheusClient) *PrometheusRepository {
	return &PrometheusRepository{
		client:   client,
		balances: make(map[string]*models.Balance),
	}
}

// SaveBalance implements Repository.SaveBalance
func (r *PrometheusRepository) SaveBalance(ctx context.Context, balance *models.Balance) error {
	// 메모리에 밸런스 정보 저장
	r.balancesMutex.Lock()
	r.balances[balance.Address] = balance
	r.balancesMutex.Unlock()

	return r.UpdateBalance(ctx, balance)
}

// GetBalance implements Repository.GetBalance
func (r *PrometheusRepository) GetBalance(ctx context.Context, address string) (*models.Balance, error) {
	r.balancesMutex.RLock()
	defer r.balancesMutex.RUnlock()

	if balance, exists := r.balances[address]; exists {
		return balance, nil
	}
	return nil, fmt.Errorf("balance not found for address: %s", address)
}

// ListBalances implements Repository.ListBalances
func (r *PrometheusRepository) ListBalances(ctx context.Context) ([]*models.Balance, error) {
	r.balancesMutex.RLock()
	defer r.balancesMutex.RUnlock()

	// 저장된 모든 밸런스 반환
	var balances []*models.Balance
	for _, balance := range r.balances {
		balances = append(balances, balance)
	}

	log.Printf("ListBalances: returning %d balances from memory", len(balances))
	return balances, nil
}

// UpdateBalance implements Repository.UpdateBalance
func (r *PrometheusRepository) UpdateBalance(ctx context.Context, balance *models.Balance) error {

	// Parse balance
	balanceValue := 0.0
	if strings.TrimSpace(balance.Balance) != "" {
		if val, err := strconv.ParseFloat(strings.TrimSuffix(balance.Balance, " DILL"), 64); err == nil {
			balanceValue = val
		}
	}

	// Parse staking balance
	stakingBalance := 0.0
	if strings.TrimSpace(balance.StakingBalance) != "" {
		if val, err := strconv.ParseFloat(balance.StakingBalance, 64); err == nil {
			stakingBalance = val
		}
	}

	// Parse staked amount
	stakedAmount := 0.0
	if strings.TrimSpace(balance.StakedAmount) != "" {
		if val, err := strconv.ParseFloat(balance.StakedAmount, 64); err == nil {
			stakedAmount = val
		}
	}

	// Parse reward
	reward := 0.0
	if strings.TrimSpace(balance.Reward) != "" {
		if val, err := strconv.ParseFloat(balance.Reward, 64); err == nil {
			reward = val
		}
	}

	// Parse daily reward
	dailyReward := 0.0
	if strings.TrimSpace(balance.DailyReward) != "" {
		if val, err := strconv.ParseFloat(balance.DailyReward, 64); err == nil {
			dailyReward = val
		}
	}

	// Parse latest income
	latestIncome := 0.0
	if strings.TrimSpace(balance.LatestIncome) != "" {
		if val, err := strconv.ParseFloat(balance.LatestIncome, 64); err == nil {
			latestIncome = val
		}
	}

	// Parse lastEpoch to float64
	lastEpoch := 0.0
	if balance.LastEpoch != "" {
		if epochValue, err := strconv.ParseFloat(balance.LastEpoch, 64); err == nil {
			lastEpoch = epochValue
		}
	}

	// Parse lastRewardTime to unix timestamp
	lastRewardTime := float64(time.Now().Unix()) // Default to current time
	if balance.LastRewardTime != "" {
		if rewardTime, err := time.Parse(time.RFC3339, balance.LastRewardTime); err == nil {
			lastRewardTime = float64(rewardTime.Unix())
		}
	}

	// Convert pool stats to float64
	poolCreatedCount := float64(balance.PoolCreatedCount)

	poolParticipatedCount := float64(balance.PoolParticipatedCount)

	// 기본 메트릭 업데이트 (validator_index 유무와 관계없이 항상 표시)
	r.client.UpdateBasicMetrics(
		balance.Address,
		balance.Label,
		balanceValue,
		stakedAmount,
		reward,
		poolCreatedCount,
		poolParticipatedCount,
		lastRewardTime,
	)

	// validator_index가 있는 경우에만 특정 메트릭 업데이트
	if balance.ValidatorIndex != "" && strings.TrimSpace(balance.ValidatorIndex) != "" {
		r.client.UpdateValidatorRelatedMetrics(
			balance.Address,
			balance.Label,
			stakingBalance,
			dailyReward,
			latestIncome,
			lastEpoch,
		)
	}

	return nil
}

// DeleteBalance implements Repository.DeleteBalance
func (r *PrometheusRepository) DeleteBalance(ctx context.Context, address string) error {
	// Prometheus doesn't support deletion of metrics
	return nil
}

// SaveValidatorReward implements Repository.SaveValidatorReward
func (r *PrometheusRepository) SaveValidatorReward(ctx context.Context, reward *models.ValidatorReward) error {
	return r.UpdateValidatorReward(ctx, reward)
}

// GetValidatorReward implements Repository.GetValidatorReward
func (r *PrometheusRepository) GetValidatorReward(ctx context.Context, validatorIdx string) (*models.ValidatorReward, error) {
	// Prometheus is not designed for data retrieval
	return nil, nil
}

// ListValidatorRewards implements Repository.ListValidatorRewards
func (r *PrometheusRepository) ListValidatorRewards(ctx context.Context) ([]*models.ValidatorReward, error) {
	// Prometheus is not designed for data retrieval
	return nil, nil
}

// UpdateValidatorReward implements Repository.UpdateValidatorReward
func (r *PrometheusRepository) UpdateValidatorReward(ctx context.Context, reward *models.ValidatorReward) error {

	// Determine if validator is active based on status or reward date
	isActive := false

	// 먼저 status로 확인
	if reward.Status != "" {
		isActive = strings.Contains(strings.ToLower(reward.Status), "active")
	} else if reward.Date != "" {
		// status 정보가 없는 경우 date로 확인
		// If we have a date, consider the validator active
		isActive = true

		// Check if the date is recent (within last 24 hours)
		if rewardTime, err := time.Parse(time.RFC3339, reward.Date); err == nil {
			// If the reward time is more than 24 hours old, validator might be inactive
			if time.Since(rewardTime) > 24*time.Hour {
				isActive = false
			}
		}
	}

	// Parse lastEpoch to float64
	lastEpoch := 0.0
	if reward.LastEpoch != "" {
		if epochValue, err := strconv.ParseFloat(reward.LastEpoch, 64); err == nil {
			lastEpoch = epochValue
		}
	}

	// Parse lastRewardTime to unix timestamp
	lastRewardTime := float64(time.Now().Unix()) // Default to current time
	if reward.Date != "" {
		if rewardTime, err := time.Parse(time.RFC3339, reward.Date); err == nil {
			lastRewardTime = float64(rewardTime.Unix())
		}
	}

	// Parse validator balance
	validatorBalance := 0.0
	if reward.Balance != "" {
		if val, err := strconv.ParseFloat(reward.Balance, 64); err == nil {
			validatorBalance = val / 1e9 // Convert to DILL (assuming balance is in Gwei)
		}
	}

	// For now, use validator index as label
	r.client.UpdateValidatorMetrics(
		reward.ValidatorIdx,
		reward.UserLabel,
		reward.LastReward,
		validatorBalance,
		isActive,
		lastEpoch,
		lastRewardTime,
		reward.Status,
	)

	return nil
}

// DeleteValidatorReward implements Repository.DeleteValidatorReward
func (r *PrometheusRepository) DeleteValidatorReward(ctx context.Context, validatorIdx string) error {
	// Prometheus doesn't support deletion of metrics
	return nil
}

// RecordBalanceMetric implements Repository.RecordBalanceMetric
func (r *PrometheusRepository) RecordBalanceMetric(balance *models.Balance) error {
	return r.UpdateBalance(context.Background(), balance)
}

// RecordValidatorRewardMetric implements Repository.RecordValidatorRewardMetric
func (r *PrometheusRepository) RecordValidatorRewardMetric(reward *models.ValidatorReward) error {
	return r.UpdateValidatorReward(context.Background(), reward)
}

// RecordAPIMetric implements Repository.RecordAPIMetric
func (r *PrometheusRepository) RecordAPIMetric(endpoint string, duration float64, status int) error {
	r.client.RecordAPIMetrics(endpoint, "GET", status, duration)
	return nil
}

// UpdateSummaryMetrics updates the summary metrics with aggregated data
func (r *PrometheusRepository) UpdateSummaryMetrics(ctx context.Context, balances []*models.Balance) error {
	addressCount := len(balances)
	var totalBalance, totalReward, totalStakedAmount float64
	validatorCount := 0
	activeValidatorCount := 0

	// 상태별 밸리데이터 수를 추적하는 맵
	statusCounts := make(map[string]int)

	for _, balance := range balances {
		// Parse balance
		if strings.TrimSpace(balance.Balance) != "" {
			if val, err := strconv.ParseFloat(strings.TrimSuffix(balance.Balance, " DILL"), 64); err == nil {
				totalBalance += val
			}
		}

		// Parse reward
		if strings.TrimSpace(balance.Reward) != "" {
			if val, err := strconv.ParseFloat(balance.Reward, 64); err == nil {
				totalReward += val
			}
		}

		// Parse staked amount
		if strings.TrimSpace(balance.StakedAmount) != "" {
			if val, err := strconv.ParseFloat(balance.StakedAmount, 64); err == nil {
				totalStakedAmount += val
			}
		}

		// Count validators
		if balance.ValidatorIndex != "" && strings.TrimSpace(balance.ValidatorIndex) != "" {
			validatorCount++

			// Count active validators
			if strings.Contains(strings.ToLower(balance.Status), "active") {
				activeValidatorCount++
			}

			// 상태별 카운트 증가
			status := strings.TrimSpace(balance.Status)
			if status == "" {
				status = "unknown"
			}
			statusCounts[status]++
		}
	}

	// Update summary metrics
	r.client.UpdateSummaryMetrics(
		addressCount,
		validatorCount,
		activeValidatorCount,
		totalBalance,
		totalReward,
		totalStakedAmount,
	)

	// 상태별 밸리데이터 수 업데이트
	r.client.UpdateValidatorStatusMetrics(statusCounts)

	return nil
}
