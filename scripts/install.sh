#!/bin/sh

BASE=/home/lvuser/dman/

if [ ! -e $BASE ]; then
	mkdir $BASE
fi
cd $BASE

echo Extracting...
mv /home/lvuser/dalek-manager.tar.gz $BASE
tar -xzf dalek-manager.tar.gz > /dev/null 2>&1

echo Making files executable...
chmod +x scripts/start_server.sh
chmod +x scripts/stop_server.sh
chmod +x dalek-manager-arm

if [ ! -e /etc/init.d/dalek-manager ]; then
    echo Installing startup script...
    su Admin -c "cp /home/lvuser/dman/scripts/dalek-manager-init.sh /etc/init.d/dalek-manager && \
    chmod 755 /etc/init.d/dalek-manager && \
    update-rc.d dalek-manager defaults"
fi
