# NCS Content Management System (CMS) & Page Builder Module Implementation Plan

> **Target File:** `/home/fidi/Projects/NCS_Intranet/Docs/plans/cms-page-builder-plan.md`  
> **Role:** Content Manager / Public Relations Officer / Systems Admin  
> **System Scope:** Public Website & Internal Portal CMS (`CMSView.vue`, `PageBuilderView.vue`)  
> **Authority Scope:** Public Website Content Publishing, Dynamic Drag-and-Drop Page Builder, Hero Banners, Sports Event News Articles, Press Release Sync, Image/Media Storage (`StorageSettingsView.vue`)  
> **Technical Stack:** Vue 3 Page Builder, Go REST API, S3/Local Storage Controller  

---

## 1. Executive Purpose & Architecture

The **CMS & Page Builder Module** gives the Public Relations and IT teams full control over the official National Council of Sports public web portal and intranet landing pages without requiring code changes.

```
+-----------------------------------------------------------------------------------+
|                        NCS CMS & DRAG-AND-DROP PAGE BUILDER                       |
+-----------------------------------------------------------------------------------+
       |                                  |                                 |
       v                                  v                                 v
+-----------------------+      +-----------------------+      +---------------------+
| 1. CMS MANAGEMENT     |      | 2. DRAG-AND-DROP PAGE |      | 3. MEDIA STORAGE &  |
|    (`CMSView.vue`)    |      |    BUILDER ENGINE     |      |    ASSET CONTROLLER |
+-----------------------+      +-----------------------+      +---------------------+
| - News Articles       |      | - Hero Slider Blocks  |      | - Local / S3 Storage|
| - Event Calendars     |      | - Dynamic Grid Layouts|      | - Image Optimization|
| - Media Photo Gallery |      | - Widget Components   |      | - Storage Quotas    |
+-----------------------+      +-----------------------+      +---------------------+
```

---

## 2. Key Modules & User Interface Specifications

### 2.1 CMS Content Publishing Workspace (`CMSView.vue`)
- **News & Articles Manager:** Create, edit, schedule, and publish sports news, competition results, press statements, and official circulars.
- **Media Library & Asset Manager (`StorageSettingsView.vue`):** High-resolution photo gallery for national games, press conferences, and stadium facility photos.

### 2.2 Drag-and-Drop Visual Page Builder (`PageBuilderView.vue`)
- **Visual Layout Canvas:** Live preview layout editor allowing administrators to drag, re-order, configure, and publish custom blocks (Hero Banners, Key Metrics Widgets, Recent News Grid, Federation Directory List, Contact Footer).

---

## 3. Database Schema & REST API Mapping

```sql
CREATE TABLE cms_pages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    layout_data JSONB NOT NULL, -- Visual block layout definitions
    is_published BOOLEAN DEFAULT FALSE,
    author_id UUID NOT NULL REFERENCES users(id),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

| Endpoint | Method | Scope | Function |
| :--- | :--- | :--- | :--- |
| `/api/v1/cms/pages` | GET/POST | Content Manager | Create and list custom visual web pages |
| `/api/v1/cms/pages/:slug` | PUT | Content Manager | Save layout_data JSON from visual Page Builder |
| `/api/v1/cms/media/upload` | POST | Content Manager | Upload and optimize web assets |
