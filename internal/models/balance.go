package models

// Balance represents the balance information for an account
type Balance struct {
	Label                 string `json:"label"`
	Address               string `json:"address"`
	ValidatorAddress      string `json:"validator_address"`
	ValidatorIndex        string `json:"validator_index"`
	Status                string `json:"status"`
	Balance               string `json:"balance"`
	StakingBalance        string `json:"staking_balance"`
	StakedAmount          string `json:"staked_amount"`
	Reward                string `json:"reward"`
	PoolCreatedCount      int    `json:"pool_created_count"`
	PoolParticipatedCount int    `json:"pool_participated_count"`
	LastEpoch             string `json:"last_epoch"`
	LastRewardTime        string `json:"last_reward_time"`
	LatestIncome          string `json:"latest_income"`
	DailyReward           string `json:"daily_reward"`
}

// ValidatorReward represents the reward information for a validator
type ValidatorReward struct {
	ValidatorIdx string  `json:"validator_idx"`
	LastEpoch    string  `json:"last_epoch"`
	LastReward   float64 `json:"last_reward"`
	Date         string  `json:"date"`
	Balance      string  `json:"balance"`
	Status       string  `json:"status"`
	UserLabel    string  `json:"user_label"`
}

// Config represents the application configuration
type Config struct {
	Addresses []Address `json:"addresses"`
}

// Address represents an address configuration
type Address struct {
	Label            string `json:"label"`
	Address          string `json:"address"`
	ValidatorAddress string `json:"validator_address"`
}

// StakerResponse represents the response from the staker API
type StakerResponse struct {
	StakedAmount          int64 `json:"stakedAmount"`
	Reward                int64 `json:"reward"`
	PoolCreatedCount      int   `json:"poolCreatedCount"`
	PoolParticipatedCount int   `json:"poolParticipatedCount"`
}

// ValidatorInfo represents validator information from the API
type ValidatorInfo struct {
	Index   string `json:"index"`
	Status  string `json:"status"`
	Balance string `json:"balance"`
}

// ValidatorDetailResponse represents detailed validator information from the API
type ValidatorDetailResponse struct {
	Result struct {
		Data struct {
			JSON struct {
				ValidatorIdx         string     `json:"validatorIdx"`
				ValidatorPublicKey   string     `json:"validatorPublicKey"`
				EpochIdx             []string   `json:"epochIdx"`
				AggEpochIdx          [][]string `json:"aggEpochIdx"`
				IncomeGWei           []string   `json:"incomeGWei"`
				IncomeGweiDaySum     []int64    `json:"incomeGweiDaySum"`
				IncomeGweiDaySumDate []string   `json:"incomeGweiDaySumDate"`
			} `json:"json"`
		} `json:"data"`
	} `json:"result"`
}
