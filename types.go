package main

type Conversation struct {
	ID    string `json:"id"`    // id
	Title string `json:"title"` // title
}

type User struct {
	ID          string `json:"id"`           // id
	Username 		string `json:"username"`		 // username
	FirstName   string `json:"first_name"`   // first_name
	LastName    string `json:"last_name"`    // last_name
	PhoneNumber string `json:"phone_number"` // phone_number
}
