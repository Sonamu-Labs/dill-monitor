package main

import (
	"context"
	"dill-monitor/internal/models"
	"dill-monitor/internal/repository"
	"dill-monitor/pkg/metrics"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// 로그 설정
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting validator status info test...")

	// 메트릭 클라이언트 및 리포지토리 초기화
	promClient := metrics.NewPrometheusClient()
	promRepo := repository.NewPrometheusRepository(promClient)

	// 메트릭 서버 시작
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		addr := ":9091"
		log.Printf("Starting metrics server on %s", addr)
		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Fatalf("Failed to start metrics server: %v", err)
		}
	}()

	// 컨텍스트 생성
	ctx := context.Background()

	// 테스트용 밸리데이터 보상 객체 생성
	validatorReward := &models.ValidatorReward{
		ValidatorIdx: "17021",
		LastEpoch:    "57006",
		LastReward:   0.0804,
		Date:         "2025-05-19T05:09:11+09:00",
		Balance:      "36000.566",
		Status:       "active_ongoing",
		UserLabel:    "Test Validator",
	}

	// 저장소에 밸리데이터 보상 정보 저장
	log.Println("Saving validator reward...")
	if err := promRepo.SaveValidatorReward(ctx, validatorReward); err != nil {
		log.Fatalf("Error saving validator reward: %v", err)
	}

	// 테스트용 밸런스 객체 생성
	balance := &models.Balance{
		Label:                 "Test Validator",
		Address:               "0xFEFCa083D55C196605ab3733BF5000c2B43861b4",
		ValidatorAddress:      "0x8e7b68fcc8813303debf815208a8fe3fe4b7fc557ce87e6bd419f6bf05064dde800056d69f11266a7b1d88cc72f5c6af",
		ValidatorIndex:        "17021",
		Status:                "active_ongoing",
		Balance:               "2.6395060217 DILL",
		StakingBalance:        "36000.566",
		StakedAmount:          "37100.0000",
		Reward:                "381.1681",
		PoolCreatedCount:      1,
		PoolParticipatedCount: 2,
		LastEpoch:             "57006",
		LastRewardTime:        "2025-05-19T05:09:11+09:00",
		LatestIncome:          "0.0804",
		DailyReward:           "17.7535",
	}

	// 저장소에 밸런스 정보 저장
	log.Println("Saving balance...")
	if err := promRepo.SaveBalance(ctx, balance); err != nil {
		log.Fatalf("Error saving balance: %v", err)
	}

	// 요약 메트릭 업데이트
	balances := []*models.Balance{balance}

	// 추가 테스트용 밸런스 객체 생성 (다른 라벨 값 사용)
	balance2 := &models.Balance{
		Label:                 "My Validator", // 다른 사용자 지정 라벨
		Address:               "0xcC195833442B8D6142DDA28a31dc4881425Ebf28",
		ValidatorAddress:      "0xb0ec80500de5ad5e8e316f71312ed6d4d2837f744f8b41950b0e570b90ecd6f8126a058ea8da03a10dd740c1e3395212",
		ValidatorIndex:        "17022",
		Status:                "active_ongoing",
		Balance:               "2.4764711682 DILL",
		StakingBalance:        "36000.322",
		StakedAmount:          "37100.0000",
		Reward:                "381.6015",
		PoolCreatedCount:      1,
		PoolParticipatedCount: 2,
		LastEpoch:             "57006",
		LastRewardTime:        "2025-05-19T05:09:11+09:00",
		LatestIncome:          "0.0804",
		DailyReward:           "17.7330",
	}

	// 저장소에 두 번째 밸런스 정보 저장
	log.Println("Saving second balance...")
	if err := promRepo.SaveBalance(ctx, balance2); err != nil {
		log.Fatalf("Error saving balance: %v", err)
	}

	// 두 번째 밸리데이터 보상 객체 생성
	validatorReward2 := &models.ValidatorReward{
		ValidatorIdx: "17022",
		LastEpoch:    "57006",
		LastReward:   0.0804,
		Date:         "2025-05-19T05:09:11+09:00",
		Balance:      "36000.322",
		Status:       "active_ongoing",
		UserLabel:    "My Validator", // 두 번째 사용자 지정 라벨
	}

	// 저장소에 두 번째 밸리데이터 보상 정보 저장
	log.Println("Saving second validator reward...")
	if err := promRepo.SaveValidatorReward(ctx, validatorReward2); err != nil {
		log.Fatalf("Error saving validator reward: %v", err)
	}

	balances = append(balances, balance2)

	// 요약 메트릭 업데이트 (두 개의 밸런스 사용)
	log.Println("Updating summary metrics with 2 validators...")
	if err := promRepo.UpdateSummaryMetrics(ctx, balances); err != nil {
		log.Fatalf("Error updating summary metrics: %v", err)
	}

	// 프로그램 종료 신호 대기
	log.Println("Test completed successfully. Server running at http://localhost:9091/metrics")
	log.Println("Press Ctrl+C to exit...")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Shutting down...")
}
