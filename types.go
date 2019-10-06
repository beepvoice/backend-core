package main

// String pointer means nullable

type UpdateMsg struct {
  Type    string `json:"type"`
  Data    string `json:"data"`
}

type Contact struct {
  UserA   string `json:"usera"`   // First user ID
  UserB   string `json:"userb"`   // Second user ID
}

type Member struct {
  User          string  `json:"user"`
  Conversation  string  `json:"conversation"`
  Pinned        bool    `json:"pinned"`
}

type Conversation struct {
	ID      string `json:"id"`      // id
	Title   string `json:"title"`   // title
	DM      bool   `json:"dm"`      // dm
	Picture string `json:"picture"` // picture
	Pinned  bool   `json:"pinned"`  // pinned
}

type User struct {
	ID          string  `json:"id"`           // id
	Username    *string `json:"username"`     // username
	Bio         string  `json:"bio"`          // bio
	ProfilePic  string  `json:"profile_pic"`  // profile_pic
	FirstName   string  `json:"first_name"`   // first_name
	LastName    string  `json:"last_name"`    // last_name
	PhoneNumber string  `json:"phone_number"` // phone_number
}

type PhoneNumber struct {
	PhoneNumber string `json:"phone_number"`
}
