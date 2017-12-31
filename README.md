# VPR
Minimalistic Point-to-Point VPN

## USAGE

```
Usage of ./vpr:
  -dev="vpr0": tun device name
  -key="": key used for authentication with length of 16/24/32 bytes
  -local="": local address of the tun interface: e.g. 10.0.1.83/24
  -loglevel="info": set the loglevel
  -port=2821: UDP port
  -remote="": Remote server (external) IP like 87.13.2.44
  -route="": adds a local route to a remote network routed through the TUN device
```

## SETUP
Goal: We want to establish a secure tunnel between two (public) endpoints: *left* and *right*. Consider the following network plan 

```
+-------------------------------------------------------------------------------------------------------------------+
|                                                                                                                   |
|                                                                                                                   |
|                                         192.168.122.0/24                                                          |
|                    +------+                                    +--------+                                         |
|   192.168.122.3/24 |      |  +-------------------------------+ |        |  192.168.122.254/24                     |
|        10.0.0.1/24 | left +--+    ENCRYPTED IP-IN-UDP        +-+  right |  10.0.1.1/24                            |
|                    |      |  +-------------------------------+ |        |                                         |
|                    +------+                                    +--------+                                         |
|                        |                                           |                                              |
|        +----------------------+                                 +------------------------+                        |
|        |   left internal net  |    10.0.0.0/24                  |   right internal net   |    10.0.1.0/24         |
|        +----------------------+                                 +------------------------+                        |
|                |                                                         |                                        |
|             +---------------+                                    +------------------+                             |
|             |               |                                    |                  |                             |
|             |   left node X |   10.0.0.10/24                     |   right node X   |     10.0.1.10/24            |
|             |               |                                    |                  |                             |
|             |               |                                    |                  |                             |
|             +---------------+                                    +------------------+                             |
|                                                                                                                   |
+-------------------------------------------------------------------------------------------------------------------+
```

We can start this application on the left node like this:
`./vpr -local 10.0.0.1/24 -key 1234123412341234 -remote 192.168.122.3 -route 10.0.1.0/24 -loglevel debug`


Similarly on the right node:
`./vpr -local 10.0.1.1/24 -key 1234123412341234 -remote 192.168.122.254 -route 10.0.0.0/24 -loglevel debug`


And then we can start pinging *from* a right node:
`$ ping 10.0.0.10`

Or a left node:
`$ ping 10.0.1.10`

## Docker
You can build&run a container with the following commands:
```
sudo docker build -t test/vpr .
sudo docker run -it -e REMOTE=87.14.33.12 --cap-add=NET_ADMIN --device=/dev/net/tun test/vpr
```