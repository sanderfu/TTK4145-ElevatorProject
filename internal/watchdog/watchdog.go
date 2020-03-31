package watchdog

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/sanderfu/TTK4145-ElevatorProject/internal/configuration"
)

type WatchdogMessage struct {
	PID        int
	UpdateTime time.Time
}

var latestMessage WatchdogMessage

const (
	connHost = ":"
	connType = "tcp"

	updateTimeMS       = 500
	timeoutMS          = 750
	noConnectionWaitMS = 500
)

func initSenderNode(port string) net.Conn {
	addr := connHost
	addr += port

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		panic(err.Error())
	}

	return conn
}

func initWatchdogNode(watchdogport string) (net.Listener, net.Conn) {
	addr := connHost
	addr += watchdogport
	fmt.Println(addr)
	l, _ := net.Listen("tcp", addr)

	conn, err := l.Accept()
	if err != nil {
		panic(err.Error())
	}
	return l, conn
}

func watchdogCheckTimeout(watchdogport string, elevport string) {
	for {
		if time.Since(latestMessage.UpdateTime).Nanoseconds()/1e6 > timeoutMS {
			fmt.Println(time.Since(latestMessage.UpdateTime).String())
			fmt.Println("The SenderNode is not responding!")
			command := "build/elevator -lastpid " + strconv.Itoa(latestMessage.PID) + " -elevport " + elevport + " -watchdogport " + watchdogport
			fmt.Println("Restarting software: ", command)
			cmd := exec.Command("gnome-terminal", "-e", command)
			cmd.Run()
		}
		time.Sleep(updateTimeMS * time.Millisecond)
	}
}

func WatchdogNode(watchdogport string, elevport string) {
	l, conn := initWatchdogNode(watchdogport)
	bytebuffer := make([]byte, 500)

	msg := new(WatchdogMessage)
	msg.PID = 0
	msg.UpdateTime = time.Now()
	latestMessage = *msg

	go watchdogCheckTimeout(watchdogport, elevport)

	for {
		_, err := conn.Read(bytebuffer)

		if err != nil {
			conn.Close()
			l.Close()
			l, conn = initWatchdogNode(watchdogport)
		}

		// Convert bytes into Buffer (implements io.Reader/io.Writer)
		buffer := bytes.NewBuffer(bytebuffer)

		//Init a new WatchdogMessage struct

		//Create a decoder object that takes in the Buffer
		gobobj := gob.NewDecoder(buffer)

		// Decode buffer and unmarshal it into a WatchdogMessage
		gobobj.Decode(msg)

		latestMessage = *msg

	}
}

func SenderNode() {
	conn := initSenderNode(configuration.Flags.WatchdogPort)
	defer conn.Close()
	msg := new(WatchdogMessage)
	msg.PID = os.Getpid()

	for {
		time.Sleep(updateTimeMS * time.Millisecond)
		msg.UpdateTime = time.Now()
		binaryBuffer := new(bytes.Buffer)
		gobobj := gob.NewEncoder(binaryBuffer)
		gobobj.Encode(msg)

		conn.Write(binaryBuffer.Bytes())
	}
}
