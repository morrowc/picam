#!/bin/sh
#
# Collect images from the pi camera, store those in a defined
# directory, for collection by the actual client software system.
# NOTE: This is used because pi-camera integration/binding in golang
#       is 'not terrific'.
#
STORAGE=/tmp/camstore
NAME=$(date '+capture-%Y-%M-%d-%H-%m-%S.jpg')
CAM=/usr/bin/libcamera-still

if [ ! -d ${STORAGE} }; then
	mkdir -p ${STORAGE}
fi

# Loop 3 times with 15 seconds break at each loop.
for d in 1 2 3; do
  NAME=$(date '+capture-%Y-%M-%d-%H-%m-%S.jpg')
  ${CAM} -n -o ${NAME} >> /tmp/capture.log 2>&1
  sleep 15
done

