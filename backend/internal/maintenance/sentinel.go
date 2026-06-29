package maintenance

import (
	"context"
	"time"
)

// EnabledFlipper is the narrow contract the schedule sentinel needs from
// the OperatorRepo — it can persist a flip of the maintenance_mode bit
// without disturbing the scheduling window.
type EnabledFlipper interface {
	SetEnabled(ctx context.Context, enabled bool, actor string) (Snapshot, error)
	SetScopeEnabled(ctx context.Context, scope string, enabled bool, actor string) (Snapshot, error)
}

// StartScheduleSentinel watches the in-memory state on a 30-second ticker
// and persists a flip when the schedule says we should be in or out of
// maintenance but the persisted Enabled bit disagrees. Returns immediately;
// the goroutine ends when ctx is cancelled.
//
// Concretely: once `expected_end` passes for an active window the sentinel
// flips Enabled → false so the UI reflects "operational" without an
// operator having to click anything.
func StartScheduleSentinel(ctx context.Context, state *State, flipper EnabledFlipper) {
	go func() {
		t := time.NewTicker(30 * time.Second)
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				reconcile(ctx, state, flipper)
			}
		}
	}()
}

func reconcile(ctx context.Context, state *State, flipper EnabledFlipper) {
	s := state.Get()
	s.Normalize()
	now := time.Now()
	// Only the "active window expired" transition is automatic. We
	// deliberately don't auto-disable when scheduled_start is in the
	// future — IsActiveAt handles that without mutating storage so
	// admins still see their planned window in the UI.
	for _, scope := range []string{ScopePublicCMS, ScopeAdminDashboard} {
		scoped := s.Scoped(scope)
		if scoped.Enabled && scoped.ExpectedEnd != nil && now.After(*scoped.ExpectedEnd) {
			updated, err := flipper.SetScopeEnabled(ctx, scope, false, "schedule-sentinel")
			if err != nil {
				return
			}
			state.Set(updated)
			s = updated
		}
	}
}
