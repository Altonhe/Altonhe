package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	log "unknwon.dev/clog/v2"
)

type project struct {
	Icon        string   `yaml:"icon"`
	Name        string   `yaml:"name"`
	Link        string   `yaml:"link"`
	Description string   `yaml:"desc"`
	Tags        []string `yaml:"tags"`
}

type profile struct {
	Projects []project `yaml:"projects"`
}

func main() {
	defer log.Stop()
	err := log.NewConsole()
	if err != nil {
		panic(err)
	}

	profileBytes, err := os.ReadFile("profile.yml")
	if err != nil {
		log.Fatal("Failed to read profile.yml: %v", err)
	}

	var profile profile
	err = yaml.Unmarshal(profileBytes, &profile)
	if err != nil {
		log.Fatal("Failed to unmarshal profile: %v", err)
	}

	readmeBytes, err := os.ReadFile("README_template.md")
	if err != nil {
		log.Fatal("Failed to read README template: %v", err)
	}

	projectsMarkdown := makeProjectMarkdown(profile.Projects)
	readmeBytes = bytes.ReplaceAll(readmeBytes, []byte("{{PROJECTS}}"), []byte(projectsMarkdown))

	tarotsMarkdown := getRandomTarot()
	readmeBytes = bytes.ReplaceAll(readmeBytes, []byte("{{TAROTS}}"), []byte(tarotsMarkdown))

	err = os.WriteFile("README.md", readmeBytes, 0644)
	if err != nil {
		log.Fatal("Failed to write README.md: %v", err)
	}

}

func makeProjectMarkdown(projects []project) string {
	var projectMarkdown string
	for _, project := range projects {
		name := project.Name
		if name == "" {
			name = path.Base(project.Link)
		}

		var tagMarkdown string
		tags := project.Tags
		if len(tags) != 0 {
			tagMarkdown += "/"
			for _, tag := range tags {
				tagMarkdown += fmt.Sprintf(" `%s`", tag)
			}
		}

		var starMarkdown string
		if strings.HasPrefix(project.Link, "https://github.com/") {
			log.Trace("Fetch %q star counts...", name)
			starCount, err := getRepoStarCount(project.Link)
			if err != nil {
				log.Error("Failed to repo's star count: %v", err)
			} else if starCount != 0 {
				starMarkdown = fmt.Sprintf("/ [★%d](%s/stargazers)", starCount, project.Link)
			}
		}

		// - 🔮 [Elaina](https://github.com/wuhan005/Elaina) - Docker-based remote code runner / [★1](https://github.com/wuhan005/Elaina/stargazers) `Docker`
		projectMarkdown += fmt.Sprintf("- %s [%s](%s) - %s %s %s\n",
			project.Icon, name, project.Link, project.Description,
			starMarkdown, tagMarkdown)
	}

	return projectMarkdown
}

func getRepoStarCount(link string) (int64, error) {
	link = strings.ReplaceAll(link, "https://github.com/", "https://api.github.com/repos/")

	resp, err := http.Get(link)
	if err != nil {
		return 0, errors.Wrap(err, "request GitHub API")
	}
	defer resp.Body.Close()

	type repoMeta struct {
		StargazersCount int64 `json:"stargazers_count"`
	}

	var meta repoMeta
	err = json.NewDecoder(resp.Body).Decode(&meta)
	if err != nil {
		return 0, errors.Wrap(err, "unmarshal")
	}
	return meta.StargazersCount, nil
}

func getRandomTarot() string {
	tarots := []string{
		"cups01.jpg",
		"cups02.jpg",
		"cups03.jpg",
		"cups04.jpg",
		"cups05.jpg",
		"cups06.jpg",
		"cups07.jpg",
		"cups08.jpg",
		"cups09.jpg",
		"cups10.jpg",
		"cups11.jpg",
		"cups12.jpg",
		"cups13.jpg",
		"cups14.jpg",
		"maj00.jpg",
		"maj01.jpg",
		"maj02.jpg",
		"maj03.jpg",
		"maj04.jpg",
		"maj05.jpg",
		"maj06.jpg",
		"maj07.jpg",
		"maj08.jpg",
		"maj09.jpg",
		"maj10.jpg",
		"maj11.jpg",
		"maj12.jpg",
		"maj13.jpg",
		"maj14.jpg",
		"maj15.jpg",
		"maj16.jpg",
		"maj17.jpg",
		"maj18.jpg",
		"maj19.jpg",
		"maj20.jpg",
		"maj21.jpg",
		"pents01.jpg",
		"pents02.jpg",
		"pents03.jpg",
		"pents04.jpg",
		"pents05.jpg",
		"pents06.jpg",
		"pents07.jpg",
		"pents08.jpg",
		"pents09.jpg",
		"pents10.jpg",
		"pents11.jpg",
		"pents12.jpg",
		"pents13.jpg",
		"pents14.jpg",
		"swords01.jpg",
		"swords02.jpg",
		"swords03.jpg",
		"swords04.jpg",
		"swords05.jpg",
		"swords06.jpg",
		"swords07.jpg",
		"swords08.jpg",
		"swords09.jpg",
		"swords10.jpg",
		"swords11.jpg",
		"swords12.jpg",
		"swords13.jpg",
		"swords14.jpg",
		"wands01.jpg",
		"wands02.jpg",
		"wands03.jpg",
		"wands04.jpg",
		"wands05.jpg",
		"wands06.jpg",
		"wands07.jpg",
		"wands08.jpg",
		"wands09.jpg",
		"wands10.jpg",
		"wands11.jpg",
		"wands12.jpg",
		"wands13.jpg",
		"wands14.jpg",
	}
	shuffle(tarots)
	var tarotsMarkdown string
	reverse := `style="transform: rotate(180deg);"`
	for _, tarot := range tarots[:3] {
		if rand.Int()%2 != 0 {
			tarotsMarkdown += fmt.Sprintf("<img src=\"https://raw.githubusercontent.com/Altonhe/Altonhe/master/tarot/%v\" width=\"25%s\" />", tarot, "%")
		} else {
			tarotsMarkdown += fmt.Sprintf("<img %v src=\"https://raw.githubusercontent.com/Altonhe/Altonhe/master/tarot/%v\" width=\"25%s\" />", reverse, tarot, "%")
		}

	}
	return tarotsMarkdown
}

func shuffle(slice []string) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for len(slice) > 0 {
		n := len(slice)
		randIndex := r.Intn(n)
		slice[n-1], slice[randIndex] = slice[randIndex], slice[n-1]
		slice = slice[:n-1]
	}
}
