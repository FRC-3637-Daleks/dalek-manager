#!/bin/sh
### BEGIN INIT INFO
# Provides:          dalek-manager
# Required-Start:    $local_fs $network $remote_fs
# Required-Stop:     $local_fs $remote_fs
# Default-Start:     2 3 4 5
# Default-Stop:      0 1 6
# Short-Description: robot manager web app
# Description:       This file should be used to start and stop dalek-manager.
### END INIT INFO

USERNAME='lvuser'
SERVICE='dalek-manager-arm'
SPATH='/home/lvuser/dman'
OPTIONS=''
LOGPATH=$SPATH'/logs'
NOW=`date`
LOGNAME="log-$NOW"
ME=`whoami`

as_user() {
if [ "$ME" = "$USERNAME" ] ; then
 bash -c "$1"
else
 su - "$USERNAME" -c "$1"
fi
}

service_start() {
if  pgrep -u $USERNAME -f $SERVICE > /dev/null
then
	echo "$SERVICE is already running!"
else
	echo "Starting $SERVICE..."
	as_user "$SPATH/$SERVICE $OPTIONS > $LOGPATH/$LOGNAME"
	sleep 1
	if pgrep -u $USERNAME -f $SERVICE > /dev/null
	then
		echo "$SERVICE is now running."
	else
		echo "Error! Could not start $SERVICE!"
	fi
fi
}

service_stop() {
if pgrep -u $USERNAME -f $SERVICE > /dev/null
then
	echo "Stopping $SERVICE"
	killall $SERVICE
else
	echo "$SERVICE was not running."
fi
if pgrep -u $USERNAME -f $SERVICE > /dev/null
then
	echo "Error! $SERVICE could not be stopped."
else
	echo "$SERVICE is stopped."
fi
}

case "$1" in
start)
service_start
;;
stop)
service_stop
;;
restart)
service_stop
service_start
;;
*)
echo "Usage: /etc/init.d/dalek-manager {start|stop|restart}"
exit 1
;;
esac

exit 0
