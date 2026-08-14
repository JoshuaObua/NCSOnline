# NCS Stores & Inventory Management Department Dashboard Implementation Plan

> **Target File:** `/home/fidi/Projects/NCS_Intranet/Docs/plans/stores-inventory-department-plan.md`  
> **Role:** Stores Officer / Inventory Manager / Warehouse Custodian  
> **Department:** Stores & Inventory Management Unit - National Council of Sports (NCS)  
> **Reporting Line:** Finance & Administration Directorate / Assistant General Secretary - Administration (AGS-A)  
> **Operational Scope:** Central Warehouse Management, Goods Received Notes (GRN), Electronic Bin Cards, Material Requisitions & Issuances, Sports Equipment Pool, Stock Takes & Reconciliation with Accounting, and Obsolescence Flagging.  
> **Compliance Standards:** Public Finance Management Act (PFMA 2015), Treasury Instructions (Stores Regulations), PPDA Asset Management Guidelines.

---

## 1. Role Overview & Stores Management Architecture

The **Stores & Inventory Management Unit** maintains physical custody and inventory accounting of all stock items, sports equipment, engineering spares, stationery, and maintenance supplies belonging to the National Council of Sports. It ensures uninterrupted supply to operational departments while preventing stock loss, leakage, and obsolescence.

```
+---------------------------------------------------------------------------------------------------+
|               NCS STORES & INVENTORY OPERATIONAL WORKFLOW                                         |
+---------------------------------------------------------------------------------------------------+
       |                                      |                                      |
       v                                      v                                      v
+-----------------------+              +-----------------------+              +-----------------------+
| 1. GOODS RECEIPT      |              | 2. INVENTORY CONTROL  |              | 3. STOCK ISSUANCE     |
| - Delivery Note Match |              | - Electronic Bin Cards|              | - Departmental Reqs   |
| - PDU PO Verification |              | - Reorder Thresholds  |              | - Store Issue Voucher |
| - Quality Inspection  |              | - Batch & Location Log|              | - Handover Signatures |
| - Issue Electronic GRN|              | - Stock Aging & Count |              | - Automated Deduction |
+-----------------------+              +-----------------------+              +-----------------------+
       |                                      |                                      |
       +--------------------------------------+--------------------------------------+
                                              |
                                              v
                       +-----------------------------------------------+
                       | RECONCILIATION & DISPOSAL PIPELINE            |
                       | - Monthly Stock Reconciliation with Accounts  |
                       | - High-Performance Sports Gear Pool Tracking  |
                       | - Damaged / Obsolete Stock Board of Survey    |
                       | - Escalation to AGS-A & GS for Write-Offs     |
                       +-----------------------------------------------+
```

---

## 2. Key Modules & User Interface Specifications

---

### 2.1 Goods Received Note (GRN) & Inward Inspection Hub
- **Procurement PO Verification:** Validates incoming supplier deliveries against approved Purchase Orders and PPDA Form 5 contracts.
- **Quality & Quantity Inspection Checklist:** Inspects technical specifications before goods acceptance.
- **Electronic GRN Generation:** Issues digital GRN with timestamp, vendor details, and inspector signature, automatically alerting Accounts for invoice matching.

### 2.2 Electronic Bin Cards & Stock Ledger
- Real-time stock balance tracking across 4 main store categories:
  1. **Sports Equipment & Competition Gear:** Balls, jerseys, nets, timing mats, boxing gloves, athletics implements.
  2. **Engineering & Maintenance Spares:** Plumbing pipes, electrical cables, LED floodlight bulbs, turf fertilizers, mower blades.
  3. **ICT Consumables & Light Hardware:** Toners, paper reams, network patch cords, power surge strips.
  4. **Office & Administrative Supplies:** Stationeries, printed registers, cleaning detergents, corporate branded collaterals.
- Automated low-stock alerts when inventory drops below safety buffer thresholds.

### 2.3 Store Requisition & Issuance System (SRV / SIV)
- Cross-departmental electronic requisition workflow: Staff submit item requests $\rightarrow$ HOD endorses $\rightarrow$ Stores Officer reviews & issues Store Issue Voucher (SIV).
- Digital handover signature capture on mobile/desktop.

### 2.4 High-Performance Gear & Federation Loan Registry
- Tracks tournament equipment loaned out to national sports federations or national teams with return condition logging.

### 2.5 Physical Stock Count & Accounting Reconciliation
- Bi-annual electronic stock count module enabling mobile tablet spot-checks.
- Automated variance report highlighting book inventory vs physical stock counts, sent to Internal Auditor and Accountant.

---

## 3. Stores Dashboard User Interface (`StoresInventoryDashboard.vue`)

