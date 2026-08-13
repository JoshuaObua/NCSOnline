# Assistant Engineer Electrical Dashboard Blueprint & Implementation Plan

> **Target File:** `/home/fidi/Projects/NCS_Intranet/Docs/plans/assistant-engineer-electrical-dashboard-plan.md`  
> **Role:** Assistant Engineer Electrical (Electrical Operations Supervisor)  
> **Reports To:** Engineering Officer Electrical  
> **Oversees:** Electricians, Generator Operators, AC/Sound Technicians  

---

## 1. Role Overview & Objectives

The **Assistant Engineer Electrical** manages field electrical execution, routine electrical safety checks, generator servicing logs, and floodlight operational maintenance across all NCS venues.

---

## 2. Shared Navigation & Sidebar Menu Items

Inherits the unified **NCS Intranet Navigation Framework**:

```
+-----------------------------------------------------------------------------------+
| [≡] NCS INTRANET | Asst Eng Electrical      [🔍 Search...] [🔄] [🌙] [🔔 3] [💬 2] [User] |
+-----------------------------------------------------------------------------------+
| [SIDEBAR MENU]  | [MAIN CONTENT DASHBOARD]                                        |
| 1. Dashboard    | +-------------------------------------------------------------+ |
| 2. Work Orders  | | STATS: Unassigned: 3 | Electricians: 3 | Gen 1 Fuel: 94%       | |
| 3. Generators   | +-------------------------------------------------------------+ |
| 4. Field Reports| | [Overview] [Create WO] [My Activities] [Leave] [Profile]     | |
| 5. My Activities| +-------------------------------------------------------------+ |
| 6. Leave        | | Active Electrical Maintenance Tasks (8)                      | |
| 7. Messages     | | WO-0881 | Arena Main Switchboard | Alex Musoke | [Assign]     | |
| 8. Profile      | +-------------------------------------------------------------+ |
| 9. Settings     |                                                                 |
+-----------------------------------------------------------------------------------+
```

### Functional Features
- **My Activities Menu Item:** Routes to `/intranet/my-activities`. Displays personal activity logs (generator logs recorded, electrical work orders assigned, pre-match floodlight checklists submitted).
- **Navbar Controls:** Theme switcher, global search, refresh, notifications drawer, and messages drawer.

---

## 3. "My Activities" Personal Audit Log for Assistant Engineer Electrical

```json
// GET /api/v1/intranet/my-activities?category=WORK_ORDER
{
  "status": "success",
  "data": [
    {
      "id": "act-128",
      "action_type": "GENERATOR_LOG_RECORDED",
      "description": "Recorded Generator Run Hours (1,420 hrs) & 200L Diesel Refill for Lugogo 500kVA Generator.",
      "target_resource_id": "NCS-ELE-GEN-002",
      "created_at": "2026-07-29T09:30:00Z"
    }
  ]
}
```

---

## 4. Go REST API & Checklist

| Endpoint | Method | Scope | Description |
| :--- | :--- | :--- | :--- |
| `/api/v1/engineering/generators/log` | POST, GET | Asst Eng Elec | Generator logging |
| `/api/v1/intranet/my-activities` | GET | Self Only | Query personal account audit log |

- [ ] **Step 1:** Add **"My Activities"** link to `AssistantEngineerElectricalDashboard.vue` sidebar.
- [ ] **Step 2:** Test personal activity audit query for Assistant Engineer Electrical account.
