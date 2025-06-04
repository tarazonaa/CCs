package models

import (
	"strconv"
	"time"
)

type KongLogEntry struct {
	RequestID         string    `json:"request_id"`
	ClientIP          string    `json:"client_ip"`
	Method            string    `json:"method"`
	URI               string    `json:"uri"`
	URL               string    `json:"url"`
	UserAgent         string    `json:"user_agent"`
	Host              string    `json:"host"`
	Status            int       `json:"status"`
	ResponseSize      int       `json:"response_size"`
	RequestSize       int       `json:"request_size"`
	UpstreamLatencyMs int       `json:"upstream_latency_ms"`
	KongLatencyMs     int       `json:"kong_latency_ms"`
	ProxyLatencyMs    int       `json:"proxy_latency_ms"`
	TotalLatencyMs    int       `json:"total_latency_ms"`
	ServiceID         string    `json:"service_id"`
	ServiceName       string    `json:"service_name"`
	RouteID           string    `json:"route_id"`
	RouteName         string    `json:"route_name"`
	StartedAt         time.Time `json:"started_at"`
}

func parseHeaderInt(value string) int {
	if value == "" {
		return 0
	}
	n, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	return n
}

func MapKongLogToEntry(raw *KongLog) KongLogEntry {
	return KongLogEntry{
		RequestID:         raw.Request.ID,
		ClientIP:          raw.ClientIP,
		Method:            raw.Request.Method,
		URI:               raw.Request.URI,
		URL:               raw.Request.URL,
		UserAgent:         raw.Request.Headers["user-agent"],
		Host:              raw.Request.Headers["host"],
		Status:            raw.Response.Status,
		ResponseSize:      raw.Response.Size,
		RequestSize:       raw.Request.Size,
		UpstreamLatencyMs: parseHeaderInt(raw.Response.Headers["x-kong-upstream-latency"]),
		KongLatencyMs:     raw.Latencies.Kong,
		ProxyLatencyMs:    raw.Latencies.Proxy,
		TotalLatencyMs:    raw.Latencies.Request,
		ServiceID:         raw.Service.ID,
		ServiceName:       raw.Service.Name,
		RouteID:           raw.Route.ID,
		RouteName:         raw.Route.Name,
		StartedAt:         time.UnixMilli(raw.StartedAt),
	}
}
