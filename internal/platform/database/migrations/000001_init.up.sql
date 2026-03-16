
/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Database: `grooster_test`
--

-- --------------------------------------------------------

--
-- Table structure for table `organs`
--

CREATE TABLE IF NOT EXISTS `organs` (
                                        `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
                                        `created_at` datetime(3) DEFAULT NULL,
                                        `updated_at` datetime(3) DEFAULT NULL,
                                        `name` varchar(255) DEFAULT NULL,
                                        PRIMARY KEY (`id`),
                                        UNIQUE KEY `idx_organs_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Table structure for table `rosters`
--

CREATE TABLE IF NOT EXISTS `rosters` (
                                         `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
                                         `created_at` datetime(3) DEFAULT NULL,
                                         `updated_at` datetime(3) DEFAULT NULL,
                                         `name` varchar(255) DEFAULT NULL,
                                         `values` longtext DEFAULT NULL,
                                         `organ_id` bigint(20) UNSIGNED DEFAULT NULL,
                                         `date` datetime(3) DEFAULT NULL,
                                         `saved` tinyint(1) DEFAULT 0,
                                         `template_id` bigint(20) UNSIGNED DEFAULT NULL,
                                         PRIMARY KEY (`id`),
                                         KEY `fk_rosters_organ` (`organ_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Table structure for table `roster_answers`
--

CREATE TABLE IF NOT EXISTS `roster_answers` (
                                                `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
                                                `created_at` datetime(3) DEFAULT NULL,
                                                `updated_at` datetime(3) DEFAULT NULL,
                                                `user_id` bigint(20) UNSIGNED DEFAULT NULL,
                                                `roster_id` bigint(20) UNSIGNED DEFAULT NULL,
                                                `roster_shift_id` bigint(20) UNSIGNED DEFAULT NULL,
                                                `value` longtext DEFAULT NULL,
                                                PRIMARY KEY (`id`),
                                                UNIQUE KEY `user_answer_idx` (`user_id`,`roster_id`,`roster_shift_id`),
                                                KEY `fk_roster_answers_roster_shift` (`roster_shift_id`),
                                                KEY `fk_rosters_roster_answer` (`roster_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Table structure for table `roster_shifts`
--

CREATE TABLE IF NOT EXISTS `roster_shifts` (
                                               `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
                                               `created_at` datetime(3) DEFAULT NULL,
                                               `updated_at` datetime(3) DEFAULT NULL,
                                               `name` varchar(255) DEFAULT NULL,
                                               `roster_id` bigint(20) UNSIGNED DEFAULT NULL,
                                               `order` bigint(20) UNSIGNED DEFAULT NULL,
                                               `shift_group_id` bigint(20) UNSIGNED DEFAULT NULL,
                                               PRIMARY KEY (`id`),
                                               KEY `fk_roster_shifts_shift_group` (`shift_group_id`),
                                               KEY `fk_rosters_roster_shift` (`roster_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Table structure for table `roster_templates`
--

CREATE TABLE IF NOT EXISTS `roster_templates` (
                                                  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
                                                  `created_at` datetime(3) DEFAULT NULL,
                                                  `updated_at` datetime(3) DEFAULT NULL,
                                                  `organ_id` bigint(20) UNSIGNED DEFAULT NULL,
                                                  `name` varchar(255) DEFAULT NULL,
                                                  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Table structure for table `roster_template_shifts`
--

CREATE TABLE IF NOT EXISTS `roster_template_shifts` (
                                                        `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
                                                        `created_at` datetime(3) DEFAULT NULL,
                                                        `updated_at` datetime(3) DEFAULT NULL,
                                                        `template_id` bigint(20) UNSIGNED DEFAULT NULL,
                                                        `shift_name` longtext DEFAULT NULL,
                                                        `shift_group_id` bigint(20) UNSIGNED DEFAULT NULL,
                                                        PRIMARY KEY (`id`),
                                                        KEY `fk_roster_templates_shifts` (`template_id`),
                                                        KEY `fk_roster_template_shifts_shift_group` (`shift_group_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Table structure for table `roster_template_shift_preferences`
--

CREATE TABLE IF NOT EXISTS `roster_template_shift_preferences` (
                                                                   `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
                                                                   `created_at` datetime(3) DEFAULT NULL,
                                                                   `updated_at` datetime(3) DEFAULT NULL,
                                                                   `roster_template_shift_id` bigint(20) UNSIGNED DEFAULT NULL,
                                                                   `user_id` bigint(20) UNSIGNED DEFAULT NULL,
                                                                   `preference` longtext DEFAULT NULL,
                                                                   PRIMARY KEY (`id`),
                                                                   KEY `fk_roster_template_shift_preferences_roster_template_shift` (`roster_template_shift_id`),
                                                                   KEY `fk_roster_template_shift_preferences_user` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Table structure for table `saved_shifts`
--

CREATE TABLE IF NOT EXISTS `saved_shifts` (
                                              `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
                                              `created_at` datetime(3) DEFAULT NULL,
                                              `updated_at` datetime(3) DEFAULT NULL,
                                              `roster_id` bigint(20) UNSIGNED DEFAULT NULL,
                                              `roster_shift_id` bigint(20) UNSIGNED DEFAULT NULL,
                                              PRIMARY KEY (`id`),
                                              KEY `fk_saved_shifts_roster_shift` (`roster_shift_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Table structure for table `shift_groups`
--

CREATE TABLE IF NOT EXISTS `shift_groups` (
                                              `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
                                              `created_at` datetime(3) DEFAULT NULL,
                                              `updated_at` datetime(3) DEFAULT NULL,
                                              `organ_id` bigint(20) UNSIGNED DEFAULT NULL,
                                              `name` varchar(255) DEFAULT NULL,
                                              PRIMARY KEY (`id`),
                                              UNIQUE KEY `organ_shift_group` (`organ_id`,`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Table structure for table `users`
--

CREATE TABLE IF NOT EXISTS `users` (
                                       `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
                                       `created_at` datetime(3) DEFAULT NULL,
                                       `updated_at` datetime(3) DEFAULT NULL,
                                       `name` varchar(255) DEFAULT NULL,
                                       `gewis_id` bigint(20) UNSIGNED DEFAULT NULL,
                                       PRIMARY KEY (`id`),
                                       UNIQUE KEY `idx_name` (`gewis_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Table structure for table `user_organs`
--

CREATE TABLE IF NOT EXISTS `user_organs` (
                                             `organ_id` bigint(20) UNSIGNED NOT NULL,
                                             `user_id` bigint(20) UNSIGNED NOT NULL,
                                             `username` varchar(25) DEFAULT NULL,
                                             PRIMARY KEY (`organ_id`,`user_id`),
                                             KEY `fk_user_organs_user` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Table structure for table `user_shift_saved`
--

CREATE TABLE IF NOT EXISTS `user_shift_saved` (
                                                  `saved_shift_id` bigint(20) UNSIGNED NOT NULL,
                                                  `user_id` bigint(20) UNSIGNED NOT NULL,
                                                  PRIMARY KEY (`saved_shift_id`,`user_id`),
                                                  KEY `fk_user_shift_saved_user` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Constraints for dumped tables
--

--
-- Constraints for table `rosters`
--
ALTER TABLE `rosters`
    ADD CONSTRAINT `fk_rosters_organ` FOREIGN KEY (`organ_id`) REFERENCES `organs` (`id`);

--
-- Constraints for table `roster_answers`
--
ALTER TABLE `roster_answers`
    ADD CONSTRAINT `fk_roster_answers_roster_shift` FOREIGN KEY (`roster_shift_id`) REFERENCES `roster_shifts` (`id`) ON DELETE CASCADE,
    ADD CONSTRAINT `fk_rosters_roster_answer` FOREIGN KEY (`roster_id`) REFERENCES `rosters` (`id`) ON DELETE CASCADE;

--
-- Constraints for table `roster_shifts`
--
ALTER TABLE `roster_shifts`
    ADD CONSTRAINT `fk_roster_shifts_shift_group` FOREIGN KEY (`shift_group_id`) REFERENCES `shift_groups` (`id`) ON DELETE SET NULL,
    ADD CONSTRAINT `fk_rosters_roster_shift` FOREIGN KEY (`roster_id`) REFERENCES `rosters` (`id`) ON DELETE CASCADE;

--
-- Constraints for table `roster_template_shifts`
--
ALTER TABLE `roster_template_shifts`
    ADD CONSTRAINT `fk_roster_template_shifts_shift_group` FOREIGN KEY (`shift_group_id`) REFERENCES `shift_groups` (`id`) ON DELETE SET NULL,
    ADD CONSTRAINT `fk_roster_templates_shifts` FOREIGN KEY (`template_id`) REFERENCES `roster_templates` (`id`) ON DELETE CASCADE;

--
-- Constraints for table `roster_template_shift_preferences`
--
ALTER TABLE `roster_template_shift_preferences`
    ADD CONSTRAINT `fk_roster_template_shift_preferences_roster_template_shift` FOREIGN KEY (`roster_template_shift_id`) REFERENCES `roster_template_shifts` (`id`) ON DELETE CASCADE,
    ADD CONSTRAINT `fk_roster_template_shift_preferences_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE;

--
-- Constraints for table `saved_shifts`
--
ALTER TABLE `saved_shifts`
    ADD CONSTRAINT `fk_saved_shifts_roster_shift` FOREIGN KEY (`roster_shift_id`) REFERENCES `roster_shifts` (`id`) ON DELETE CASCADE;

--
-- Constraints for table `shift_groups`
--
ALTER TABLE `shift_groups`
    ADD CONSTRAINT `fk_shift_groups_organ` FOREIGN KEY (`organ_id`) REFERENCES `organs` (`id`) ON DELETE CASCADE;

--
-- Constraints for table `user_organs`
--
ALTER TABLE `user_organs`
    ADD CONSTRAINT `fk_user_organs_organ` FOREIGN KEY (`organ_id`) REFERENCES `organs` (`id`),
    ADD CONSTRAINT `fk_user_organs_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);

--
-- Constraints for table `user_shift_saved`
--
ALTER TABLE `user_shift_saved`
    ADD CONSTRAINT `fk_user_shift_saved_saved_shift` FOREIGN KEY (`saved_shift_id`) REFERENCES `saved_shifts` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
    ADD CONSTRAINT `fk_user_shift_saved_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;