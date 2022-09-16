package requests

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spaam/fanslyio/internal/context"
)

func Request(crawl *context.Context, url string) []byte {
	client := crawl.Client
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic("dududu")
	}
	req.Header.Set("Authorization", crawl.Authorization)
	req.Header.Set("User-Agent", crawl.UserAgent)
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode > 400 {
		fmt.Println("Wrong user-agent or authorization")
		return []byte{}
	}
	body, _ := ioutil.ReadAll(resp.Body)
	return body
}
