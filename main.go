package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"tidy/cloudinit"

	"github.com/digitalocean/godo"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
)

type Params struct {
	name              string
	dryRun            bool
	dockerComposePath string
}

func main() {
	var typeb string
	var name string
	var action string
	var dryrun bool
	var dockercompose string
	flag.StringVar(&typeb, "type", "do", "Cloud that will be used")
	flag.StringVar(&name, "name", "launchlab", "Name that will be used in Cloud Instance")
	flag.StringVar(&action, "action", "create", "Name that will be used in Cloud Instance")
	flag.StringVar(&dockercompose, "file", "cloudinit/examples/elasticsearch.yml", "Docker compose file to be used.")
	flag.BoolVar(&dryrun, "dry-run", false, "Dry run command to be created.")

	flag.Parse()
	log.Debug().Msg("Args parsed")

	param := Params{
		name:              name,
		dryRun:            dryrun,
		dockerComposePath: dockercompose,
	}

	switch typeb {
	case "do":
		launchDo(param)
	default:
		fmt.Println("Type not supported:", typeb)
	}
}

type DigitalOceanToken struct {
	AccessToken string `yaml:"access-token"`
}

var configurationLocation = os.Getenv("HOME") + "/.config/doctl/config.yaml"

func loadDoClient(path string) *godo.Client {
	token := DigitalOceanToken{}
	f, _ := os.Open(path)
	content, _ := ioutil.ReadAll(f)
	yaml.Unmarshal(content, &token)
	return godo.NewFromToken(token.AccessToken)
}

func launchDo(param Params) {
	client := loadDoClient(configurationLocation)

	base64Content, err := cloudinit.GetFileAsBase64(param.dockerComposePath)
	if err != nil {
		panic(err)
	}

	dc := cloudinit.CloudInitConfig{
		Base64: base64Content,
		Raw:    base64Content,
	}

	createRequest := &godo.DropletCreateRequest{
		Name:     param.name,
		Region:   "nyc3",
		Size:     "s-1vcpu-1gb",
		UserData: cloudinit.GenerateCloudInit(dc),
		SSHKeys: []godo.DropletCreateSSHKey{
			{0, "43:7d:f6:a5:2e:15:78:4e:58:8a:f8:1a:ae:47:bf:5f"},
		},
		Image: godo.DropletCreateImage{
			Slug: "ubuntu-20-04-x64",
		},
	}
	ctx := context.TODO()

	if param.dryRun {
		return
	}

	newDroplet, _, err := client.Droplets.Create(ctx, createRequest)
	if err != nil {
		fmt.Print("Error:", err)
		os.Exit(1)
	}

	fmt.Println("Droplet Created!!!")
	fmt.Println(newDroplet)
}
