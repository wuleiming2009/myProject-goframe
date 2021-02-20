CREATE TABLE `account` (
                           `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
                           `user_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '用户雪花ID',
                           `password` varchar(200) NOT NULL DEFAULT '' COMMENT '用户密码',
                           `email` varchar(100) NOT NULL DEFAULT '' COMMENT '注册邮箱',
                           `facebook_id` varchar(100) NOT NULL DEFAULT '' COMMENT 'facebook的三方账号',
                           `gooleplus_id` varchar(100) NOT NULL DEFAULT '' COMMENT 'google+的三方账号',
                           `is_deleted` tinyint(4) NOT NULL DEFAULT '0',
                           `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
                           `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                           PRIMARY KEY (`id`),
                           KEY `user_id` (`user_id`),
                           KEY `email` (`email`),
                           KEY `facebook_id` (`facebook_id`),
                           KEY `gooleplus_id` (`gooleplus_id`)
) ENGINE=InnoDB AUTO_INCREMENT=13 DEFAULT CHARSET=utf8mb4;