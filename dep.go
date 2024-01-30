package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

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
	host, ok := DefaultConfig.Repositories[d.Repository]
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

func (d Dependency) DirectDependencies() []Dependency {
	baseUrl, err := d.AsUrl()
	if err != nil {
		log.Println("failed to resolve url")
		return nil
	}

	u, _ := url.JoinPath(baseUrl, d.Pom())

	log.Printf("fetching url %s", u)

	res, err := http.Get(u)
	if err != nil {
		log.Println("failed to fetch dependencies")
		return nil
	}

	defer res.Body.Close()

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("failed to read response body")
		return nil
	}

	file := &PomFile{}
	if err := xml.Unmarshal(bytes, file); err != nil {
		log.Println("failed to unmarsharl xml")
		log.Println(err.Error())
		return nil
	}

	var deps []Dependency

	for _, d := range file.Project.Dependencies {
		if d.Scope != "provided" {
			deps = append(deps, Dependency{
				Group:      d.GroupID,
				Artifact:   d.ArtifactID,
				Version:    d.Version,
				Type:       "pkg",
				Repository: "maven",
			})
		}
	}

	return deps
}
