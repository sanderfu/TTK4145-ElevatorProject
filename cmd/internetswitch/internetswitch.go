package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
)

var internet bool

func internetOn(on bool, networkCard string) {
	var command string
	if on {
		fmt.Println("Turning on internet")
		command = "sudo ifconfig " + networkCard + " up"
	} else {
		fmt.Println("Turning off internet")
		command = "sudo ifconfig " + networkCard + " down"
	}
	out, err := exec.Command("/bin/bash", "-c", command).Output()
	if err != nil {
		fmt.Printf("%s", err)
	}
	fmt.Println("Command Successfully Executed")
	output := string(out[:])
	fmt.Println(output)

}
func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("What is the name of the network card? (typically eno1 or enp4s0)")
	networkCard, _ := reader.ReadString('\n')
	networkCard = networkCard[:len(networkCard)-1]
	internet = true
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("Hit enter to toggle internet")
		text, _ := reader.ReadString('\n')
		fmt.Println(text)
		internet = !internet
		internetOn(internet, networkCard)
	}
}
