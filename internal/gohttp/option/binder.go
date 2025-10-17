package option

type BinderOptions struct {
	Headers     map[string]any
	ContentType ContentType
	Operation   string
	RequestID   string
}

type BinderOption func(*BinderOptions)

func NewBinderOptions(options ...BinderOption) *BinderOptions {
	o := BinderOptions{
		Headers:     make(map[string]any),
		ContentType: ContentTypeApplicationJson,
	}

	for _, option := range options {
		option(&o)
	}

	return &o
}

func WithHeader(key string, val any) BinderOption {
	return func(o *BinderOptions) {
		o.Headers[key] = val
	}
}

func WithHeaders(headers map[string]any) BinderOption {
	return func(o *BinderOptions) {
		o.Headers = headers
	}
}

func WithContentType(contentType ContentType) BinderOption {
	return func(o *BinderOptions) {
		o.ContentType = contentType
	}
}

func WithOperation(operation string) BinderOption {
	return func(o *BinderOptions) {
		o.Operation = operation
	}
}

func WithRequestID(requestID string) BinderOption {
	return func(o *BinderOptions) {
		o.RequestID = requestID
	}
}
