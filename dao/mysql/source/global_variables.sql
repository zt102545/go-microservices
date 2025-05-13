CREATE TABLE `global_variables`
(
    `id`          int(10) unsigned                        NOT NULL AUTO_INCREMENT,
    `name`        text COLLATE utf8mb4_unicode_ci         NOT NULL COMMENT '配置名称',
    `integer`     int(11)                                 NOT NULL,
    `string`      varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
    `text`        text COLLATE utf8mb4_unicode_ci         NOT NULL,
    `float`       double(8, 2)                            NOT NULL,
    `data`        text COLLATE utf8mb4_unicode_ci         NOT NULL COMMENT '配置数据',
    `description` text COLLATE utf8mb4_unicode_ci         NOT NULL COMMENT '描述',
    `created_at`  timestamp                               NULL DEFAULT NULL,
    `updated_at`  timestamp                               NULL DEFAULT NULL,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  AUTO_INCREMENT = 22
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci COMMENT ='全局配置表';