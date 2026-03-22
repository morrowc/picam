# picam

Raspberry Pi Camera image collector and processing server

## Introduction

Collect images from a pi-camera on the pi, send those images
to a remote server for storage. Maximum resolution for a
Logitech QuickCam Pro is 800 x 600, the camera software used,
however, is not limited to the QuickCam Pro.

To enable the oled screen, enable I2C in raspi-config.

## Architecture

The system consists of two primary components:

1. **Client**: Runs on the Raspberry Pi. A shell script (`imager.sh`) captures images periodically. A Go binary (`client_main/client`) watches the local storage directory and streams new images to the remote server over gRPC.
1. **Server**: Runs on a central machine. It accepts incoming images from authenticated clients and saves them to a designated storage directory per client.

## Dependencies

Software to install includes:

* [golang](https://golang.org)
* [grpc libraries](https://google.golang.org/grpc)

For the OLED display (optional):

* apt-get install build-essential python3-pip python3-dev python3-smbus python3-pil
* pip3 install Adafruit-Blinka adafruit-circuitpython-ssd1306

Hardware required:

* [raspberry pi](https://www.raspberrypi.org) - Tested on at least pi4.
* [usb webcam](https://amzn.com/dp/B00006LIOM) - Tested, but probably any usb cam.

## Server Operations

The server expects a text protobuf configuration file (`config.textpb`) that defines its listening port and the client mappings.

**Example `config.textpb`:**

```textpb
port: 9987
client:{
  id: "abc12345"
  store: "/tmp/store1"
}
```

**Running the server:**

```shell
go run server/main.go -config /path/to/config.textpb
```

The server process typically runs at system startup from a systemd file.

## Client Operations

The client watches a local directory and uploads new files to the server.

**Starting the client:**

```shell
go run client_main/main.go -id "abc12345" -store /tmp/camstore -server docs.as701.net:9987
```

Actuation of the image capture is done through cron. For example:

```shell
  # m h  dom mon dow   command
  # Run the image collection script.
  * * * * * /home/pi/scripts/git/picam/client/imager.sh >> /tmp/capture.log 2>&1
  # Remove stale/old images from the collection bin (/tmp/camstore)
  1 * * * * find /tmp/camstore -type f -ctime +1 -exec rm {} \; >> /tmp/cleanup.log 2>&1
```

For OLED display usage, simply make a cron entry to update the screen:

```shell
0 * * * *  /home/pi/scripts/git/picam/client/oled.py >> /tmp/oled.log 2>&1
```

## Systemd and Unique Client IDs

The `client_main` typically runs at system startup via a systemd script. Every client device communicating with the server must have a unique ID that matches the server configuration.

To preserve client privacy while maintaining a unique identifier per device, we recommend two simple approaches for your client-side systemd service (`img_client.service`):

### Option 1: Static Environment File (Recommended)

Generate a unique UUID at provisioning time and store it in `/etc/default/picam-client`:

```shell
echo "PICAM_CLIENT_ID=$(uuidgen)" | sudo tee /etc/default/picam-client
```

Then reference it in your systemd service file:

```ini
[Service]
EnvironmentFile=-/etc/default/picam-client
ExecStart=/home/pi/scripts/git/picam/client_main/client -id ${PICAM_CLIENT_ID} -store /tmp/camstore -server docs.as701.net:9987
```

### Option 2: Dynamically Derived ID

If you prefer not to manage static files on each client, you can derive a hashed identifier from the machine ID dynamically on execution using `ExecStartPre`:

```ini
[Service]
ExecStartPre=/bin/sh -c 'echo PICAM_CLIENT_ID=$(cat /etc/machine-id | sha256sum | head -c 16) > /run/picam-client.env'
EnvironmentFile=/run/picam-client.env
ExecStart=/home/pi/scripts/git/picam/client_main/client -id ${PICAM_CLIENT_ID} -store /tmp/camstore -server docs.as701.net:9987
```

## Be good
