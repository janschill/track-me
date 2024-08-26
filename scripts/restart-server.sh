#!/bin/bash

echo "Running new binary test..."
~/track-me-temp/trackme --ping
if [ $? -eq 0 ]; then
    echo "New binary test passed. Proceeding with deployment."
    sudo systemctl stop trackme.service
    mv ~/track-me-temp/trackme /root/track-me/trackme
    rsync -av ~/track-me-temp/web/assets/* /root/track-me/web/assets/
    rsync -av ~/track-me-temp/web/templates/* /root/track-me/web/templates/
    sudo systemctl start trackme.service
else
    echo "New binary test failed. Aborting deployment."
    exit 1
fi
rm -rf ~/track-me-temp/*
