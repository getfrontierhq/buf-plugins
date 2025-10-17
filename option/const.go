package option

type (
	ContentType string
)

const (
	ContentTypeApplicationJson ContentType = "application/json"

	ContentTypeHeader   = "Content-Type"
	AuthorizationHeader = "Authorization"
	UserAgentHeader     = "User-Agent"
	XRequestIDHeader    = "X-Request-ID"
)

func (c ContentType) String() string {
	return string(c)
}
