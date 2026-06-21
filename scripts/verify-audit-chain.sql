WITH ordered AS (
  SELECT chain_sequence, previous_hash, entry_hash,
         lag(entry_hash) OVER (ORDER BY chain_sequence) AS actual_previous
  FROM audit_logs WHERE entry_hash IS NOT NULL
)
SELECT chain_sequence, previous_hash, actual_previous
FROM ordered
WHERE previous_hash <> COALESCE(actual_previous, repeat('0',64));
