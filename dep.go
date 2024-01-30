package main

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

var RepoUrls = map[string]string{
	"maven": "https://repo.maven.apache.org/maven2",
}

type Dependency struct {
	Type       string
	Repository string
	Group      string
	Artifact   string
	Version    string
}

func (a Dependency) AsPath() string {
	return strings.Join(
		[]string{
			strings.ReplaceAll(a.Group, ".", "/"),
			strings.ReplaceAll(a.Artifact, ".", "/"),
			a.Version,
		},
		"/")
}

func (d Dependency) AsUrl() (string, error) {
	host, ok := RepoUrls[d.Repository]
	if !ok {
		return "", errors.New("Not found")
	}

	return url.JoinPath(host, d.AsPath())
}

func (a Dependency) Jar() string {
	return fmt.Sprintf("%s-%s.jar", a.Artifact, a.Version)
}

func (a Dependency) Pom() string {
	return fmt.Sprintf("%s-%s.pom", a.Artifact, a.Version)
}
