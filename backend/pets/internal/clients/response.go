package clients

type UserNamesResponse struct {
	Users []User `json:"users"`
}

type User struct {
	ID       string `json:"id"`
	UserName string `json:"username"`
}
