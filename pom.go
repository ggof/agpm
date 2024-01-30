package main

import (
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"net/url"
)

type PomDependency struct {
	GroupID    string `xml:"groupId"`
	ArtifactID string `xml:"artifactId"`
	Version    string `xml:"version"`
	Scope      string `xml:"scope"`
}

type PomProject struct {
	Dependencies []PomDependency `xml:"dependencies"`
}

type PomFile struct {
	Project PomProject `xml:"project"`
}

func GetDirectDependencies(dep Dependency) []Dependency {
	baseUrl, err := dep.AsUrl()
	if err != nil {
		log.Println("failed to resolve url")
		return nil
	}

  u, _ := url.JoinPath(baseUrl, dep.Pom())

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
        Group: d.GroupID,
        Artifact: d.ArtifactID,
        Version: d.Version,
        Type: "pkg",
        Repository: "maven",
      })
    }
  }

  return deps
}
