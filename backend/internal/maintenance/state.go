package maintenance

import (
	"encoding/json"
	"sync/atomic"
	"time"
)

const (
	ScopePublicCMS      = "public_cms"
	ScopeAdminDashboard = "admin_dashboard"
)

type DisplayMeta struct {
	CustomTitle            string `json:"custom_title"`
	CustomMessage          string `json:"custom_message"`
	ShowCountdown          bool   `json:"show_countdown"`
	AllowEmailNotification bool   `json:"allow_email_notification"`
}

type BypassRules struct {
	AllowedRoles     []string `json:"allowed_roles"`
	AllowedIPRanges  []string `json:"allowed_ip_ranges"`
	SecretQueryParam string   `json:"secret_query_param"`
}

type ScopedSnapshot struct {
	Scope             string      `json:"scope"`
	Enabled           bool        `json:"maintenance_mode"`
	Reason            string      `json:"reason"`
	ScheduledStart    *time.Time  `json:"scheduled_start,omitempty"`
	ExpectedEnd       *time.Time  `json:"expected_end,omitempty"`
	AutoStartEnforced bool        `json:"auto_start_enforced"`
	AutoEndEnforced   bool        `json:"auto_end_enforced"`
	DisplayMeta       DisplayMeta `json:"display_meta"`
	BypassRules       BypassRules `json:"bypass_rules"`
	ChangedAt         time.Time   `json:"changed_at"`
	ChangedBy         string      `json:"changed_by,omitempty"`
}

type Snapshot struct {
	Enabled                 bool           `json:"maintenance_mode"`
	Reason                  string         `json:"reason"`
	ScheduledStart          *time.Time     `json:"scheduled_start,omitempty"`
	ExpectedEnd             *time.Time     `json:"expected_end,omitempty"`
	ChangedAt               time.Time      `json:"changed_at"`
	ChangedBy               string         `json:"changed_by,omitempty"`
	PublicCMS               ScopedSnapshot `json:"public_cms"`
	AdminDashboard          ScopedSnapshot `json:"admin_dashboard"`
	GlobalAuthInvalidBefore *time.Time     `json:"global_auth_invalid_before,omitempty"`
	CacheGeneration         uint64         `json:"cache_generation"`
}

func DefaultScopedSnapshot(scope string) ScopedSnapshot {
	now := time.Now().UTC()
	s := ScopedSnapshot{
		Scope:             scope,
		AutoStartEnforced: true,
		AutoEndEnforced:   true,
		ChangedAt:         now,
		DisplayMeta: DisplayMeta{
			CustomTitle:            "Scheduled Maintenance",
			CustomMessage:          "The platform is temporarily unavailable while scheduled maintenance is in progress.",
			ShowCountdown:          true,
			AllowEmailNotification: false,
		},
	}
	switch scope {
	case ScopeAdminDashboard:
		s.DisplayMeta.CustomTitle = "Admin Dashboard Maintenance"
		s.DisplayMeta.CustomMessage = "Back-office tools are temporarily unavailable while maintenance is in progress."
		s.BypassRules.AllowedRoles = []string{"super_admin", "admin", "content_manager"}
		s.BypassRules.SecretQueryParam = "admin_maintenance_bypass"
	default:
		s.Scope = ScopePublicCMS
		s.DisplayMeta.CustomTitle = "We'll be right back"
		s.DisplayMeta.CustomMessage = "The National Council of Sports platform is undergoing scheduled maintenance."
		s.BypassRules.AllowedRoles = []string{"super_admin", "admin", "content_manager"}
		s.BypassRules.SecretQueryParam = "public_maintenance_bypass"
	}
	return s
}

func (s ScopedSnapshot) Normalize(scope string) ScopedSnapshot {
	if s.Scope == "" {
		s.Scope = scope
	}
	if s.ChangedAt.IsZero() {
		s.ChangedAt = time.Now().UTC()
	}
	if !s.AutoStartEnforced {
		s.AutoStartEnforced = true
	}
	if !s.AutoEndEnforced {
		s.AutoEndEnforced = true
	}
	defaults := DefaultScopedSnapshot(scope)
	if s.DisplayMeta.CustomTitle == "" {
		s.DisplayMeta.CustomTitle = defaults.DisplayMeta.CustomTitle
	}
	if s.DisplayMeta.CustomMessage == "" {
		s.DisplayMeta.CustomMessage = defaults.DisplayMeta.CustomMessage
	}
	if len(s.BypassRules.AllowedRoles) == 0 {
		s.BypassRules.AllowedRoles = defaults.BypassRules.AllowedRoles
	} else if scope == ScopeAdminDashboard {
		s.BypassRules.AllowedRoles = ensureRoles(s.BypassRules.AllowedRoles, "super_admin", "admin", "content_manager")
	}
	if s.BypassRules.SecretQueryParam == "" {
		s.BypassRules.SecretQueryParam = defaults.BypassRules.SecretQueryParam
	}
	return s
}

