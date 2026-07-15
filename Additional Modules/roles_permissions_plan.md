# NAMIS Security Roles and Permissions Plan

This plan outlines the integration of system-wide administrative roles, security groups, and resource permission levels to manage sports data within NCS.

## 1. Role Definitions & Permissions Mapping

We assign standard permission scopes to roles to dictate who can view or modify records:

### Role: `ncs_general_secretary`
- Full cross-federation reporting access, board report compilation, and council dashboards.
- Permissions: `reports:review:any`, `federations:read:any`, `dashboard:governance:read`, `dashboard:athletes:read`, `dashboard:performance:read`, `dashboard:finance:read`, `dashboard:talent:read`.

### Role: `technical_department`
- Federation audits, tracking performances, talent, and coaching certifications.
- Permissions: `reports:review:any`, `federations:read:any`, `dashboard:athletes:read`, `dashboard:performance:read`, `dashboard:talent:read`.

### Role: `finance_department`
- Monitoring disbursements, financial audits, and budget approvals.
- Permissions: `reports:review:any`, `federations:read:any`, `dashboard:finance:read`.

### Role: `federation_president`
- Submitting final approvals for their own federation data and reporting.
- Permissions: `reports:read:own`, `reports:approve:own`, `federations:read:own`.

### Role: `federation_general_secretary`
- Creating, updating, and draft-saving athlete lists, coaches, talent pathways, and competitions.
- Permissions: `reports:read:own`, `reports:write:own`, `federations:read:own`, `federations:write:own`.

### Role: `safeguarding_officer`
- Restricted access to medical conditions, parent details, and child safety cases.
- Permissions: `safeguarding:cases:manage`.

### Role: `auditor`
- Read-only lookup permissions for compliance audits.
- Permissions: `federations:read:own`, `reports:read:own`.

## 2. Integration with User Manager Panel

To allow administrators to assign these roles:
- The backend uses dynamic fetching of roles from the database (`List` method in `RolesHandler`).
- When seeding roles using `016_nsmis_reporting_foundation.sql`, all these roles are marked as `is_system = TRUE`, making them immediately returned by `/api/v1/admin/roles`.
- We will verify that the frontend User list and details edit form in the admin panel fetches this updated role list and lists them as select options, allowing admins to map users to their official department or federation role.
