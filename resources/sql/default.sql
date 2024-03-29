INSERT INTO
  organizations (name, parent_id)
VALUES
  ('管理系统', 0);

INSERT INTO
  permissions (name, description)
VALUES
  ('dashboard', '概览'),
  ('setting', '系统设置'),
  ('permission', '权限管理'),
  ('role', '角色管理'),
  ('user', '用户管理'),
  ('log', '日志管理'),
  ("audit_log", "审计日志"),
  ("security_log", "安全日志"),
  ("run_log", "运行日志"),
  ("log_config", "日志设置"),
  ("alert", "告警管理"),
  ("alert_list", "告警列表"),
  ("alert_config", "告警配置"),
  ("alert_access", "告警接入"),
  ("alert_push", "告警推送"),
  ("report", "报告管理"),
  ("report_list", "报告列表"),
  ("report_config", "报告配置"),
  ("sched_job", "定时任务管理"),
  ("resource", "资源管理"),
  ("dict", "字典管理"),
  ("system", "系统设置"),
  ("security", "系统安全")
  ;

INSERT INTO
  roles (name, remark, is_default)
VALUES
  ('管理员', '', 1),
  ('操作员', '', 1),
  ('审计员', '', 1);

INSERT INTO
  role_permissions (role_id, permission_id)
VALUES
  (1, 1),
  (1, 2),
  (1, 3),
  (1, 4),
  (1, 5),
  (1, 6),
  (1, 7),
  (1, 8),
  (1, 9),
  (1, 10),
  (1, 11),
  (1, 12),
  (1, 13),
  (1, 14),
  (1, 15),
  (1, 16),
  (1, 17),
  (1, 18),
  (1, 19),
  (1, 20),
  (1, 21),
  (1, 22),
  (1, 23),
  (2, 1),
  (2, 2),
  (2, 3),
  (2, 4),
  (2, 5),
  (2, 6),
  (2, 9),
  (2, 10),
  (2, 11),
  (2, 12),
  (2, 13),
  (2, 14),
  (2, 15),
  (2, 16),
  (2, 17),
  (2, 18),
  (2, 19),
  (2, 20),
  (2, 21),
  (2, 22),
  (2, 23),
  (3, 1),
  (3, 2),
  (3, 6),
  (3, 7),
  (3, 8);

INSERT INTO
  users (
    username,
    nickname,
    password,
    email,
    phone,
    avatar,
    sex,
    age,
    organization_id
  )
VALUES
  (
    'admin',
    '管理员',
    'a76eb12d1a7fe3e1530c83c7eb683afb2695bc3d6500fff826280a3c7ed24a34',
    '',
    '',
    '',
    0,
    18,
    1
  );

INSERT INTO
  user_roles (user_id, role_id)
VALUES
  (1, 1);

INSERT INTO
  system_settings (name, description)
VALUES
  ("后台管理系统", "一款开箱即用的后台管理系统");

INSERT INTO
  log_configs (id, keep_time, archive)
VALUES
  ("1", 30, 1);

INSERT INTO
  dicts (type, description, label, value, value_type)
VALUES
  ("report", "报表", "默认", "default", "string"),
  ("schedJob", "定时任务", "报告", "report", "string"),
  ("schedJob", "定时任务", "日志", "log", "string"),
  ("alertType", "告警类型", "主机", "host", "string"),
  ("alertType", "告警类型", "服务", "service", "string");
