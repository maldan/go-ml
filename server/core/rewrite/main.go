package rewrite

import (
	"net/url"
	"regexp"
	"strings"
)

type Rewrite struct {
	RegExp     *regexp.Regexp
	HttpMethod string
	Url        string
}

// HandleRewrite can remap url. For example GET /api/me -> GET /api/user/profile
// It also pass url parameters /user/(?P<id>[0-9]+) -> /api/user/id?=$id
func HandleRewrite(rewriteList []Rewrite, httpMethod string, url *url.URL) {
	for _, r := range rewriteList {
		if !strings.Contains(r.HttpMethod, httpMethod) {
			continue
		}

		// Check regex
		match := r.RegExp.FindStringSubmatch(url.Path)
		if len(match) == 0 {
			continue
		}

		// Replace path groups to query
		for i, name := range r.RegExp.SubexpNames() {
			if i != 0 && name != "" {
				url.RawQuery += "&" + name + "=" + match[i]
			}
		}

		// Set redirect
		url.Path = r.Url
		return
	}
}
