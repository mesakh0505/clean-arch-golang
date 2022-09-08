BEGIN;
CREATE TABLE IF NOT EXISTS todo(
  id uuid NOT NULL PRIMARY KEY,
  text varchar(255) NOT NULL,
  status varchar(50) NOT NULL,
  created_at timestamptz NOT NULL
);
COMMIT;
