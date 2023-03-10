# CCVT Question of the Day
This is an app to host the thursday night question of the day
Tools:

- GO
- Nginx
- MariaDB
- Docker (optional)

# Local Setup
For running without Docker
1. Build binary
```
cd src/
go mod init webserver
go mod tidy
go build
```
2. Configure nginx: This moves the domain_conf file to the sites_enabled directory
```
./install_and_configure_nginx.sh qotd.ccvt-home.com
```
3. Copy service script to `/lib/systemd/system/web-app.service`

4. Set up Mariadb
```
./install_mariadb.sh {SQL PASSWORD}
./initialize_mariadb.sh {SQL PASSWORD}
```
5. Test Mariadb connection:
```
mysql -u "qotd" "-ppassword" -h localhost qotd
```

# Docker Setup
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