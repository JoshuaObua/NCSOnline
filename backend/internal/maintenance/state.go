package maintenance

import (
	"sync/atomic"
	"time"
)

type Snapshot struct {
	Enabled                 bool       `json:"maintenance_mode"`
	Reason                  string     `json:"reason"`
	ExpectedEnd             *time.Time `json:"expected_end,omitempty"`
	ChangedAt               time.Time  `json:"changed_at"`
	ChangedBy               string     `json:"changed_by,omitempty"`
	GlobalAuthInvalidBefore *time.Time `json:"global_auth_invalid_before,omitempty"`
	CacheGeneration         uint64     `json:"cache_generation"`
}

type State struct{ value atomic.Pointer[Snapshot] }

func New(initial Snapshot) *State { s := &State{}; s.value.Store(&initial); return s }
func (s *State) Get() Snapshot {
	v := s.value.Load()
	if v == nil {
		return Snapshot{}
	}
	return *v
}
func (s *State) Set(next Snapshot) { copy := next; s.value.Store(&copy) }
func (s *State) FlushCache() Snapshot {
	for {
		old := s.value.Load()
		if old == nil {
			initial := Snapshot{}
			s.value.CompareAndSwap(nil, &initial)
			continue
		}
		next := *old
		next.CacheGeneration++
		if s.value.CompareAndSwap(old, &next) {
			return next
		}
	}
}
