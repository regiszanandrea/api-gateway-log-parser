package apigateway

type Log struct {
	Request             Request             `json:"request"`
	UpstreamURI         string              `json:"upstream_uri"`
	Response            Response            `json:"response"`
	AuthenticatedEntity AuthenticatedEntity `json:"authenticated_entity"`
	Route               Route               `json:"route"`
	Service             Service             `json:"service"`
	Latencies           Latencies           `json:"latencies"`
	ClientIP            string              `json:"client_ip"`
	StartedAt           int64               `json:"started_at"`
	ServiceID           string              `json:"service_id"`
	CustomerID          string              `json:"customer_id"`
}

type Request struct {
	Method  string         `json:"method"`
	URI     string         `json:"uri"`
	URL     string         `json:"url"`
	Size    int            `json:"size"`
	Headers RequestHeaders `json:"headers"`
}

type RequestHeaders struct {
	Accept    string `json:"accept"`
	Host      string `json:"host"`
	UserAgent string `json:"user-agent"`
}

type Response struct {
	Status  int             `json:"status"`
	Size    int             `json:"size"`
	Headers ResponseHeaders `json:"headers"`
}

type ResponseHeaders struct {
	ContentLength                 string `json:"Content-Length"`
	Via                           string `json:"via"`
	Connection                    string `json:"Connection"`
	AccessControlAllowCredentials string `json:"access-control-allow-credentials"`
	ContentType                   string `json:"Content-Type"`
	Server                        string `json:"server"`
	AccessControlAllowOrigin      string `json:"access-control-allow-origin"`
}

type AuthenticatedEntity struct {
	ConsumerID struct {
		UUID string `json:"uuid"`
	} `json:"consumer_id"`
}

type Route struct {
	CreatedAt int `json:"created_at"`
	Hosts     interface {
	} `json:"hosts"`
	ID            string   `json:"id"`
	Methods       []string `json:"methods"`
	Paths         []string `json:"paths"`
	PreserveHost  bool     `json:"preserve_host"`
	Protocols     []string `json:"protocols"`
	RegexPriority int      `json:"regex_priority"`
	Service       struct {
		ID string `json:"id"`
	} `json:"service"`
	StripPath bool `json:"strip_path"`
	UpdatedAt int  `json:"updated_at"`
}

type Service struct {
	ConnectTimeout int    `json:"connect_timeout"`
	CreatedAt      int    `json:"created_at"`
	Host           string `json:"host"`
	ID             string `json:"id"`
	Name           string `json:"name"`
	Path           string `json:"path"`
	Port           int    `json:"port"`
	Protocol       string `json:"protocol"`
	ReadTimeout    int    `json:"read_timeout"`
	Retries        int    `json:"retries"`
	UpdatedAt      int    `json:"updated_at"`
	WriteTimeout   int    `json:"write_timeout"`
}

type Latencies struct {
	Proxy   int `json:"proxy"`
	Gateway int `json:"gateway"`
	Request int `json:"request"`
}

type ApiGatewayLogService interface {
	Parse(path string) error
}
