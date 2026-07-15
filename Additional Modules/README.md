# National Athlete Management Information System (NAMIS) - Implementation Plans

This folder contains the modular implementation plans for adding the comprehensive athlete management capabilities to the NCS Portal, based on the requirements defined in [Atheletes.md](file:///c:/NCSPortal/Additional%20Modules/Atheletes.md).

## Module Architecture & Plans Index

We have structured the implementation plans into separate, highly-detailed documents for each architectural area:

1. **[Clubs Module Plan](file:///c:/NCSPortal/Additional%20Modules/clubs_module.md)**
   - Database schema details for athletic clubs, academies, and school sports associations.
   - Core CRUD mappings under the generic repository layer.

2. **[Federations Module Plan](file:///c:/NCSPortal/Additional%20Modules/federations_module.md)**
   - Managing physical addresses, contact information, and licensing status.
   - Assigning roles like President, General Secretary, Treasurer, and Arbitrator.

3. **[Athletes Module Plan](file:///c:/NCSPortal/Additional%20Modules/athletes_module.md)**
   - Master profiles, NIN/Passport records, personal details, and emergency contacts.
   - Dual-career educational pathways, sports scholarships, and occupation tracking.
   - Medical logs (restricted blood group and injury states), child safeguarding controls, and consent.
   - Anti-doping compliance tracking (testing pools, history, WADA certifications).

4. **[Coaches Module Plan](file:///c:/NCSPortal/Additional%20Modules/coaches_module.md)**
   - Coaches registry, license numbers, expiration status, and certification levels.
   - Training entourage support staff mapping (Strength & Conditioning, Physiotherapists, Nutritionists, Doctors, Managers).

5. **[Performance & Medals Module Plan](file:///c:/NCSPortal/Additional%20Modules/performance_module.md)**
   - Tracking competition details (from local District up to Olympic level).
   - Fine-grained athlete results (Time, Distance, Weight, Score, Ranking, PB/SB/NR flags).
   - Medal registry (Gold, Silver, Bronze), funding types, and official NCS Recognition.

6. **[Talent Identification Module Plan](file:///c:/NCSPortal/Additional%20Modules/talent_identification_module.md)**
   - Talent scouting profiles, age at identification, schools, and scout credentials.
   - Training pathway recommendations and talent scholarships tracking.

7. **[National Team Module Plan](file:///c:/NCSPortal/Additional%20Modules/national_team_module.md)**
   - National squad selections, category tracking, appearances count, and call-up dates.

8. **[Analytics Module Plan](file:///c:/NCSPortal/Additional%20Modules/analytics_module.md)**
   - KPI aggregate query routines and reporting models for council leaders.
   - Dashboard layouts for geographic regional sports development reports and gender stats.

9. **[Security Roles & Permissions Plan](file:///c:/NCSPortal/Additional%20Modules/roles_permissions_plan.md)**
   - Access control mapping for `ncs_general_secretary`, `technical_department`, `finance_department`, `federation_president`, `federation_general_secretary`, `safeguarding_officer`, and `auditor`.
   - User panel role assignments.

---

## Technical Approach

All plans align with the existing Go backend design patterns (using the `NSMISRepo` generic repository mapping framework) and reuse the pre-seeded DB tables inside Postgres for maximum performance and minimum friction.
