ALTER TABLE user_organs
ADD COLUMN role VARCHAR(20) not null DEFAULT 'member';

ALTER TABLE user_organs
    ADD "CONSTRAINT" chk_user_organs_role
        CHECK (role IN ('owner', 'admin', 'member'));