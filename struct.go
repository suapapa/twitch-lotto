package main

// User reperesnt genenal user information of meetup and twitch
type User struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	ProfileImageURL string `json:"profile_image_url"`
}
