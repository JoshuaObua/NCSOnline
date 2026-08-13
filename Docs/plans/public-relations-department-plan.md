# NCS Public Relations & Corporate Communications Department Implementation Plan

> **Target File:** `/home/fidi/Projects/NCS_Intranet/Docs/plans/public-relations-department-plan.md`  
> **Role:** Head of Public Relations / Corporate Communications Officer / Media Relations Lead  
> **Department:** Public Relations & Communications Department - National Council of Sports (NCS)  
> **Authority Scope:** Press Releases, Media Accreditation, Social Media Channels, Official Announcements, Public Event Briefings, CMS Content Publishing (`CMSView.vue`, `PageBuilderView.vue`)  
> **Compliance Standards:** Government Communications Guidelines, Access to Information Act (Uganda), Uganda Media Council Regulations  

---

## 1. Executive Role Overview & Mission

The **Public Relations & Corporate Communications Department** is responsible for public information dissemination, brand management, media accreditation for international sports events, press release broadcasting, crisis communications, and website/CMS portal publishing for the National Council of Sports.

This plan details the **Public Relations & Media Command Portal** inside `NCS_Intranet`. It integrates directly with the **Content Management System (CMS)** (`CMSView.vue`, `PageBuilderView.vue`) and public announcement bulletin board.

---

## 2. Key Modules & User Interface Specifications

```
+-----------------------------------------------------------------------------------+
| [≡] NCS INTRANET | Public Relations & Media Command Workspace [➕ Press Release] [📰 CMS Post] |
+-----------------------------------------------------------------------------------+
| [PR & MEDIA ENGAGEMENT KPI CARDS]                                                 |
| +--------------------+ +--------------------+ +-------------------+ +-------------------+ |
| | Press Releases     | | Media Accreditations| | Website Traffic   | | Social Engagement | |
| | 24 Published YTD   | | 145 Journalists   | | 85.4K Visits / Mo | | 94.2% Positive    | |
| +--------------------+ +--------------------+ +-------------------+ +-------------------+ |
+-----------------------------------------------------------------------------------+
| [PR WORKSPACE: (1) Press Releases | (2) Media Accreditation | (3) CMS & Website | (4) Bulletins] |
| +-------------------------------------------------------------------------------+ |
| | [PRESS RELEASE & PUBLIC STATEMENT PUBLISHER]                                  | |
| | Title / Headline               | Target Audience | Release Date | Status      | |
| | NCS Statement on Lugogo Turf   | General Public  | Jul 25, 2026 | Published   | |
| | AFCON 2027 Stadium Readiness   | International   | Jul 18, 2026 | Published   | |
| | Federation Elections Circular  | Sports Media    | Aug 02, 2026 | Draft Review| |
| | [ Create Press Release ] [ Schedule Broadcast ] [ Submit to GS for Sign-off ] | |
| +-------------------------------------------------------------------------------+ |
| +-----------------------------------------------+ +-------------------------------+ |
| | Media Accreditation & Journalist Registry    | | CMS Web Portal Page Builder   | |
| | Approved Journalists : 145 Accredited        | | Homepage Banner  : Updated    | |
| | Pending Applications : 12 Journalists        | | News Articles    : 48 Live    | |
| | [ Review Media Badges ] [ Export Press List ] | | [ Open CMS Page Builder ]     | |
| +-----------------------------------------------+ +-------------------------------+ |
+-----------------------------------------------------------------------------------+
```

### Key Workstation Tabs & Features

1. **Press Release & Official Statement Publisher:**
   - Drafting, review, and publication of official press releases, government sports advisories, and event briefings.
   - Approval routing to General Secretary (GS) for official statements prior to media broadcast.

2. **Media Accreditation & Press Pass Management:**
   - Online journalist registration portal for national and international sports reporters covering events at Lugogo, Namboole, and national championships.
   - Badge generation, media pass clearance, and media tribunal credentials.

3. **Content Management System (CMS) & Portal Builder (`CMSView.vue`, `PageBuilderView.vue`):**
   - Direct integration with the public NCS web portal for publishing news articles, photo galleries, event calendars, sports results, and board policies.

4. **Media Monitoring & Public Information Broadcast:**
   - Tracking print, television, radio, and digital media coverage regarding Ugandan sports governance.
   - Emergency crisis communication messaging and public notice distribution.

---

## 3. Database Schema & Technical Architecture

```sql
CREATE TABLE press_releases (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    headline VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    category VARCHAR(64) NOT NULL, -- OFFICIAL_STATEMENT, EVENT_BRIEF, MEDIA_ADVISORY
    body_content TEXT NOT NULL,
    featured_image_url TEXT,
    author_id UUID NOT NULL REFERENCES users(id),
    status VARCHAR(32) DEFAULT 'DRAFT', -- DRAFT, PENDING_GS_APPROVAL, PUBLISHED
    published_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE media_accreditations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    applicant_name VARCHAR(128) NOT NULL,
    media_house VARCHAR(128) NOT NULL,
    media_type VARCHAR(32) NOT NULL, -- TV, RADIO, PRINT, ONLINE, PHOTOGRAPHER
    press_card_number VARCHAR(64) NOT NULL,
    event_name VARCHAR(255) NOT NULL,
    status VARCHAR(32) DEFAULT 'PENDING', -- PENDING, APPROVED, REJECTED
    issued_badge_code VARCHAR(64) UNIQUE,
    approved_by UUID REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

---

## 4. REST API Mapping (`backend/internal/handlers/cms.go`)

| Endpoint | Method | Description |
| :--- | :--- | :--- |
| `/api/v1/cms/press-releases` | GET/POST | Create and query press releases |
| `/api/v1/cms/press-releases/:id/approve` | PUT | GS sign-off and publish press release |
| `/api/v1/cms/media-accreditation` | GET/POST | Review and approve media pass applications |
| `/api/v1/cms/pages` | GET/POST | Manage public portal web pages via PageBuilder |

---

## 5. Implementation Verification Roadmap

- [ ] Connect `CMSView.vue` and `PageBuilderView.vue` to Go backend handler `backend/internal/handlers/cms.go`.
- [ ] Implement press release approval workflow routing to General Secretary.
