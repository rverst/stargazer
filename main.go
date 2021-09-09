package main

import (
	"bufio"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/integrii/flaggy"
)

const (
	appName = "stargazer"
	appDesc = ""

	defaultOutput = "README.md"
	defaultFormat = "list"

	envUser   = "GITHUB_USER"
	envToken  = "GITHUB_TOKEN"
	envOutput = "OUTPUT_FILE"
	envFormat = "OUTPUT_FORMAT"
	envIgnore = "IGNORE_REPOS"

	envStars   = "WITH_STARS"
	envLicense = "WITH_LICENSE"
)

var (
	version = ""
	ignored []string
)

func main() {
	if version == "" {
		version = "dev"
	}
	flaggy.SetName(appName)
	flaggy.SetDescription(appDesc)
	flaggy.SetVersion(version)

	var user, token, output, format string
	var test, wStars, wLicense bool
	wStars, wLicense = true, true
	flaggy.String(&output, "o", "output-file", "the file to create (default:"+defaultOutput+" )")
	flaggy.String(
		&format,
		"f",
		"output-format",
		"the format of the output ["+strings.Join(availableFormats, ", ")+"] (default:"+defaultFormat+" )",
	)
	flaggy.String(&user, "u", "github-user", "github user name")
	flaggy.String(&token, "", "github-token", "github access token")
	flaggy.StringSlice(&ignored, "i", "ignore", "repositories to ignore (flag can be specified multiple times)")
	flaggy.Bool(&test, "t", "test", "just put out some test data")
	flaggy.Bool(&wStars, "", "with-stars", "print starcount of repositories")
	flaggy.Bool(&wLicense, "", "with-license", "print license of repositories")

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
	if format == "" {
		format = os.Getenv(envFormat)
		if v, ok := env[envFormat]; ok {
			format = v
		}
		if format == "" {
			format = defaultFormat
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

	var e string
	var ok bool
	if e, ok = env[envStars]; !ok {
		e = os.Getenv(envStars)
	}
	if b, err := strconv.ParseBool(e); err == nil {
		wStars = b
	}

	if e, ok = env[envLicense]; !ok {
		e = os.Getenv(envLicense)
	}
	if b, err := strconv.ParseBool(e); err == nil {
		wLicense = b
	}

	if token == "" {
		log.Fatal("github token is required")
	}

	if err := initTemplate(TemplateType(format)); err != nil {
		log.Fatal(err)
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

	for k, v := range stars {
		sort.Slice(v, func(i, j int) bool {
			return strings.ToLower(v[i].NameWithOwner) < strings.ToLower(v[j].NameWithOwner)
		})
		stars[k] = v
	}

	err = writeList(output, stars, total, wLicense, wStars)
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
