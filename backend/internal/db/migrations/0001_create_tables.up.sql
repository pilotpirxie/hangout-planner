CREATE TABLE IF NOT EXISTS calendars (
  id text PRIMARY KEY,
  title text NOT NULL,
  description text,
  dates jsonb NOT NULL,
  created_at timestamptz DEFAULT now()
);

CREATE TABLE IF NOT EXISTS votes (
  id text PRIMARY KEY,
  calendar_id text REFERENCES calendars(id) ON DELETE CASCADE,
  username text NOT NULL,
  available jsonb NOT NULL,
  created_at timestamptz DEFAULT now()
);

