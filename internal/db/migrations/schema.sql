-- Users
CREATE TABLE users (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  email TEXT NOT NULL UNIQUE,
  password_hash TEXT NOT NULL,
  full_name TEXT
);

-- Classes
CREATE TABLE classes (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,
  description TEXT,
  capacity INTEGER DEFAULT 20
);

-- Class signups
CREATE TABLE class_signups (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL,
  class_id INTEGER NOT NULL,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  UNIQUE(user_id, class_id)
);

-- Activities
CREATE TABLE activities (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL,
  type TEXT NOT NULL,
  duration_minutes INTEGER,
  notes TEXT,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Data annotations
CREATE TABLE data_annotations (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  data_field TEXT NOT NULL,
  category TEXT,
  description TEXT,
  policy_type TEXT,
  scope TEXT,
  annotation_id TEXT UNIQUE,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
