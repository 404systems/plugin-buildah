package main

import (
	"codeberg.org/woodpecker-plugins/go-plugin"
	. "github.com/404systems/plugin-buildah/internal"
)

func main() {
	p := &Plugin{
		Settings: &Settings{},
	}

	p.Plugin = plugin.New(plugin.Options{
		Name:        "buildah",
		Description: "build oci images with buildah",
		Flags:       p.Flags(),
		Execute:     p.Execute,
	})

	p.Run()
}
