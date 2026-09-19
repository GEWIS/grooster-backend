CREATE TABLE `roster_responsibles` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `created_at` datetime(3) DEFAULT NULL,
    `updated_at` datetime(3) DEFAULT NULL,
    `roster_id` BIGINT UNSIGNED NOT NULL,
    `user_id` BIGINT UNSIGNED NOT NULL,

    UNIQUE KEY `idx_roster_responsible` (`roster_id`, `user_id`),

    CONSTRAINT `fk_roster_responsibles_roster`
        FOREIGN KEY (`roster_id`)
            REFERENCES `rosters`(`id`)
            ON DELETE CASCADE,

    CONSTRAINT `fk_roster_responsibles_user`
        FOREIGN KEY (`user_id`)
            REFERENCES `users`(`id`)
            ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
