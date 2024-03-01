CREATE TABLE users (
    id         uuid         NOT NULL PRIMARY KEY
	, login    varchar(50)  NOT NULL UNIQUE
	, password varchar(500) NOT NULL 
);

CREATE INDEX login_idx on users(login);

CREATE TABLE text (
	id         uuid          NOT NULL PRIMARY KEY
	, user_id  uuid          NOT NULL
	, content  bytea         NOT NULL
);

ALTER TABLE text
	ADD FOREIGN KEY (user_id) REFERENCES users(id);