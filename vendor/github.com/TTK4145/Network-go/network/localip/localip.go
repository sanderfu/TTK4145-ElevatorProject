package localip

import (
	"fmt"
	"net"
	"strings"
	"time"
)

var localIP string

func LocalIP() (string, error) {
	//addr := net.TCPAddr{IP: []byte{8, 8, 8, 8}, Port: 53}
	addr := "8.8.8.8:53"
	//fmt.Printf("%+v\n", addr)
	conn, err := net.DialTimeout("tcp4", addr, 250*time.Millisecond)

	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer conn.Close()
	localIP = strings.Split(conn.LocalAddr().String(), ":")[0]
	return localIP, nil
}
