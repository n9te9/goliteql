package generated

type Post struct {
	Id        string `json:"id"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	Published bool   `json:"published"`
}
type User struct {
	Id    string  `json:"id"`
	Name  string  `json:"name"`
	Email string  `json:"email"`
	Age   *int    `json:"age"`
	Posts []*Post `json:"posts"`
}
