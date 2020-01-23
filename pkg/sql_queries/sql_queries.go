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
	InsertForumUsers = `INSERT INTO forum_users(nickname, forum) VALUES `

	UpdateUserByNickname = "UPDATE persons SET about = COALESCE(nullif($2, ''), about)," +
		"email = COALESCE(nullif($3, ''), email), fullname = COALESCE(nullif($4, ''), fullname)" +
		"WHERE nickname = $1 RETURNING nickname, fullname, email, about;"
	UpdatePost = "UPDATE posts " +
		"SET message = COALESCE(nullif($1, ''), message), is_edited = ($1 <> message AND $1 <> '')" + // not null, maybe need trim()
		"WHERE id = $2" +
		"RETURNING author, forum, created, id, is_edited, message, parent, thread;"

	UpdateThreadByID = "UPDATE threads SET title = COALESCE(nullif($2, ''), title)," +
		"message = COALESCE(nullif($1, ''), message) WHERE id = $3 " +
		"RETURNING author, created, forum, id, message, slug, title, votes"

	UpdateThreadBySlug = "UPDATE threads SET title = COALESCE(nullif($2, ''), title)," +
		"message = COALESCE(nullif($1, ''), message) WHERE slug = $3 " +
		"RETURNING author, created, forum, id, message, slug, title, votes"

	UpdateVoteByThreadID   = "UPDATE votes SET voice = $1 WHERE nickname = $2 AND thread = $3;"
	UpdateVoteByThreadSlug = "UPDATE votes SET voice = $1 WHERE nickname = $2 AND " +
		"thread = (SELECT id FROM threads WHERE slug = $3);"

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
		`JOIN forum_users ON p.nickname = forum_users.nickname ` +
		`WHERE forum_users.forum = :forum `

	// select курильщика
	// SelectUsersWithParams = "SELECT p.about, p.email, p.fullname, p.nickname " +
	// 	`FROM forum_users as p ` +
	// 	`WHERE p.forum = :forum `

	SelectDBStatus = "SELECT " +
		"(SELECT COUNT(*) FROM posts) AS posts, " +
		"(SELECT COUNT(*) FROM threads) AS threads, " +
		"(SELECT COUNT(*) FROM persons) AS persons, " +
		"(SELECT COUNT(*) FROM forums) AS forums;"

	SelectPostsSorted = `SELECT author, forum, created, posts.id, is_edited, message, COALESCE(parent, 0), thread 
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
