#! /bin/sh

### BEGIN INIT INFO
# Provides:     blinky
# Required-Start: $remote_fs $syslog
# Required-Stop:  $remote_fs $syslog
# Default-Start:  2 3 4 5
# Default-Stop:
# Short-Description:	Runs a BlinkM over I2C
### END INIT INFO

. /lib/lsb/init-functions

MAIN=/home/pi/Desktop/blinky/files/blinky.run
PIDFILE=/var/run/blinky/blinky.pid

case "$1" in
    start)
	echo -n "Starting blinky "
	if start-stop-daemon --start --quiet --oknodo --pidfile $PIDFILE --exec $MAIN; then
	    log_end_msg 0 || true
	else
	    log_end_msg 1 || true
	fi
	;;
    stop)
	echo -n "Stopping blinky "
	if start-stop-daemon --stop --quiet --oknodo --make-pidfile --pidfile $PIDFILE; then
	    log_end_msg 0 || true
	else
	    log_end_msg 1 || true
	fi
	;;
    *)
	echo "Usage: /etc/init.d/blinky {start|stop}"
	exit 1
	;;
esac

exit 0
