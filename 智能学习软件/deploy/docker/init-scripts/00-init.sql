-- 智能学习软件 — 数据库初始化脚本
-- 调试环境下，GORM AutoMigrate 会自动创建表结构
-- 如需种子数据，在此文件添加 INSERT 语句
-- 例如：
-- INSERT INTO subjects (id, name, created_at, updated_at) VALUES (1, '数学', NOW(), NOW());
-- INSERT INTO subjects (id, name, created_at, updated_at) VALUES (2, '英语', NOW(), NOW());

-- 说明：本目录下的 .sql 文件会在 PostgreSQL 容器首次启动时按文件名顺序执行。
-- 仅用于调试环境，生产环境请使用版本化迁移工具。