```
+-----------------------------------------------------------------------------------+
| [≡] NCS INTRANET | Stores & Inventory Control Portal  [➕ New GRN] [📦 Issue Stock] [Stores]|
+-----------------------------------------------------------------------------------+
| [STORES & INVENTORY KPI CARDS]                                                    |
| +--------------------+ +--------------------+ +-------------------+ +-------------------+ |
| | Total Stock Value  | | Active SKU Items   | | Low Stock Alerts  | | Pending Issuances | |
| | UGX 412.8 Million  | | 642 Stock Items    | | 7 Items to Order  | | 4 Dept Vouchers   | |
| +--------------------+ +--------------------+ +-------------------+ +-------------------+ |
+-----------------------------------------------------------------------------------+
| [STORES WORKSPACE: (1) Stock Ledger | (2) Inward GRN | (3) Issuance Vouchers | (4) Federation Gear] |
| +-------------------------------------------------------------------------------+ |
| | [REAL-TIME BIN CARD & STOCK INVENTORY TABLE]                                  | |
| | SKU Code   | Item Description         | Category     | Qty on Hand | Status   | |
| | SKU-SP-042 | Match Soccer Balls (FIFA)| Sports Gear  | 145 Units   | In Stock | |
| | SKU-ENG-19 | Floodlight Bulbs 1000W   | Engineering  | 4 Units     | LOW STOCK| |
| | SKU-ICT-08 | HP Laserjet Toner 85A    | ICT Supplies | 18 Units    | In Stock | |
| | SKU-ADM-31 | A4 Paper Reams (Box)     | Admin Supplies| 32 Boxes    | In Stock | |
| | [ Export Monthly Stores Ledger PDF ] [ Generate Stock Reconciliation Brief ]  | |
| +-------------------------------------------------------------------------------+ |
+-----------------------------------------------------------------------------------+
```

---

## 4. Database Schema & Technical Architecture

```sql
CREATE TABLE store_inventory_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sku_code VARCHAR(64) UNIQUE NOT NULL,
    item_name VARCHAR(255) NOT NULL,
    category VARCHAR(64) NOT NULL, -- SPORTS_GEAR, ENGINEERING_SPARES, ICT_SUPPLIES, OFFICE_SUPPLIES
    unit_of_measure VARCHAR(32) NOT NULL, -- Units, Boxes, Meters, Pairs, Litres
    unit_cost_ugx NUMERIC(18,2) NOT NULL DEFAULT 0.00,
    quantity_on_hand INT NOT NULL DEFAULT 0,
    minimum_reorder_level INT NOT NULL DEFAULT 5,
    warehouse_bin_location VARCHAR(64) NOT NULL, -- e.g. 'Rack B-04', 'Locker 12'
    status VARCHAR(32) DEFAULT 'IN_STOCK', -- IN_STOCK, LOW_STOCK, OUT_OF_STOCK, OBSOLETE
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE goods_received_notes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    grn_number VARCHAR(64) UNIQUE NOT NULL,
    po_reference VARCHAR(64) NOT NULL,
    supplier_name VARCHAR(255) NOT NULL,
    received_by UUID NOT NULL REFERENCES users(id),
    total_received_value_ugx NUMERIC(18,2) NOT NULL,
    items_received JSONB NOT NULL, -- Array of {sku_code, qty, unit_cost, condition}
    delivery_note_url TEXT,
    received_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE store_issuances (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    siv_number VARCHAR(64) UNIQUE NOT NULL,
    requesting_department VARCHAR(64) NOT NULL,
    requested_by UUID NOT NULL REFERENCES users(id),
    approved_by_hod UUID NOT NULL REFERENCES users(id),
    issued_by UUID NOT NULL REFERENCES users(id),
    items_issued JSONB NOT NULL, -- Array of {sku_code, qty_issued, purpose}
    recipient_signature_url TEXT,
    issued_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

---

## 5. Go REST API Endpoints (`backend/internal/handlers/stores.go`)

| Endpoint | Method | Scope | Description |
| :--- | :--- | :--- | :--- |
| `/api/v1/stores/inventory` | GET/POST | `stores_officer` | List, search, filter, and create stock inventory items |
| `/api/v1/stores/inventory/:id` | GET/PUT | `stores_officer` | View bin card history and update stock attributes |
| `/api/v1/stores/grn` | GET/POST | `stores_officer` | Generate new Goods Received Notes and view receiving logs |
| `/api/v1/stores/issuances` | GET/POST | All Staff / Stores | Submit stock requisition and process Store Issue Vouchers (SIV) |
| `/api/v1/stores/reconciliation` | GET/POST | Stores / Accounts | Generate physical stock count reconciliation reports |
