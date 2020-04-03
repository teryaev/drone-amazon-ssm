[![Build Status](https://cloud.drone.io/api/badges/teryaev/drone-amazon-ssm/status.svg)](https://cloud.drone.io/teryaev/drone-amazon-ssm)

A secret extension to Drone secrets extension for integration with AWS SSM Parameter store. _Please note this project requires Drone server version 1.4 or higher._

Docker image -- https://hub.docker.com/r/reptiloid666/drone-amazon-ssm

## Installation

Create a shared secret:

```console
$ openssl rand -hex 16
bea26a2221fd8090ea38720fc445eca6
```

Download and run the plugin:

```console
$ docker run -d \
  --publish=3000:3000 \
  --env=DRONE_DEBUG=true \
  --env=DRONE_SECRET=bea26a2221fd8090ea38720fc445eca6 \
  --restart=always \
  --name=secrets reptiloid666/drone-amazon-ssm
```

Update your runner configuration to include the plugin address and the shared secret.

```text
DRONE_SECRET_PLUGIN_ENDPOINT=http://1.2.3.4:3000
DRONE_SECRET_PLUGIN_TOKEN=bea26a2221fd8090ea38720fc445eca6
