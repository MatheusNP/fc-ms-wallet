CREATE TABLE clients (
  id varchar(255),
  name varchar(255),
  email varchar(255),
  created_at timestamp default current_timestamp,
  updated_at timestamp default current_timestamp on update current_timestamp
);

INSERT INTO clients (id, name, email, created_at) 
VALUES ('8d34617d-81ac-4e5e-ada3-777c33087c46', 'John', 'john@doe.com', now());

INSERT INTO clients (id, name, email, created_at) 
VALUES ('ca3dea49-59df-4c0c-abcb-d2106300c3bf', 'Jane', 'jane@doe.com', now());
