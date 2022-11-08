# picam
Raspberry Pi Camera image collector and processing server

# Introduction

Collect images from a pi-camera on the pi, send those images
to a remote server for storage. Maximum resolution for a
Logitech QuickCam Pro is 800 x 600, the camera software used,
however, is not limited to the QuickCam Pro.

To enable the oled screen, enable I2C in raspi-config.

# Dependencies

Software to install includes:

   * [golang](https://golang.org)
   * [grpc libraries](https://google.golang.org/grpc)
   * python libraries
      * apt-get install build-essential python3-pip python3-dev python3-smbus
      * apt-get install python3-pil
      * pip3 install Adafruit-Blinka
      * pip3 install adafruit-circuitpython-ssd1306

Hardware required:

   * [raspberry pi](https://www.raspberrypi.org) - Tested on at least pi4.
   * [usb webcam](https://amzn.com/dp/B00006LIOM) - Tested, but probably any usb cam.

# Normal Operations

Golang support to query the pi camera is not terrific, so for now, poll the camera with:
   
   * imager.sh

and collect images from the storage location with the client-server:

   * client_main

Actuation of the imager.sh is done through cron:

```shell
  # m h  dom mon dow   command
  # Run the image collection script.
  * * * * * /home/pi/scripts/git/picam/client/imager.sh >> /tmp/capture.log 2>&1
  # Remove stale/old images from the collection bin (/tmp/camstore)
  1 * * * * find /tmp/camstore -type f -ctime +1 -exec rm {} \; >> /tmp/cleanup.log 2>&1 
```

The client_main runs at system startup from a systemd file, included in the client_main.
The server process similarly runs at system startup on the server, using a systemd file.
