package repository

import (
	"context"
	"dill-monitor/internal/models"
)

// Repository defines the interface for data storage and retrieval
type Repository interface {
	// Balance operations
	SaveBalance(ctx context.Context, balance *models.Balance) error
	GetBalance(ctx context.Context, address string) (*models.Balance, error)
	ListBalances(ctx context.Context) ([]*models.Balance, error)
	UpdateBalance(ctx context.Context, balance *models.Balance) error
	DeleteBalance(ctx context.Context, address string) error

	// Validator reward operations
	SaveValidatorReward(ctx context.Context, reward *models.ValidatorReward) error
	GetValidatorReward(ctx context.Context, validatorIdx string) (*models.ValidatorReward, error)
	ListValidatorRewards(ctx context.Context) ([]*models.ValidatorReward, error)
	UpdateValidatorReward(ctx context.Context, reward *models.ValidatorReward) error
	DeleteValidatorReward(ctx context.Context, validatorIdx string) error

	// Metrics operations
	RecordBalanceMetric(balance *models.Balance) error
	RecordValidatorRewardMetric(reward *models.ValidatorReward) error
	RecordAPIMetric(endpoint string, duration float64, status int) error

	// Summary metrics operations
	UpdateSummaryMetrics(ctx context.Context, balances []*models.Balance) error
}
