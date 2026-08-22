# NAMIS NCS Management & Executive Analytics Dashboard Specification

## 1. Overview
The **NCS Management & Executive Analytics Dashboard** provides high-level executive decision support for the Minister of Sports, NCS Board, General Secretary, and Technical Department. It aggregates national sports statistics, measures key performance indicators (KPIs), tracks government funding return-on-investment, monitors international medal yields, and visualizes grass-roots sports development across all sub-regions of Uganda.

---

## 2. Ten Key NCS KPI Reports

NAMIS provides real-time automated reports across 10 critical governance indicators:

```
+---------------------------------------------------------------------------------+
|                          10 MANDATORY NCS KPI REPORTS                           |
+---------------------------------------------------------------------------------+
| 1. Total Registered Athletes in Uganda   | 6. Athletes Identified via Talent Prog|
| 2. Total Active Athletes Roster          | 7. Athletes Supported by Government  |
| 3. Total National Team Athletes          | 8. School -> Club -> National Team   |
| 4. Total International Medals Won        | 9. Gender Participation Statistics   |
| 5. Medal Table by Federation             | 10. Regional Sports Development Stats |
+---------------------------------------------------------------------------------+
```

---

## 3. UI Layout & Visual Dashboard Mockup

```
+-----------------------------------------------------------------------------------+
|  [NCS Logo]  NCS EXECUTIVE ANALYTICS COMMAND CENTER  [Year: 2026] [Export PDF]    |
+-----------------------------------------------------------------------------------+
| [Sidebar Navigation Menu]                                                         |
|  Executive Command Center Overview                                              |
|  NCS Council KPI Analytics [v]                                                   |
|    ├── 10 Mandatory Governance KPI Command Board                                  |
|    ├── Athlete Demographics & Regional Heatmap                                    |
|    └── Gender Equity & Inclusivity (Para-Sports)                                  |
|  Federation Performance Audits [v]                                               |
|    ├── Category A/B/C Federation Standings                                        |
|    ├── National Medal Table Rankings                                              |
|    └── Statutory Grant Allocation ROI Reports                                     |
|  Talent Development Pipeline [v]                                                |
|    ├── School -> Club -> National Team Progression Funnel                         |
|    ├── Regional Talent Scouting Centers Report                                    |
|    └── Government Sports Scholarship Audit                                        |
|  Dynamic Form Builder & Management [v]                                           |
|    ├── Form Schemas & Form Builder                                                |
|    ├── Form Publication & Categories                                              |
|    └── Submissions & Approval Review Desk                                         |
+-----------------------------------------------------------------------------------+

```

---

## 4. Analytical Breakdowns & Reports

### 4.1 Athletes Analytics Breakdown
* **Athletes by Federation**: Bar charts comparing athlete volume across Priority vs Established federations.
* **Athletes by District & Region**: Geographical distribution map of athlete origins across Northern, Eastern, Central, and Western Uganda.
* **Gender Participation Ratio**: Target tracking for 50/50 gender equality initiatives.
* **Age Group Demographics**: Pyramid distribution (*U10, U12, U15, U17, U20, Senior*).

---

### 4.2 Elite Athlete & Government Support Tracking
* **Government Support ROI**: Tracking performance outcomes of athletes receiving government stipends or equipment grants.
* **Athletes on Employment**: Tracking employment opportunities provided to elite athletes by government agencies (Police, Prisons, UPDF) and private sponsors.
* **Para-Athletes & Inclusivity**: Tracking athlete participation and medal conversion rates for athletes with disabilities.

---

### 4.3 Super Admin Dynamic Form Builder & Form Management Engine
* **Form Builder Interface**: Drag-and-drop / JSON dynamic form schema creator for Super Admins.
* **Supported Field Types**:
  - `TEXT` / `TEXTAREA`: Single-line & multi-line text input fields.
  - `NUMBER`: Integer / decimal financial & count inputs.
  - `DATE`: Single date or date-range pickers.
  - `DROPDOWN` / `RADIO`: Single-select options from configured lists.
  - `MULTI_SELECT` / `CHECKBOX`: Multi-select choices & compliance agreement checks.
  - `FILE_UPLOAD`: Document upload dropzone with MIME type validation (`.pdf`, `.png`, `.jpeg`) and maximum size constraints (up to 25MB).
* **Form Categories**: `LICENSE_RENEWAL`, `CLUB_ACCREDITATION`, `INTERNATIONAL_TOUR_SANCTION`, `TALENT_GRANT`, `GENERAL_COMPLIANCE`.
* **Submissions Review & Approval Desk**:
  - Inspect raw dynamic JSON response payloads alongside rendered attachments.
  - One-click actions: **Approve Application**, **Reject**, or **Request Revision** with mandatory feedback comments.
  - Auto-trigger email notifications and status update Webhooks to submitting federations.

---


### 4.2 Elite Athlete & Government Support Tracking
* **Government Support ROI**: Tracking performance outcomes of athletes receiving government stipends or equipment grants.
* **Athletes on Employment**: Tracking employment opportunities provided to elite athletes by government agencies (Police, Prisons, UPDF) and private sponsors.
* **Para-Athletes & Inclusivity**: Tracking athlete participation and medal conversion rates for athletes with disabilities.

---

## 5. SQL Analytics Queries Specification

```sql
-- Query: Medal Table by Federation
SELECT 
    f.name AS federation_name,
    COUNT(CASE WHEN m.medal_type = 'Gold' THEN 1 END) AS gold_count,
    COUNT(CASE WHEN m.medal_type = 'Silver' THEN 1 END) AS silver_count,
    COUNT(CASE WHEN m.medal_type = 'Bronze' THEN 1 END) AS bronze_count,
    COUNT(m.id) AS total_medals
FROM federations f
LEFT JOIN medals m ON f.id = m.federation_id
GROUP BY f.id, f.name
ORDER BY gold_count DESC, silver_count DESC, bronze_count DESC;

-- Query: Talent Progression Pipeline
SELECT 
    a.region,
    COUNT(DISTINCT ti.id) AS identified_talent,
    COUNT(DISTINCT c.id) AS club_athletes,
    COUNT(DISTINCT nta.id) AS national_team_caps
FROM athletes a
LEFT JOIN talent_identification ti ON a.id = ti.athlete_id
LEFT JOIN clubs c ON a.club_id = c.id
LEFT JOIN national_team_appearances nta ON a.id = nta.athlete_id
GROUP BY a.region;
```
