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
	PID        int // Process ID
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

// TODO: denne filen blander ms og ns (det er ikke nødvendig) skal vi holde oss til en?

////////////////////////////////////////////////////////////////////////////////
// Public functions
////////////////////////////////////////////////////////////////////////////////

func WatchdogNode(watchdogport string, elevport string) {
	listener, conn := initWatchdogNode(watchdogport)
	bytebuffer := make([]byte, 500)

	msg := new(WatchdogMessage)
	msg.PID = 0
	msg.UpdateTime = time.Now()
	latestMessage = *msg

	go watchdogCheckTimeout(watchdogport, elevport)

	// TODO: klarer ikke helt å følge for løkken
	for {
		_, err := conn.Read(bytebuffer)

		if err != nil {
			conn.Close()
			listener.Close()
			listener, conn = initWatchdogNode(watchdogport)
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

// TODO: Hvor blir denne kalt?
func SenderNode(port string) {
	conn := initSenderNode(port)
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

////////////////////////////////////////////////////////////////////////////////
// Private functions
////////////////////////////////////////////////////////////////////////////////

func initSenderNode(port string) net.Conn {
	addr := connHost
	addr += port

	conn, err := net.Dial(connType, addr)
	if err != nil {
		panic(err.Error())
	}

	return conn
}

func initWatchdogNode(watchdogport string) (net.Listener, net.Conn) {
	addr := connHost
	addr += watchdogport
	fmt.Println(addr)  // TODO: nødvendig?
	listener, _ := net.Listen(connType, addr)

	conn, err := listener.Accept()
	if err != nil {
		panic(err.Error())
	}
	return listener, conn
}

// TODO: denne gjør vel mer enn bare å sjekke? endre navn eller lage en ekstra func?
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
