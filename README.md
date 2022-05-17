# picam
Raspberry Pi Camera image collector and processing server

# Introduction

Collect images from a pi-camera on the pi, send those images
to a remote server for storage. Maximum resolution for a
Logitech QuickCam Pro is 800 x 600.

# Dependencies

Software to install includes:

   * [golang](https://golang.org)
   * [grpc libraries](https://google.golang.org/grpc)

Hardware required:

   * [raspberry pi](https://www.raspberrypi.org) - Tested on at least pi4.
   * [usb webcam](https://amzn.com/dp/B00006LIOM) - Tested, but probably any usb cam.
