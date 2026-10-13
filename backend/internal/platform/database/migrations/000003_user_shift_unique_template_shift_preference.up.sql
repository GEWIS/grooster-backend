DELETE p1 FROM roster_template_shift_preferences p1
    INNER JOIN roster_template_shift_preferences p2
    WHERE
       p1.id < p2.id AND
       p1.roster_template_shift_id = p2.roster_template_shift_id AND
       p1.user_id = p2.user_id;

ALTER TABLE  `roster_template_shift_preferences`
ADD UNIQUE INDEX `idx_user_shift` (`roster_template_shift_id`, `user_id`);