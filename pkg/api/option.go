package api

import (
	"fmt"
	"time"

	"github.com/wakatime/wakatime-cli/pkg/heartbeat"
)

// Option is a functional option for Client.
type Option func(*Client)

// WithAuth adds authentication via Authorization header.
func WithAuth(auth BasicAuth) (Option, error) {
	authHeaderValue, err := auth.HeaderValue()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve auth header value: %w", err)
	}

	return func(c *Client) {
		c.authHeader = authHeaderValue
	}, nil
}

// WithHostname sets the X-Machine-Name header to the passed in hostname.
func WithHostname(hostname string) Option {
	return func(c *Client) {
		c.machineNameHeader = hostname
	}
}

// WithTimeout configures a timeout for all requests.
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.client.Timeout = timeout
	}
}

// WithUserAgentUnknownPlugin sets the User-Agent header on all requests,
// including default value for plugin.
func WithUserAgentUnknownPlugin() Option {
	return WithUserAgent("Unknown/0")
}

// WithUserAgent sets the User-Agent header on all requests, including the passed
// in value for plugin.
func WithUserAgent(plugin string) Option {
	userAgent := heartbeat.UserAgent(plugin)

	return func(c *Client) {
		c.userAgentHeader = userAgent
	}
}
