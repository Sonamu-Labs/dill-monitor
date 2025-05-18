package models

// ServerConfig represents server specific configuration
type ServerConfig struct {
	MetricsPort int    `json:"metricsPort"`
	LogLevel    string `json:"logLevel"`
	Host        string `json:"host"`
	// 기타 서버 관련 설정 추가 가능
}
