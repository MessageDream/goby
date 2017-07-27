UPDATE user SET status = 1 WHERE is_active = 1;
UPDATE user SET status = 0 WHERE is_active = 0;
ALTER TABLE user DROP is_active;