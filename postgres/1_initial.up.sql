
CREATE TABLE IF NOT EXISTS "user" (
	id BYTEA PRIMARY KEY,
	username VARCHAR(63555) UNIQUE,
  bio VARCHAR(63535),
  profile_pic VARCHAR(63535),
	first_name VARCHAR(65535),
	last_name VARCHAR(65535),
	phone_number VARCHAR(32) UNIQUE
);

CREATE TABLE IF NOT EXISTS "conversation" (
	id BYTEA PRIMARY KEY,
  dm BOOLEAN,
	title VARCHAR(65535),
  picture VARCHAR(63535)
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

CREATE OR REPLACE FUNCTION notify_permissions_new () RETURNS TRIGGER AS $$
  BEGIN
    PERFORM pg_notify('member_new', CONCAT(NEW."user", '+', NEW."conversation"));
    RETURN NULL;
  END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION notify_permissions_delete () RETURNS TRIGGER AS $$
  BEGIN
    PERFORM pg_notify('member_delete', CONCAT(OLD."user", '+', OLD."conversation"));
    RETURN NULL;
  END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER notify_permissions_new
  AFTER INSERT OR UPDATE
  ON "member"
  FOR EACH ROW
    EXECUTE PROCEDURE notify_permissions_new();

CREATE TRIGGER notify_permissions_delete
  AFTER DELETE
  ON "member"
  FOR EACH ROW
    EXECUTE PROCEDURE notify_permissions_delete();
