# NCSMS v1 Design

All design documentation for NCSMS v1 is stored under the `Modular Design` folder.

## Folder structure

- `Modular Design/`
  - `ncsms-v1-module-design.md` — main architecture, security, middleware, and shared guidance
  - `super-admin-v1.md` — Super Admin module design
  - `admin-v1.md` — Admin module design
  - `federations-v1.md` — Federations module design
  - `super-admin.md`, `admin.md` — legacy or supplementary module notes

## Editing workflow

1. Update the shared architecture and security guidance in `Modular Design/ncsms-v1-module-design.md`.
2. Modify module-specific details in the corresponding `*-v1.md` file.
3. Keep the root `README.md` as a simple entrypoint and navigation guide.

## Notes

- The main repository root is intentionally kept minimal.
- The `Modular Design` folder contains all the working design documents.
