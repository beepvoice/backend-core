
CREATE TABLE IF NOT EXISTS "user" (
	id BYTEA PRIMARY KEY,
	first_name VARCHAR(65535),
	last_name VARCHAR(65535),
	phone_number VARCHAR(32) UNIQUE
);

CREATE TABLE IF NOT EXISTS "conversation" (
	id BYTEA PRIMARY KEY,
	title VARCHAR(65535)
);

CREATE TABLE IF NOT EXISTS member (
	"user" BYTEA REFERENCES "user"(id),
	"conversation" BYTEA REFERENCES "conversation"(id),
	UNIQUE ("user", "conversation")
);

CREATE TABLE IF NOT EXISTS contact (
	"user" BYTEA REFERENCES "user"(id),
	contact BYTEA REFERENCES "user"(id),
	UNIQUE ("user", contact)
);
