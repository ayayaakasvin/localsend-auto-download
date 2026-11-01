package domain

import "time"

type ShareLinkConfig struct {
	Host      string
	SessionID string
	PinCode   string
}

func NewShareLinkConfig(
	host string,
	sessionID string,
	pinCode string,
) *ShareLinkConfig {
	return &ShareLinkConfig{
		Host:      host,
		SessionID: sessionID,
		PinCode:   pinCode,
	}
}

type TimeoutConfig struct {
	Timeout     time.Duration
	IdleTimeout time.Duration
}

func NewTimeoutConfig(
	t time.Duration,
	it time.Duration,
) *TimeoutConfig {
	return &TimeoutConfig{
		Timeout:     t,
		IdleTimeout: it,
	}
}
