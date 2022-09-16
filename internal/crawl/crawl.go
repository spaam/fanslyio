package crawl

import (
	"fmt"
	"strings"

	"github.com/spaam/fanslyio/internal/context"
	"github.com/spaam/fanslyio/internal/requests"
	"github.com/tidwall/gjson"
)

type User struct {
	Username string
	UserID   string
	GroupID  string
}

type Post struct {
	Created     int64
	Attachments []string
}

type UserContent struct {
	Username string
	Posts    []Post
	Total    int
}

type Items struct {
	After string
	Posts []Post
	Total int
}

func AccountMedia(postst *Post, id string, body []byte) {
	res := gjson.Get(string(body), fmt.Sprintf("response.accountMedia.#(id==%s)", id))
	locations := gjson.Get(res.String(), "media.variants.#.locations")
	if len(locations.Array()) > 0 {
		location := locations.Array()[0]
		if len(location.Array()) > 0 {
			picurl := gjson.Get(location.Array()[0].String(), "location")
			postst.Attachments = append(postst.Attachments, picurl.String())
		}
	}
}

func GetItems(posts gjson.Result, body []byte) (items Items) {
	lastpost := "0"
	var posts_array []Post
	total := 0
	for _, post := range posts.Array() {
		attachments := gjson.Get(post.String(), "attachments")
		created := int64(gjson.Get(post.String(), "createdAt").Int())
		after := gjson.Get(post.String(), "id").String()
		postst := Post{Created: created}
		for _, attachment := range attachments.Array() {
			contentType := gjson.Get(attachment.String(), "contentType").Int()
			id := gjson.Get(attachment.String(), "contentId")
			if contentType == 1 {
				AccountMedia(&postst, id.String(), body)
			} else if contentType == 2 {
				mediaIDS := gjson.Get(string(body), fmt.Sprintf("response.accountMediaBundles.#(id==%s).accountMediaIds", id.String()))
				for _, mediaID := range mediaIDS.Array() {
					AccountMedia(&postst, mediaID.String(), body)
				}
			}
		}
		lastpost = after

		if len(postst.Attachments) > 0 {
			total = total + len(postst.Attachments)
			posts_array = append(posts_array, postst)
		}
	}

	items.After = lastpost
	items.Posts = posts_array
	items.Total = total
	return items
}

func Crawl(contxt *context.Context) ([]UserContent, error) {
	var content []UserContent
	posts := GetPosts(contxt)
	messages := GetMessages(contxt)

	for _, post := range posts {
		content = append(content, post)
	}

	for _, mess := range messages {
		isthere := false
		for _, post := range posts {
			if post.Username == mess.Username {
				isthere = true
				post.Posts = append(post.Posts, mess.Posts...)
				post.Total = post.Total + mess.Total
			}
		}
		if !isthere {
			posts = append(posts, mess)
		}
	}

	return posts, nil
}

func GetMessages(contxt *context.Context) (usercontent []UserContent) {
	allow := contxt.Allowlistset
	block := contxt.Blocklistset
	adduser := true
	if contxt.Allowlistset {
		adduser = false
	}
	var url = "https://apiv3.fansly.com/api/v1/messaging/groups?sortOrder=1&flags=0&subscriptionTierId=&search=&limit=100&offset=0"
	body := requests.Request(contxt, url)
	messages := gjson.Get(string(body), "response.data")
	var users []User
	for _, i := range messages.Array() {
		username := strings.ToLower(i.Get("partnerUsername").String())
		if allow {
			adduser = false
		}
		if allow && context.Checkuser(username, contxt.Allowlist) {
			adduser = true
		}
		if block && context.Checkuser(username, contxt.Blocklist) {
			continue
		}
		if adduser {
			user := User{Username: username, UserID: i.Get("partnerAccountId").String(), GroupID: i.Get("groupId").String()}
			users = append(users, user)
		}
	}
	for _, i := range users {
		fmt.Printf("Getting messages %s\n", i.Username)
		//url = fmt.Sprintf("https://apiv3.fansly.com/api/v1/message?groupId=%d&limit=100", i.GroupID)
		user_content := UserContent{Username: i.Username, Total: 0}
		more := true
		after := "0"
		for more {
			url = fmt.Sprintf("https://apiv3.fansly.com/api/v1/message?groupId=%s&before=%s&limit=25", i.GroupID, after)
			body = requests.Request(contxt, url)
			posts := gjson.Get(string(body), "response.messages")
			if len(posts.Array()) == 0 {
				more = false
			}
			items := GetItems(posts, body)
			after = items.After
			user_content.Posts = append(user_content.Posts, items.Posts...)
			user_content.Total = user_content.Total + items.Total
		}
		usercontent = append(usercontent, user_content)
	}

	return usercontent
}

func GetPosts(contxt *context.Context) []UserContent {
	var url = "https://apiv3.fansly.com/api/v1/account/settings"
	allow := contxt.Allowlistset
	block := contxt.Blocklistset
	adduser := true
	if contxt.Allowlistset {
		adduser = false
	}
	var usercontent []UserContent
	body := requests.Request(contxt, url)
	myid := gjson.Get(string(body), "response.accountId").String()
	url = fmt.Sprintf("https://apiv3.fansly.com/api/v1/account/%s/following?before=0&after=0&limit=200&offset=0", myid)

	body = requests.Request(contxt, url)
	following := gjson.Get(string(body), "response.#.accountId")
	var accounts string
	for _, follow := range following.Array() {
		accounts = fmt.Sprintf("%s,%s", accounts, follow.String())
	}
	url = fmt.Sprintf("https://apiv3.fansly.com/api/v1/account?ids=%s", accounts)
	body = requests.Request(contxt, url)
	follow := gjson.Get(string(body), "response")
	var users []User
	for _, i := range follow.Array() {
		username := strings.ToLower(gjson.Get(i.String(), "username").String())
		if allow {
			adduser = false
		}
		if allow && context.Checkuser(username, contxt.Allowlist) {
			adduser = true
		}
		if block && context.Checkuser(username, contxt.Blocklist) {
			continue
		}
		if adduser {
			user := User{Username: username, UserID: gjson.Get(i.String(), "id").String()}
			users = append(users, user)
		}
	}
	if len(users) == 0 {
		return usercontent
	}

	for _, i := range users {
		fmt.Printf("Getting posts for %s\n", i.Username)
		user_content := UserContent{Username: i.Username, Total: 0}
		more := true
		after := "0"
		for more {
			url := fmt.Sprintf("https://apiv3.fansly.com/api/v1/timeline/%s?before=%s&after=0", i.UserID, after)
			body := requests.Request(contxt, url)
			posts := gjson.Get(string(body), "response.posts")
			if len(posts.Array()) == 0 {
				more = false
			}
			items := GetItems(posts, body)
			after = items.After
			user_content.Posts = append(user_content.Posts, items.Posts...)
			user_content.Total = user_content.Total + items.Total
		}
		usercontent = append(usercontent, user_content)
	}
	return usercontent
}
