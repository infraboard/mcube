CREATE TABLE IF NOT EXISTS `books` (
  `id` char(64) COLLATE utf8mb4_general_ci NOT NULL COMMENT '对象Id',
  `name` varchar(255) COLLATE utf8mb4_general_ci NOT NULL COMMENT '书名',
  `author` varchar(255) COLLATE utf8mb4_general_ci NOT NULL COMMENT '作者',
  `create_at` bigint NOT NULL COMMENT '创建时间(13位时间戳)',
  PRIMARY KEY (`id`),
  KEY `idx_name` (`name`) USING BTREE COMMENT '用于书名搜索',
  KEY `idx_author` (`author`) USING BTREE COMMENT '用于作者搜索'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;