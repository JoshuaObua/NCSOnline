# NAMIS Athlete Dashboard Specification

## 1. Overview
The **Athlete Self-Service Dashboard** is an intuitive, mobile-responsive portal designed for national athletes across Uganda. It empowers athletes to inspect their master profiles, track official competition results, view verified medals and honors, manage dual-career educational data, review government equipment support, and access verified digital identity credentials.

---

## 2. UI Layout & View Structure

```
+-----------------------------------------------------------------------------------+
|  [NCS Logo]  NAMIS ATHLETE DASHBOARD        [Notifications] [Dark Mode] [Profile] |
+-----------------------------------------------------------------------------------+
| [Sidebar Navigation Menu]                                                         |
|  Athlete Overview                                                               |
|  My Profile [v]                                                                  |
|    ├── Personal Info & Identity (NIN/Passport)                                    |
|    ├── Classification & Event Specialization                                      |
|    ├── Parents & Emergency Contacts                                               |
|    └── Digital Athlete Credentials / Pass                                         |
|  Performance & Results [v]                                                      |
|    ├── Sanctioned Competition Results Log                                         |
|    ├── Personal Bests (PBs) & Records                                             |
|    └── National Team Appearances & Caps                                           |
|  Medals & Honors Vault [v]                                                      |
|    ├── Podium Medals (Gold/Silver/Bronze)                                         |
|    ├── NCS Recognition & Prize Money Log                                         |
|    └── Trophies & National Certificates                                           |
|  Dual Career & Education [v]                                                     |
|    ├── Academic Institutions & Qualifications                                     |
|    ├── Sports Scholarships & Grants                                               |
|    └── Employment & Career Pathways                                               |
|  Government & Equipment Support [v]                                              |
|    ├── Issued Equipment & Kit Logs                                                |
|    └── Valuation & Receipt Confirmations                                          |
| Health & Safeguarding (Restricted) [v]                                         |
|    ├── Medical History & Injury Log                                               |
|    ├── Guardian Consent & Minor Safeguarding                                      |
|    └── WADA Anti-Doping Education & Testing Log                                   |
+-----------------------------------------------------------------------------------+

```

---

## 3. Detailed Data Models & Specifications

### 3.1 Athlete Master Profile Fields (15 Attributes)
1. **Athlete ID**: System auto-generated string (e.g. `NAMIS-UG-2026-08492`).
2. **National ID (NIN) / Passport Number**: Encrypted national identification string.
3. **Full Name**: Legal full name matching official identification.
4. **Gender**: Binary selection (`Male`, `Female`).
5. **Date of Birth**: Verified birth date.
6. **Age Category**: Dynamic computed category (`U10`, `U12`, `U15`, `U17`, `U20`, `Senior`).
7. **Nationality**: Default `Ugandan`.
8. **District of Origin**: District dropdown (e.g., *Kampala*, *Gulu*, *Mbale*, *Mbarara*, *Kapchorwa*).
9. **Region**: Regional grouping (*Central*, *North*, *East*, *West*).
10. **Current Residence**: Postal/physical residence address.
11. **Parents Names & Contacts**: Primary parental contact registry.
12. **Phone Contact**: Primary mobile phone number (+256 format).
13. **Email Address**: Registered email contact.
14. **Next of Kin**: Emergency next-of-kin full name & relation.
15. **Emergency Contact**: 24/7 emergency phone contact number.

---

### 3.2 Athlete Classification & Registration Attributes (7 Attributes)
1. **Name of Sport**: Primary sports discipline (e.g., *Athletics*, *Football*, *Netball*, *Boxing*, *Basketball*, *Rugby*, *Swimming*).
2. **Discipline / Event**: Specific event specialization (e.g., *5000m*, *Welterweight*, *Point Guard*, *100m Freestyle*).
3. **Federation / Association**: Primary recognized national body.
4. **Club / Academy / Institution**: Registered local club or university team.
5. **Registration Number**: Official federation license/registration code.
6. **Athlete Status**: State indicator (`Active`, `Injured`, `Retired`, `Suspended`).
7. **Date Registered with Federation**: Date of initial federation accreditation.

---

### 3.3 Dual Career & Education Module
* **School / Higher Education Institution**: Current or highest institution attended.
* **Highest Education Level**: `Primary`, `O-Level`, `A-Level`, `Diploma`, `Bachelors`, `Masters`.
* **Sports Scholarship Status**: Indicator of government or institutional sponsorship (`Active`, `Graduated`, `None`).
* **Employment Status & Employer**: Dual career employment tracking (e.g., *Uganda Police Forces*, *Uganda Prisons*, *Private Corporate*).

---

### 3.4 Government & NCS Support Log
* **Equipment Support Issued**: Running shoes, spike suits, boxing gloves, bicycles, wheelchairs.
* **Date Issued & Valuation**: Financial valuation (UGX) provided by government grant allocations.
* **Scholarship Financial Grants**: Tuition fee grants and living stipends.

---

## 4. Frontend Component Design (Vue 3 / Composition API)

```vue
<!-- Component: AthleteDashboardView.vue -->
<template>
  <div class="athlete-dashboard p-4">
    <!-- Hero Profile Header -->
    <div class="card mb-4 bg-primary text-white shadow-sm border-0">
      <div class="card-body d-flex align-items-center justify-content-between">
        <div>
          <h2 class="mb-1 fw-bold">{{ athlete.full_name }}</h2>
          <p class="mb-0 opacity-75">
            ID: {{ athlete.athlete_id }} | {{ athlete.sport_name }} ({{ athlete.discipline_event }})
          </p>
        </div>
        <span :class="['badge fs-6', statusBadgeClass(athlete.athlete_status)]">
          {{ athlete.athlete_status }}
        </span>
      </div>
    </div>

    <!-- Master Tabs -->
    <ul class="nav nav-pills mb-3" id="athlete-tabs">
      <li class="nav-item">
        <button class="nav-link active" @click="activeTab = 'profile'">Master Profile</button>
      </li>
      <li class="nav-item">
        <button class="nav-link" @click="activeTab = 'performance'">Performance & PBs</button>
      </li>
      <li class="nav-item">
        <button class="nav-link" @click="activeTab = 'education'">Dual Career & Support</button>
      </li>
    </ul>
  </div>
</template>
```
