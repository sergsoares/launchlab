package main

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/digitalocean/godo"
	"gopkg.in/yaml.v2"
)

type Params struct {
	name              string
	dryRun            bool
	dockerComposePath string
	sshPath           string
	image             string
	userData          string
	fingerprint       string
	region            string
	size              string
}

var configurationLocation = os.Getenv("HOME") + "/.config/doctl/config.yaml"
var baseSshPath = os.Getenv("HOME") + "/.ssh/id_rsa.pub"

func main() {
	var typeb string
	var name string
	var action string
	var dryrun bool
	var dockercompose string
	var sshpath string
	var image string
	var userData string
	var region string
	var size string

	flag.StringVar(&typeb, "type", "do", "Cloud that will be used")
	flag.StringVar(&name, "name", "launchlab", "Name that will be used in Cloud Instance")
	flag.StringVar(&action, "action", "create", "Name that will be used in Cloud Instance")
	flag.StringVar(&dockercompose, "file", "docker-compose.yml", "Docker compose file to be used.")
	flag.BoolVar(&dryrun, "dry-run", false, "Dry run command to be created.")
	flag.StringVar(&sshpath, "ssh", baseSshPath, "SSH public path to be used.")
	flag.StringVar(&image, "image", "ubuntu-20-04-x64", "Imase used as base: Default ubuntu 20.04")
	flag.StringVar(&userData, "userData", "", "Default command as userdata")
	flag.StringVar(&region, "region", "nyc3", "Default command as userdata")
	flag.StringVar(&size, "size", "s-1vcpu-1gb", "Default command as userdata")

	flag.Parse()

	key, err := ioutil.ReadFile(sshpath)
	if err != nil {
		log.Fatal(err)
	}
	sshFingerPrint := getFingerPrintFromKey(string(key))
	fmt.Println("> Fingerprint generated from", sshpath, ":", sshFingerPrint)

	if userData == "" {
		fmt.Println("> Using default userdata")
		command, err := GetFileAsCommandBase64(dockercompose)
		if err != nil {
			fmt.Println("> Invalid path in dockerfile")
			os.Exit(1)
		}

		userData = getUserdataWithDockerCompose(command)
	}

	params := Params{
		name:              name,
		dryRun:            dryrun,
		dockerComposePath: dockercompose,
		sshPath:           sshpath,
		image:             image,
		userData:          userData,
		fingerprint:       sshFingerPrint,
		region:            region,
		size:              size,
	}

	fmt.Println("> Params initialized")

	switch typeb {
	case "do":
		fmt.Println("> Type deteted: digital ocean")
		launchDo(params)
	default:
		fmt.Println("Type not supported:", typeb)
	}
}

type DigitalOceanToken struct {
	AccessToken string `yaml:"access-token"`
}

func loadDoClient(path string) *godo.Client {
	token := DigitalOceanToken{}
	f, _ := os.Open(path)
	content, _ := ioutil.ReadAll(f)
	yaml.Unmarshal(content, &token)
	return godo.NewFromToken(token.AccessToken)
}

func getFingerPrintFromKey(key string) string {
	parts := strings.Fields(string(key))
	if len(parts) < 2 {
		log.Fatal("bad key")
	}

	k, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		log.Fatal(err)
	}

	fp := md5.Sum([]byte(k))
	result := ""
	for i, b := range fp {
		result += fmt.Sprintf("%02x", b)
		if i < len(fp)-1 {
			result += fmt.Sprint(":")
		}
	}

	return result
}

func getUserdataWithDockerCompose(base64compose string) string {
	return fmt.Sprintf(`#!/bin/bash
sudo apt-get update
sudo apt-get install -y docker.io docker-compose

echo %s | base64 -d > /root/docker-compose.yml
sudo docker-compose -f /root/docker-compose.yml up -d`, base64compose)
}

func launchDo(param Params) {
	client := loadDoClient(configurationLocation)

	createRequest := &godo.DropletCreateRequest{
		Name:     param.name,
		Region:   param.region,
		Size:     param.size,
		UserData: param.userData,
		SSHKeys: []godo.DropletCreateSSHKey{
			{0, param.fingerprint},
		},
		Image: godo.DropletCreateImage{
			Slug: param.image,
		},
	}
	ctx := context.TODO()

	if param.dryRun {
		fmt.Println("> Dry run Activate")
		return
	}
	fmt.Println("> Creating Droplet")

	dropletObj, _, err := client.Droplets.Create(ctx, createRequest)
	if err != nil {
		fmt.Print("Error:", err)
		os.Exit(1)
	}

	fmt.Println("> Droplet Created:", dropletObj.Name)
}

func GetFileAsCommandBase64(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("failure with path: %s", err)
	}

	content, err := ioutil.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failure with file content: %s", err)
	}

	return base64.StdEncoding.EncodeToString(content), nil
}
