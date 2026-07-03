Here is a comprehensive, production-grade system prompt designed to instruct an AI or developer to rebuild the lost operational modules. It includes the server **Command Center**, a robust database-backed **Maintenance Mode & DB Tooling** system, and a git-integrated **Smart Updates** pipeline.

---

## System Prompt: DevOps Command Center, Maintenance Mode, & Automated Git Deployer

### Objective

Re-architect and implement three missing mission-critical administrative modules within the system backend and sidebar wrapper: the server **Command Center**, the dynamic **Maintenance Mode & Database Management Engine**, and the GitHub-integrated **Smart Updates Pipeline**. These modules must feature tight hardware/Docker abstraction, robust fail-safes, and granular permission checks.

---

### 1. Module A: System Command Center (Server & Docker Orchestration)

Build a real-time system monitoring and container management interface.

* **Metrics Engine (Backend):** Fetch hardware statistics natively via server APIs or system utilities.
* *Metrics:* CPU Core Utilization %, Memory Consumption (Used/Total), Disk I/O, and Network Throughput.


* **Docker Container Orchestration:**
* Expose a secure bridge to the host machine's Docker daemon socket (or equivalent container architecture).
* Provide safe admin actions to **Start**, **Stop**, and **Restart** specific application container services directly from the web interface.


* **Log Streaming Viewer:**
* Implement an asynchronous backend stream (e.g., WebSockets or Server-Sent Events) to tail container logs (`docker logs --tail 100 -f`).
* **UI Requirement:** A scrollable, dark-mode terminal window with auto-scroll toggles, copy-to-clipboard functionality, and keyword log filtering (e.g., `ERROR`, `WARN`).



---

### 2. Module B: Maintenance Mode & Database Tooling Engine

Create a bulletproof system lock and data redundancy toolkit.

#### Maintenance Control & Routing Middleware

* **Global Toggle:** An instant toggle switch to enter/exit Maintenance Mode.
* **Bypass Whitelists:** Implement fields to whitelist specific IP addresses (e.g., developer IPs) and logged-in Admin User IDs so they can still browse and test public pages during downtime.
* **User Experience:** When active, all non-whitelisted traffic must be intercepted at the global middleware layer and served a lightweight, responsive static `503 Service Unavailable` page.

#### Database Operations Matrix

Provide a specialized sub-panel with full transactional execution for database management:

* **Schema & Data Backup:** Trigger a background database dump process (e.g., native SQL dumps or binary backups). Save backups to a protected storage directory with automated timestamped names (`backup_YYYYMMDD_HHMMSS.sql`).
* **Restore & Import:** Allow administrators to upload an existing `.sql` or zipped dump file, or select a previously saved local backup to restore.
* *Safety Catch:* The UI must display an explicit text-confirmation modal requiring the admin to type `CONFIRM RESTORE` before overwriting the active live database.


* **Verification Tool:** Execute an optimization/integrity check command against the active relational tables (e.g., checking indexes, resolving fragmentation) and display a green/red health status log.
* **Delete/Purge:** Enable explicit purging of specific stale database snapshots to free up disk space.

---

### 3. Module C: Smart Updates & GitHub Deployer

Implement a secure, token-driven continuous deployment pipeline utilizing git.

* **Secure Fine-Grained Access Token Storage:**
* Store a GitHub Personal Access Token (PAT) securely using standard environmental encryption practices at rest. Never expose this token in plain text anywhere in the frontend UI.


* **Automated Remote Update Checker:**
* Run a background routine or cron check comparing the current local commit hash against the production target branch on GitHub via the GitHub API.
* **Notification System:** If a mismatch is found (new commits available), dispatch an immediate admin notification alert across the dashboard UI indicating an update is ready.


* **One-Click Pull & Deploy Execution Engine:**
* When the admin clicks "Deploy Update", trigger a controlled deployment sequence:
1. **Auto-Backup:** Spin up an automated, pre-update database dump.
2. **Git Execution:** Execute a clean remote fetch and pull (`git fetch && git pull origin <branch>`) leveraging the stored secure access token.
3. **Dependency & Migration Resolution:** Programmatically trigger application dependency syncs (e.g., `go mod tidy`, `npm install`, or `pip install -r`) followed by any outstanding database schema migrations.
4. **Graceful Reload:** Restart the application service worker to apply changes smoothly.
5. **Audit Logs:** Write the complete shell execution output (success or failure trace logs) directly into an admin update deployment ledger for tracing.





---

### 4. Security & Architecture Constraints

* **Execution Isolation:** Ensure all shell processes run under highly restricted permission limits to prevent arbitrary code execution vulnerabilities.
* **State Protection:** If an update pull or migration fails, the engine must halt gracefully, log the exact stack trace, and alert the admin without dropping the database state.
* **Performance:** Terminal logs and metric polling must use efficient streaming or debounced interval fetching to prevent client-side or server memory exhaustion.