package main

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/cli/go-gh/pkg/tableprinter"
	"github.com/cli/go-gh/pkg/term"
	"gopkg.in/yaml.v3"
)

const languagesDefinitionYAMLURL string = "https://api.github.com/repos/github/linguist/contents/lib/linguist/languages.yml"

func main() {
	if len(os.Args) > 1 {
		log.Fatal("This command does not take any arguments at the moment")
	}
	res, err := http.Get(languagesDefinitionYAMLURL)
	if err != nil {
		log.Fatal(err)
	}

	responseBodyRaw, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	responseBody := map[string]any{}
	if err := json.Unmarshal(responseBodyRaw, &responseBody); err != nil {
		log.Fatal(err)
	}

	languagesDefinitionYAMLbytes, err := base64.StdEncoding.DecodeString(
		responseBody["content"].(string),
	)
	if err != nil {
		log.Fatal(err)
	}

	languagesDefinitions := map[string]map[string]any{}
	if err := yaml.Unmarshal(
		languagesDefinitionYAMLbytes, &languagesDefinitions,
	); err != nil {
		log.Fatal(err)
	}

	termwidth, _, err := term.FromEnv().Size()
	if err != nil {
		log.Fatal(err)
	}

	tp := tableprinter.New(os.Stdout, term.IsTerminal(os.Stdout), termwidth)
	tp.AddField("LANGUAGE")
	tp.AddField("TYPE")
	tp.AddField("EXTENSIONS")
	tp.EndRow()

	for lang, def := range languagesDefinitions {
		tp.AddField(lang)
		tp.AddField(def["type"].(string))
		tp.AddField(func() string {
			if def["extensions"] == nil {
				return ""
			}

			extsRaw := def["extensions"].([]any)
			if len(extsRaw) > 0 {
				exts := []string{}
				for _, ext := range extsRaw {
					exts = append(exts, ext.(string))
				}
				return strings.Join(exts, ",")
			}
			return ""
		}())
		tp.EndRow()
	}

	if err := tp.Render(); err != nil {
		log.Fatal(err)
	}
}
