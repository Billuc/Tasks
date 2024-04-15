package tasks

import (
	"bufio"
	"errors"
	"os"
	"regexp"
	"strings"
)

var templVarRegex = regexp.MustCompile(`---\w+---`)

type template struct {
	templateDir  string
	templateName string
}

type loadedTemplate struct {
	templateContent string
	paramNames      []string
}

func NewTemplate(dir string, name string) template {
	return template{
		templateDir:  dir,
		templateName: name,
	}
}

func (t *template) Load() (loadedTemplate, error) {
	file, err := os.Open(t.templateDir + "/" + t.templateName)
	if err != nil {
		return loadedTemplate{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var strBuilder strings.Builder
	params := make([]string, 0)

	for scanner.Scan() {
		line := scanner.Text()
		paramsInLine := templVarRegex.FindAllString(line, -1)

		for _, p := range paramsInLine {
			params = append(params, p[3:len(p)-3])
		}

		strBuilder.WriteString(line)
	}

	return loadedTemplate{templateContent: strBuilder.String(), paramNames: params}, nil
}

func (lt *loadedTemplate) GetEmptyParamsMap() map[string]string {
	data := make(map[string]string)

	for _, name := range lt.paramNames {
		data[name] = ""
	}

	return data
}

func (lt *loadedTemplate) Render(data map[string]string) (string, error) {
	result := lt.templateContent

	for _, name := range lt.paramNames {
		value, ok := data[name]

		if !ok {
			return "", errors.New("Missing parameter : " + name)
		}

		result = strings.ReplaceAll(result, "---"+name+"---", value)
	}

	return result, nil
}
