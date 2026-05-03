-- Users
CREATE TABLE users (
  id INTEGER PRIMARY KEY AUTOINCREMENT, -- @data_type: uuid
  email TEXT NOT NULL UNIQUE, -- @data_type: ugc-uid
  password_hash TEXT NOT NULL, -- @data_type: ugc-nid
  full_name TEXT -- @data_type: ugc-uid
);

-- Classes
CREATE TABLE classes (
  id INTEGER PRIMARY KEY AUTOINCREMENT, -- @data_type: uuid
  name TEXT NOT NULL, -- @data_type: ugc-nid
  description TEXT, -- @data_type: ugc-nid
  capacity INTEGER DEFAULT 20 -- @data_type: ugc-nid
);

-- Class signups
CREATE TABLE class_signups (
  id INTEGER PRIMARY KEY AUTOINCREMENT, -- @data_type: uuid
  user_id INTEGER NOT NULL, -- @data_type: uuid
  class_id INTEGER NOT NULL, -- @data_type: uuid
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP, -- @data_type: system
  UNIQUE(user_id, class_id)
);

-- Activities
CREATE TABLE activities (
  id INTEGER PRIMARY KEY AUTOINCREMENT, -- @data_type: uuid
  user_id INTEGER NOT NULL, -- @data_type: uuid
  type TEXT NOT NULL, -- @data_type: ugc-nid
  duration_minutes INTEGER, -- @data_type: ugc-nid
  notes TEXT, -- @data_type: ugc-nid
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP -- @data_type: system
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
