#!/bin/bash
DOMAIN_NAME=$1
if [[ $1 -eq 0 ]]
then
	echo "Please Specify Hostname"
	exit
fi
sudo apt install nginx -y
cp domain_conf $DOMAIN_NAME
sudo mv $DOMAIN_NAME /etc/nginx/sites-available/ 
sudo ln -s /etc/nginx/sites-available/$DOMAIN_NAME /etc/nginx/sites-enabled/$DOMAIN_NAME
sudo nginx -s reload
