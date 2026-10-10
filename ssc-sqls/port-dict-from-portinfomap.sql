-- Migrate PortInfoMap into 港口字典 (base_dicts)
ALTER TABLE base_dicts
  ADD COLUMN IF NOT EXISTS port_name VARCHAR(100) NULL COMMENT '口岸名称' AFTER `type`,
  ADD COLUMN IF NOT EXISTS extra_pay VARCHAR(100) NULL COMMENT '补差项' AFTER port_name;

-- Update existing active rows by key+type
UPDATE base_dicts SET value='上海', port_name='上海口岸', extra_pay='洋山补差', updated_at=NOW(3)
 WHERE `type`='港口字典' AND `key`='上海洋山' AND deleted_at IS NULL;
UPDATE base_dicts SET value='上海', port_name='上海口岸', extra_pay='', updated_at=NOW(3)
 WHERE `type`='港口字典' AND `key`='上海外高桥' AND deleted_at IS NULL;
UPDATE base_dicts SET value='天津', port_name='天津口岸', extra_pay='补差', updated_at=NOW(3)
 WHERE `type`='港口字典' AND `key`='天津东疆' AND deleted_at IS NULL;

-- Insert missing keys
INSERT INTO base_dicts (`key`, `value`, `type`, port_name, extra_pay, created_at, updated_at)
SELECT * FROM (
  SELECT '上海洋山保税' AS `key`, '上海' AS value, '港口字典' AS type, '上海口岸' AS port_name, '洋山补差' AS extra_pay, NOW(3) AS created_at, NOW(3) AS updated_at
  UNION ALL SELECT '天津新港', '天津', '港口字典', '天津口岸', '', NOW(3), NOW(3)
  UNION ALL SELECT '广州黄埔', '广州', '港口字典', '广州口岸', '补差', NOW(3), NOW(3)
  UNION ALL SELECT '重庆果园', '重庆', '港口字典', '重庆口岸', '', NOW(3), NOW(3)
  UNION ALL SELECT '重庆果园港', '重庆', '港口字典', '重庆口岸', '', NOW(3), NOW(3)
) AS t
WHERE NOT EXISTS (
  SELECT 1 FROM base_dicts d
  WHERE d.`type`='港口字典' AND d.`key`=t.`key` AND d.deleted_at IS NULL
);
