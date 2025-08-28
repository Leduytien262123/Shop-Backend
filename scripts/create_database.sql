-- Tạo database nếu chưa tồn tại
CREATE DATABASE IF NOT EXISTS backend CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- Sử dụng database
USE backend;

-- Tạo user admin mặc định (chạy sau khi ứng dụng đã khởi động)
-- INSERT INTO users (username, email, password, full_name, role, is_active, created_at, updated_at) 
-- VALUES ('admin', 'admin@example.com', '$2a$10$hash_password_here', 'Administrator', 'admin', true, NOW(), NOW());
