-- name: CreateUser
INSERT INTO users (email, password_hash, full_name) VALUES (:email, :password_hash, :full_name);

-- name: GetUserByEmail
SELECT id, email, password_hash, full_name FROM users WHERE email = :email;

-- name: CreateClass
INSERT INTO classes (name, description, capacity) VALUES (:name, :description, :capacity);

-- name: ListClasses
SELECT id, name, description, capacity FROM classes ORDER BY id DESC LIMIT :limit OFFSET :offset;

-- name: SignupClass
INSERT INTO class_signups (user_id, class_id) VALUES (:user_id, :class_id);

-- name: LogActivity
INSERT INTO activities (user_id, type, duration_minutes, notes) VALUES (:user_id, :type, :duration_minutes, :notes);

-- name: ListActivitiesByUser
SELECT id, user_id, type, duration_minutes, notes, created_at FROM activities WHERE user_id = :user_id ORDER BY created_at DESC LIMIT :limit OFFSET :offset;
