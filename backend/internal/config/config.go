package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config 集中管理应用配置，全部通过环境变量注入。
type Config struct {
	AppEnv      string
	ServerPort  string
	MongoURI    string
	MongoDBName string
	JWTSecret   string
	TokenExpire time.Duration

	MongoConnectTimeout time.Duration

	CORSAllowedOrigins []string
	RateLimitPerMin    int
	LoginRateLimit     int
	LoginRateWindow    time.Duration

	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
	ShutdownWait time.Duration
}

// Load 从环境变量读取配置并提供安全默认值。
func Load() *Config {
	return &Config{
		AppEnv:      getEnv("APP_ENV", "development"),
		ServerPort:  getEnv("SERVER_PORT", "8080"),
		MongoURI:    getEnv("MONGO_URI", "mongodb://artvault:artvault_pwd@127.0.0.1:27017/artvault?authSource=admin"),
		MongoDBName: getEnv("MONGO_DB_NAME", "artvault"),
		JWTSecret:   getEnv("JWT_SECRET", ""),
		TokenExpire: time.Duration(getEnvInt("TOKEN_EXPIRE_HOURS", 72)) * time.Hour,

		MongoConnectTimeout: time.Duration(getEnvInt("MONGO_CONNECT_TIMEOUT_SEC", 30)) * time.Second,

		CORSAllowedOrigins: splitCSV(getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:18808")),
		RateLimitPerMin:    getEnvInt("RATE_LIMIT_PER_MIN", 120),
		LoginRateLimit:     getEnvInt("LOGIN_RATE_LIMIT", 10),
		LoginRateWindow:    time.Duration(getEnvInt("LOGIN_RATE_WINDOW_SEC", 60)) * time.Second,

		ReadTimeout:  time.Duration(getEnvInt("SERVER_READ_TIMEOUT_SEC", 10)) * time.Second,
		WriteTimeout: time.Duration(getEnvInt("SERVER_WRITE_TIMEOUT_SEC", 20)) * time.Second,
		IdleTimeout:  time.Duration(getEnvInt("SERVER_IDLE_TIMEOUT_SEC", 60)) * time.Second,
		ShutdownWait: time.Duration(getEnvInt("SHUTDOWN_WAIT_SEC", 10)) * time.Second,
	}
}

// Validate 校验配置；生产/预发环境强制安全密钥与 CORS 白名单。
func (c *Config) Validate() error {
	if c.AppEnv == "production" || c.AppEnv == "staging" {
		if c.JWTSecret == "" || len(c.JWTSecret) < 32 || c.JWTSecret == "change-me-artvault" {
			return fmt.Errorf("JWT_SECRET 过短或仍为默认值，生产环境必须设置为至少 32 字符的随机字符串")
		}
		for _, origin := range c.CORSAllowedOrigins {
			if origin == "*" {
				return fmt.Errorf("生产环境禁止将 CORS_ALLOWED_ORIGINS 配置为 *")
			}
		}
	}
	if c.ServerPort == "" {
		return fmt.Errorf("SERVER_PORT 不能为空")
	}
	return nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}

func splitCSV(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