func ensureRoles(current []string, required ...string) []string {
	seen := map[string]bool{}
	for _, role := range current {
		seen[role] = true
	}
	for _, role := range required {
		if !seen[role] {
			current = append(current, role)
		}
	}
	return current
}

func (s ScopedSnapshot) IsActiveAt(now time.Time) bool {
	if s.ExpectedEnd != nil && now.After(*s.ExpectedEnd) {
		return false
	}
	if s.ScheduledStart != nil && now.Before(*s.ScheduledStart) {
		return false
	}
	return s.Enabled
}

func (s ScopedSnapshot) MarshalJSON() ([]byte, error) {
	type alias ScopedSnapshot
	now := time.Now()
	return json.Marshal(&struct {
		alias
		IsActive    bool `json:"is_active"`
		IsScheduled bool `json:"is_scheduled"`
		Scheduling  struct {
			StartTime         *time.Time `json:"start_time,omitempty"`
			EndTime           *time.Time `json:"end_time,omitempty"`
			AutoStartEnforced bool       `json:"auto_start_enforced"`
			AutoEndEnforced   bool       `json:"auto_end_enforced"`
		} `json:"scheduling"`
	}{
		alias:       alias(s),
		IsActive:    s.IsActiveAt(now),
		IsScheduled: s.Enabled && s.ScheduledStart != nil && s.ScheduledStart.After(now),
		Scheduling: struct {
			StartTime         *time.Time `json:"start_time,omitempty"`
			EndTime           *time.Time `json:"end_time,omitempty"`
			AutoStartEnforced bool       `json:"auto_start_enforced"`
			AutoEndEnforced   bool       `json:"auto_end_enforced"`
		}{s.ScheduledStart, s.ExpectedEnd, s.AutoStartEnforced, s.AutoEndEnforced},
	})
}

func (s *Snapshot) Normalize() {
	if s.PublicCMS.Scope == "" && s.AdminDashboard.Scope == "" {
		s.PublicCMS = DefaultScopedSnapshot(ScopePublicCMS)
		s.PublicCMS.Enabled = s.Enabled
		s.PublicCMS.Reason = s.Reason
		s.PublicCMS.ScheduledStart = s.ScheduledStart
		s.PublicCMS.ExpectedEnd = s.ExpectedEnd
		s.PublicCMS.ChangedAt = s.ChangedAt
		s.PublicCMS.ChangedBy = s.ChangedBy
	}
	s.PublicCMS = s.PublicCMS.Normalize(ScopePublicCMS)
	s.AdminDashboard = s.AdminDashboard.Normalize(ScopeAdminDashboard)
	s.Enabled = s.PublicCMS.Enabled
	s.Reason = s.PublicCMS.Reason
	s.ScheduledStart = s.PublicCMS.ScheduledStart
	s.ExpectedEnd = s.PublicCMS.ExpectedEnd
	s.ChangedAt = s.PublicCMS.ChangedAt
	s.ChangedBy = s.PublicCMS.ChangedBy
}

// MarshalJSON adds two derived fields the frontend uses to decide what to
// render: `is_active` is true only when the window is currently live, and
// `is_scheduled` is true when a future window has been planned.
func (s Snapshot) MarshalJSON() ([]byte, error) {
	type alias Snapshot
	now := time.Now()
	s.Normalize()
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
	s.Normalize()
	return s.PublicCMS.IsActiveAt(now)
}

func (s Snapshot) Scoped(scope string) ScopedSnapshot {
	s.Normalize()
	if scope == ScopeAdminDashboard {
		return s.AdminDashboard
	}
	return s.PublicCMS
}

func (s Snapshot) WithScoped(scope string, scoped ScopedSnapshot) Snapshot {
	s.Normalize()
	scoped = scoped.Normalize(scope)
	if scope == ScopeAdminDashboard {
		s.AdminDashboard = scoped
	} else {
		s.PublicCMS = scoped
	}
	s.Normalize()
	return s
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
