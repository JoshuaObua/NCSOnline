# NAMIS Competitions & Medals Dashboard Specification

## 1. Overview
The **Competitions & Medals Dashboard** is the official results ledger and honors registry for the National Council of Sports Uganda. It tracks national and international competition results, national team appearances, podium medal achievements, prize money allocations, and personal/national record badges.

---

## 2. Eight-Tier Competition Hierarchy

NAMIS categorizes all official sporting competitions into an 8-tier hierarchy:

```
+---------------------------------------------------------------------------------+
|                         COMPETITION LEVEL HIERARCHY                             |
+---------------------------------------------------------------------------------+
| Level I   : District Championships (e.g., Kampala District Games)               |
| Level II  : Regional Games (e.g., Northern Uganda Regional Athletics)           |
| Level III : National Championships (e.g., Uganda National Track & Field)        |
| Level IV  : East African Games (e.g., FEASSSA Games, East Africa University)    |
| Level V   : African Championships / Games (e.g., All Africa Games)            |
| Level VI  : Commonwealth Games                                                  |
| Level VII : Olympic Games / Paralympic Games                                    |
| Level VIII: World Championships (e.g., World Athletics, FIFA World Cup)         |
+---------------------------------------------------------------------------------+
```

---

## 3. UI Layout & View Structure

```
+-----------------------------------------------------------------------------------+
|  [NCS Logo]  COMPETITIONS & MEDALS DESK        [Event Filter] [Search] [Export]   |
+-----------------------------------------------------------------------------------+
| [Sidebar Navigation Menu]                                                         |
| Competitions Desk Overview                                                     |
|  Competitions Registry [v]                                                      |
|    ├── All Sanctioned Events (8 Tiers)                                            |
|    ├── National Championships (Level III)                                         |
|    ├── International Games (Level IV - VIII)                                      |
|    └── Register Sanctioned Competition                                            |
|  Results & Scoring Engine [v]                                                    |
|    ├── Live Official Scoring Log                                                  |
|    ├── Result Verification Desk                                                   |
|    └── National Record (NR) Flags                                                 |
|  Medals & Honors Desk [v]                                                       |
|    ├── Medals Registry (Gold, Silver, Bronze)                                     |
|    ├── NCS Recognition Approvals                                                  |
|    └── Prize Money Disbursements                                                  |
|  National Squads & Caps [v]                                                      |
|    ├── National Team Call-Up Registries                                           |
|    ├── International Caps & Appearances Ledger                                    |
|    └── Debut Milestones & Honors                                                  |
+-----------------------------------------------------------------------------------+

```

---

## 4. Key Core Modules

### 4.1 Performance & Results Engine
* **Performance Value Formatting**: Supports time (*seconds/minutes*), distance (*meters*), weight (*kg*), score (*points/goals*), and ranking.
* **Automated Record Badges**:
  - **National Record (NR)**: Triggered when a result surpasses the current Ugandan national record for the event.
  - **Personal Best (PB)**: Triggered when an athlete beats their historical best result.
  - **Seasonal Best (SB)**: Fastest/highest result recorded in the current calendar year.

---

### 4.2 Medals and Awards Module
```
Medal Record Structure Example:
+-----------------------------------------------------------------------------------+
| Athlete   | Medal | Event | Competition           | Date        | Coach   | Prize |
+-----------------------------------------------------------------------------------+
| Athlete A | Gold  | 5000m | African Championships | 12 Jun 2026 | Coach X | \$5,000|
+-----------------------------------------------------------------------------------+
```
* **Medal Attributes**: Medal ID, Athlete ID, Medal Type (`Gold`, `Silver`, `Bronze`), Event Won, Competition Level, Date Won, Host Country, Represented Federation, Coach at Time of Winning, Team or Individual Event, Prize Money Awarded (UGX/USD), NCS Recognition Status.

---

### 4.3 National Team Representation Tracker
* **Team Roster Call-Up History**: Tracks national selection call-ups (*Uganda Cranes*, *She Cranes*, *Uganda Olympic Team*, *Rugby Cranes*).
* **Cap Counters**: Total number of official international caps/appearances.
* **Debut & Milestone Tracking**: Records date of first call-up and latest appearance.

---

## 5. Sample API Output (Medal Ledger)

```json
{
  "competition": "African Championships 2026",
  "level": "African",
  "host_city": "Nairobi",
  "host_country": "Kenya",
  "medals_summary": {
    "gold": 4,
    "silver": 2,
    "bronze": 3,
    "total": 9
  },
  "medalists": [
    {
      "medal_id": "MED-2026-001",
      "athlete_name": "Athlete A",
      "medal_type": "Gold",
      "event": "5000m",
      "performance": "12:45.10",
      "is_national_record": true,
      "coach_at_time": "Coach X",
      "prize_money_ugx": 20000000,
      "ncs_recognition": "Official Government Award"
    }
  ]
}
```
