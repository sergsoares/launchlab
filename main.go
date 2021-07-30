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
	name string
}

func main() {
	var typeb string
	var name string
	var action string
	flag.StringVar(&typeb, "type", "do", "Cloud that will be used")
	flag.StringVar(&name, "name", "launchlab", "Name that will be used in Cloud Instance")
	flag.StringVar(&action, "action", "create", "Name that will be used in Cloud Instance")
	flag.Parse()
	log.Debug().Msg("Args parsed")

	param := Params{
		name: name,
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

func launchDo(param Params) {
	token := DigitalOceanToken{}
	f, _ := os.Open(configurationLocation)
	content, _ := ioutil.ReadAll(f)
	yaml.Unmarshal(content, &token)

	base64Content, err := cloudinit.GetFileAsBase64(string(content))
	if err != nil {
		panic(err)
	}

	dc := cloudinit.DockerComposeConfig{
		Base64: base64Content,
		Raw:    base64Content,
	}

	client := godo.NewFromToken(token.AccessToken)
	createRequest := &godo.DropletCreateRequest{
		Name:     param.name,
		Region:   "nyc3",
		Size:     "s-1vcpu-1gb",
		UserData: cloudinit.GenerateDockerCompose(dc),
		SSHKeys: []godo.DropletCreateSSHKey{
			{0, "43:7d:f6:a5:2e:15:78:4e:58:8a:f8:1a:ae:47:bf:5f"},
		},
		Image: godo.DropletCreateImage{
			Slug: "ubuntu-20-04-x64",
		},
	}
	ctx := context.TODO()

	newDroplet, _, err := client.Droplets.Create(ctx, createRequest)

	if err != nil {
		fmt.Print("Error:", err)
		os.Exit(1)
	}

	fmt.Println("Droplet Created!!!")
	fmt.Println(newDroplet)
}
