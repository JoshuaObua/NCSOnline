# Content Manager Module — NCSMS v1.0

## Overview

The Content Manager module provides a secure CMS-style experience for website content and public page management. It supports blog publishing, page authoring, content sections, widgets, testimonials, project showcases, contact page content, and SEO metadata.

This module is intended for non-technical content staff, website editors, and communications teams who need a structured interface to manage public-facing content without direct database or system configuration access.

## Domain Scope

- Blog posts and article content
- Blog categories and subcategories
- Public website pages and sections
- Contact page content and metadata
- Testimonials and featured quotes
- Footer widgets, banners, and promotional cards
- Project entries and portfolio showcases
- Page composition and SEO metadata
- Content publishing workflows and status

## Key Entities

- `content_pages`
  - `id`, `slug`, `title`, `summary`, `body`, `status`, `locale`, `seo_title`, `seo_description`, `seo_keywords`, `published_at`, `author_id`, `updated_by`
  - Supports static pages, landing sections, policy pages, and homepage sections.

- `blog_posts`
  - `id`, `slug`, `title`, `excerpt`, `body`, `status`, `category_id`, `tags`, `author_id`, `published_at`, `feature_image`, `read_time_minutes`, `seo_title`, `seo_description`

- `blog_categories`
  - `id`, `name`, `slug`, `description`, `parent_category_id`, `status`

- `testimonials`
  - `id`, `name`, `role`, `quote`, `organization`, `photo_url`, `status`

- `widgets`
  - `id`, `name`, `type`, `payload`, `position`, `status`

- `projects`
  - `id`, `title`, `summary`, `details`, `status`, `start_date`, `end_date`, `media_urls`, `project_url`

- `content_media`
  - `id`, `object_type`, `object_id`, `media_type`, `storage_path`, `caption`, `alt_text`, `uploaded_by`

## API Surface

### Content Manager endpoints

- `GET /api/v1/content/pages` — list pages for editing and preview
- `GET /api/v1/content/pages/:slug` — fetch single page details
- `POST /api/v1/content/pages` — create a new page
- `PUT /api/v1/content/pages/:id` — update page content and metadata
- `DELETE /api/v1/content/pages/:id` — remove outdated pages (explicit permission)
- `POST /api/v1/content/pages/:id/publish` — publish page to public website

- `GET /api/v1/content/blog` — list blog posts, filters by category/status
- `GET /api/v1/content/blog/:slug` — fetch blog details
- `POST /api/v1/content/blog` — create blog draft
- `PUT /api/v1/content/blog/:id` — update blog post
- `DELETE /api/v1/content/blog/:id` — delete blog post
- `POST /api/v1/content/blog/:id/publish` — publish blog post

- `GET /api/v1/content/categories` — list blog categories
- `POST /api/v1/content/categories` — create category
- `PUT /api/v1/content/categories/:id` — update category
- `DELETE /api/v1/content/categories/:id` — archive category

- `GET /api/v1/content/testimonials`
- `POST /api/v1/content/testimonials`
- `PUT /api/v1/content/testimonials/:id`
- `DELETE /api/v1/content/testimonials/:id`

- `GET /api/v1/content/widgets`
- `POST /api/v1/content/widgets`
- `PUT /api/v1/content/widgets/:id`
- `DELETE /api/v1/content/widgets/:id`

- `GET /api/v1/content/projects`
- `POST /api/v1/content/projects`
- `PUT /api/v1/content/projects/:id`
- `DELETE /api/v1/content/projects/:id`

- `POST /api/v1/content/media` — upload or link media assets
- `DELETE /api/v1/content/media/:id`

## Access Control & RBAC

- Primary role: `CONTENT_MANAGER`
- `CONTENT_MANAGER` can `VIEW`, `CREATE`, `UPDATE`, and optionally `DELETE` public content entities.
- Publishing actions may require an additional `publish` permission flag or reviewer workflow for sensitive content.
- `SYSTEM_ADMIN` retains full content module control.
- `GENERAL_SECRETARY` can view and approve website content when configured as a content reviewer.
- Public users can access `VIEW`-only published content through the front-end website.

### Role guidelines

- `CONTENT_MANAGER` should not have access to system or configuration modules such as `roles`, `audit`, `assets`, or federation financial data.
- A separate custom role such as `CONTENT_EDITOR` can be defined for users who need draft access but not publishing rights.
- `CUSTOM_ROLE` can be used to create scoped editors with `CREATE`, `UPDATE`, and `VIEW` on website content without broader system privileges.

## Supabase & Storage

- Store uploaded content media in a private Supabase bucket with read access granted only to published front-end endpoints or signed URLs.
- Use structured storage paths like `content/blog/:post_id/:filename`, `content/pages/:page_id/:filename`, and `content/media/:media_id`.
- Maintain media metadata in `content_media` for image alt text, captions, and ownership audit.
- Signed URLs for public content assets should expire quickly and be renewed on-demand by the front-end.

## Security and Validation

- Validate all HTML or rich text input before saving. Prefer sanitized Markdown or a structured rich text schema.
- Strip unsafe tags and JavaScript from `body`, `excerpt`, `summary`, and widget payloads.
- Enforce strict slug uniqueness and prevent path traversal in page slugs.
- Protect public page endpoints from unauthorized preview access by requiring a valid admin/session token for draft routes.
- Log content publish and edit events in audit logs with actor, entity type, entity id, action, and timestamp.

## Audit and Change Tracking

- Track `created_by`, `updated_by`, `published_by`, and `published_at` on major content entities.
- Record content status transitions such as `draft`, `review`, `published`, `archived`.
- Use audit events for content publish approvals, editor actions, and deleted content removals.
- Allow rollbacks to prior versions or archived copies via audit-backed content history if required by governance.

## Public Website Workflow

- Draft content is saved privately and only visible to authenticated content managers.
- Published content becomes available on public pages and blog feeds via public API or static page generation.
- Page composition allows component-driven sections such as hero banners, featured posts, testimonials, and CTA widgets.
- Contact page content includes form copy, address details, map embed metadata, and phone/email display fields.
