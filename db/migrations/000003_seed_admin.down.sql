BEGIN;
DELETE FROM users WHERE username = 'admin';
COMMIT;