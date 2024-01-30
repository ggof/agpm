package main

import (
	"errors"
	"net/url"
	"strings"
)

var RepoUrls = map[string]string{
	"maven": "https://repo.maven.apache.org/maven2",
}

type AgpmDependency struct {
	Type       string
	Repository string
	Group      string
	Artifact   string
	Version    string
}

func (a AgpmDependency) AsPath() string {
	return strings.ReplaceAll(a.Group, ".", "/") +
		strings.ReplaceAll(a.Artifact, ".", "/") +
		a.Version
}

func (a AgpmDependency) AsUrl() (string, error) {
	host, ok := RepoUrls[a.Repository]
	if !ok {
		return "", errors.New("Not found")
	}

	return url.JoinPath(host, a.AsPath())
}
