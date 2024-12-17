CREATE TABLE accounts (
  id varchar(255),
  client_id varchar(255),
  balance float,
  created_at timestamp default current_timestamp,
  updated_at timestamp default current_timestamp on update current_timestamp
);

INSERT INTO accounts (id, client_id, balance, created_at)
VALUES ('912ca3ff-59d4-4640-acd8-76eb1faf0dbe', '8d34617d-81ac-4e5e-ada3-777c33087c46', 200, now());

INSERT INTO accounts (id, client_id, balance, created_at)
VALUES ('31fb673a-012a-4cfc-8dc9-7a9ad0219437', 'ca3dea49-59df-4c0c-abcb-d2106300c3bf', 200, now());
