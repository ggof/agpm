package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"

	"gopkg.in/yaml.v3"
)

type YamlConfig struct {
	Version      string            `yaml:"version"`
	Toolchain    int               `yaml:"toolchain"`
	Classpath    string            `yaml:"classpath"`
	Repositories map[string]string `yaml:"repositories"`
	Dependencies []string          `yaml:"dependencies"`
	Src          string            `yaml:"src"`
}

const (
	cfgName = "module.yaml"
)

var DefaultConfig = &YamlConfig{
	Version:      "1.9.22",
	Toolchain:    21,
	Classpath:    "cp/",
	Src:          "src/",
	Repositories: map[string]string{"maven": "https://repo.maven.apache.org/maven2"},
}

func main() {
	loadConfig()

	for _, d := range DefaultConfig.Dependencies {
		if err := Download(Parse(d)); err != nil {
			panic(err)
		}
	}
}

func loadConfig() {
	cfgFile, err := os.Open(cfgName)
	defer cfgFile.Close()

	if err != nil {
		log.Println("Config file not found, using default config as-is")
		return
	}
	bytes, err := io.ReadAll(cfgFile)
	if err != nil {
		log.Printf("Failed to read the config file: %s", err.Error())
		log.Println("using default config as-is")
		return
	}

	cfg := &YamlConfig{}
	if err := yaml.Unmarshal(bytes, cfg); err != nil {
		log.Printf("Failed to parse the config file: %s", err.Error())
		log.Println("using default config as-is")
		return
	}

	if cfg.Toolchain > 0 {
		DefaultConfig.Toolchain = cfg.Toolchain
	}

	if len(cfg.Src) > 0 {
		DefaultConfig.Src = cfg.Src
	}

	DefaultConfig.Dependencies = append(DefaultConfig.Dependencies, cfg.Dependencies...)
	for k, v := range cfg.Repositories {
		if _, ok := DefaultConfig.Repositories[k]; !ok {
			DefaultConfig.Repositories[k] = v
		}
	}
}

func Parse(s string) Dependency {
	r := regexp.MustCompile(`([^:]+):([^\/]+)\/([^\/]+)\/([^@]+)@(.+)`)
	matches := r.FindStringSubmatch(s)
	if matches == nil {
		panic("malformed dependency " + s)
	}

	matches = matches[1:]

	return Dependency{
		Type:       matches[0],
		Repository: matches[1],
		Group:      matches[2],
		Artifact:   matches[3],
		Version:    matches[4],
	}
}

func Download(d Dependency) error {
	// download dependency
	baseUrl, err := d.AsUrl()
	if err != nil {
		return err
	}

	jarUrl, _ := url.JoinPath(baseUrl, d.Jar())

	log.Printf("downloading jar from url %s", jarUrl)

	res, err := http.Get(jarUrl)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	os.Mkdir(DefaultConfig.Classpath, 0777)
	p := path.Join(DefaultConfig.Classpath, d.Jar())
	f, err := os.OpenFile(p, os.O_CREATE|os.O_RDWR, 0664)
	if err != nil {
		return err
	}

	_, err = io.Copy(f, res.Body)
	if err != nil {
		return err
	}

	f.Close()

	// find all direct dependencies
	for _, dd := range d.DirectDependencies() {
		if err := Download(dd); err != nil {
			return err
		}
	}

	return nil
}
