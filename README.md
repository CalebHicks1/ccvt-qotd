# Caleb's Webserver
This is the start of my personal webserver. 
Uses:
- docker compose
    - nginx
    - webapp

# Setup
Create `.env` file with the following values:
```
# Random key to encrypt your session
SESSION-KEY="<key-to-encrypt-session>"
# Set to true if you are running on a local machine and not a container
RUN-LOCAL="false"
```
## Install
1. Install docker: 
https://www.digitalocean.com/community/tutorials/how-to-install-and-use-docker-on-ubuntu-18-04
2. Install docker-compose
## Run
1. Run nginx proxy: 
```
docker-compose -f nginx-proxy-compose.yaml up -d
```
2. Run webapp and force a rebuild if the source code was changed:
```
docker-compose up -d --build
```

```
sudo service web-app restart
```

# Credits
Using the gorilla/mux toolkit
https://github.com/gorilla/mux#static-files

# Tools
- nginx
- systemd
- certbot
- go executable