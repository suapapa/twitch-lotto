package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"golang.org/x/oauth2/clientcredentials"
	"golang.org/x/oauth2/twitch"
)

var (
	clientID     = os.Getenv("TWITCH_CLIENT_ID")
	clientSecret = os.Getenv("TWITCH_CLIENT_SECRET")
	oauth2Config *clientcredentials.Config
	httpClient   *http.Client
)

// TwitchFollowerTo returns follows of loging user
func TwitchFollowerTo(login string) ([]*User, error) {
	oauth2Config = &clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     twitch.Endpoint.TokenURL,
		Scopes:       []string{"user:read:email"},
	}

	// token, err := oauth2Config.Token(context.TODO())
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("Access token: %s\n", token.AccessToken)

	httpClient = oauth2Config.Client(context.TODO())

	u := getTwitchUserInfo([]string{login}, nil)
	followTo := u.Data[0].ID
	members := getTwitchUserFollowTo(followTo).Data
	followerIDs := make([]string, len(members))
	for i, f := range members {
		followerIDs[i] = f.FromID
	}
	fu := getTwitchUserInfo(nil, followerIDs).Data
	users := make([]*User, len(fu))
	for i, f := range fu {
		users[i] = &User{
			ID:              f.Login,
			Name:            f.DisplayName,
			ProfileImageURL: f.ProfileImageURL,
		}
	}

	return users, nil
}

// TwitchUser reperesents users info
type TwitchUser struct {
	Data []struct {
		ID              string `json:"id"`
		Login           string `json:"login"`
		DisplayName     string `json:"display_name"`
		Type            string `json:"type"`
		BroadcasterType string `json:"broadcaster_type"`
		Description     string `json:"description"`
		ProfileImageURL string `json:"profile_image_url"`
		OfflineImageURL string `json:"offline_image_url"`
		ViewCount       int    `json:"view_count"`
	} `json:"data"`
}

func getTwitchUserInfo(login []string, id []string) *TwitchUser {
	values := make(url.Values)
	for _, l := range login {
		values.Add("login", l)
	}
	for _, i := range id {
		values.Add("id", i)
	}
	req, err := http.NewRequest("GET", "https://api.twitch.tv/helix/users?"+values.Encode(), nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Client-Id", clientID)

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	var user TwitchUser
	json.NewDecoder(resp.Body).Decode(&user)
	return &user
}

// TwitchFollow reperesents users follows
type TwitchFollow struct {
	Total int `json:"total"`
	Data  []struct {
		FromID     string    `json:"from_id"`
		FromName   string    `json:"from_name"`
		ToID       string    `json:"to_id"`
		ToName     string    `json:"to_name"`
		FollowedAt time.Time `json:"followed_at"`
	} `json:"data"`
	Pagination struct {
		Cursor string `json:"cursor"`
	} `json:"pagination"`
}

func getTwitchUserFollowTo(id string) *TwitchFollow {
	values := make(url.Values)
	values.Add("to_id", id)
	values.Add("first", "100")
	// values.Add("after", "TODO")
	req, err := http.NewRequest("GET", "https://api.twitch.tv/helix/users/follows?"+values.Encode(), nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Client-Id", clientID)

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	var follow TwitchFollow
	json.NewDecoder(resp.Body).Decode(&follow)
	return &follow
}
