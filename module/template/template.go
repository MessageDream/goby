package template

import (
	"container/list"
	"fmt"
	"html/template"
	"runtime"
	"strings"
	"time"

	"github.com/MessageDream/goby/module/infrastructure"
)

var (
	appName, appURL, appSubURL, appDomain, appVersion string
)

func init() {
	appName = ""
	appURL = ""
	appSubURL = ""
	appDomain = ""
	appVersion = ""
}

func InitTemplate(appname, appurl, appsuburl, appdomain, appversion string) {
	appName = appname
	appURL = appurl
	appSubURL = appsuburl
	appVersion = appversion
}

func Str2html(raw string) template.HTML {
	return template.HTML(raw)
}

func Range(l int) []int {
	return make([]int, l)
}

func List(l *list.List) chan interface{} {
	e := l.Front()
	c := make(chan interface{})
	go func() {
		for e != nil {
			c <- e.Value
			e = e.Next()
		}
		close(c)
	}()
	return c
}

func ShortSha(sha1 string) string {
	if len(sha1) == 40 {
		return sha1[:10]
	}
	return sha1
}

var mailDomains = map[string]string{
	"gmail.com": "gmail.com",
}

func NewFuncMap() []template.FuncMap {
	return []template.FuncMap{map[string]interface{}{
		"GoVer": func() string {
			return strings.Title(runtime.Version())
		},
		"UseHTTPS": func() bool {
			return strings.HasPrefix(appURL, "https")
		},
		"AppName": func() string {
			return appName
		},
		"AppSubURL": func() string {
			return appSubURL
		},
		"AppURL": func() string {
			return appURL
		},
		"AppVer": func() string {
			return appVersion
		},
		"AppDomain": func() string {
			return appDomain
		},
		"Str2html":  Str2html,
		"MD5":       infrastructure.EncodeMd5,
		"TimeSince": infrastructure.TimeSince,
		"ShowFooterTemplateLoadTime": func() bool {
			return false
		},
		"LoadTimes": func(startTime time.Time) string {
			return fmt.Sprint(time.Since(startTime).Nanoseconds()/1e6) + "ms"
		},
		"ThemeColorMetaTag": func() string {
			return "#ff5343"
		},
		"FileSize": infrastructure.FileSize,
		"Subtract": infrastructure.Subtract,
		"Add": func(a, b int) int {
			return a + b
		},
		"DateFmtLong": func(t time.Time) string {
			return t.Format(time.RFC1123Z)
		},
		"DateFmtShort": func(t time.Time) string {
			return t.Format("Jan 02, 2006")
		},
		"List": List,
		"SubStr": func(str string, start, length int) string {
			if len(str) == 0 {
				return ""
			}
			end := start + length
			if length == -1 {
				end = len(str)
			}
			if len(str) < end {
				return str
			}
			return str[start:end]
		},
	}}
}

func Oauth2Icon(t int) string {
	switch t {
	case 1:
		return "fa-github-square"
	case 2:
		return "fa-google-plus-square"
	case 3:
		return "fa-twitter-square"
	case 4:
		return "fa-qq"
	case 5:
		return "fa-weibo"
	}
	return ""
}

func Oauth2Name(t int) string {
	switch t {
	case 1:
		return "GitHub"
	case 2:
		return "Google+"
	case 3:
		return "Twitter"
	case 4:
		return "腾讯 QQ"
	case 5:
		return "Weibo"
	}
	return ""
}
