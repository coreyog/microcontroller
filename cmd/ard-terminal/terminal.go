package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/coreyog/board-discovery"
	"github.com/coreyog/microcontroller"
)

var scanner = bufio.NewReader(os.Stdin)

func main() {
	devices, _ := discovery.DiscoverNow(true, false)

	targetDevice := ""
	switch len(devices) {
	case 0:
		fmt.Println("no devices discovered")
		return
	case 1:
		targetDevice = devices[0]
		fmt.Printf("device discovered: %s\n", targetDevice)
	default:
		choice, err := chooseOption("Please choose a device:", devices)
		if err != nil {
			fmt.Printf("unable to receive choice: %s", err)
			panic(err)
		}

		targetDevice = devices[choice]
	}

	fmt.Printf("opening connection to device at %s...", targetDevice)

	ard, err := microcontroller.NewArduino(targetDevice, 9600)
	if err != nil {
		fmt.Println("ERROR")
		fmt.Printf("unable to connect to device at %s: %s\n", targetDevice, err)
		panic(err)
	}

	time.Sleep(time.Second * 2)
	fmt.Println("SUCCESS\nctrl-c to stop")

	for {
		// Read-Eval-Print Loop = REPL

		// prompt and read
		fmt.Printf("> ")
		line, err := scanner.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("unable to read input: %s\n", err)
			panic(err)
		}

		// trim
		line = strings.TrimSpace(line)
		fmt.Printf("-> %s : {%s}\n", spacedHex([]byte(line)), line)

		resp, err := ard.Request([]byte(line))
		if err != nil {
			fmt.Printf("unable to receive response: %s\n", err)
			panic(err)
		}

		if bytes.HasSuffix(resp, []byte("\n")) {
			resp = resp[:len(resp)-1]
		}

		if len(resp) == 0 {
			fmt.Println("<- <nil response>")
		} else {
			fmt.Printf("<- %s : {%s}\n", spacedHex(resp), string(resp))
		}
	}
}

func chooseOption(prompt string, options []string) (choice int, err error) {
	invalid := true
	for invalid {
		invalid = false

		fmt.Println(prompt)
		fmt.Println()

		for i, o := range options {
			fmt.Printf("%d) %s\n", i, o)
		}

		fmt.Printf("Choice [1-%d]: ", len(options))
		line, err := scanner.ReadString('\n')
		if err != nil {
			return -1, err
		}

		choice, err = strconv.Atoi(line)
		if err != nil {
			fmt.Printf("invalid choice: %s\n", err)
			invalid = true
		}

		choice-- // the list started counting at 1

		if choice < 0 || choice >= len(options) {
			fmt.Printf("choice outside of range")
			invalid = true
		}
	}

	return choice, nil
}

func spacedHex(data []byte) string {
	sb := strings.Builder{}
	for i, b := range data {
		h := fmt.Sprintf("%x", b)
		if len(h) < 2 {
			h = "0" + h
		}

		sb.WriteString(h)

		if i != len(data)-1 {
			sb.WriteRune(' ')
		}
	}

	return sb.String()
}
