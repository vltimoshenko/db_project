package sql_queries

const (
	InsertForum = `INSERT INTO forums(slug, title, person)
		VALUES($1, $2, $3) RETURNING id;`
	InsertThread = `INSERT INTO threads(author, message, title, forum, slug)
		VALUES($1,$2,$3,$4,$5) RETURNING id;`
	InsertUser = `INSERT INTO persons(about, email, fullname, nickname)
		VALUES($1,$2,$3,$4);`
	InsertPost = "INSERT INTO posts(author, message, parent, thread, forum, created) " +
		"VALUES ($1,$2,$3,$4,$5,$6) RETURNING id;"
	InsertVote = "INSERT INTO votes (nickname, voice, thread)VALUES($1,$2,$3);"

	SelectForumBySlug = `SELECT f.posts, f.slug, f.threads, f.title, f.person
		FROM forums as f WHERE lower(f.slug) = lower($1);`
	SelectThreadBySlug = `SELECT t.author, t.created, t.forum, t.id, t.message, t.slug, t.title, t.votes ` +
		`FROM threads as t WHERE lower(t.slug) = lower($1);`
	SelectThreadsWithParams = `SELECT t.author, t.created, t.forum, t.id, t.message, t.slug, t.title, t.votes ` +
		`FROM threads as t WHERE lower(t.forum) = lower(:forum) `

	SelectUserByNickname = "SELECT p.about, p.email, p.fullname, p.nickname FROM persons as p WHERE p.nickname = $1"
	SelectUserByEmail    = "SELECT p.about, p.email, p.fullname, p.nickname FROM persons as p WHERE p.email = $1"

	UpdateUserByNickname = "UPDATE persons SET about = $1, email = $2, fullname = $3 WHERE nickname = $4 RETURNING id;"
	UpdateThreadByID     = "UPDATE threads SET message = $1, title = $2 WHERE id = $3;"
	UpdateThreadRating   = "UPDATE threads SET votes = votes + $1 WHERE id = $2;"
	UpdateVote           = "UPDATE votes SET voice = $1 WHERE nickname = $2 AND thread = $3;"
	UpdatePost           = "UPDATE posts SET message = $1, is_edited = $2 WHERE id = $3;"

	SelectThreadsBySlug = `SELECT t.author, t.created, t.forum, t.id, t.message, t.slug, t.title, t.votes ` +
		`FROM threads as t WHERE lower(t.slug) = lower($1);`
	SelectThreadsByID = `SELECT t.author, t.created, t.forum, t.id, t.message, t.slug, t.title, t.votes ` +
		`FROM threads as t WHERE t.id = $1;`
	SelectPostByID = "SELECT p.author, p.created, p.forum, p.id, p.is_edited, p.message, p.parent, p.thread " +
		"FROM posts as p WHERE p.id = $1;"
	SelectVote = "SELECT nickname, voice FROM votes WHERE nickname = $1 AND thread = $2;"

	SelectUsersWithParams = "SELECT p.about, p.email, p.fullname, p.nickname " +
		`FROM persons as p ` +
		"WHERE p.nickname IN ( " +
		"SELECT t.author AS nickname " +
		"FROM threads as t " +
		"WHERE lower(t.forum) = lower($1) " +
		"UNION " +
		"SELECT pos.author AS nickname " +
		"FROM posts as pos " +
		"WHERE lower(pos.forum) = lower(:forum) ) "

	SelectDBStatus = "SELECT " +
		"(SELECT COALESCE(SUM(posts), 0) FROM forums WHERE posts > 0) AS posts, " +
		"(SELECT COALESCE(SUM(threads), 0) FROM forums WHERE threads > 0) AS threads, " +
		"(SELECT COUNT(*) FROM persons) AS persons, " +
		"(SELECT COUNT(*) FROM forums) AS forums;"

	Clear = "TRUNCATE votes, posts, threads, forums, persons RESTART IDENTITY CASCADE;"

	// SELECTPostsFlat = "SELECT p.author, p.created, p.forum, p.id, p.is_edited, p.message, p.parent, p.thread " +
	// 	"FROM posts as p WHERE p.thread = $1 AND p.id > $3 ORDER BY p.id LIMIT $2"
	// SELECTPostsFlatDesc = "SELECT p.author, p.created, p.forum, p.id, p.is_edited, p.message, p.parent, p.thread " +
	// 	"FROM posts as p WHERE p.thread = $1 AND p.id < $3 ORDER BY p.id DESC LIMIT $2"

	// SelectPostsFlat = "SELECT p.author, p.created, p.forum, p.id, p.is_edited, p.message, p.parent, p.thread " +
	// 	"FROM posts as p WHERE p.thread = $1 "

	//AND p.id > $3 ORDER BY p.id LIMIT $2"

	////////////////////////////////////////////////////////////////
	SelectPostsFlat = "SELECT p.author, p.created, p.forum, p.id, p.is_edited, p.message, p.parent, p.thread " +
		"FROM posts as p WHERE p.thread = $1 AND p.id > $3 ORDER BY p.id LIMIT $2"
	SelectPostsFlatDesc = "SELECT p.author, p.created, p.forum, p.id, p.is_edited, p.message, p.parent, p.thread " +
		"FROM posts as p WHERE p.thread = $1 AND p.id < $3 ORDER BY p.id DESC LIMIT $2"

	SelectPostsTree = "WITH RECURSIVE temp1 (author, created, forum, id, is_edited, message, parent, thread, PATH, LEVEL, root ) AS ( " +
		"SELECT T1.author, T1.created, T1.forum, T1.id, T1.is_edited, T1.message, T1.parent, T1.thread, CAST (10000 + T1.id AS VARCHAR (50)) as PATH, 1, T1.id as root " +
		"FROM posts as T1 WHERE T1.parent = 0 and T1.thread = $1 " +
		"union " +
		"select T2.author, T2.created, T2.forum, T2.id, T2.is_edited, T2.message, T2.parent, T2.thread, CAST ( temp1.PATH ||'->'|| 10000 + T2.id AS VARCHAR(50)), LEVEL + 1, root " +
		"FROM posts T2 INNER JOIN temp1 ON( temp1.id = T2.parent) " +
		") " +
		"select author, created, forum, id, is_edited, message, parent, thread from temp1 ORDER BY root, PATH LIMIT $2;"

	SelectPostsTreeSince = "WITH RECURSIVE temp1 (author, created, forum, id, is_edited, message, parent, thread, PATH, LEVEL ) AS ( " +
		"SELECT T1.author, T1.created, T1.forum, T1.id, T1.is_edited, T1.message, T1.parent, T1.thread, CAST (1000000 + T1.id AS VARCHAR (50)) as PATH, 1 " +
		"FROM posts as T1 WHERE T1.parent = 0 AND T1.thread = $1" +
		"union " +
		"select T2.author, T2.created, T2.forum, T2.id, T2.is_edited, T2.message, T2.parent, T2.thread, CAST ( temp1.PATH ||'->'|| T2.id AS VARCHAR(50)), LEVEL + 1 " +
		"FROM posts T2 INNER JOIN temp1 ON( temp1.id = T2.parent) " +
		") " +
		"select author, created, forum, id, is_edited, message, parent, thread from temp1 WHERE id > $3 ORDER BY PATH LIMIT $2;"

	SelectPostsTreeDesc = "WITH RECURSIVE temp1 (author, created, forum, id, is_edited, message, parent, thread, PATH, LEVEL, root ) AS ( " +
		"SELECT T1.author, T1.created, T1.forum, T1.id, T1.is_edited, T1.message, T1.parent, T1.thread, CAST (1000000 + T1.id AS VARCHAR (50)) as PATH, 1, T1.id as root " +
		"FROM posts as T1 WHERE T1.parent = 0 AND T1.thread = $1" +
		"union " +
		"select  T2.author, T2.created, T2.forum, T2.id, T2.is_edited, T2.message, T2.parent, T2.thread, CAST (temp1.PATH ||'->'|| T2.id AS VARCHAR(50)), LEVEL + 1, root " +
		"FROM posts as T2 INNER JOIN temp1 ON (temp1.id = T2.parent) " +
		") " +
		"select author, created, forum, id, is_edited, message, parent, thread from temp1 WHERE id < $3 ORDER BY PATH DESC LIMIT $2;"
	SelectPostsTreeSinceDesc = "WITH RECURSIVE temp1 (author, created, forum, id, is_edited, message, parent, thread, PATH, LEVEL, root ) AS ( " +
		"SELECT T1.author, T1.created, T1.forum, T1.id, T1.is_edited, T1.message, T1.parent, T1.thread, CAST (1000000 + T1.id AS VARCHAR (50)) as PATH, 1, T1.id as root " +
		"FROM posts as T1 WHERE T1.parent = 0 AND T1.thread = $1" +
		"union " +
		"select  T2.author, T2.created, T2.forum, T2.id, T2.is_edited, T2.message, T2.parent, T2.thread, CAST (temp1.PATH ||'->'|| T2.id AS VARCHAR(50)), LEVEL + 1, root " +
		"FROM posts as T2 INNER JOIN temp1 ON (temp1.id = T2.parent) " +
		") " +
		"select author, created, forum, id, is_edited, message, parent, thread from temp1 ORDER BY PATH;"

	SelectPostsParentTree = "WITH RECURSIVE temp1 (author, created, forum, id, is_edited, message, parent, thread, PATH, LEVEL, root ) AS ( " +
		"SELECT T1.author, T1.created, T1.forum, T1.id, T1.is_edited, T1.message, T1.parent, T1.thread, CAST (10000 + T1.id AS VARCHAR (50)) as PATH, 1, T1.id as root " +
		"FROM posts as T1 WHERE T1.parent = 0 AND T1.thread = $1" +
		"union " +
		"select T2.author, T2.created, T2.forum, T2.id, T2.is_edited, T2.message, T2.parent, T2.thread, CAST ( temp1.PATH ||'->'|| 10000 + T2.id AS VARCHAR(50)), LEVEL + 1, root " +
		"FROM posts T2 INNER JOIN temp1 ON( temp1.id = T2.parent) " +
		") " +
		"select author, created, forum, id, is_edited, message, parent, thread from temp1 ORDER BY root, PATH;"

	SelectPostsParentTreeDesc = "WITH RECURSIVE temp1 (author, created, forum, id, is_edited, message, parent, thread, PATH, LEVEL, root ) AS ( " +
		"SELECT T1.author, T1.created, T1.forum, T1.id, T1.is_edited, T1.message, T1.parent, T1.thread, CAST (10000 + T1.id AS VARCHAR (50)) as PATH, 1, T1.id as root " +
		"FROM posts as T1 WHERE T1.parent = 0 AND T1.thread = $1" +
		"union " +
		"select  T2.author, T2.created, T2.forum, T2.id, T2.is_edited, T2.message, T2.parent, T2.thread, CAST ( temp1.PATH ||'->'|| 10000 + T2.id AS VARCHAR(50)), LEVEL + 1, root " +
		"FROM posts as T2 INNER JOIN temp1 ON( temp1.id = T2.parent) " +
		") " +
		"select author, created, forum, id, is_edited, message, parent, thread  from temp1 ORDER BY root desc, PATH;"
)
