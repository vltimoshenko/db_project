package sql_queries

const (
	InsertForum = `INSERT INTO forums(slug, title, person)
		VALUES($1, $2, $3)`
	InsertThread                       = `INSERT INTO threads (author, created, message, title, forum) values ($1,$2,$3,$4,$5) RETURNING id;`
	InsertThreadWithoutCreated         = `INSERT INTO threads (author, message, title, forum) values ($1,$2,$3,$4) RETURNING id;`
	InsertThreadWithSlugWithoutCreated = `INSERT INTO threads (author, message, title, forum, slug) values ($1,$2,$3,$4,$5) RETURNING id;`
	InsertThreadWithSlug               = `INSERT INTO threads (author, created, message, title, forum, slug) values ($1,$2,$3,$4,$5,$6) RETURNING id;`
	InsertUser                         = `INSERT INTO persons(about, email, fullname, nickname)VALUES($1,$2,$3,$4);`
	InsertPosts                        = `INSERT INTO posts(author, message, parent, thread, created, forum)VALUES `
	InsertVoteByThreadID               = "INSERT INTO votes (nickname, voice, thread)VALUES($1,$2,$3);"
	InsertVoteByThreadSlug             = "INSERT INTO votes (nickname, voice, thread)VALUES" +
		"($1,$2,(SELECT id FROM threads WHERE slug = $3));"

	SelectForumBySlug = `SELECT f.posts, f.slug, f.threads, f.title, f.person
		FROM forums as f WHERE f.slug = $1;`
	SelectThreadsWithParams = `SELECT t.author, t.created, t.forum, t.id, t.message, t.slug, t.title, t.votes ` +
		`FROM threads as t WHERE t.forum = :forum `
	SelectUserByNickname = "SELECT p.about, p.email, p.fullname, p.nickname FROM persons as p WHERE p.nickname = $1"
	SelectUserByEmail    = "SELECT p.about, p.email, p.fullname, p.nickname FROM persons as p WHERE p.email = $1"

	// UpdateUserByNickname   = "UPDATE persons SET about = $1, email = $2, fullname = $3 WHERE nickname = $4;"

	UpdateUserByNickname = `
				UPDATE persons
				SET about = coalesce(nullif($2, ''), about),
					email = coalesce(nullif($3, ''), email),
					fullname = coalesce(nullif($4, ''), fullname)
				WHERE nickname = $1
				RETURNING nickname, fullname, email, about;`

	UpdateThreadByID       = "UPDATE threads SET message = $1, title = $2 WHERE id = $3;"
	UpdateVoteByThreadID   = "UPDATE votes SET voice = $1 WHERE nickname = $2 AND thread = $3;"
	UpdateVoteByThreadSlug = "UPDATE votes SET voice = $1 WHERE nickname = $2 AND thread = (SELECT id FROM threads WHERE slug = $3);"
	UpdatePost             = "UPDATE posts SET message = $1, is_edited = $2 WHERE id = $3;"

	SelectThreadBySlug = `SELECT t.author, t.created, t.forum, t.id, t.message, t.slug, t.title, t.votes ` +
		`FROM threads as t WHERE t.slug = $1;`
	SelectThreadByID = `SELECT t.author, t.created, t.forum, t.id, t.message, t.slug, t.title, t.votes ` +
		`FROM threads as t WHERE t.id = $1;`
	SelectPostByID = "SELECT p.author, p.created, p.forum, p.id, p.is_edited, p.message, p.parent, p.thread " +
		"FROM posts as p WHERE p.id = $1;"
	SelectVoteByThreadID   = "SELECT nickname, voice FROM votes WHERE nickname =$1 AND thread = $2;"
	SelectVoteByThreadSlug = "SELECT nickname, voice FROM votes WHERE nickname = $1 AND thread = (SELECT id FROM threads WHERE slug = $2);"

	SelectUsersWithParams = "SELECT p.about, p.email, p.fullname, p.nickname " +
		`FROM persons as p ` +
		`JOIN forum_users ON p.nickname = forum_users.person ` +
		`WHERE forum_users.forum = :forum `

		// SelectDBStatus = "SELECT " +
		// 	"(SELECT COALESCE(SUM(posts), 0) FROM forums WHERE posts > 0) AS posts, " +
		// 	"(SELECT COALESCE(SUM(threads), 0) FROM forums WHERE threads > 0) AS threads, " +
		// 	"(SELECT COUNT(*) FROM persons) AS persons, " +
		// 	"(SELECT COUNT(*) FROM forums) AS forums;"
	SelectDBStatus = "SELECT " +
		"(SELECT COUNT(*) FROM posts) AS posts, " +
		"(SELECT COUNT(*) FROM threads) AS threads, " +
		"(SELECT COUNT(*) FROM persons) AS persons, " +
		"(SELECT COUNT(*) FROM forums) AS forums;"

	Clear = "TRUNCATE votes, posts, threads, forums, persons RESTART IDENTITY CASCADE;"
)
