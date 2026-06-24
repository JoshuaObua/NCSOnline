package maintenance

import (
	"encoding/json"
	"sync/atomic"
	"time"
)

type Snapshot struct {
	Enabled                 bool       `json:"maintenance_mode"`
	Reason                  string     `json:"reason"`
	ScheduledStart          *time.Time `json:"scheduled_start,omitempty"`
	ExpectedEnd             *time.Time `json:"expected_end,omitempty"`
	ChangedAt               time.Time  `json:"changed_at"`
	ChangedBy               string     `json:"changed_by,omitempty"`
	GlobalAuthInvalidBefore *time.Time `json:"global_auth_invalid_before,omitempty"`
	CacheGeneration         uint64     `json:"cache_generation"`
}

// MarshalJSON adds two derived fields the frontend uses to decide what to
// render: `is_active` is true only when the window is currently live, and
// `is_scheduled` is true when a future window has been planned.
func (s Snapshot) MarshalJSON() ([]byte, error) {
	type alias Snapshot
	now := time.Now()
	return json.Marshal(&struct {
		alias
		IsActive    bool `json:"is_active"`
		IsScheduled bool `json:"is_scheduled"`
	}{
		alias:       alias(s),
		IsActive:    s.IsActiveAt(now),
		IsScheduled: s.Enabled && s.ScheduledStart != nil && s.ScheduledStart.After(now),
	})
}

// IsActiveAt reports whether maintenance mode should be considered ON at the
// supplied instant, taking the schedule window into account. A start in the
// past with no end (or an end in the future) means active; an end in the past
// always means inactive even if the persisted Enabled bit is still true.
func (s Snapshot) IsActiveAt(now time.Time) bool {
	if s.ExpectedEnd != nil && now.After(*s.ExpectedEnd) {
		return false
	}
	if s.ScheduledStart != nil && now.Before(*s.ScheduledStart) {
		return false
	}
	return s.Enabled
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
