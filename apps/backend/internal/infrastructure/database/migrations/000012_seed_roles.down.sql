-- 000012_seed_roles.down.sql

DELETE FROM roles
WHERE name IN ('admin','moderator','editor','author','premium','user');