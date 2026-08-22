# NAMIS Talent Identification, Medical & Safeguarding Dashboard Specification

## 1. Overview
The **Talent Identification, Medical & Safeguarding Dashboard** is a secure, role-restricted module within NAMIS. It governs the National Council of Sports grass-roots talent development pipeline (School → District → Talent Centre → Elite National Squad), protects confidential athlete medical records, manages minor safeguarding/guardian consent, and ensures WADA anti-doping compliance.

> [!CAUTION]
> Medical, Injury, Safeguarding, and Anti-Doping data contain sensitive health and minor protection records. Access is strictly restricted to certified Medical Officers, Safeguarding Officers, and NCS Super Administrators with full audit log tracking.

---

## 2. System Security Architecture & Permission Scoping

```
+---------------------------------------------------------------------------------+
|               RESTRICTED SECURITY & DATA ACCESSIBILITY MATRIX                   |
+---------------------------------------------------------------------------------+
| Domain                | Athlete   | Coach       | Medical/Safeguarding | NCS Admin |
| --------------------- | --------- | ----------- | -------------------- | --------- |
| Talent Scouting       | Read Own  | Write Squad | Read                 | Full      |
| Medical Records       | Read Own  | Restricted  | Read/Write           | Audited   |
| Safeguarding Consent  | Read Own  | Restricted  | Read/Write           | Audited   |
| Anti-Doping WADA      | Read Own  | Read Squad  | Read/Write           | Full      |
+---------------------------------------------------------------------------------+
```

---

## 3. UI Layout & View Structure

```
+-----------------------------------------------------------------------------------+
|  [NCS Logo]  SAFEGUARDING & TALENT DEVELOPMENT DESK   [Encrypted Session] [Officer] |
+-----------------------------------------------------------------------------------+
| [Sidebar Navigation Menu]                                                         |
|  Safeguarding & Talent Overview                                                 |
|  Talent Identification Pipeline [v]                                             |
|    ├── Scouting Candidates Registry                                               |
|    ├── School & Regional Talent Centers                                           |
|    ├── Recommended Development Pathways                                           |
|    └── Government Talent Scholarships                                             |
|  Medical Desk (Restricted) [v]                                                   |
|    ├── Athlete Medical Records Vault                                              |
|    ├── Injury & Rehabilitation Tracker                                            |
|    └── Insurance & Medical Clearances                                             |
|  Minor Safeguarding Vault (Restricted) [v]                                       |
|    ├── Guardian Consent Forms Registry                                            |
|    ├── Assigned Welfare Officers                                                  |
|    └── Incident Logging & Case Management                                         |
|  WADA Anti-Doping Engine [v]                                                    |
|    ├── National Testing Pool Registry                                             |
|    ├── Test Results & Lab Logs                                                    |
|    ├── WADA Digital Education Certificates                                        |
|    └── Suspension & Disciplinary Sanctions                                        |
+-----------------------------------------------------------------------------------+

```

---

## 4. Functional Modules & Data Fields

### 4.1 Talent Identification Module
Core component for fulfilling the NCS mandate of national grass-roots talent development:
* **Identification Date & Scout**: Date of scouting event and accredited scout identifier.
* **School & District**: Primary school, secondary school, or tertiary institution.
* **Talent Centre**: Accredited regional development center (*Gulu*, *Mbale*, *Fort Portal*, *Masaka*, *Kampala*).
* **Talent Category**: Age & skill assessment rating.
* **Recommended Pathway**: Development roadmap (*School Games → Regional Academy → Federation Club → National Team*).
* **Scholarship Status**: Status of government sports sponsorship (`Recommended`, `Approved`, `Active`).

---

### 4.2 Medical Records Desk (Restricted Access)
* **Blood Group & Allergies**: Critical emergency medical parameters.
* **Injury History & Current Status**: Comprehensive injury tracking log (*Fit*, *Rehabilitation*, *Surgery Required*).
* **Medical Insurance Provider & Policy**: Valid policy coverage for national team tours.

---

### 4.3 Safeguarding & Minor Protection Desk
* **Guardian Details & Consent**: Required legal consent for athletes under 18 years of age.
* **Manager / Agent Info**: Registered representative contact details.
* **Safeguarding Officer Assigned**: Designated welfare officer responsible for the athlete.
* **Anti-Doping Education Status**: WADA digital learning completion certificate.

---

### 4.4 Anti-Doping Compliance Engine
* **Testing Status**: In-competition / Out-of-competition testing log.
* **Test Results**: Official lab findings (`Negative`, `Pending`, `AAF`).
* **WADA Clearance**: Date of completed anti-doping educational module.
* **Suspension History**: Formal disciplinary suspension logging.
