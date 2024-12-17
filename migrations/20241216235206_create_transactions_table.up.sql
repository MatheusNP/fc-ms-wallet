CREATE TABLE transactions (
  id varchar(255),
  account_from_id varchar(255),
  account_to_id varchar(255),
  amount float,
  created_at timestamp default current_timestamp,
  updated_at timestamp default current_timestamp on update current_timestamp
);
