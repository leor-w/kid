package elasticsearch

import (
	"net/http"
	"time"
)

type Options struct {
	Addresses []string // A list of Elasticsearch nodes to use.
	Username  string   // Username for HTTP Basic Authentication.
	Password  string   // Password for HTTP Basic Authentication.

	CloudID                string // Endpoint for the Elastic Service (https://elastic.co/cloud).
	APIKey                 string // Base64-encoded token for authorization; if set, overrides username/password and service token.
	ServiceToken           string // Service token for authorization; if set, overrides username/password.
	CertificateFingerprint string // SHA256 hex fingerprint given by Elasticsearch on first launch.

	Header http.Header // Global HTTP request header.

	// PEM-encoded certificate authorities.
	// When set, an empty certificate pool will be created, and the certificates will be appended to it.
	// The option is only valid when the transport is not specified, or when it's http.Transport.
	CACert string

	RetryOnStatus []int                           // List of status codes for retry. Default: 502, 503, 504.
	DisableRetry  bool                            // Default: false.
	MaxRetries    int                             // Default: 3.
	RetryOnError  func(*http.Request, error) bool // Optional function allowing to indicate which error should be retried. Default: nil.

	CompressRequestBody bool // Default: false.

	DiscoverNodesOnStart  bool          // Discover nodes when initializing the client. Default: false.
	DiscoverNodesInterval time.Duration // Discover nodes periodically. Default: disabled.

	EnableMetrics           bool // Enable the metrics collection.
	EnableDebugLogger       bool // Enable the debug logging.
	EnableCompatibilityMode bool // Enable sends compatibility header

	DisableMetaHeader bool // Disable the additional "X-Elastic-Client-Meta" HTTP header.
}

type Option func(*Options)

func WithAddresses(addresses []string) Option {
	return func(o *Options) {
		o.Addresses = addresses
	}
}

func WithUsername(username string) Option {
	return func(o *Options) {
		o.Username = username
	}
}

func WithPassword(password string) Option {
	return func(o *Options) {
		o.Password = password
	}
}

func WithCloudID(cloudId string) Option {
	return func(o *Options) {
		o.CloudID = cloudId
	}
}

func WithApiKey(apiKey string) Option {
	return func(o *Options) {
		o.APIKey = apiKey
	}
}

func WithServiceToken(serviceToken string) Option {
	return func(o *Options) {
		o.ServiceToken = serviceToken
	}
}

func WithCertificateFingerprint(certificateFingerprint string) Option {
	return func(o *Options) {
		o.CertificateFingerprint = certificateFingerprint
	}
}

func WithHeader(header map[string][]string) Option {
	return func(o *Options) {
		o.Header = header
	}
}

func WithCACert(caCert string) Option {
	return func(o *Options) {
		o.CACert = caCert
	}
}

func WithRetryOnStatus(retryOnStatus []int) Option {
	return func(o *Options) {
		o.RetryOnStatus = retryOnStatus
	}
}

func WithDisableRetry(disableRetry bool) Option {
	return func(o *Options) {
		o.DisableRetry = disableRetry
	}
}

func WithMaxRetries(maxRetries int) Option {
	return func(o *Options) {
		o.MaxRetries = maxRetries
	}
}

func WithCompressRequestBody(compressRequestBody bool) Option {
	return func(o *Options) {
		o.CompressRequestBody = compressRequestBody
	}
}

func WithDiscoverNodesOnStart(discoverNodesOnStart bool) Option {
	return func(o *Options) {
		o.DiscoverNodesOnStart = discoverNodesOnStart
	}
}

func WithDiscoverNodesInterval(discoverNodesInterval time.Duration) Option {
	return func(o *Options) {
		o.DiscoverNodesInterval = discoverNodesInterval * time.Second
	}
}

func WithEnableMetrics(enableMetrics bool) Option {
	return func(o *Options) {
		o.EnableMetrics = enableMetrics
	}
}

func WithEnableDebugLogger(enableDebugLogger bool) Option {
	return func(o *Options) {
		o.EnableDebugLogger = enableDebugLogger
	}
}

func WithEnableCompatibilityMode(enableCompatibilityMode bool) Option {
	return func(o *Options) {
		o.EnableCompatibilityMode = enableCompatibilityMode
	}
}

func WithDisableMetaHeader(disableMetaHeader bool) Option {
	return func(o *Options) {
		o.DisableMetaHeader = disableMetaHeader
	}
}
