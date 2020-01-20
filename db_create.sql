DROP TABLE IF EXISTS votes;
DROP TABLE IF EXISTS posts;
DROP TABLE IF EXISTS threads;
DROP TABLE IF EXISTS forums;
DROP TABLE IF EXISTS persons;
CREATE EXTENSION IF NOT EXISTS citext;

CREATE UNLOGGED TABLE persons(
    about text,
    -- email text NOT NULL UNIQUE,
    -- fullname text NOT NULL,
    -- nickname text NOT NULL UNIQUE
    email CITEXT NOT NULL UNIQUE CONSTRAINT persons_email_right CHECK(email ~ '^.*@[A-Za-z0-9\-_\.]*$'),
    fullname TEXT NOT NULL DEFAULT '',
    nickname CITEXT COLLATE "POSIX" PRIMARY KEY CONSTRAINT persons_nick_right CHECK(nickname ~ '^[A-Za-z0-9_\.]*$')
);

CREATE UNLOGGED TABLE forums(
    id serial,
    posts integer DEFAULT 0 NOT NULL,
    slug text PRIMARY KEY,
    threads integer DEFAULT 0 NOT NULL,
    title text NOT NULL,
    person CITEXT REFERENCES persons (nickname) ON DELETE RESTRICT ON UPDATE RESTRICT NOT NULL
    -- person text NOT NULL REFERENCES persons(nickname) ON DELETE CASCADE NOT NULL
);

CREATE UNLOGGED TABLE threads(
    id SERIAL PRIMARY KEY,
    author CITEXT REFERENCES persons (nickname) ON DELETE CASCADE ON UPDATE CASCADE NOT NULL,
    -- author text NOT NULL REFERENCES persons(nickname) ON DELETE CASCADE NOT NULL,
    created timestamptz DEFAULT now(),
    forum text NOT NULL,
   -- forum text REFERENCES forums(slug) ON DELETE CASCADE NOT NULL,
    message text NOT NULL,
    -- slug text UNIQUE,
    slug CITEXT UNIQUE CONSTRAINT slug_correct CHECK(slug ~ '^(\d|\w|-|_)*(\w|-|_)(\d|\w|-|_)*$'),
    title text NOT NULL,
    votes integer DEFAULT 0 NOT NULL
);

CREATE UNLOGGED TABLE posts(
    id SERIAL PRIMARY KEY, --maybe need to switch serial to int
    -- author text NOT NULL REFERENCES persons(nickname) ON DELETE CASCADE,
    author CITEXT REFERENCES persons (nickname) NOT NULL,

    -- author text NOT NULL,
    created timestamp with time zone DEFAULT '1970-01-01 03:00:00+03'::timestamp with time zone NOT NULL,
    forum text REFERENCES forums(slug) ON DELETE CASCADE NOT NULL,
    is_edited boolean DEFAULT false NOT NULL,
    message text NOT NULL,
    parent integer DEFAULT 0 NOT NULL, --reference on parent id
    thread integer REFERENCES threads(id) ON DELETE CASCADE NOT NULL
    -- thread integer NOT NULL

);

CREATE UNLOGGED TABLE votes(
    -- nickname text NOT NULL REFERENCES persons(nickname) ON DELETE CASCADE,
    nickname CITEXT REFERENCES persons (nickname) ON DELETE CASCADE ON UPDATE CASCADE NOT NULL,

    voice smallint NOT NULL,
    thread integer NOT NULL,
    -- thread integer REFERENCES threads(id) ON DELETE CASCADE NOT NULL,
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
        IF TG_OP='INSERT' THEN
            UPDATE threads SET votes=votes+NEW.voice WHERE id=NEW.thread;
            RETURN NEW;
        ELSIF TG_OP='UPDATE' THEN
            UPDATE threads SET votes=votes+(NEW.voice-OLD.voice) WHERE id=NEW.thread;
            RETURN NEW;
        ELSE
            RAISE EXCEPTION 'Invalid call update_thread_votes_counter()';
        end if;
    END
$BODY$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS update_thread_vote ON votes;

CREATE TRIGGER update_thread_vote AFTER INSERT OR UPDATE ON votes
    FOR EACH ROW EXECUTE PROCEDURE update_thread_votes_counter();


-- CREATE INDEX IF NOT EXISTS idx_persons_email ON persons(email);
-- CREATE INDEX  IF NOT EXISTS idx_threads_slug ON threads(lower(slug));
-- -- create index IF NOT EXISTS thread_forum ON threads(forum);
-- CREATE INDEX IF NOT EXISTS idx_posts_parent ON posts(parent);
-- CREATE INDEX IF NOT EXISTS idx_posts_author ON posts(lower(author));
-- CREATE INDEX IF NOT EXISTS idx_posts_thread ON posts(thread);
-- CREATE INDEX IF NOT EXISTS idx_thread_author ON threads(lower(author));


----------------------------------------------------------------
CREATE UNIQUE INDEX idx_persons_nickname ON persons(lower(nickname));
CREATE UNIQUE INDEX idx_forums_slug ON forums(lower(slug));

CREATE INDEX IF NOT EXISTS idx_posts_thread ON posts(thread);
CREATE INDEX IF NOT EXISTS idx_posts_id_thread ON posts(id, thread);
CREATE INDEX IF NOT EXISTS idx_posts_forum_author ON posts(forum, author);
CREATE INDEX IF NOT EXISTS idx_posts_author ON posts(lower(author));

CREATE UNIQUE INDEX idx_threads_slug ON threads(lower(slug));
CREATE INDEX IF NOT EXISTS idx_threads_author ON threads(lower(author));
CREATE INDEX IF NOT EXISTS idx_threads_forum ON threads(forum);
CREATE INDEX IF NOT EXISTS idx_threads_forum_created ON threads(lower(forum), created);

CREATE INDEX IF NOT EXISTS idx_votes_coverage ON votes(thread, lower(nickname), voice);