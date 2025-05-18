package service

import (
	"context"
	"dill-monitor/internal/models"
	"dill-monitor/internal/repository"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// BalanceService handles balance-related business logic
type BalanceService struct {
	repo repository.Repository
}

// NewBalanceService creates a new balance service
func NewBalanceService(repo repository.Repository) *BalanceService {
	return &BalanceService{
		repo: repo,
	}
}

// ProcessAddress processes a single address and updates its balance information
func (s *BalanceService) ProcessAddress(ctx context.Context, addr models.Address) (*models.Balance, error) {
	// Get wallet balance
	balance, err := s.getWalletBalance(addr.Address)
	if err != nil {
		return nil, fmt.Errorf("error getting wallet balance: %v", err)
	}

	// Get staker info
	stakerInfo, err := s.getStakerInfo(addr.Address)
	if err != nil {
		return nil, fmt.Errorf("error getting staker info: %v", err)
	}

	// Create balance object
	balanceObj := &models.Balance{
		Label:                 addr.Label,
		Address:               addr.Address,
		ValidatorAddress:      addr.ValidatorAddress,
		Balance:               balance,
		StakingBalance:        "0",
		StakedAmount:          fmt.Sprintf("%.4f", float64(stakerInfo.StakedAmount)/1e9),
		Reward:                fmt.Sprintf("%.4f", float64(stakerInfo.Reward)/1e9),
		PoolCreatedCount:      stakerInfo.PoolCreatedCount,
		PoolParticipatedCount: stakerInfo.PoolParticipatedCount,
		LastEpoch:             "0",
		LastRewardTime:        time.Now().Format(time.RFC3339),
		LatestIncome:          "0",
		DailyReward:           "0",
	}

	// If validator address exists, get validator info
	if addr.ValidatorAddress != "" {
		validatorInfo, err := s.getValidatorInfo(addr.ValidatorAddress)
		if err != nil {
			return nil, fmt.Errorf("error getting validator info: %v", err)
		}

		if validatorInfo != nil {
			balanceObj.ValidatorIndex = validatorInfo.Index
			balanceObj.Status = validatorInfo.Status

			// Convert balance from string to float64
			validatorBalance, err := strconv.ParseFloat(validatorInfo.Balance, 64)
			if err == nil {
				balanceObj.StakingBalance = fmt.Sprintf("%.3f", validatorBalance/1e9)
			} else {
				balanceObj.StakingBalance = "0"
			}

			// Get validator details
			if balanceObj.ValidatorIndex != "" {
				details, err := s.getValidatorDetails(balanceObj.ValidatorIndex)
				if err == nil && details != nil {
					if len(details.Result.Data.JSON.EpochIdx) > 0 {
						balanceObj.LastEpoch = details.Result.Data.JSON.EpochIdx[len(details.Result.Data.JSON.EpochIdx)-1]
					}

					if len(details.Result.Data.JSON.IncomeGWei) > 0 {
						latestIncome := details.Result.Data.JSON.IncomeGWei[len(details.Result.Data.JSON.IncomeGWei)-1]
						income, err := strconv.ParseFloat(latestIncome, 64)
						if err == nil {
							balanceObj.LatestIncome = fmt.Sprintf("%.4f", income/1e9)
						}
					}

					if len(details.Result.Data.JSON.IncomeGweiDaySum) >= 2 {
						dailyReward := float64(details.Result.Data.JSON.IncomeGweiDaySum[0] + details.Result.Data.JSON.IncomeGweiDaySum[1])
						balanceObj.DailyReward = fmt.Sprintf("%.4f", dailyReward/1e9)
					} else if len(details.Result.Data.JSON.IncomeGweiDaySum) == 1 {
						// 하루 데이터만 있는 경우 그 값의 2배를 일일 보상 예상치로 사용
						dailyReward := float64(details.Result.Data.JSON.IncomeGweiDaySum[0]) * 2
						balanceObj.DailyReward = fmt.Sprintf("%.4f", dailyReward/1e9)
					} else if len(details.Result.Data.JSON.IncomeGWei) > 0 {
						// 일별 합계 데이터가 없지만 수입 데이터가 있는 경우, 마지막 수입 값에 기반하여 추정
						lastIncome, err := strconv.ParseFloat(details.Result.Data.JSON.IncomeGWei[len(details.Result.Data.JSON.IncomeGWei)-1], 64)
						if err == nil {
							// 하루 동안 약 225개의 epoch가 발생한다고 가정 (예상치)
							dailyReward := lastIncome * 225
							balanceObj.DailyReward = fmt.Sprintf("%.4f", dailyReward/1e9)
						} else {
							balanceObj.DailyReward = balanceObj.LatestIncome // 단일 수입을 일일 보상으로 설정
						}
					} else {
						log.Printf("Warning: No data available to calculate DailyReward for validator %s", balanceObj.ValidatorIndex)
						balanceObj.DailyReward = "0"
					}

					balanceObj.LastRewardTime = time.Now().Format(time.RFC3339)

				}
			}
		}
	} else {
		// For addresses without validator, set current time for LastRewardTime
		balanceObj.LastRewardTime = time.Now().Format(time.RFC3339)
		log.Printf("Non-validator address %s: setting LastRewardTime to %s", addr.Address, balanceObj.LastRewardTime)
	}

	// Save balance to repository
	if err := s.repo.SaveBalance(ctx, balanceObj); err != nil {
		return nil, fmt.Errorf("error saving balance: %v", err)
	}

	return balanceObj, nil
}

// getWalletBalance retrieves the wallet balance from the API
func (s *BalanceService) getWalletBalance(address string) (string, error) {
	url := fmt.Sprintf("https://alps.dill.xyz/api/trpc/stats.getBalance?input={\"json\":{\"address\":\"%s\"}}", address)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var response struct {
		Result struct {
			Data struct {
				JSON struct {
					Balance string `json:"balance"`
				} `json:"json"`
			} `json:"data"`
		} `json:"result"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return "", err
	}

	balance, err := strconv.ParseFloat(response.Result.Data.JSON.Balance, 64)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%.10f DILL", balance/1e18), nil
}

// getStakerInfo retrieves staker information from the API
func (s *BalanceService) getStakerInfo(address string) (*models.StakerResponse, error) {
	url := "https://staker.dill.xyz/api?Action=GetUserInfo"
	requestBody := fmt.Sprintf(`{"Action":"GetUserInfo","Address":"%s"}`, address)

	req, err := http.NewRequest("POST", url, strings.NewReader(requestBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var stakerResponse models.StakerResponse
	if err := json.Unmarshal(body, &stakerResponse); err != nil {
		return &models.StakerResponse{
			StakedAmount:          0,
			Reward:                0,
			PoolCreatedCount:      0,
			PoolParticipatedCount: 0,
		}, nil
	}

	return &stakerResponse, nil
}

// getValidatorInfo retrieves validator information from the API
func (s *BalanceService) getValidatorInfo(validatorAddress string) (*models.ValidatorInfo, error) {
	url := fmt.Sprintf("https://alps.dill.xyz/api/trpc/stats.getAllValidators?input={\"json\":{\"page\":1,\"limit\":25,\"pubkey\":\"%s\"}}", validatorAddress)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response struct {
		Result struct {
			Data struct {
				JSON struct {
					Data []models.ValidatorInfo `json:"data"`
				} `json:"json"`
			} `json:"data"`
		} `json:"result"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	if len(response.Result.Data.JSON.Data) > 0 {
		return &response.Result.Data.JSON.Data[0], nil
	}

	return nil, nil
}

// getValidatorDetails retrieves detailed validator information from the API
func (s *BalanceService) getValidatorDetails(validatorIdx string) (*models.ValidatorDetailResponse, error) {
	// Use UTC time to avoid timezone issues
	now := time.Now().UTC()
	endTime := now.UnixNano() / int64(time.Millisecond)
	startTime := now.Add(-24*time.Hour).UnixNano() / int64(time.Millisecond)

	inputJSON := fmt.Sprintf(`{"json":{"item":"only to meet the parameter requirements of tRPC","validatorKey":"%s","validatorIdx":"%s","validatorIsStr":false,"startTime":%d,"endTime":%d}}`,
		validatorIdx, validatorIdx, startTime, endTime)

	// URL 인코딩 적용
	encodedInput := url.QueryEscape(inputJSON)

	url := fmt.Sprintf("https://alps.dill.xyz/api/trpc/stats.getValidatorDetailByKeyOrIdx?input=%s", encodedInput)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making HTTP request for validator details %s: %v", validatorIdx, err)
	}
	defer resp.Body.Close()

	// HTTP 상태 코드 확인
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid HTTP status code %d for validator %s", resp.StatusCode, validatorIdx)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body for validator details %s: %v", validatorIdx, err)
	}

	// Check for empty response
	if len(body) == 0 {
		return nil, fmt.Errorf("empty response body for validator details %s", validatorIdx)
	}

	// Check if the response is valid JSON
	if !json.Valid(body) {
		// 유효하지 않은 JSON이면 처음 100자만 로그로 출력
		invalidJSON := string(body)
		if len(invalidJSON) > 100 {
			invalidJSON = invalidJSON[:100] + "..."
		}
		return nil, fmt.Errorf("invalid JSON response for validator details %s: %s", validatorIdx, invalidJSON)
	}

	var detailResponse models.ValidatorDetailResponse
	if err := json.Unmarshal(body, &detailResponse); err != nil {
		return nil, fmt.Errorf("error parsing JSON response for validator details %s: %v", validatorIdx, err)
	}

	return &detailResponse, nil
}

