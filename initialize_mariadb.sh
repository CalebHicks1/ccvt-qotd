SQL_PASSWORD=$1
if [[ -z $SQL_PASSWORD ]]
then
	echo "Please Specify DB password"
	exit
fi
mysql -u "qotd" "-p${SQL_PASSWORD}" -h localhost < init.sql