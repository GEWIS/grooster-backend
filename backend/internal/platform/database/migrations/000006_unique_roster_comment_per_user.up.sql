ALTER TABLE `roster_comments`
    ADD UNIQUE KEY `idx_roster_comment_user` (`roster_id`, `user_id`);
