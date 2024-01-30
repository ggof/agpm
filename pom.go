package main

import (
	"log"
	"net/http"
)

type Dependency struct {
	GroupID    string `xml:"groupId"`
	ArtifactID string `xml:"artifactId"`
	Version    string `xml:"version"`
	Scope      string `xml:"scope"`
}

type Project struct {
	Dependencies []Dependency `xml:"dependencies"`
}

type PomFile struct {
	Project Project `xml:"project"`
}

func GetDependenciesFor(dep AgpmDependency) []AgpmDependency {
	u, err := dep.AsUrl()
	if err != nil {
		log.Println("failed to resolve url")
		return nil
	}
	res, err := http.Get(u)
	defer res.Body.Close()
	if err != nil {
		log.Println("failed to fetch dependencies")
		return nil
	}
}
