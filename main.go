package main

import (
	"bufio"
	"log"
	"os"
	"strings"
	"time"

	"github.com/integrii/flaggy"
)

const (
	appName = "stargazer"
	appDesc = ""

	defaultOutput = "README.md"
	envUser       = "GITHUB_USER"
	envToken      = "ACCESS_TOKEN"
	envOutput     = "OUTPUT_FILE"
	envIgnore     = "IGNORE_REPOS"
)

var (
	version = ""
	commit  = "none"
	date    = "unknown"
	builtBy = "unknown"
	ignored []string
)

func main() {
	if version == "" {
		version = "dev"
	}
	flaggy.SetName(appName)
	flaggy.SetDescription(appDesc)
	flaggy.SetVersion(version)

	var user, token, output string
	var test bool
	flaggy.String(&output, "o", "output-file", "the file to create (default:"+defaultOutput+" )")
	flaggy.String(&user, "u", "github-user", "github user name")
	flaggy.String(&token, "", "access-token", "github access token")
	flaggy.StringSlice(&ignored, "i", "ignore", "repositories to ignore (flag can be specified multiple times)")
	flaggy.Bool(&test, "t", "test", "just put out some test data")

	flaggy.Parse()

	var env map[string]string
	if exists(".env") {
		env = parseEnvFile(".env")
	}

	// flag > .env > environment
	if output == "" {
		output = os.Getenv(envOutput)
		if v, ok := env[envOutput]; ok {
			output = v
		}
		if output == "" {
			output = defaultOutput
		}
	}
	if user == "" {
		user = os.Getenv(envUser)
		if v, ok := env[envUser]; ok {
			user = v
		}
	}
	if token == "" {
		token = os.Getenv(envToken)
		if v, ok := env[envToken]; ok {
			token = v
		}
	}
	if len(ignored) == 0 {
		ig := os.Getenv(envIgnore)
		if v, ok := env[envIgnore]; ok {
			ig = v
		}
		sp := strings.Split(ig, ",")
		ignored = make([]string, 0)
		for _, s := range sp {
			s := strings.Trim(s, " ")
			if s != "" {
				ignored = append(ignored, s)
			}
		}
	}

	var stars map[string][]Star
	var total int
	var err error
	if test {
		stars, total = testStars()
	} else {
		if stars, total, err = fetchStars(user, token); err != nil {
			log.Fatal(err)
		}
	}

	err = writeList(output, stars, total)
	if err != nil {
		log.Fatal(err)
	}
}

func isIgnored(name string) bool {
	if len(ignored) == 0 {
		return false
	}
	for _, i := range ignored {
		if strings.ToLower(i) == strings.ToLower(name) {
			return true
		}
	}
	return false
}

func testStars() (stars map[string][]Star, total int) {

	stars = make(map[string][]Star)
	stars["go"] = make([]Star, 1)
	s := Star{
		Url:           "https://github.com/rverst/stargazer",
		Name:          "stargazer",
		NameWithOwner: "rverst/stargazer",
		Description:   "Creates awesome lists of your starred repositories",
		License:       "MIT License",
		Archived:      false,
		StarredAt:     time.Now(),
	}

	if !isIgnored(s.NameWithOwner) {
		stars["go"][0] = s
	}

	total = 1
	return
}

func parseEnvFile(file string) map[string]string {

	env := make(map[string]string)
	f, err := os.Open(file)
	if err != nil {
		return env
	}
	s := bufio.NewScanner(f)
	for s.Scan() {
		if s.Err() != nil {
			continue
		}
		l := strings.Trim(s.Text(), " ")
		if l == "" {
			continue
		}
		if strings.HasPrefix(l, "#") {
			continue
		}
		s := strings.Split(l, "=")

		if len(s) == 2 {
			env[strings.TrimRight(s[0], " ")] = strings.TrimLeft(s[1], " ")
		}
	}
	return env
}

func exists(file string) bool {
	fi, err := os.Stat(file)
	if os.IsNotExist(err) {
		return false
	}
	return !fi.IsDir()
}
