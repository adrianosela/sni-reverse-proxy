package hostmgr

import (
	"net/http/httputil"
)

// Host represents a target host of a proxy
type Host struct {
	Address      string // e.g. "host:port" or ":port" for localhost port
	ReverseProxy *httputil.ReverseProxy
}

// HostManager represents the host-management functionality required of a multiple-host proxy
type HostManager interface {
	PutHost(sni string, addr string, rp *httputil.ReverseProxy) error
	GetHost(sni string) (*Host, bool, error)
	RemoveHost(sni string) error
}
