# /etc/init.d/

case "$1" in
  start)
    echo "Starting dalek-manager"
    su lvuser -l -c "bash /home/lvuser/dman/scripts/start_server.sh &" || bash /home/lvuser/dman/scripts/start_server.sh &
    ;;
  stop)
    echo "Stopping dalek-manager"
    su lvuser -l -c "bash /home/lvuser/dman/scripts/stop_server.sh" || bash /home/lvuser/dman/scripts/stop_server.sh
    ;;
  *)
    echo "Usage: /etc/init.d/dalek-manager-init.sh {start|stop}"
    exit 1
    ;;
esac

exit 0
