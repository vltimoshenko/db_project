DROP TABLE IF EXISTS votes;
DROP TABLE IF EXISTS posts;
DROP TABLE IF EXISTS threads;
DROP TABLE IF EXISTS forums;
DROP TABLE IF EXISTS persons;
CREATE EXTENSION IF NOT EXISTS citext;

CREATE UNLOGGED TABLE persons(
    nickname CITEXT COLLATE "POSIX" PRIMARY KEY, 
    about TEXT NOT NULL DEFAULT '',
    email CITEXT NOT NULL UNIQUE,
    fullname TEXT NOT NULL DEFAULT ''
);

CREATE UNLOGGED TABLE forums(
    posts integer DEFAULT 0 NOT NULL,
    slug citext PRIMARY KEY,
    threads integer DEFAULT 0 NOT NULL,
    title text DEFAULT '' NOT NULL,
    person CITEXT REFERENCES persons (nickname) ON DELETE RESTRICT ON UPDATE RESTRICT NOT NULL
);

CREATE UNLOGGED TABLE threads(
    id SERIAL PRIMARY KEY,
    author CITEXT REFERENCES persons (nickname) ON DELETE CASCADE ON UPDATE CASCADE  NOT NULL,
    created timestamptz DEFAULT now(),
    forum CITEXT NOT NULL,
    message text NOT NULL,
    slug CITEXT UNIQUE,
    title text NOT NULL,
    votes integer DEFAULT 0 NOT NULL
);

CREATE OR REPLACE FUNCTION get_thread_by_post(post_ BIGINT) RETURNS INTEGER AS 
$BODY$
    BEGIN
        RETURN (SELECT thread FROM posts WHERE id = post_);
    END;
$BODY$ 
LANGUAGE plpgsql;

CREATE UNLOGGED TABLE posts(
    id SERIAL PRIMARY KEY,
    author CITEXT REFERENCES persons (nickname) NOT NULL,
    created timestamp with time zone DEFAULT '1970-01-01 03:00:00+03'::timestamp with time zone NOT NULL,
    forum CITEXT,
    is_edited boolean DEFAULT false NOT NULL,
    message text NOT NULL,
    parent BIGINT REFERENCES posts(id) ON DELETE CASCADE ON UPDATE RESTRICT
        CONSTRAINT post_parent_constraint CHECK (get_thread_by_post(parent) = thread),
    thread integer,
    path INTEGER[] not null
);

CREATE OR REPLACE FUNCTION change_path() RETURNS TRIGGER AS
$BODY$
    BEGIN
        NEW.path = (SELECT path FROM posts WHERE id = NEW.parent) || NEW.id;
        RETURN NEW;
    END;
$BODY$ 
LANGUAGE plpgsql;


CREATE TRIGGER change_path BEFORE INSERT ON posts
    FOR EACH ROW EXECUTE PROCEDURE change_path();

CREATE UNLOGGED TABLE votes(
    nickname CITEXT REFERENCES persons (nickname) ON DELETE CASCADE ON UPDATE CASCADE NOT NULL,
    voice smallint NOT NULL,
    thread integer NOT NULL,
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


CREATE OR REPLACE FUNCTION update_thread_votes_counter() RETURNS TRIGGER AS 
$BODY$
    BEGIN
        IF TG_OP = 'INSERT' THEN
            UPDATE threads SET votes = votes + NEW.voice WHERE id = NEW.thread;
            RETURN NEW;
        ELSIF TG_OP = 'UPDATE' THEN
            UPDATE threads SET votes = votes + NEW.voice - OLD.voice WHERE id = NEW.thread;
            RETURN NEW;
        ELSE
            RAISE EXCEPTION 'Invalid call update_thread_votes_counter()';
        END IF;
    END
$BODY$
LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS update_thread_vote ON votes;

CREATE TRIGGER update_thread_vote AFTER INSERT OR UPDATE ON votes
    FOR EACH ROW EXECUTE PROCEDURE update_thread_votes_counter();

-- CREATE UNLOGGED TABLE forum_users (
--     forum CITEXT NOT NULL,
--     nickname CITEXT NOT NULL,
--     about TEXT NOT NULL DEFAULT '',
--     email CITEXT NOT NULL UNIQUE,
--     fullname TEXT NOT NULL DEFAULT '',
--     PRIMARY KEY(forum, nickname)
-- );

CREATE UNLOGGED TABLE forum_users (
    forum CITEXT NOT NULL,
    nickname CITEXT NOT NULL
);

-- CREATE OR REPLACE FUNCTION add_user_to_forum() RETURNS TRIGGER AS
-- $BODY$
--     BEGIN
--         INSERT INTO forum_users (forum, nickname, email, fullname, about)
--         SELECT NEW.forum, nickname, email, fullname, about
--         FROM persons
--         WHERE nickname = NEW.author
--         ON CONFLICT DO NOTHING;
--         RETURN NULL;
--     END;
-- $BODY$
-- LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION  add_user_to_forum() RETURNS trigger AS
$BODY$
BEGIN
    INSERT INTO forum_users(forum, nickname) VALUES (NEW.forum, NEW.author) ON conflict do nothing;
    RETURN NEW;
END;
$BODY$
LANGUAGE plpgsql;

CREATE TRIGGER forum_user_after_thread AFTER INSERT ON threads
    FOR EACH ROW EXECUTE PROCEDURE add_user_to_forum();

-- CREATE TRIGGER forum_user_after_post AFTER INSERT ON posts
--     FOR EACH ROW EXECUTE PROCEDURE add_user_to_forum();
    

CREATE UNIQUE INDEX idx_persons_nickname ON persons(nickname); --+
CREATE INDEX IF NOT EXISTS idx_persons_email ON persons(email); --+

CREATE UNIQUE INDEX idx_forums_slug ON forums(slug); --+

CREATE INDEX IF NOT EXISTS idx_posts_thread ON posts(thread); -- -
CREATE INDEX IF NOT EXISTS idx_posts_id_thread ON posts(id, thread); -- +
CREATE INDEX IF NOT EXISTS idx_posts_forum_author ON posts(forum, author); -- +
CREATE INDEX IF NOT EXISTS idx_posts_author ON posts(author); -- -
CREATE INDEX IF NOT EXISTS idx_post_path_first ON posts((path[1])); -- +
CREATE INDEX IF NOT EXISTS idx_post_parent_thread_path_id ON posts(thread, (path[1]), id) WHERE parent IS NUll; -- +

CREATE INDEX IF NOT EXISTS idx_threads_id ON threads(id); --primary key -- +
CREATE UNIQUE INDEX idx_threads_slug ON threads(slug) INCLUDE (id); -- +
CREATE INDEX IF NOT EXISTS idx_threads_author ON threads(author);  -- +
CREATE INDEX IF NOT EXISTS idx_threads_forum_created ON threads(forum, created); --+

CREATE INDEX IF NOT EXISTS idx_votes_coverage ON votes(thread, nickname) INCLUDE (voice); --+

CREATE UNIQUE INDEX IF NOT EXISTS idx_forum_users ON forum_users(forum, nickname);
