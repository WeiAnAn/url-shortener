CREATE TABLE short_urls (
  short_url varchar(7) PRIMARY KEY NOT NULL,
  original_url varchar(2048) NOT NULL,
  expireAt timestamp NOT NULL 
);