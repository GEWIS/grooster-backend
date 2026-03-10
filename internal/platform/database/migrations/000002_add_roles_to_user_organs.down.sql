ALTER TABLE user_organs
    ADD COLUMN role VARCHAR(20) NOT NULL DEFAULT 'member',
    ADD CONSTRAINT chk_user_organs_role
        CHECK (role IN ('owner', 'admin', 'member'));