package main

import "gopkg.in/guregu/null.v3"

type UpdateMsg struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

type Contact struct {
	UserA string `json:"usera"` // First user ID
	UserB string `json:"userb"` // Second user ID
}

type Member struct {
	User         string `json:"user"`
	Conversation string `json:"conversation"`
	Pinned       bool   `json:"pinned"`
}

type Conversation struct {
	ID      string      `json:"id"`      // id
	Title   null.String `json:"title"`   // title
	Picture null.String `json:"picture"` // picture
	Pinned  bool        `json:"pinned"`  // pinned
}

type User struct {
	ID          string      `json:"id"`           // id
	Username    null.String `json:"username"`     // username
	Bio         string      `json:"bio"`          // bio
	ProfilePic  string      `json:"profile_pic"`  // profile_pic
	FirstName   string      `json:"first_name"`   // first_name
	LastName    string      `json:"last_name"`    // last_name
	PhoneNumber string      `json:"phone_number"` // phone_number
}

type PhoneNumber struct {
	PhoneNumber string `json:"phone_number"`
}
