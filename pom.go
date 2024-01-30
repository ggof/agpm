package main

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

