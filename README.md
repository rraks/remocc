# remocc
Remote command and control framework
Manage remote devices residing on private networks

<p align="center">
<img src="https://github.com/rraks/remocc/blob/master/docs/logo.png">
<img docs/logo.png >
</p>

*Still under development, open to collaborations*

## Overview
It's often difficult to manage a remote device (on LTE/Private networks). remocc provides 
a mechanism to access a shell to the remote device and in addition allows you to run apps (remote cron tasks) on your devices.
<p align="center">
<img src="https://github.com/rraks/remocc/blob/master/docs/overview.png">
</p>

## How? Reverse SSH Tunnels
<p align="center">
<img src="https://github.com/rraks/remocc/blob/master/docs/reverse.png" >
</p>


## Running it
### Production 
1. Bring up db and web containers \
`docker-compose up web db`
2. Use test device if necessary \
`docker-compose up device`
3. Exec and run test device \
`docker exec -it remocc_device_1 /bin/bash `\
`CGO_ENABLED=0 go build ./` \
`./device`

### Development
1. Bring up db and webdev containers \
`docker-compose up webdev db`
2. Connect to the docker instance
`docker exec -it remocc_webdev_1 /bin/bash`
3. Start the sshd daemon 
`/docker-entrypoint-initdb.d/init.sh`
3. Compile and run \
` CGO_ENABLED=0 go build ./` \
`./remocc`
4. Exec and run test device \
`docker exec -it remocc_device_1 /bin/bash `\
`CGO_ENABLED=0 go build ./` \
`./device`
