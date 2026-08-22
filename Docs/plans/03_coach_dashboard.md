# NAMIS Coach & Technical Staff Dashboard Specification

## 1. Overview
The **Coach & Technical Staff Dashboard** enables head coaches, assistant coaches, strength & conditioning specialists, physiotherapists, team doctors, and sports scientists to manage national squad rosters, log athlete performance metrics, monitor injury/rehabilitation status, submit competition entries, and verify professional technical certifications under recognized national sports federations.

---

## 2. Technical Support Team Roster & Roles (11 Disciplines)

NAMIS structures technical support staff around 11 core roles to ensure high-performance elite sports science and multidisciplinary care:

```
+---------------------------------------------------------------------------------+
|                       TECHNICAL STAFF ROSTER & STRUCTURE                        |
+---------------------------------------------------------------------------------+
| 1. Primary Coach (Head Coach)        | 7. Sports Scientist / Technical Director |
| 2. Coach ID & License Number         | 8. Physiotherapist                       |
| 3. Coaching Level & Certification    | 9. Team Doctor                           |
| 4. Primary Coach Contact             | 10. Sports Nutritionist                  |
| 5. Assistant Coach                   | 11. Team Manager                         |
| 6. Strength & Conditioning (S&C)     |                                          |
+---------------------------------------------------------------------------------+
```

---

## 3. UI Layout & View Structure

```
+-----------------------------------------------------------------------------------+
|  [NCS Logo]  NAMIS COACH & TECHNICAL DASHBOARD    [Squad Selector] [Dark] [User] |
+-----------------------------------------------------------------------------------+
| [Sidebar Navigation Menu]                                                         |
|  Coach Overview                                                                 |
|  Squad Management [v]                                                            |
|    ├── National & Club Squad Roster                                               |
|    ├── Competition Performance Logger                                             |
|    ├── Athlete Personal Best (PB) & Record Triggers                               |
|    └── Training Load & GPS Tracking Data                                          |
|  Multidisciplinary Entourage [v]                                                |
|    ├── Technical Support Staff Roster                                             |
|    ├── Strength & Conditioning (S&C) Assignments                                  |
|    ├── Physiotherapy & Team Doctor Injury Logs                                    |
|    └── Sports Nutrition & Dietary Plans                                           |
|  Coaching Credentials & License [v]                                             |
|    ├── NCS Coach License Status                                                   |
|    ├── License History & Renewals                                                 |
|    ├── Certification Upload Vault (CAF, WA, FIBA)                                 |
|    └── CPD & Coaching Clinics Log                                                 |
+-----------------------------------------------------------------------------------+

```

---

## 4. Key Functional Modules

### 4.1 Squad Roster & Performance Logger
* **Squad Roster View**: List of all registered athletes under the coach's federation or club.
* **Performance Entry Modal**: Log times, distances, weights, scores, and rankings achieved during sanctioned competitions.
* **Personal Best (PB) & National Record (NR) Trigger**: Auto-flag when a logged result beats an athlete's historical record.

### 4.2 Technical Team Staff Assignment
* Ability for Head Coach or Team Manager to attach accredited support personnel to specific national team tours:
  - **S&C Specialist**: Monitor load management and GPS tracking data.
  - **Physiotherapist & Doctor**: Update injury status (`Fit`, `Light Training`, `Restricted`, `Out`) and clearance certificates.
  - **Nutritionist**: Formulate dietary plans and hydration targets.

### 4.3 Certification & License Tracker
* **Certification Upload**: Submit copies of diplomas (*CAF A/B/C*, *World Athletics*, *FIBA Level 1/2/3*, *World Boxing*).
* **NCS Coaching License**: Real-time license verification status issued by the National Council of Sports.

---

## 5. API Endpoints Specification

```http
GET /api/v1/coaches/squad
Headers: Authorization: Bearer <coach_token>
Response: 200 OK
{
  "coach_id": "COACH-UG-9912",
  "full_name": "Coach X",
  "certification": "World Athletics Level 3",
  "athletes": [
    {
      "id": "usr_ath_001",
      "full_name": "Athlete A",
      "discipline": "5000m",
      "status": "Active",
      "personal_best": "12:45.10"
    }
  ]
}

POST /api/v1/coaches/results
Payload:
{
  "athlete_id": "usr_ath_001",
  "competition_id": "comp_009",
  "event": "5000m",
  "result_value": "12:45.10",
  "position": 1,
  "is_personal_best": true
}
```
