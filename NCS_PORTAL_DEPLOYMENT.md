# NCS Cloud Infrastructure Credentials & Deployment Reference

This document contains production connection details, seeded administrative accounts, database credentials, and deployment tokens for the entire NCS multi-branch cluster on VPS `169.58.210.57`.

---

## 1. VPS Server & SSH Access

| Property | Value |
|---|---|
| **Primary IPv4** | `169.58.210.57` |
| **Server Type** | Cloud VPS 4 (2026) |
| **Location** | Hub Europe |
| **Operating System** | Ubuntu 24.04.4 LTS (Noble Numbat) |
| **SSH Port** | `22` |
| **Root Username** | `root` |
| **Root Password** | `EnJk1eeto3R4E3QPPbh` *(SSH Password login disabled for security)* |
| **Admin User** | `fidi` |
| **Admin Password** | `@Fr1caObuaObali` |
| **Sudo Privileges** | `ALL=(ALL) NOPASSWD:ALL` |
| **SSH Key Used** | `~/.ssh/ncs_online_vps` (`ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOTSg/DLIaUrmAUZm7bqHiaBVTRpU0hgXnQ184t8lss6`) |

### SSH Connection Command
```bash
ssh -i ~/.ssh/ncs_online_vps fidi@169.58.210.57
```

---

## 2. Deployed Applications & Public Endpoints

| System | Branch | Production Path | Public Endpoint | Status |
|---|---|---|---|---|
| **NCS Website** | `website` | `/opt/ncs-website` | [https://ncsweb.atenimedia.com](https://ncsweb.atenimedia.com) | **Live (SSL TLSv1.3)** |
| **NCS Portal** | `portal` | `/opt/ncs-portal` | [https://ncsportal.atenimedia.com](https://ncsportal.atenimedia.com) | **Live (SSL TLSv1.3)** |
| **NCS Intranet** | `intranet` | `/opt/ncs-intranet` | [https://ncsintranet.atenimedia.com](https://ncsintranet.atenimedia.com) | **Live (SSL TLSv1.3)** |
| **NCS Bot** | `ncsbot` | `/opt/ncs-bot` | `http://169.58.210.57:3100` | **Live (Port 3100)** |

---

## 3. Seeded Administrative Accounts

### NCS Website, NCS Portal & NCS Intranet (Unified Super Admin)
The Go-based platforms share the initial seeded superadmin schema across their independent databases:

| Field | Value |
|---|---|
| **Email / Username** | `admin@ncs.go.ug` |
| **Password** | `NCS@Admin2026!` |
| **Assigned Role** | `super_admin` (`role_super_admin`) |
| **User ID** | `usr_super_admin_001` |
| **Permissions** | Unrestricted system-wide access (CMS, NAMIS, Departmental modules) |

### NCS Bot (Chatwoot Service)
| Field | Value |
|---|---|
| **Onboarding URL** | `http://169.58.210.57:3100/installation/onboarding` |
| **Recommended Email** | `admin@ncs.go.ug` |
| **Recommended Password** | `NCS@Admin2026!` |
| **Account Type** | SuperAdmin / Platform Owner |

---

## 4. Git Repository & Deployment Credentials

| Property | Value |
|---|---|
| **Repository URL** | `https://github.com/JoshuaObua/NCSOnline.git` |
| **Git Username** | `JoshuaObua` |
| **GitHub Token (PAT)** | `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` |
| **Authenticated Git Remote (VPS)** | `https://JoshuaObua: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@github.com/JoshuaObua/NCSOnline.git` |

---

## 5. Database Credentials

Each system operates on its own dedicated PostgreSQL instance with persistent data volumes:

### NCS Website Database
- **Host (Inside Docker)**: `postgres:5432`
- **Host (Local VPS port)**: `127.0.0.1:5432`
- **Database Name**: `ncswebsite`
- **Username**: `ncswebsite_user`
- **Password**: Configured in `/opt/ncs-website/.env`
- **Data Volume**: `website_postgres_data`

### NCS Portal Database
- **Host (Inside Docker)**: `postgres:5432`
- **Host (Local VPS port)**: `127.0.0.1:5436`
- **Database Name**: `ncsportal`
- **Username**: `ncsportal_user`
- **Password**: Configured in `/opt/ncs-portal/.env`
- **Data Volume**: `portal_postgres_data`

### NCS Intranet Database
- **Host (Inside Docker)**: `postgres:5432`
- **Host (Local VPS port)**: `127.0.0.1:5437`
- **Database Name**: `ncsintranet`
- **Username**: `ncsintranet_user`
- **Password**: Configured in `/opt/ncs-intranet/.env`
- **Data Volume**: `intranet_postgres_data`

### NCS Bot Database & Redis
- **PostgreSQL Host**: `ncsbot-postgres:5432`
- **Database Name**: `ncsbot`
- **Username**: `ncsbot`
- **Password**: Configured in `/opt/ncs-bot/.env`
- **Data Volume**: `ncsbot_postgres`
- **Redis Host**: `ncsbot-redis:6379`
- **Data Volume**: `ncsbot_redis`
- **Storage Volume**: `ncsbot_storage`

---

## 6. SSL / TLS Certificate Details

- **Certificate Authority**: Let's Encrypt
- **Covered Domains**: `ncsweb.atenimedia.com`, `ncsportal.atenimedia.com`, `ncsintranet.atenimedia.com`
- **Certificate Path**: `/etc/letsencrypt/live/ncsweb.atenimedia.com/fullchain.pem`
- **Private Key Path**: `/etc/letsencrypt/live/ncsweb.atenimedia.com/privkey.pem`
- **Auto-Renewal Cron**: Configured at `03:00` daily
- **Renewal Command**:
  ```bash
  cd /opt/ncs-website && sudo docker compose run --rm --entrypoint certbot certbot renew --quiet && sudo docker compose exec nginx nginx -s reload
  ```

---
*Created: August 20, 2026*

