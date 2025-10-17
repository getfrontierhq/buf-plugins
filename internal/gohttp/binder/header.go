package binder

import (
	"context"
	"fmt"
	"strings"
)

func (d *RequestDecoder) BindHeader() {
	params := d.Request.Header

	ctx := d.Request.Context()
	for key := range params {
		val := params.Get(key)
		key = strings.ToLower(key)
		if strings.HasPrefix(key, "x-") {
			key = strings.TrimPrefix(key, "x-")
			ctx = context.WithValue(ctx, key, val)
			continue
		}

		key = fmt.Sprintf("%s%s", contextHeaderPrefix, key)
		ctx = context.WithValue(ctx, key, val)
	}

	*d.Request = *d.Request.WithContext(ctx)
}

func (e *RequestEncoder) BindHeader() {
	for key, val := range e.Opts.Headers {
		e.Request.Header.Set(key, fmt.Sprintf("%v", val))
	}

	if shouldHaveBody(e.Request.Method) {
		e.Request.Header.Set("Content-Type", e.Opts.ContentType.String())
	} else {
		e.Request.Header.Del("Content-Type")
	}

	if e.Opts.Operation != "" {
		e.Request.Header.Set("X-Operation", e.Opts.Operation)
	}

	if e.Opts.RequestID != "" {
		e.Request.Header.Set("X-Request-ID", e.Opts.RequestID)
	}
}
