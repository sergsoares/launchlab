package cloudinit

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

var baseyaml string = `#cloud-config
groups:
  - docker
package_upgrade: true
packages:
  - apt-transport-https
  - ca-certificates
  - curl
  - gnupg-agent
  - software-properties-common`

var runcmd string = `runcmd:
  # install docker following the guide: https://docs.docker.com/install/linux/docker-ce/ubuntu/
  - curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
  - sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
  - sudo apt-get -y update
  - sudo apt-get -y install docker-ce docker-ce-cli containerd.io
  - sudo systemctl enable docker
  # install docker-compose following the guide: https://docs.docker.com/compose/install/
  - sudo curl -L "https://github.com/docker/compose/releases/download/1.25.4/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
  - sudo chmod +x /usr/local/bin/docker-compose`

type CloudInitConfig struct {
	CommandBase64 string
	Raw           string
	Users         string
}

func GenerateCloudInit(dc CloudInitConfig) string {
	result := fmt.Sprint(baseyaml, "\n", runcmd, `
  - echo `, dc.CommandBase64, ` | base64 -d > /root/docker-compose.yml
  - docker-compose -f /root/docker-compose.yml up -d
`, dc.Users)

	return result
}

type sshKeys string

type userCloudInit struct {
	Name              string    `yaml:"name"`
	SshAuthorizedKeys []sshKeys `yaml:"ssh-authorized-keys"`
	Sudo              string    `yaml:"sudo"`
	Groups            string    `yaml:"groups"`
	Shell             string    `yaml:"shell"`
}
type usersCloudInit struct {
	Users []userCloudInit `yaml:"users"`
}

func GetConfiguredUser(path string) string {
	f, _ := os.Open(path)
	fileContent, _ := ioutil.ReadAll(f)

	users := usersCloudInit{
		Users: []userCloudInit{
			{
				Name:              "launchlab",
				SshAuthorizedKeys: []sshKeys{sshKeys(fileContent)},
				Sudo:              "['ALL=(ALL) NOPASSWD:ALL']",
				Groups:            "sudo",
				Shell:             "/bin/bash",
			},
		},
	}
	content, err := yaml.Marshal(users)
	if err != nil {
		fmt.Print(err)
	}

	return string(content)
}

func GetFileAsCommandBase64(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Errorf("Failure with path: %s", err)
	}

	content, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Errorf("Failure with file content: %s", err)
	}

	return base64.StdEncoding.EncodeToString(content), nil
}
