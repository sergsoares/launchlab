package cloudinit

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/rs/zerolog/log"
)

var baseyaml string = `#cloud-config
groups:
  - docker
users:
  - default
  # the docker service account
  - name: docker-service
    groups: docker
package_upgrade: true
packages:
  - apt-transport-https
  - ca-certificates
  - curl
  - gnupg-agent
  - software-properties-common`

var runcmd string = `# power_state:
#   mode: reboot
#   message: Restarting after installing docker & docker-compose
runcmd:
  # install docker following the guide: https://docs.docker.com/install/linux/docker-ce/ubuntu/
  - curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
  - sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
  - sudo apt-get -y update
  - sudo apt-get -y install docker-ce docker-ce-cli containerd.io
  - sudo systemctl enable docker
  # install docker-compose following the guide: https://docs.docker.com/compose/install/
  - sudo curl -L "https://github.com/docker/compose/releases/download/1.25.4/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
  - sudo chmod +x /usr/local/bin/docker-compose`

type DockerComposeConfig struct {
	Base64 string
	Raw    string
}

func GenerateCloudInit(dc DockerComposeConfig) string {
	result := fmt.Sprint(baseyaml, "\n", runcmd, `
  - echo `, dc.Base64, ` | base64 -d > /root/docker-compose.yml
  - docker-compose -f /root/docker-compose.yml up -d
`)
	log.Debug().Msg(result)
	return result
}

func GetFileAsBase64(path string) (string, error) {
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
