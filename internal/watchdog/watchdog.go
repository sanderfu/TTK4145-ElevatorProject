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
)

type WatchdogMessage struct {
	PID        int
	UpdateTime time.Time
}

var latestMessage WatchdogMessage
var connPort int = 15579

const (
	connHost = ":"
	connType = "tcp"

	updateTimeMS       = 500
	timeoutMS          = 750
	noConnectionWaitMS = 500
)

func initSenderNode() net.Conn {
	addr := connHost
	addr += strconv.Itoa(connPort)

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		panic(err.Error())
	}

	return conn
}

func initWatchdogNode() (net.Listener, net.Conn) {
	addr := connHost
	addr += strconv.Itoa(connPort)
	fmt.Println(addr)
	l, _ := net.Listen("tcp", addr)

	conn, err := l.Accept()
	if err != nil {
		panic(err.Error())
	}
	return l, conn
}

func watchdogCheckTimeout() {
	for {
		if time.Since(latestMessage.UpdateTime).Milliseconds() > timeoutMS {
			fmt.Println(time.Since(latestMessage.UpdateTime).String())
			fmt.Println("The SenderNode is not responding!")
			cmd := exec.Command("gnome-terminal", "-e", "./main")
			cmd.Run()
		}
		time.Sleep(updateTimeMS * time.Millisecond)
	}
}

func WatchdogNode() {
	l, conn := initWatchdogNode()
	bytebuffer := make([]byte, 500)

	msg := new(WatchdogMessage)
	msg.PID = 0
	msg.UpdateTime = time.Now()
	latestMessage = *msg

	go watchdogCheckTimeout()

	for {
		_, err := conn.Read(bytebuffer)

		if err != nil {
			conn.Close()
			l.Close()
			l, conn = initWatchdogNode()
		}

		// Convert bytes into Buffer (implements io.Reader/io.Writer)
		buffer := bytes.NewBuffer(bytebuffer)

		//Init a new WatchdogMessage struct

		//Create a decoder object that takes in the Buffer
		gobobj := gob.NewDecoder(buffer)

		// Decode buffer and unmarshal it into a WatchdogMessage
		gobobj.Decode(msg)

		fmt.Println("Receiver: ", msg.UpdateTime.String())
		latestMessage = *msg

	}
}

func SenderNode() {
	conn := initSenderNode()
	defer conn.Close()
	msg := new(WatchdogMessage)
	msg.PID = os.Getpid()

	for {
		time.Sleep(updateTimeMS * time.Millisecond)
		msg.UpdateTime = time.Now()
		binaryBuffer := new(bytes.Buffer)
		gobobj := gob.NewEncoder(binaryBuffer)
		fmt.Println("Sender: ", msg.UpdateTime.String())
		gobobj.Encode(msg)

		conn.Write(binaryBuffer.Bytes())
	}
}
