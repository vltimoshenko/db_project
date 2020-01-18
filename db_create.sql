-- DROP TABLE votes;
-- DROP TABLE posts;
-- DROP TABLE threads;
-- DROP TABLE forums;
-- DROP TABLE persons;


CREATE TABLE persons(
    id SERIAL PRIMARY KEY,
    about text,
    email text NOT NULL UNIQUE,
    fullname text NOT NULL,
    nickname text NOT NULL UNIQUE
);

CREATE TABLE forums(
    id serial,
    posts integer DEFAULT 0 NOT NULL,
    slug text PRIMARY KEY,
    threads integer DEFAULT 0 NOT NULL,
    title text NOT NULL,
    person text NOT NULL REFERENCES persons(nickname) ON DELETE CASCADE NOT NULL
);

CREATE TABLE threads(
    id SERIAL PRIMARY KEY,
    author text NOT NULL REFERENCES persons(nickname) ON DELETE CASCADE NOT NULL,
    created timestamptz DEFAULT now(),
    forum text REFERENCES forums(slug) ON DELETE CASCADE NOT NULL,
    message text NOT NULL,
    slug text UNIQUE,
    title text NOT NULL,
    votes integer DEFAULT 0 NOT NULL
);

CREATE TABLE posts(
    id SERIAL PRIMARY KEY, --maybe need to switch serial to int
    author text NOT NULL REFERENCES persons(nickname) ON DELETE CASCADE NOT NULL,
    created timestamp with time zone DEFAULT '1970-01-01 03:00:00+03'::timestamp with time zone NOT NULL,
    forum text REFERENCES forums(slug) ON DELETE CASCADE NOT NULL,
    is_edited boolean DEFAULT false NOT NULL,
    message text NOT NULL,
    parent integer DEFAULT 0 NOT NULL, --reference on parent id
    thread integer REFERENCES threads(id) ON DELETE CASCADE NOT NULL
);

CREATE TABLE votes(
    nickname text NOT NULL,
    voice smallint NOT NULL,
    thread integer REFERENCES threads(id) ON DELETE CASCADE NOT NULL,
    PRIMARY KEY(nickname, thread)
);

CREATE OR REPLACE FUNCTION update_forum_after_post() RETURNS TRIGGER AS 
$BODY$
    BEGIN
    UPDATE forums
    SET posts = posts + 1
    WHERE slug = NEW.forum;
    RETURN NULL;
    END;
$BODY$ 
LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS insert_post on posts;

CREATE TRIGGER insert_post
  AFTER INSERT
  ON posts
  FOR EACH ROW EXECUTE PROCEDURE update_forum_after_post();


CREATE OR REPLACE FUNCTION update_forum_after_thread() RETURNS TRIGGER AS 
$BODY$
    BEGIN
    UPDATE forums
    SET threads = threads + 1
    WHERE slug = NEW.forum;
    RETURN NULL;
    END;
$BODY$ 
LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS insert_thread on threads;

CREATE TRIGGER insert_thread
  AFTER INSERT
  ON threads
  FOR EACH ROW EXECUTE PROCEDURE update_forum_after_thread();
