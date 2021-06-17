package main

import (
	"bufio"
	"github.com/integrii/flaggy"
	"log"
	"os"
	"strings"
)

const (
	appName = "stargazer"
	appDesc = ""

	envUser  = "GITHUB_USER"
	envToken = "ACCESS_TOKEN"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
	builtBy = "unknown"
)

func main() {
	flaggy.SetName(appName)
	flaggy.SetDescription(appDesc)
	flaggy.SetVersion(version)

	var user, token string
	flaggy.String(&user, "u", "github-user", "github user name")
	flaggy.String(&token, "", "access-token", "github access token")

	flaggy.Parse()

	var env map[string]string
	if exists(".env") {
		env = parseEnvFile(".env")
	}

	// flag > .env > environment
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

	stars, total,  err := fetchStars(user, token)
	if err != nil {
		log.Fatal(err)
	}

	err = writeReadme(stars, total)
	if err != nil {
		log.Fatal(err)
	}
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
