package sql_queries

const (
	InsertForum = `INSERT INTO forums(slug, title, person)
		VALUES($1, $2, $3)`
	InsertThread = "INSERT INTO threads (author, created, message, title, forum)VALUES" +
		"($1,$2,$3,$4,$5) RETURNING id;"
	InsertThreadWithoutCreated = "INSERT INTO threads (author, message, title, forum)VALUES" +
		"($1,$2,$3,$4) RETURNING id;"
	InsertThreadWithSlugWithoutCreated = "INSERT INTO threads (author, message, title, forum, slug)VALUES" +
		"($1,$2,$3,$4,$5) RETURNING id;"
	InsertThreadWithSlug = "INSERT INTO threads (author, created, message, title, forum, slug)VALUES" +
		"($1,$2,$3,$4,$5,$6) RETURNING id;"
	InsertUser = "INSERT INTO persons(about, email, fullname, nickname)" +
		"VALUES($1,$2,$3,$4);"
	InsertPosts            = "INSERT INTO posts(author, message, parent, thread, created, forum)VALUES "
	InsertVoteByThreadID   = "INSERT INTO votes (nickname, voice, thread)VALUES($1,$2,$3);"
	InsertVoteByThreadSlug = "INSERT INTO votes (nickname, voice, thread)VALUES" +
		"($1,$2,(SELECT id FROM threads WHERE slug = $3));"

	UpdateUserByNickname = `
				UPDATE persons
				SET about = coalesce(nullif($2, ''), about),
					email = coalesce(nullif($3, ''), email),
					fullname = coalesce(nullif($4, ''), fullname)
				WHERE nickname = $1
				RETURNING nickname, fullname, email, about;`

	UpdateThreadByID       = "UPDATE threads SET message = $1, title = $2 WHERE id = $3;"
	UpdateVoteByThreadID   = "UPDATE votes SET voice = $1 WHERE nickname = $2 AND thread = $3;"
	UpdateVoteByThreadSlug = "UPDATE votes SET voice = $1 WHERE nickname = $2 AND " +
		"thread = (SELECT id FROM threads WHERE slug = $3);"

	UpdatePost = "UPDATE posts SET message = $1, is_edited = $2 WHERE id = $3;"

	SelectForumBySlug = `SELECT f.posts, f.slug, f.threads, f.title, f.person
		FROM forums as f WHERE f.slug = $1;`
	SelectThreadsWithParams = `SELECT t.author, t.created, t.forum, t.id, t.message, t.slug, t.title, t.votes ` +
		`FROM threads as t WHERE t.forum = :forum `
	SelectUserByNickname = "SELECT p.about, p.email, p.fullname, p.nickname FROM persons as p WHERE p.nickname = $1"
	SelectUserByEmail    = "SELECT p.about, p.email, p.fullname, p.nickname FROM persons as p WHERE p.email = $1"

	SelectThreadBySlug = `SELECT t.author, t.created, t.forum, t.id, t.message, t.slug, t.title, t.votes ` +
		`FROM threads as t WHERE t.slug = $1;`
	SelectThreadByID = `SELECT t.author, t.created, t.forum, t.id, t.message, t.slug, t.title, t.votes ` +
		`FROM threads as t WHERE t.id = $1;`
	SelectPostByID = "SELECT p.author, p.created, p.forum, p.id, p.is_edited, p.message, p.parent, p.thread " +
		"FROM posts as p WHERE p.id = $1;"
	SelectVoteByThreadID   = "SELECT nickname, voice FROM votes WHERE nickname =$1 AND thread = $2;"
	SelectVoteByThreadSlug = "SELECT nickname, voice FROM votes WHERE nickname = $1 AND " +
		"thread = (SELECT id FROM threads WHERE slug = $2);"

	SelectUsersWithParams = "SELECT p.about, p.email, p.fullname, p.nickname " +
		`FROM persons as p ` +
		`JOIN forum_users ON p.nickname = forum_users.person ` +
		`WHERE forum_users.forum = :forum `

	SelectDBStatus = "SELECT " +
		"(SELECT COUNT(*) FROM posts) AS posts, " +
		"(SELECT COUNT(*) FROM threads) AS threads, " +
		"(SELECT COUNT(*) FROM persons) AS persons, " +
		"(SELECT COUNT(*) FROM forums) AS forums;"

	SelectPostsSorted = `SELECT author, forum, created, posts.id, is_edited, message, coalesce(parent, 0), thread 
					FROM posts
					{{.Condition}}
					ORDER BY {{.OrderBy}}
					{{.Limit}}`
	SelectPostsParentTree = `JOIN (
						SELECT parents.id FROM posts AS parents
						WHERE parents.thread=$1 AND parents.parent IS NULL
							{{- if .Since}} AND {{.Since}}{{- end}}
						ORDER BY parents.path[1] {{.Desc}}
						{{.Limit}}
						) as p ON path[1]=p.id`

	Clear = "TRUNCATE votes, posts, threads, forums, persons RESTART IDENTITY CASCADE;"
)
