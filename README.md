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

## Start playing with launchlab (Postgresql + Adminer)

```bash
sergsoares@host:~/projects$ cat docker-compose.yml
# Use postgres/example user/password credentials
version: '3.1'

services:
  # Connection string
  # postgresql://postgres:postgres@localhost:5432/db?sslmode=disable
  db:
    image: postgres
    restart: always
    ports: 
        - "5432:5432"
    environment:
      POSTGRES_USER: postgres 
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: db

  adminer:
    image: adminer
    restart: always
    ports:
      - 8080:8080

# Launchlab use default file docker-compose.yml
sergsoares@host:~/projects/$ launchlab 
> Fingerprint generated from /home/sergsoares/.ssh/id_rsa.pub: 27:9d:f6:b5:2e:49:78:4e:52:8e:f8:1b:ae:47:ba:5f
> Using default userdata
> Params initialized
> Type deteted: digital ocean
> Creating Droplet
> Droplet Created: launchlab


sergsoares@host:~/projects/$ doctl compute droplet list
ID           Name         Public IPv4       Private IPv4    Public IPv6    Memory    VCPUs    Disk    Region    Image                     VPC UUID                                Status    Tags    Features              Volumes
258546822    launchlab    165.227.190.21    10.108.0.2                     1024      1        25      nyc3      Ubuntu 20.04 (LTS) x64    74c17ef7-e5fb-4525-a0b7-447740cf58cf    active            private_networking 


sergsoares@host:~/projects/$ curl 165.227.190.21:8080 -I
HTTP/1.1 200 OK
Host: 165.227.190.21:8080
Date: Sat, 07 Aug 2021 19:34:00 GMT
Connection: close
X-Powered-By: PHP/7.4.22
Set-Cookie: adminer_sid=a87958d89b03bdb5bfe11189e1e558f7; path=/; HttpOnly
Set-Cookie: adminer_key=f7826141933755f867d74857425c64e6; path=/; HttpOnly; SameSite=lax
Content-Type: text/html; charset=utf-8
Cache-Control: no-cache
X-Frame-Options: deny
X-XSS-Protection: 0
X-Content-Type-Options: nosniff
Referrer-Policy: origin-when-cross-origin
Content-Security-Policy: script-src 'self' 'unsafe-inline' 'nonce-ZjA5OTk3MmZmZmNkNTFlZjNjMGYwMTVhYmRkNjE2NWQ=' 'strict-dynamic'; connect-src 'self'; frame-src https://www.adminer.org; object-src 'none'; base-uri 'none'; form-action 'self'
```