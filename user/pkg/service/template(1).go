package service

import (
	"bytes"
	"html/template"
)

// Template ...
type Template struct {
	tm *template.Template
}

// Parse ...
func (t *Template) Parse() error {
	var err error
	t.tm, err = template.ParseFiles("./data/email-activation.html")
	if err != nil {
		return err
	}
	return nil
}

// GetActivationMessage ...
func (t Template) GetActivationMessage(url string) string {

	b := bytes.Buffer{}
	t.tm.Execute(&b, url)

	return b.String()
}
