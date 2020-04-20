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
	PID       int
	Timestamp time.Time
}

var latestMessage WatchdogMessage

const (
	updateTimeMS = 500
	timeoutMS    = 750
)

////////////////////////////////////////////////////////////////////////////////
// Public functions
////////////////////////////////////////////////////////////////////////////////

func WatchdogNode(watchdogport string, elevport string, bcastlocalport string) {
	fmt.Println("Starting watchdog")
	listener, conn := initWatchdogNode(watchdogport)
	bytebuffer := make([]byte, 500)
	msg := new(WatchdogMessage)
	msg.PID = 0
	msg.Timestamp = time.Now()
	latestMessage = *msg

	go watchdogTimeoutHandler(watchdogport, elevport, bcastlocalport)

	for {
		_, err := conn.Read(bytebuffer)

		// If read fails, ElevatorNode is gone and WatchdogNode must be reinitialized
		if err != nil {
			conn.Close()
			listener.Close()
			listener, conn = initWatchdogNode(watchdogport)
		}

		// Convert bytes into Buffer (implements io.Reader/io.Writer)
		buffer := bytes.NewBuffer(bytebuffer)

		//Create a decoder object that takes in the Buffer
		gobDecoder := gob.NewDecoder(buffer)

		// Decode buffer and unmarshal it into a WatchdogMessage
		gobDecoder.Decode(msg)
		latestMessage = *msg
	}
}

func ElevatorNode(port string) {
	conn := initElevatorNode(port)
	defer conn.Close()
	msg := new(WatchdogMessage)
	msg.PID = os.Getpid()
	for {
		time.Sleep(updateTimeMS * time.Millisecond)
		msg.Timestamp = time.Now()
		binaryBuffer := new(bytes.Buffer)
		gobEncoder := gob.NewEncoder(binaryBuffer)
		gobEncoder.Encode(msg)
		conn.Write(binaryBuffer.Bytes())
	}
}

////////////////////////////////////////////////////////////////////////////////
// Private functions
////////////////////////////////////////////////////////////////////////////////

func initElevatorNode(port string) net.Conn {
	addr := ":" + port
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		panic(err.Error())
	}
	return conn
}

func initWatchdogNode(watchdogport string) (net.Listener, net.Conn) {
	addr := ":" + watchdogport
	listener, _ := net.Listen("tcp", addr)
	conn, err := listener.Accept()
	if err != nil {
		panic(err.Error())
	}
	return listener, conn
}

func watchdogTimeoutHandler(watchdogport string, elevport string, bcastlocalport string) {
	for {
		if time.Since(latestMessage.Timestamp).Nanoseconds()/1e6 > timeoutMS {
			fmt.Println(time.Since(latestMessage.Timestamp).String())
			fmt.Println("The ElevatorNode is not responding!")
			command := "build/elevator -lastpid " + strconv.Itoa(latestMessage.PID) + " -elevport " + elevport + " -watchdogport " + watchdogport + " -bcastlocalport " + bcastlocalport
			fmt.Println("Restarting software: ", command)
			cmd := exec.Command("gnome-terminal", "-e", command)
			cmd.Run()
		}
		time.Sleep(updateTimeMS * time.Millisecond)
	}
}
