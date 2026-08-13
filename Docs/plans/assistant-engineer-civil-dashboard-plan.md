# Assistant Engineer Civil Dashboard Blueprint & Implementation Plan

> **Target File:** `/home/fidi/Projects/NCS_Intranet/Docs/plans/assistant-engineer-civil-dashboard-plan.md`  
> **Role:** Assistant Engineer Civil (Operational Field Supervisor)  
> **Reports To:** Engineering Officer Civil  
> **Oversees:** Plumbers, Masons, Carpenters, Grounds Technicians  

---

## 1. Role Overview & Objectives

The **Assistant Engineer Civil** serves as the primary operational supervisor for civil maintenance, plumbing systems, building structures, and sports ground maintenance across NCS facilities.

---

## 2. Shared Navigation & Sidebar Menu Items

Inherits the unified **NCS Intranet Navigation Framework**:

```
+-----------------------------------------------------------------------------------+
| [≡] NCS INTRANET | Asst Engineer Civil       [🔍 Search...] [🔄] [🌙] [🔔 5] [💬 2] [User] |
+-----------------------------------------------------------------------------------+
| [SIDEBAR MENU]  | [MAIN CONTENT DASHBOARD]                                        |
| 1. Dashboard    | +-------------------------------------------------------------+ |
| 2. Work Orders  | | STATS: Unassigned: 5 | Active Plumbers: 4 | Verifications: 7  | |
| 3. Plumbers     | +-------------------------------------------------------------+ |
| 4. Field Reports| | [Overview] [Create WO] [My Activities] [Leave] [Profile]     | |
| 5. My Activities| +-------------------------------------------------------------+ |
| 6. Leave        | | Active Civil & Plumbing Work Orders (12)                        | |
| 7. Messages     | | WO-0880 | Arena Washroom Leak | Plumber Kato | [Assign]         | |
| 8. Profile      | +-------------------------------------------------------------+ |
| 9. Settings     |                                                                 |
+-----------------------------------------------------------------------------------+
```

### Functional Features
- **My Activities Menu Item:** Routes to `/intranet/my-activities`. Displays personal activity logs (work orders created, plumber task assignments, material requisitions drafted, site verifications cleared).
- **Navbar Controls:** Theme switcher, global search, refresh, notifications drawer, and messages drawer.

---

## 3. "My Activities" Personal Audit Log for Assistant Engineer Civil

```json
// GET /api/v1/intranet/my-activities?category=WORK_ORDER
{
  "status": "success",
  "data": [
    {
      "id": "act-120",
      "action_type": "WORK_ORDER_ASSIGNED",
      "description": "Assigned Plumbing Work Order WO-0880 to Plumber John Kato.",
      "target_resource_id": "WO-0880",
      "created_at": "2026-07-29T08:15:00Z"
    }
  ]
}
```

---

## 4. Go REST API & Checklist

| Endpoint | Method | Scope | Description |
| :--- | :--- | :--- | :--- |
| `/api/v1/engineering/work-orders` | POST, GET | Asst Eng Civil | Work order management |
| `/api/v1/intranet/my-activities` | GET | Self Only | Query personal account audit log |

- [ ] **Step 1:** Add **"My Activities"** link to `AssistantEngineerCivilDashboard.vue` sidebar.
- [ ] **Step 2:** Test personal activity audit query for Assistant Engineer Civil account.
