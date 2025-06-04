/* Kong Logs

Contains go types to unmarshall kong logs

Joaquin Badillo
2025-06-04
*/


package models

type KongLog struct {
	Response       Response  `json:"response"`
	Route          Route     `json:"route"`
	Workspace      string    `json:"workspace"`
	WorkspaceName  string    `json:"workspace_name"`
	Tries          []Try     `json:"tries"`
	ClientIP       string    `json:"client_ip"`
	Request        Request   `json:"request"`
	UpstreamURI    string    `json:"upstream_uri"`
	StartedAt      int64     `json:"started_at"`
	Source         string    `json:"source"`
	UpstreamStatus string    `json:"upstream_status"`
	Latencies      Latencies `json:"latencies"`
	Service        Service   `json:"service"`
}

type Response struct {
	Size    int               `json:"size"`
	Headers map[string]string `json:"headers"`
	Status  int               `json:"status"`
}

type Route struct {
	UpdatedAt               int64      `json:"updated_at"`
	Tags                    []string   `json:"tags"`
	ResponseBuffering       bool       `json:"response_buffering"`
	PathHandling            string     `json:"path_handling"`
	Protocols               []string   `json:"protocols"`
	Service                 ServiceRef `json:"service"`
	HTTPSRedirectStatusCode int        `json:"https_redirect_status_code"`
	RegexPriority           int        `json:"regex_priority"`
	Name                    string     `json:"name"`
	ID                      string     `json:"id"`
	StripPath               bool       `json:"strip_path"`
	PreserveHost            bool       `json:"preserve_host"`
	CreatedAt               int64      `json:"created_at"`
	RequestBuffering        bool       `json:"request_buffering"`
	WSID                    string     `json:"ws_id"`
	Paths                   []string   `json:"paths"`
}

type ServiceRef struct {
	ID string `json:"id"`
}

type Try struct {
	BalancerStart     int64   `json:"balancer_start"`
	BalancerStartNS   float64 `json:"balancer_start_ns"`
	IP                string  `json:"ip"`
	BalancerLatency   int     `json:"balancer_latency"`
	Port              int     `json:"port"`
	BalancerLatencyNS int     `json:"balancer_latency_ns"`
}

type Request struct {
	ID          string            `json:"id"`
	Headers     map[string]string `json:"headers"`
	URI         string            `json:"uri"`
	Size        int               `json:"size"`
	Method      string            `json:"method"`
	QueryString map[string]string `json:"querystring"`
	URL         string            `json:"url"`
}

type Latencies struct {
	Kong    int `json:"kong"`
	Proxy   int `json:"proxy"`
	Request int `json:"request"`
	Receive int `json:"receive"`
}

type Service struct {
	WriteTimeout   int    `json:"write_timeout"`
	ReadTimeout    int    `json:"read_timeout"`
	UpdatedAt      int64  `json:"updated_at"`
	Host           string `json:"host"`
	Name           string `json:"name"`
	ID             string `json:"id"`
	Port           int    `json:"port"`
	Enabled        bool   `json:"enabled"`
	CreatedAt      int64  `json:"created_at"`
	Protocol       string `json:"protocol"`
	WSID           string `json:"ws_id"`
	ConnectTimeout int    `json:"connect_timeout"`
	Retries        int    `json:"retries"`
}


