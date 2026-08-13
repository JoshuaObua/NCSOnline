# Seamless System Integration Blueprint & Implementation Plan

> **Target File:** `/home/fidi/Projects/NCS_Intranet/Docs/plans/seamless-dashboard-integration-plan.md`  
> **Goal:** Seamlessly integrate all 6 Engineering Department dashboards, shared modules (Notifications, Messages, Leave, Profile, Settings, My Activities Personal Audit Log, **PPDA Procurement Form 5**), functional responsive navigation (Theme Toggle, Search, Refresh, Notifications, Messages), and **Cascading Hierarchical Reporting** into `NCS_Intranet` with zero friction.

---

## 1. Shared Vue 3 Route Declarations (`frontend/src/router/index.js`)

```javascript
const sharedRoutes = [
  {
    path: '/intranet',
    component: () => import('@/layouts/EngineeringLayout.vue'),
    meta: { requiresAuth: true },
    children: [
      {
        path: 'procurement/form-5',
        name: 'ProcurementForm5',
        component: () => import('@/views/procurement/ProcurementForm5View.vue')
      },
      {
        path: 'my-activities',
        name: 'MyActivities',
        component: () => import('@/views/shared/MyActivitiesView.vue')
      },
      {
        path: 'profile',
        name: 'UserProfile',
        component: () => import('@/views/shared/UserProfileView.vue')
      },
      {
        path: 'settings',
        name: 'UserSettings',
        component: () => import('@/views/shared/UserSettingsView.vue')
      },
      {
        path: 'leave',
        name: 'LeaveManagement',
        component: () => import('@/views/shared/LeaveManagementView.vue')
      },
      {
        path: 'messages',
        name: 'MessageCenter',
        component: () => import('@/views/shared/MessageCenterView.vue')
      },
      {
        path: 'reports',
        name: 'HierarchicalReports',
        component: () => import('@/views/shared/HierarchicalReportsView.vue')
      }
    ]
  }
]
```

---

## 2. Complete Step-by-Step Implementation Roadmap

- [ ] **Task 1: Database Migration Execution**
  Run SQL migrations creating `procurement_form5` and `procurement_form5_items` tables alongside `intranet_user_activities`.

- [ ] **Task 2: Go Backend Procurement Form 5 Handlers**
  Implement Form 5 CRUD, HOD confirmation, Accounting Officer approval, and PDF generation in `backend/internal/procurement/form5_handler.go`.

- [ ] **Task 3: Sidebar Menu Link Integration**
  Add **"Procurement Form 5"** link (`/intranet/procurement/form-5`) to `AppSidebar.vue` and `AppBottomNav.vue` across all department portals.

- [ ] **Task 4: Vue 3 Form 5 Interface (`ProcurementForm5View.vue`)**
  Build 3-step wizard with auto-generating reference numbers, itemized rows, auto-cost calculation, and PDF preview.

- [ ] **Task 5: End-to-End Test**
  Verify complete procurement lifecycle: User creates Form 5 -> HOD confirms -> Accounting Officer approves -> PDF exported.
