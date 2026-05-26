# NCSMS v1 — National Council of Sports Management System

Design documentation for NCSMS v1. All working design documents live under `Modular Design/`.

---

## Application Submission Workflow

All user-facing applications follow a shared **6-stage wizard workflow**:

```
STAGE 1          STAGE 2         STAGE 3                STAGE 4              STAGE 5       STAGE 6
FILL FORM   →   REVIEW    →   DOWNLOAD & SIGN   →   PAYMENT       →   SUBMIT    →   TRACK
(Wizard)        (Summary)      (PDF + Upload)        (Digital or              (Confirm)     (Dashboard)
                                                      Proof Upload)
```

- **Fill Form** — Multi-step wizard with auto-save, Save Draft, and Resume from last step
- **Review** — Read-only summary with per-section edit links
- **Download & Sign** — System generates a pre-filled PDF; applicant signs digitally in-browser or prints, signs, and scans for upload
- **Payment** — Two methods accepted for all application and renewal types:
  - Online digital payment (card / mobile money)
  - Upload scanned proof of payment (bank slip, receipt)
- **Submit** — Checklist confirmation; unique `application_reference` assigned on submit
- **Track** — Status dashboard with reviewer notes, NEEDS_INFORMATION response flow, and notifications

Full workflow architecture: [`Modular Design/applications/application-workflow-v1.md`](Modular Design/applications/application-workflow-v1.md)

---

## Folder Structure

```
Modular Design/
├── ncsms-v1-module-design.md         Main architecture, security, middleware, shared guidance
├── super-admin-v1.md                 Super Admin module
├── general-secretary-v1.md           Application review and licence processing
├── roles-v1.md                       System roles and custom role permissions
├── user-profiles-v1.md               User profile and account categories
├── athletes-v1.md                    Athlete registration and credentialing
├── assets-v1.md                      Asset and inventory management
├── audit-v1.md                       Audit events and compliance
├── content-manager-v1.md             Website content and public pages
└── applications/
    ├── application-workflow-v1.md    Shared 6-stage wizard workflow architecture
    ├── applications-v1.md            Applications domain overview and index
    ├── national-sport-recognition-application-v1.md      Form 1 — Declaration of National Sport
    ├── federation-registration-renewal-v1.md             Form 3 — NSA / NSF Registration or Renewal
    ├── national-sports-association-transformation-v1.md  Form 5 — NSA Transformation (applicant view)
    ├── national-sports-federation-transformation-v1.md   Form 5 — NSA→NSF Transformation (full form)
    ├── sports-competition-organization-application-v1.md Form 7 — Organise a Sports Competition
    ├── sports-facility-operation-application-v1.md       Form 8 — Operate a Sports Facility
    ├── community-academy-registration-renewal-v1.md      Form 10 — Community Sports Club Reg/Renewal
    └── sports-academy-operation-application-v1.md        Form 11 — Operate a Sports Academy
```

---

## Design Document Conventions

1. Update shared architecture and security guidance in `ncsms-v1-module-design.md`
2. Module-specific changes go in the corresponding `*-v1.md` file
3. Application form changes go in the appropriate file under `applications/`
4. Workflow changes that affect all applications go in `application-workflow-v1.md`
5. Keep this `README.md` as a navigation entrypoint — detailed design stays in the module files
