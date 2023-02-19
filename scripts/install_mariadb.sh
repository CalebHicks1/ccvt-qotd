#!/bin/bash
SQL_PASSWORD=$1
if [[ -z $SQL_PASSWORD ]]
then
	echo "Please Specify DB password"
	exit
fi

# Get required packages
sudo apt install -y software-properties-common
# Next, import the GPG signing key.
sudo apt-key adv --fetch-keys 'https://mariadb.org/mariadb_release_signing_key.asc'
# add the MariaDB APT repository
sudo add-apt-repository 'deb [arch=amd64,arm64,ppc64el] https://mariadb.mirror.liquidtelecom.com/repo/10.6/ubuntu focal main'
# install mariadb
sudo apt update && sudo apt install -y mariadb-server mariadb-client


sudo systemctl enable mariadb

echo "drop schema if exists qotd;create schema qotd;" | sudo mysql -u root -h localhost
echo "CREATE USER 'qotd'@'localhost' IDENTIFIED BY '${SQL_PASSWORD}'; GRANT ALL PRIVILEGES ON qotd.* TO 'qotd'@'localhost';" | sudo mysql -u root -h localhost "qotd"