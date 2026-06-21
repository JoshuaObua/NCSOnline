# NCS-Online — Public URL Routes

## Public Pages (PublicLayout)

| URL | Route Name | Component | Description |
|-----|-----------|-----------|-------------|
| `/` | Home | `HomeView.vue` | Homepage — hero slides, fun facts, featured content |
| `/news` | News | `BlogView.vue` | News & blog listing with dynamic category filter chips |
| `/news/:slug` | NewsPost | `BlogPostView.vue` | Individual article — HTML content, author name, view count |
| `/pages/:slug` | Page | `PageView.vue` | Standalone CMS-managed static pages (About, Privacy, etc.) |
| `/events` | Events | `EventsView.vue` | Upcoming and past sports events listing |
| `/events/:slug` | EventDetail | `EventDetailView.vue` | Individual event detail page |
| `/careers` | Careers | `CareersView.vue` | Jobs, tenders and internship opportunities |
| `/careers/:id` | CareerDetail | `CareerDetailView.vue` | Individual job/tender posting detail |
| `/projects` | Projects | `ProjectsView.vue` | Key projects and development initiatives |
| `/case-studies` | CaseStudies | `CaseStudiesView.vue` | In-depth impact stories and case studies |
| `/apply` | Apply | `LicensePortalView.vue` | Sports licence online application portal |
| `/resource-centre` | ResourceCentre | `ResourceCentreView.vue` | Downloadable guidelines, reports and rules |
| `/resources` | — | redirect | Redirects to `/resource-centre` |
| `/facilities` | Facilities | `FacilitiesView.vue` | Sports facilities and venues directory |
| `/facilities/:slug` | FacilityDetail | `FacilityDetailView.vue` | Individual facility detail page |
| `/associations` | Associations | `AssociationsView.vue` | Sports associations and federations |
| `/invest` | Invest | `InvestView.vue` | Investment opportunities in the sports sector |
| `/faqs` | FAQs | `FAQsView.vue` | Frequently asked questions |
| `/contact-us` | ContactUs | `ContactUsView.vue` | Contact form and office information |

## Auth Routes (no layout wrapper)

| URL | Route Name | Component | Notes |
|-----|-----------|-----------|-------|
| `/login` | Login | `LoginView.vue` | Redirects authenticated users to their portal |
| `/register` | Register | `RegisterView.vue` | Self-registration for applicants; redirects if already logged in |

## Authenticated Routes

| URL | Route Name | Component | Roles |
|-----|-----------|-----------|-------|
| `/my-portal` | ApplicantPortal | `ApplicantPortalView.vue` | `user` / applicant |
| `/dashboard` | Dashboard | `DashboardView.vue` | All authenticated users |
| `/applications` | Applications | `ApplicationsView.vue` | All authenticated users |
| `/profile` | Profile | `ProfileView.vue` | All authenticated users |
| `/users` | Users | `UsersView.vue` | `super_admin`, `admin` |
| `/audit-logs` | AuditLogs | `AuditLogsView.vue` | `super_admin`, `admin` |
| `/roles` | Roles | `RolesView.vue` | `super_admin` |
| `/cms` | CMS | `CMSView.vue` | `super_admin`, `admin`, `content_manager` |

## Dynamic Segments

| Segment | Used in | Format |
|---------|---------|--------|
| `:slug` | `/news/:slug`, `/events/:slug`, `/facilities/:slug`, `/pages/:slug` | URL-safe string (e.g. `ncs-launches-new-programme`) |
| `:id` | `/careers/:id` | UUID string |

## Catch-all

`/:pathMatch(.*)* → /` — any unmatched path redirects to the homepage.
