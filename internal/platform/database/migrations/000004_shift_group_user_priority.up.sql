CREATE TABLE `shift_group_priorities` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `created_at` datetime(3) DEFAULT NULL,
    `updated_at` datetime(3) DEFAULT NULL,
    `shift_group_id` BIGINT UNSIGNED NOT NULL,
    `user_id` BIGINT UNSIGNED NOT NULL,
    `priority` TINYINT DEFAULT 1 COMMENT '0: Low, 1: Default, 2: High',
    UNIQUE INDEX `idx_shift_group` (`shift_group_id`, `user_id`),

    CONSTRAINT `fk_shift_group_priorities_shift_group`
        FOREIGN KEY (`shift_group_id`)
            REFERENCES `shift_groups`(`id`)
            ON DELETE CASCADE,

    CONSTRAINT `fk_shift_group_priorities_user`
        FOREIGN KEY (`user_id`)
            REFERENCES `users`(`id`)
            ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO shift_group_priorities (user_id, shift_group_id, priority, created_at, updated_at)
SELECT
    u.id as user_id,
    sg.id as shift_group_id,
    2 AS priority,
    NOW(),
    NOW()
FROM users u
CROSS JOIN shift_groups sg
WHERE NOT EXISTS (
    SELECT 2 FROM shift_group_priorities sgp
    WHERE sgp.user_id = u.id AND sgp.shift_group_id = sg.id
);
