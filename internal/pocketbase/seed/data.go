//go:build dev

package seed

type seedSuperuser struct{ Email, Password string }
type seedUser struct{ Email, Password, Name string }
type seedPost struct{ Title, Body, AuthorEmail string }

var superusers = []seedSuperuser{
	{Email: "root@dev.com", Password: "root1234"},
}

var users = []seedUser{
	{Email: "alice@dev.com", Password: "password1234", Name: "Alice"},
	{Email: "bob@dev.com", Password: "password1234", Name: "Bob"},
}

var posts = []seedPost{
	{Title: "Hello World", Body: "First test post.", AuthorEmail: "alice@dev.com"},
}