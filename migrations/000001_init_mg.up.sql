CREATE TABLE users (
    id         uuid         NOT NULL PRIMARY KEY
	, login    varchar(50)  NOT NULL UNIQUE
	, password varchar(500) NOT NULL 
);

CREATE INDEX login_idx on users(login);

CREATE TABLE text_data (
	id        uuid  NOT NULL PRIMARY KEY
	, user_id uuid  NOT NULL
	, content bytea NOT NULL
);

ALTER TABLE text_data
	ADD FOREIGN KEY (user_id) REFERENCES users(id);

CREATE TABLE binary_data (
	id        uuid  NOT NULL PRIMARY KEY
	, user_id uuid  NOT NULL
	, content bytea NOT NULL
);

ALTER TABLE binary_data
	ADD FOREIGN KEY (user_id) REFERENCES users(id);

CREATE TABLE credentials_data (
	id         uuid  NOT NULL PRIMARY KEY
	, user_id  uuid  NOT NULL
	, name     bytea NOT NULL
	, login    bytea NOT NULL
	, password bytea NOT NULL
	, meta     bytea NOT NULL
);

ALTER TABLE credentials_data
	ADD FOREIGN KEY (user_id) REFERENCES users(id);

CREATE TABLE bank_card_data (
	id            uuid  NOT NULL PRIMARY KEY
	, user_id     uuid  NOT NULL
	, number      bytea NOT NULL
	, valid_thru  bytea NOT NULL
	, cvv         bytea NOT NULL
	, card_holder bytea NOT NULL
	, meta        bytea NOT NULL
);

ALTER TABLE bank_card_data
	ADD FOREIGN KEY (user_id) REFERENCES users(id);
