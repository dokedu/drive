package mail

import (
	"bytes"
	"embed"
	"html/template"
)

//go:embed templates/*
var templateFiles embed.FS

func LoginLinkMailTemplate(name string, link string) (string, error) {
	t, err := template.ParseFS(templateFiles, "templates/*.gohtml")
	if err != nil {
		return "", err
	}

	data := struct {
		Name string
		Link string
	}{
		Name: name,
		Link: link,
	}

	out := new(bytes.Buffer)
	err = t.ExecuteTemplate(out, "login_link.gohtml", data)

	if err != nil {
		return "", err
	}

	return out.String(), nil
}
