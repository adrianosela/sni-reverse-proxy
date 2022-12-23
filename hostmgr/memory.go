package hostmgr

import (
	"net/http/httputil"
	"sync"
)

// InMemoryHostManager is an in-memory implementation of the HostManager iface
type InMemoryHostManager struct {
	sync.RWMutex // inherit read/write lock behavior
	hosts        map[string]*Host
}

// ensure InMemoryTargetStorage implements HostManager at compile-time
var _ HostManager = (*InMemoryHostManager)(nil)

// NewInMemoryHostManager is the InMemoryHostManager constructor
func NewInMemoryHostManager() *InMemoryHostManager {
	return &InMemoryHostManager{
		hosts: make(map[string]*Host),
	}
}

// PutHost adds/updates the target host for a given server name
func (s *InMemoryHostManager) PutHost(
	sni string,
	addr string,
	rp *httputil.ReverseProxy,
) error {
	s.Lock()
	defer s.Unlock()

	s.hosts[sni] = &Host{
		Address:      addr,
		ReverseProxy: rp,
	}
	return nil
}

// GetHost looks-up the target host for a given server name
func (s *InMemoryHostManager) GetHost(sni string) (*Host, bool, error) {
	s.RLock()
	defer s.RUnlock()

	t, ok := s.hosts[sni]
	return t, ok, nil
}

// RemoveHost removes the target host for a given server name
func (s *InMemoryHostManager) RemoveHost(sni string) error {
	s.Lock()
	defer s.Unlock()

	delete(s.hosts, sni)
	return nil
}
