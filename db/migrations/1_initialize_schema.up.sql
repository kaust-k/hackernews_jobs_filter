DROP TABLE IF EXISTS hn_jobs CASCADE;

CREATE TABLE hn_jobs (
  _id integer PRIMARY KEY,
  _parent integer ,
  _by varchar(32),
  _text varchar,
  _time timestamp,
  _type varchar(12)
);