// UpdateSummaryMetrics updates the summary metrics with aggregated data from all balances
func (s *BalanceService) UpdateSummaryMetrics(ctx context.Context, balances []*models.Balance) error {
	return s.repo.UpdateSummaryMetrics(ctx, balances)
}

// ProcessValidators processes validator information and updates metrics
func (s *BalanceService) ProcessValidators(ctx context.Context) error {
	log.Printf("Starting to process validator information...")

	// Get all validators with balances
	balances, err := s.repo.ListBalances(ctx)
	if err != nil {
		log.Printf("Error listing balances: %v", err)
		return err
	}

	validatorCount := 0

	// Process each validator
	for _, balance := range balances {
		// Skip non-validators
		if balance.ValidatorIndex == "" || strings.TrimSpace(balance.ValidatorIndex) == "" {
			continue
		}

		validatorCount++

		// Create validator reward object
		validatorReward := &models.ValidatorReward{
			ValidatorIdx: balance.ValidatorIndex,
			LastEpoch:    balance.LastEpoch,
			Date:         balance.LastRewardTime,
			Status:       balance.Status,
			Balance:      balance.StakingBalance,
			UserLabel:    balance.Label,
		}

		// Parse LastReward
		lastReward := 0.0
		if balance.LatestIncome != "" {
			if val, err := strconv.ParseFloat(balance.LatestIncome, 64); err == nil {
				lastReward = val
			}
		}
		validatorReward.LastReward = lastReward

		// Update validator reward
		if err := s.repo.SaveValidatorReward(ctx, validatorReward); err != nil {
			log.Printf("Error saving validator reward: %v", err)
		}
	}

	return nil
}
