# Engineering Officer Civil Dashboard Blueprint & Implementation Plan

> **Target File:** `/home/fidi/Projects/NCS_Intranet/Docs/plans/engineering-officer-civil-dashboard-plan.md`  
> **Role:** Engineering Officer Civil (Division Head - Civil Engineering)  
> **Reports To:** Senior Engineer  
> **Oversees:** Assistant Engineer Civil (and indirectly Plumbers, Masons, Groundsmen)  

---

## 1. Role Overview & Objectives

The **Engineering Officer Civil** leads the Civil Engineering Division at the National Council of Sports (NCS). This role is responsible for the physical infrastructure integrity of all NCS sports complexes, pitches, stadium seating, roofing structures, drainage networks, and civil asset lifecycles.

---

## 2. Shared Navigation & Sidebar Menu Items

Inherits the unified **NCS Intranet Navigation Framework**:

```
+-----------------------------------------------------------------------------------+
| [≡] NCS INTRANET | Eng Officer Civil        [🔍 Search...] [🔄] [🌙] [🔔 4] [💬 1] [User] |
+-----------------------------------------------------------------------------------+
| [SIDEBAR MENU]  | [MAIN CONTENT DASHBOARD]                                        |
| 1. Dashboard    | +-------------------------------------------------------------+ |
| 2. Civil WOs    | | STATS: Active Civil WOs: 24 | Reviews: 6 | Pitch Score: 92% | |
| 3. Pitch Health | +-------------------------------------------------------------+ |
| 4. Civil Reports| | [Overview] [Requisitions] [My Activities] [Leave] [Profile] | |
| 5. My Activities| +-------------------------------------------------------------+ |
| 6. Leave        | | Civil Material Requisitions Pending Approval (3)             | |
| 7. Messages     | | REQ-042 | 50mm PVC Pipes | 1.8M | [Approve]                 | |
| 8. Profile      | +-------------------------------------------------------------+ |
| 9. Settings     |                                                                 |
+-----------------------------------------------------------------------------------+
```

### Functional Features
- **My Activities Menu Item:** Directly routes to `/intranet/my-activities`. Queries personal action audit history (civil requisitions authorized, pitch inspections logged, leave requests reviewed, password changes).
- **Navbar Controls:** Theme switcher, global search, instant refresh, notifications drawer, and direct messages drawer.

---

## 3. "My Activities" Personal Audit Log for Civil Officer

```json
// GET /api/v1/intranet/my-activities?category=CIVIL
{
  "status": "success",
  "data": [
    {
      "id": "act-104",
      "action_type": "CIVIL_REQUISITION_AUTHORIZED",
      "description": "Authorized Civil Requisition REQ-042 (UGX 1,800,000) for 50mm PVC Pipes.",
      "created_at": "2026-07-29T11:20:00Z"
    }
  ]
}
```

---

## 4. Go REST API & Checklist

| Endpoint | Method | Scope | Description |
| :--- | :--- | :--- | :--- |
| `/api/v1/engineering/civil-officer/dashboard` | GET | Civil Officer | Civil dashboard KPIs & queue |
| `/api/v1/intranet/my-activities` | GET | Self Only | Query personal account audit log |

- [ ] **Step 1:** Register **"My Activities"** in `EngineeringOfficerCivilDashboard.vue` sidebar.
- [ ] **Step 2:** Test personal activity query for Civil Officer account.
