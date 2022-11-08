INSERT INTO
  organizations (name, parent_id)
VALUES
  ('管理系统', 0);

INSERT INTO
  permissions (name, description, enable)
VALUES
  ('dashboard', '概览', 1),
  ('setting', '系统设置', 1),
  ('permission', '权限管理', 1),
  ('role', '角色管理', 1),
  ('user', '用户管理', 1),
  ('log', '日志管理', 1);

INSERT INTO
  roles (name, remark, default)
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
  (2, 1),
  (2, 2),
  (3, 1),
  (3, 6);

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
  (1, 1)
