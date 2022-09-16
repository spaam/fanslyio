package context

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/tidwall/gjson"
)

type Context struct {
	Client        http.Client
	UserAgent     string
	Authorization string
	Allowlist     []string
	Blocklist     []string
	Allowlistset  bool
	Blocklistset  bool
}

func ParseConfig(file string) (Context, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return Context{}, fmt.Errorf("Can't find config file")
	}
	useragent := gjson.Get(string(data), "user-agent")
	authorization := gjson.Get(string(data), "authorization")
	if useragent.Type == gjson.Null {
		return Context{}, fmt.Errorf("user-agent is empty")
	}
	if authorization.Type == gjson.Null {
		return Context{}, fmt.Errorf("authoriztion is empty")
	}

	cc := Context{Client: http.Client{}, UserAgent: useragent.String(), Authorization: authorization.String(), Allowlistset: false, Blocklistset: false}

	allowlist := gjson.Get(string(data), "allowlist")
	if allowlist.Type != gjson.Null {
		for _, allow := range allowlist.Array() {
			cc.Allowlist = append(cc.Allowlist, strings.ToLower(allow.String()))
			cc.Allowlistset = true
		}
	}
	blocklist := gjson.Get(string(data), "blocklist")
	if blocklist.Type != gjson.Null {
		for _, block := range blocklist.Array() {
			cc.Blocklist = append(cc.Blocklist, strings.ToLower(block.String()))
			cc.Blocklistset = true
		}
	}
	return cc, nil
}

func Checkuser(username string, users []string) bool {
	for _, i := range users {
		if i == username {
			return true
		}
	}
	return false
}
