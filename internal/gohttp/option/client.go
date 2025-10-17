package option

import "time"

type ClientOptions struct {
	BaseURL string
	Timeout time.Duration
}

type ClientOption func(*ClientOptions)

func NewClientOptions(options ...ClientOption) *ClientOptions {
	o := ClientOptions{
		BaseURL: "",
		Timeout: 10 * time.Second,
	}

	for _, option := range options {
		option(&o)
	}

	return &o
}

func WithBaseURL(baseURL string) ClientOption {
	return func(o *ClientOptions) {
		o.BaseURL = baseURL
	}
}

func WithTimeout(timeout time.Duration) ClientOption {
	return func(o *ClientOptions) {
		o.Timeout = timeout
	}
}
