# Launchlab

App to launch docker-compose apps in a simple way in Digital Ocean.

When you can use launchlab:

- Launch a Database for validate how your app behavior with high latency.
- Play with infrastructure services through a Docker compose file.
- Pair programming pet projects with friends.

## Install

1. [Installing doctl](https://github.com/digitalocean/doctl#installing-doctl)
1. [Using `doctl auth init` to configure default credentails](https://ldocs.digitalocean.com/reference/doctl/reference/auth/init/) 
1. Get last release in [Launchlab releases](https://github.com/sergsoares/launchlab/releases/)

## Start playing with launchlab 

```bash
# Minimal command using digital ocean config.
launchlab -file docker-compose.yml
```