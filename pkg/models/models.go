package models

type NewForum struct {
	Slug  string `json:"slug"     db:"slug"`
	Title string `json:"title"    db:"title"`
	User  string `json:"user"     db:"person"`
}

type Forum struct {
	Posts  int    `json:"posts"   db:"posts"`
	Slug   string `json:"slug"    db:"slug"`
	Thread int    `json:"threads" db:"threads"`
	Title  string `json:"title"   db:"title"`
	User   string `json:"user"    db:"person"`
}

type NewPost struct {
	Author  string `json:"author"`
	Message string `json:"message"`
	Parent  int64  `json:"parent"`
}

type PostUpdate struct {
	Message string `json:"message"`
}

type Post struct {
	Author   string `json:"author"     db:"author"`
	Created  string `json:"created"    db:"created"`
	Forum    string `json:"forum"      db:"forum"`
	ID       int64  `json:"id"         db:"id"`
	IsEdited bool   `json:"isEdited"  db:"is_edited"`
	Message  string `json:"message"    db:"message"`
	Parent   int64  `json:"parent"     db:"parent"`
	Thread   int    `json:"thread"     db:"thread"`
}

type NewThread struct {
	Author  string `json:"author"`
	Created string `json:"created"`
	Message string `json:"message"`
	Slug    string `json:"slug"`
	Title   string `json:"title"`
}

type ThreadUpdate struct {
	Message string `json:"message"`
	Title   string `json:"title"`
}

type Thread struct {
	Author  string `json:"author"`
	Created string `json:"created"`
	Forum   string `json:"forum"`
	ID      int    `json:"id"`
	Message string `json:"message"`
	Slug    string `json:"slug"`
	Title   string `json:"title"`
	Votes   int    `json:"votes"`
}

type NewUser struct {
	About    string `json:"about"    db:"about"`
	Email    string `json:"email"    db:"email"`
	Fullname string `json:"fullname" db:"fullname"`
}

type User struct {
	About    string `json:"about"    db:"about"`
	Email    string `json:"email"    db:"email"`
	Fullname string `json:"fullname" db:"fullname"`
	Nickname string `json:"nickname" db:"nickname"`
}

// // type Users []*User

type Status struct {
	Post   int `json:"post"          db:"posts"`
	Thread int `json:"thread"        db:"threads"`
	User   int `json:"user"          db:"persons"`
	Forum  int `json:"forum"         db:"forums"`
}

type Vote struct {
	Nickname string `json:"nickname" db:"nickname"`
	Voice    int    `json:"voice"    db:"voice"`
}
