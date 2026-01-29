package internal

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"codeberg.org/woodpecker-plugins/go-plugin"
	"github.com/BurntSushi/toml"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v3"
)

type Settings struct {
	Repo           string
	Registry       string
	Username       string
	Password       string
	RegistriesFile string
	AuthsFile      string
	Mirror         string
	Containerfile  string
	Context        string
	Target         string
	SkipPush       bool
	SkipBuild      bool
	AutoTag        bool
	Tags           string
	SearchRegistry string
	CacheRepo      string
	Retries        int
}

type Plugin struct {
	*plugin.Plugin
	Settings *Settings
}

func (p *Plugin) Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "oci.repo",
			Sources:     cli.EnvVars("PLUGIN_REPO"),
			Destination: &p.Settings.Repo,
		},
		&cli.StringFlag{
			Name:        "oci.registry",
			Sources:     cli.EnvVars("PLUGIN_REGISTRY"),
			Usage:       "the registry to push the image to",
			Destination: &p.Settings.Registry,
			Value:       "docker.io",
		},
		&cli.StringFlag{
			Name:        "oci.registry.username",
			Sources:     cli.EnvVars("PLUGIN_USERNAME"),
			Destination: &p.Settings.Username,
		},
		&cli.StringFlag{
			Name:        "oci.registry.password",
			Sources:     cli.EnvVars("PLUGIN_PASSWORD"),
			Destination: &p.Settings.Password,
		},
		&cli.StringFlag{
			Name:        "registries.file",
			Sources:     cli.EnvVars("REGISTRIES_FILE"),
			Destination: &p.Settings.RegistriesFile,
		},
		&cli.StringFlag{
			Name:        "auths.file",
			Sources:     cli.EnvVars("AUTHS_FILE"),
			Destination: &p.Settings.AuthsFile,
		},
		&cli.StringFlag{
			Name:        "containerfile",
			Sources:     cli.EnvVars("PLUGIN_CONTAINERFILE"),
			Destination: &p.Settings.Containerfile,
		},
		&cli.StringFlag{
			Name:        "context",
			Sources:     cli.EnvVars("PLUGIN_CONTEXT"),
			Destination: &p.Settings.Context,
		},
		&cli.StringFlag{
			Name:        "target",
			Sources:     cli.EnvVars("PLUGIN_TARGET"),
			Destination: &p.Settings.Target,
		},
		&cli.BoolFlag{
			Name:        "skip.push",
			Sources:     cli.EnvVars("PLUGIN_SKIP_PUSH"),
			Destination: &p.Settings.SkipPush,
		},
		&cli.BoolFlag{
			Name:        "skip.build",
			Sources:     cli.EnvVars("PLUGIN_SKIP_BUILD"),
			Destination: &p.Settings.SkipBuild,
		},
		&cli.StringFlag{
			Name:        "tags",
			Sources:     cli.EnvVars("PLUGIN_TAGS"),
			Destination: &p.Settings.Tags,
		},
		&cli.BoolFlag{
			Name:        "tags.auto",
			Sources:     cli.EnvVars("PLUGIN_AUTO_TAG"),
			Destination: &p.Settings.AutoTag,
		},
		&cli.StringFlag{
			Name:        "oci.searchregistry",
			Sources:     cli.EnvVars("PLUGIN_SEARCH_REGISTRY"),
			Usage:       "the registry to use for pulling images",
			Value:       "docker.io",
			Destination: &p.Settings.SearchRegistry,
		},
		&cli.StringFlag{
			Name:        "oci.searchregistry.mirror",
			Sources:     cli.EnvVars("PLUGIN_MIRROR"),
			Destination: &p.Settings.Mirror,
		},
		&cli.StringFlag{
			Name:        "cacherepo",
			Sources:     cli.EnvVars("PLUGIN_CACHE_REPO"),
			Destination: &p.Settings.CacheRepo,
		},
		&cli.IntFlag{
			Name:        "retries",
			Sources:     cli.EnvVars("PLUGIN_RETRIES"),
			Destination: &p.Settings.Retries,
			Value:       1,
		},
	}
}

type RegistriesConfig struct {
	SearchRegistries []string   `toml:"unqualified-search-registries"`
	Registry         []Registry `toml:"registry"`
}

type Registry struct {
	Location string     `toml:"location"`
	Insecure bool       `toml:"insecure"`
	Mirror   []Registry `toml:"mirror"`
}

type AuthsConfig struct {
	Auths map[string]Auth `json:"auths"`
}

type Auth struct {
	Auth string `json:"auth"`
}

