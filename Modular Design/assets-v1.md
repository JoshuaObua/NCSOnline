# Assets Module — NCSMS v1.0

## Purpose

- Manage inventory, equipment, and asset tracking for federations and programs
- Record asset assignment, status, and accountability
- Support audit-grade tracking of asset lifecycle events

## Domain Responsibilities

- Define asset categories, locations, and ownership metadata
- Create and update asset records for federation equipment, facilities, and supplies
- Track asset transfers, disposals, and condition updates
- Support reconciliation and inventory reporting

## High-level Code Layout

- `/internal/domain/assets/models.go`
  - `Asset`, `AssetCategory`, `AssetAssignment`, `AssetCondition`
- `/internal/domain/assets/repository.go`
  - `CreateAsset(ctx, asset)`, `FindAssetByID(ctx, id)`, `ListAssets(ctx, filter)`, `UpdateAsset(ctx, asset)`
- `/internal/domain/assets/service.go`
  - `RegisterAsset(ctx, req)`, `TransferAsset(ctx, req)`, `UpdateAssetCondition(ctx, req)`, `DeactivateAsset(ctx, req)`
- `/internal/domain/assets/handler.go`
  - HTTP endpoints: `POST /api/v1/assets`, `GET /api/v1/assets/{id}`, `PATCH /api/v1/assets/{id}`, `GET /api/v1/assets`

## Access and Policy

- `ADMIN` roles manage assets across federations or globally depending on scope
- `FEDERATION_USER` may manage assets within their own federation if allowed
- Asset data is protected with RLS and audit trails for every mutation

## Example Routes

- `POST /api/v1/assets` -> register a new asset
- `GET /api/v1/assets/{id}` -> retrieve asset details
- `PATCH /api/v1/assets/{id}` -> update asset metadata or status
- `GET /api/v1/assets` -> list and filter assets
