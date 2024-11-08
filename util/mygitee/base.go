package mygitee

import "time"

type (
	ErrMsg struct {
		Messages []string `json:"messages"`
	}

	TagInfo struct {
		Name    string `json:"name"`
		Message string `json:"message"`
		Commit  struct {
			Sha  string    `json:"sha"`
			Date time.Time `json:"date"`
		} `json:"commit"`
		Tagger struct {
			Name  string    `json:"name"`
			Email string    `json:"email"`
			Date  time.Time `json:"date"`
		} `json:"tagger"`
	}
)