func (p *Plugin) Execute(ctx context.Context) error {
	log.Debug().Msg("starting execution")
	// setup registries.conf
	cfg := RegistriesConfig{
		SearchRegistries: []string{p.Settings.SearchRegistry},
		Registry: []Registry{
			{
				Location: p.Settings.SearchRegistry,
			},
		},
	}
	log.Info().Msg(fmt.Sprintf("using %q as the default search registry", p.Settings.SearchRegistry))

	if p.Settings.Mirror != "" {
		cfg.Registry[0].Mirror = []Registry{
			{
				Location: p.Settings.Mirror,
				Insecure: true,
			},
		}
		log.Info().Msg(fmt.Sprintf("using %q as the search registry mirror", p.Settings.Mirror))
	}

	f, err := os.OpenFile(p.Settings.RegistriesFile, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to create registries.conf: %w", err)
	}

	if err := toml.NewEncoder(f).Encode(cfg); err != nil {
		return fmt.Errorf("failed to encode registries file %q: %w", p.Settings.RegistriesFile, err)
	}

	f.Close()

	//generate auth file

	authCfg := &AuthsConfig{Auths: make(map[string]Auth)}
	authFileExists := false
	if p.Settings.Registry != "" {

		if p.Settings.Username == "" {
			log.Warn().Msg("username was not set!")
		}

		if p.Settings.Password == "" {
			log.Warn().Msg("password was not set!")
		}

		authStr := p.Settings.Username + ":" + p.Settings.Password

		authEnc := base64.StdEncoding.EncodeToString([]byte(authStr))

		authCfg.Auths[p.Settings.Registry] = Auth{Auth: authEnc}

		data, err := json.MarshalIndent(authCfg, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal auth config: %w", err)
		}

		if err := os.WriteFile(p.Settings.AuthsFile, data, 0644); err != nil {
			return fmt.Errorf("failed to write auths file %q: %w", p.Settings.AuthsFile, err)
		}
		authFileExists = true
	} else {
		log.Info().Msg("no registry configured")
	}

	// set tags

	tags := strings.Split(p.Settings.Tags, ",")
	tags = append(tags, p.generateTags()...)
	if len(tags) == 0 {
		tags = append(tags, "latest")
	}

	dests := []string{}
	for _, tag := range tags {
		dests = append(dests, p.Settings.Registry+"/"+p.Settings.Repo+":"+tag)
	}

	// run build

	buildCmd := exec.Command(
		"buildah",
		"build",
		"--registries-conf="+p.Settings.RegistriesFile,
		fmt.Sprintf("--retry=%d", p.Settings.Retries),
	)
	if p.Settings.Containerfile != "" {
		buildCmd.Args = append(buildCmd.Args, "-f="+p.Settings.Containerfile)
	}
	if authFileExists {
		buildCmd.Args = append(buildCmd.Args, "--authfile="+p.Settings.AuthsFile)
	}
	if p.Settings.Target != "" {
		buildCmd.Args = append(buildCmd.Args, "--target="+p.Settings.Target)
	}
	if p.Settings.CacheRepo != "" {
		cacheUri := p.Settings.Registry + "/" + p.Settings.CacheRepo
		buildCmd.Args = append(buildCmd.Args, "--cache-from="+cacheUri)
		if !p.Settings.SkipPush {
			buildCmd.Args = append(buildCmd.Args, "--cache-to="+cacheUri)
		}
	}
	for _, d := range dests {
		buildCmd.Args = append(buildCmd.Args, "-t="+d)
	}
	// we need to set build context last
	if p.Settings.Context != "" {
		buildCmd.Args = append(buildCmd.Args, p.Settings.Context)
	} else {
		buildCmd.Args = append(buildCmd.Args, ".")
	}

	log.Info().Msg(fmt.Sprintf("build args: %s", strings.Join(buildCmd.Args, ", ")))

	if !p.Settings.SkipBuild {
		buildCmd.Stdout = os.Stdout
		buildCmd.Stderr = os.Stderr

		if err := buildCmd.Run(); err != nil {
			return fmt.Errorf("failed to build image: %w", err)
		}
	} else {
		log.Info().Msg("skipping build")
	}

	// push image

	for _, dest := range dests {
		if err := p.pushImage(dest, authFileExists); err != nil {
			return err
		}
	}

	return nil
}

func (p *Plugin) pushImage(dest string, useAuthFile bool) error {
	pushCmd := exec.Command(
		"buildah",
		"push",
	)
	if useAuthFile {
		pushCmd.Args = append(pushCmd.Args, "--authfile="+p.Settings.AuthsFile)
	}
	pushCmd.Args = append(pushCmd.Args, dest)

	log.Info().Msg(fmt.Sprintf("push args: %s", strings.Join(pushCmd.Args, ", ")))

	if p.Settings.SkipPush {
		log.Info().Msg("skipping push")
		return nil
	}

	pushCmd.Stdout = os.Stdout
	pushCmd.Stderr = os.Stderr

	if err := pushCmd.Run(); err != nil {
		return fmt.Errorf("failed to push image: %w", err)
	}

	return nil
}
