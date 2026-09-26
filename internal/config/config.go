package config

import "os"

type Config struct {
	HTTPAddr      string
	TemporalHost  string
	TemporalTaskQ string
	TemporalNS    string
}

func Load() Config {
	return Config{
		HTTPAddr:      getEnv("HTTP_ADDR", ":8080"),
		TemporalHost:  getEnv("TEMPORAL_HOST", "localhost:7233"),
		TemporalTaskQ: getEnv("TEMPORAL_TASK_QUEUE", "reservation-task-queue"),
		TemporalNS:    getEnv("TEMPORAL_NAMESPACE", "default"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
