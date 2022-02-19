package main

import (
	"context"

	"github.com/go-redis/redis/extra/redisotel/v8"
	"github.com/go-redis/redis/v8"
)

// for the sake simplicity only single redis client connection
var (
	redisClient    *redis.Client
	redisTelemetry bool
)

func RedisConnContext(ctx context.Context) *redis.Conn {
	conn := redisClient.Conn(ctx)

	// add tracing hook when enabled
	if redisTelemetry {
		conn.AddHook(redisotel.NewTracingHook())
	}

	return conn
}

func OpenRedis(s *Settings) error {
	if !s.Conn.Redis.Enabled {
		return nil
	}

	redisTelemetry = s.Telemetry.Enabled
	redisClient = redis.NewClient(&redis.Options{
		Addr:            s.Conn.Redis.Addr,
		Username:        s.Conn.Redis.Username,
		Password:        s.Conn.Redis.Password,
		DB:              s.Conn.Redis.DB,
		MaxRetries:      s.Conn.Redis.MaxRetries,
		MinRetryBackoff: s.Conn.Redis.MinRetryBackoff,
		MaxRetryBackoff: s.Conn.Redis.MaxRetryBackoff,
		DialTimeout:     s.Conn.Redis.DialTimeout,
		ReadTimeout:     s.Conn.Redis.ReadTimeout,
		WriteTimeout:    s.Conn.Redis.WriteTimeout,
		PoolFIFO:        false,
		PoolSize:        s.Conn.Redis.PoolSize,
		MinIdleConns:    s.Conn.Redis.MinIdleConns,
		// MaxConnAge:         0,
		// PoolTimeout:        0,
		// IdleTimeout:        0,
		// IdleCheckFrequency: 0,
		// TLSConfig:          &tls.Config{},
		// Limiter:            nil,
	})

	return nil
}

func CloseRedis() error {
	if redisClient == nil {
		return nil
	}

	return redisClient.Close()
}
