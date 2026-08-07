CREATE UNIQUE INDEX IF NOT EXISTS ux_integrations_store_type
ON integrations (store_id, integration_type_id)
WHERE deleted_at IS NULL
  AND store_id IS NOT NULL
  AND store_id <> '';
