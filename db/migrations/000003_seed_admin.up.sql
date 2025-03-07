BEGIN;

WITH inserted_user AS (
  INSERT INTO users (username, email, password_hash)
  VALUES ('admin', 'admin@example.com','$2a$10$nV8wXSf1Z0Qx/Vc99MPCVujGnkRl/dqciMM/.j5tr3xRieJrK6KlG')
  RETURNING id
)
INSERT INTO user_roles (user_id, role_id)
VALUES ((SELECT id FROM inserted_user), 1);

COMMIT;