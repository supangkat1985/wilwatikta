package main

import "time"

type JaegerSettings struct {
	URL string `mapstructure:"url"`
}

type PrometheusSettings struct {
	Addr string `mapstructure:"addr"`
}

type TelemetrySettings struct {
	Enabled    bool               `mapstructure:"enabled"`
	Jaeger     JaegerSettings     `mapstructure:"jaeger"`
	Prometheus PrometheusSettings `mapstructure:"prometheus"`
}

type RedisSettings struct {
	Enabled         bool          `mapstructure:"enabled"`
	Addr            string        `mapstructure:"addr"`
	Username        string        `mapstructure:"username"`
	Password        string        `mapstructure:"password"`
	DB              int           `mapstructure:"db"`
	MaxRetries      int           `mapstructure:"max_retries"`
	MinRetryBackoff time.Duration `mapstructure:"min_retry_backoff"`
	MaxRetryBackoff time.Duration `mapstructure:"max_retry_backoff"`
	DialTimeout     time.Duration `mapstructure:"dial_timeout"`
	ReadTimeout     time.Duration `mapstructure:"read_timeout"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout"`
	PoolSize        int           `mapstructure:"pool_size"`
	MinIdleConns    int           `mapstructure:"min_idle_conns"`
	// TODO: expose more configuration
}

type ConnectionSettings struct {
	Redis RedisSettings `mapstructure:"redis"`
}

type Settings struct {
	Conn      ConnectionSettings `mapstrucutre:"conn"`
	Telemetry TelemetrySettings  `mapstructure:"telemetry"`
}
