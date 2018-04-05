package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type Thread struct {
	Abbrev    bool
	NonLive   bool
	PostCtr   uint
	ImageCtr  uint
	ReplyTime uint
	BumpTime  uint
	Id        uint
	Subject   string
	Board     string
	Body      string
	Posts     []Post
	Banned    bool
	Time      uint
	Image     Image
}

type Post struct {
	Editing  bool
	Deleted  bool
	Banned   bool
	Sage     bool
	Time     uint
	Id       uint
	Body     string
	Flag     string
	PosterID string
	Name     string
	Trip     string
	Auth     string
	//links 	PostLinks
	//commands 	[]Command
	Image Image
}

type Image struct {
	Name string
}

type Aaaaa struct {
	Threads []Thread
}

const url = "https://meguca.org/json/boards/a/"

func main() {
	posts := make(chan Post)

	go checkPosts(posts)

	// check every so often
	fetchPosts(posts)
	for _ = range time.Tick(5 * time.Second) {
		fetchPosts(posts)
	}
}

func fetchPosts(posts chan<- Post) {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()

	result := Aaaaa{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		panic(err)
	}
	threads := result.Threads
	if len(threads) == 0 {
		panic("no threads!")
	}

	for _, thread := range threads {
		// handle OP posts
		p := Post{Banned: thread.Banned, Body: thread.Body, Id: thread.Id, Image: thread.Image}
		posts <- p
		// handle replies
		for _, post := range thread.Posts {
			posts <- post
		}
	}
}

var haikus = make(map[uint]bool)

func checkPosts(posts <-chan Post) {
	for post := range posts {
		if post.Editing {
			// let them finish
			continue
		}
		if haikus[post.Id] {
			// already got it
			continue
		}
		if maybeHaiku(post.Body) {
			haikus[post.Id] = true
			fmt.Printf("Post %v:\n%v\n\n", post.Id, post.Body)
		}
	}
}

var haiku = []int{5, 7, 5}

func maybeHaiku(post string) bool {
	lines := strings.Split(post, "\n")
	pos := 0
	for _, line := range lines {
		words := strings.Fields(line)
		if len(words) <= haiku[pos] && len(words) > 0 {
			pos++
		} else {
			pos = 0
		}
		if pos == len(haiku) {
			return true
		}
	}
	return false
}
