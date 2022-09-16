package download

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/lestrrat-go/strftime"

	"github.com/spaam/fanslyio/internal/context"
	"github.com/spaam/fanslyio/internal/crawl"
	"github.com/spaam/fanslyio/internal/requests"
)

func Filename(timestamp int64) string {
	tm := time.Unix(timestamp, 0)
	s, err := strftime.Format("%Y-%m-%d_%H-%M-%S", tm)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%s_UTC", s)
}

func findExt(url string) string {
	r, _ := http.NewRequest("GET", url, nil)
	return filepath.Ext(path.Base(r.URL.Path))
}

func Download(context *context.Context, usercontent []crawl.UserContent) {
	for _, content := range usercontent {
		fmt.Printf("Downloading content for %s\n", content.Username)
		username := content.Username
		_, err := os.Stat(username)
		if err != nil {
			os.Mkdir(username, os.ModePerm)
		}
		nr := 1
		total := content.Total
		for _, post := range content.Posts {
			count := 1
			for _, attachment := range post.Attachments {
				var name string
				if len(post.Attachments) > 1 {
					name = fmt.Sprintf("%s_%d", Filename(post.Created), count)
					count++
				} else {
					name = Filename(post.Created)
				}
				_, err := os.Stat(fmt.Sprintf("%s/%s%s", username, name, findExt(attachment)))
				if err == nil {
					fmt.Printf("[ %d / %d ] %s, already exists, Skipping\n", nr, total, fmt.Sprintf("%s/%s%s", username, name, findExt(attachment)))
					nr++
					continue
				} else {
					fmt.Printf("[ %d / %d ] %s, Downloading\n", nr, total, fmt.Sprintf("%s/%s%s", username, name, findExt(attachment)))
				}

				body := requests.Request(context, attachment)
				fd, err := os.Create(fmt.Sprintf("%s/%s%s", username, name, findExt(attachment)))
				if err != nil {
					fmt.Println(err)
					continue
				}
				fd.Write(body)
				fd.Close()
				nr++
			}
		}
	}
}
