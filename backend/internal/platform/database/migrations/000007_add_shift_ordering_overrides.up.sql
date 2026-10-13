CREATE TABLE `shift_ordering_overrides` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `created_at` datetime(3) DEFAULT NULL,
    `updated_at` datetime(3) DEFAULT NULL,
    `user_id` BIGINT UNSIGNED NOT NULL,
    `shift_group_id` BIGINT UNSIGNED NOT NULL,
    `set_at` datetime(3) NOT NULL,

    KEY `fk_shift_ordering_overrides_user` (`user_id`),
    KEY `fk_shift_ordering_overrides_shift_group` (`shift_group_id`),

    CONSTRAINT `fk_shift_ordering_overrides_user`
        FOREIGN KEY (`user_id`)
            REFERENCES `users`(`id`)
            ON DELETE CASCADE,

    CONSTRAINT `fk_shift_ordering_overrides_shift_group`
        FOREIGN KEY (`shift_group_id`)
            REFERENCES `shift_groups`(`id`)
            ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
