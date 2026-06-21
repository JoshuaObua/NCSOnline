# Database Restore Drill

Restore drills use an isolated database and never overwrite production. Evidence records backup name, checksum, PostgreSQL version, timing, restore result, critical table counts, errors and cleanup confirmation. A backup is verified only after the restored database accepts connections and counts are compared.
