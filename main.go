package main

import (
	"codeberg.org/woodpecker-plugins/go-plugin"
	. "github.com/404systems/plugin-buildah/internal"
)

var (
	version string
	commit  string
)

func getVersion() string {
	var v = "unknown"
	var c = "unknown"
	if version != "" {
		v = version
	}
	if commit != "" {
		c = commit
	}
	return v + "-" + c
}

func main() {
	p := &Plugin{
		Settings: &Settings{},
	}

	p.Plugin = plugin.New(plugin.Options{
		Name:        "plugin-buildah",
		Description: "build oci images with buildah",
		Flags:       p.Flags(),
		Execute:     p.Execute,
		Version:     getVersion(),
	})

	p.Run()
}
