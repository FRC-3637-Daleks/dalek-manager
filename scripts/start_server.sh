#!/bin/sh
bash /home/lvuser/dman/scripts/stop_server.sh
cd /home/lvuser/dman/
/home/lvuser/dman/dalek-manager-arm &
echo Server started!
