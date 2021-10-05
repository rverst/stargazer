package main

import (
	"bufio"
	"fmt"
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

	defaultOutput      = "README.md"
	defaultFormat      = "list"
	defaultWithToc     = true
	defaultWithStars   = true
	defaultWithLicense = true
  defaultWithBtt     = false

	envUser   = "GITHUB_USER"
	envToken  = "GITHUB_TOKEN"
	envOutput = "OUTPUT_FILE"
	envFormat = "OUTPUT_FORMAT"
	envIgnore = "IGNORE_REPOS"

	envToc     = "WITH_TOC"
	envStars   = "WITH_STARS"
	envLicense = "WITH_LICENSE"
	envBttLink = "WITH_BACK_TO_TOP"
)

var (
	version = ""
	ignored []string
	env     map[string]string
)

func main() {
	if version == "" {
		version = "dev"
	}
	flaggy.SetName(appName)
	flaggy.SetDescription(appDesc)
	flaggy.SetVersion(version)

	var user, token, output, format string
	var test, wToc, wStars, wLicense, wBtt bool
	wToc, wStars, wLicense = true, true, true
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
	flaggy.Bool(&wToc, "", "with-toc", "print table of contents")
	flaggy.Bool(&wStars, "", "with-stars", "print starcount of repositories")
	flaggy.Bool(&wLicense, "", "with-license", "print license of repositories")
	flaggy.Bool(&wBtt, "", "with-back-to-top", "generate 'back to top' links for each language")

	flaggy.Parse()

	if exists(".env") {
		env = parseEnvFile(".env")
	}

	// flag > .env > environment
	if output == "" {
		output = getEnv(envOutput, defaultOutput)
	}
	if format == "" {
		format = getEnv(envFormat, defaultFormat)
	}
	if user == "" {
		user = getEnv(envUser, "")
	}
	if token == "" {
		token = getEnv(envToken, "")
	}
	if len(ignored) == 0 {
		ig := getEnv(envIgnore, "")
		sp := strings.Split(ig, ",")
		ignored = make([]string, 0)
		for _, s := range sp {
			s := strings.Trim(s, " ")
			if s != "" {
				ignored = append(ignored, s)
			}
		}
	}

	if wToc == defaultWithToc {
		v := getEnv(envToc, fmt.Sprintf("%t", defaultWithToc))
		if b, err := strconv.ParseBool(v); err == nil {
			wToc = b
		}
	}
	if wStars == defaultWithStars {
		v := getEnv(envStars, fmt.Sprintf("%t", defaultWithStars))
		if b, err := strconv.ParseBool(v); err == nil {
			wStars = b
		}
	}
	if wLicense == defaultWithLicense {
		v := getEnv(envLicense, fmt.Sprintf("%t", defaultWithLicense))
		if b, err := strconv.ParseBool(v); err == nil {
			wLicense = b
		}
	}
  if wBtt == defaultWithBtt {
    v := getEnv(envBttLink, fmt.Sprintf("%t", defaultWithBtt))
    if b, err := strconv.ParseBool(v); err == nil {
      wBtt = b
    }
  }

	if token == "" && !test {
		log.Fatal("github token is required")
	}

	if err := initTemplate(format); err != nil {
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

	err = writeList(output, stars, total, wToc, wLicense, wStars, wBtt)
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
		Stars:         1,
		Archived:      false,
		StarredAt:     time.Now(),
	}
	if !isIgnored(s.NameWithOwner) {
		stars["go"][0] = s
	}
	stars["markdown"] = make([]Star, 1)
	s = Star{
		Url:           "https://github.com/rverst/stars",
		Name:          "stars",
		NameWithOwner: "rverst/stars",
		Description:   "A list of awesome repositories I starred",
		License:       "MIT License",
		Stars:         1,
		Archived:      false,
		StarredAt:     time.Now(),
	}
	if !isIgnored(s.NameWithOwner) {
		stars["markdown"][0] = s
	}

	stars["C#"] = make([]Star, 0)
	stars["C++"] = make([]Star, 0)
	stars["C##"] = make([]Star, 0)

	stars["C#"] = append(stars["C#"], Star{
		Url:           "https://github.com/rverst/test",
		Name:          "test",
		NameWithOwner: "rverst/test",
		Description:   "",
		License:       "MIT License",
		Stars:         1,
		StarredAt:     time.Now(),
	})
	stars["C++"] = append(stars["C++"], Star{
		Url:           "https://github.com/rverst/test_2",
		Name:          "test_2",
		NameWithOwner: "rverst/test_2",
		Description:   "Some description",
		License:       "",
		Stars:         1,
		StarredAt:     time.Now(),
	})

	stars["C##"] = append(stars["C##"], Star{
		Url:           "https://github.com/rverst/test_3",
		Name:          "test_3",
		NameWithOwner: "rverst/test_3",
		Description:   "",
		License:       "",
		Stars:         1,
		StarredAt:     time.Now(),
	})

	total = 4
	return
}

// .env > environment
func getEnv(key, defVal string) string {
	val := os.Getenv(key)
	if v, ok := env[key]; ok {
		val = v
	}
	if val == "" {
		return defVal
	}
	return val
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
