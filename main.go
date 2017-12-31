package main

import (
	"fmt"
	"os"
	"os/exec"

	flag "github.com/namsral/flag"

	log "github.com/Sirupsen/logrus"

	"github.com/songgao/water"
)

const (
	BUFFERSIZE = 1500
	MTU        = "1300"
)

var (
	localIP  = flag.String("local", "", "local address of the tun interface: e.g. 10.0.1.83/24")
	remoteIP = flag.String("remote", "", "Remote server (external) IP like 87.13.2.44")
	key      = flag.String("key", "", "key used for authentication with length of 16/24/32 bytes")
	route    = flag.String("route", "", "adds a local route to a remote network routed through the VPN")
	devName  = flag.String("dev", "vpr0", "tun device name")
	port     = flag.Int("port", 2821, "UDP port")
	loglevel = flag.String("loglevel", "info", "set the loglevel")
)

func main() {
	flag.Parse()
	var err error
	if "" == *localIP {
		flag.Usage()
		log.Fatalln("\nlocal ip is not specified")
	}
	if "" == *remoteIP {
		flag.Usage()
		log.Fatalln("\nremote server is not specified")
	}
	if "" == *key {
		flag.Usage()
		log.Fatalln("\nkey not specified")
	}
	lvl, err := log.ParseLevel(*loglevel)
	if err != nil {
		log.Fatalln(err)
	}
	log.SetLevel(lvl)

	tunCfg := water.Config{
		DeviceType: water.TUN,
	}
	tunCfg.Name = *devName
	iface, err := water.New(tunCfg)
	if nil != err {
		log.Fatalf("Unable to create TUN: %s", err)
	}
	log.Printf("created TUN: %s", iface.Name())
	execIP("link", "set", "dev", iface.Name(), "mtu", MTU)
	execIP("addr", "add", *localIP, "dev", iface.Name())
	execIP("link", "set", "dev", iface.Name(), "up")
	if "" != *route {
		execIP("route", "add", *route, "dev", iface.Name())
	}

	conn, err := newConnection(*port, []byte(*key))
	if err != nil {
		log.Fatalln("unable to create connection", err)
	}
	defer conn.Close()
	go func() {
		log.Infof("starting udp read loop at :%d", *port)
		buf := make([]byte, BUFFERSIZE)
		for {
			n, err := conn.Read(buf)
			if err != nil {
				log.Fatalln("error reading:", err)
				continue
			}
			if err != nil || n == 0 {
				fmt.Println("Error: ", err)
				continue
			}
			log.Debugf("fwd packet(%d) to tun", n)
			iface.Write(buf[:n])
		}
	}()
	packet := make([]byte, BUFFERSIZE)
	log.Infof("starting tun read loop on %s", iface.Name())
	for {
		n, err := iface.Read(packet)
		if err != nil {
			break
		}
		log.Debugf("fwd packet(%d) to udp", n)
		_, err = conn.Write(packet[:n])
		if err != nil {
			fmt.Println("Error writing to conn: ", err)
		}
	}
}

func execIP(args ...string) {
	cmd := exec.Command("/sbin/ip", args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if nil != err {
		log.Fatalf("Error running /sbin/ip: %s, %v", err, args)
	}
}
