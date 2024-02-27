package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/micmonay/keybd_event"
	"go.bug.st/serial"
)

var buttonsToKeys = map[int]int{
	52: 49,
	7:  17,
	9:  18,
	5:  50,
	50: 20,
	6:  21,
	8:  22,
	2:  23,
	68: 24,
	4:  59,
	3:  60,
	14: 61,
	43: 32,
	64: 33,
	65: 34,
	46: 35,
	44: 36,
	42: 37,
	69: 38,
	67: 44,
	63: 45,
	40: 46,
	49: 47,
	22: 48,
}

var boundingsForKeys = map[int]keybd_event.KeyBonding{}

func setKeyForkB(kb *keybd_event.KeyBonding, button int) int {
	switch button {
	case 69:
		kb.SetKeys(keybd_event.VK_SP12)
		kb.HasSHIFT(true)
	case 67:
		kb.SetKeys(keybd_event.VK_7)
		kb.HasSHIFT(true)
	case 63:
		kb.SetKeys(keybd_event.VK_SP5)
		kb.HasSHIFT(true)
	case 40:
		kb.SetKeys(keybd_event.VK_SP11)
	case 38:
		kb.SetKeys(keybd_event.VK_7)
	case 39:
		kb.SetKeys(keybd_event.VK_8)
	case 66:
		kb.SetKeys(keybd_event.VK_9)
	case 48:
		kb.SetKeys(keybd_event.VK_KPPLUS)
	case 62:
		kb.SetKeys(keybd_event.VK_6)
	case 10:
		kb.SetKeys(keybd_event.VK_1)
	case 28:
		kb.SetKeys(keybd_event.VK_2)
	case 26:
		kb.SetKeys(keybd_event.VK_3)
	case 27:
		kb.SetKeys(keybd_event.VK_Q)
		kb.HasALTGR(true)
	case 25:
		kb.SetKeys(keybd_event.VK_0)
	case 18:
		kb.SetKeys(keybd_event.VK_DOT)
	case 24:
		kb.SetKeys(keybd_event.VK_ENTER)
	// her starts the arrow pad
	case 16:
		kb.SetKeys(keybd_event.VK_BACKSPACE)
	case 17:
		kb.HasCTRL(true)
		return 1
	case 19:
		kb.SetKeys(keybd_event.VK_UP)
	case 51:
		kb.HasSHIFT(true)
		return 1
	case 20:
		kb.SetKeys(keybd_event.VK_LEFT)
	case 15:
		kb.SetKeys(keybd_event.VK_DOWN)
	case 21:
		kb.SetKeys(keybd_event.VK_RIGHT)
	default:
		key, ok := buttonsToKeys[button]
		if !ok {
			return 0
		}
		kb.HasCTRL(true)
		kb.HasSHIFT(true)
		kb.SetKeys(key)
	}
	return 0
}

func convertSerialToKeystrokes(serialData []byte) {
	if len(serialData) != 3 {
		return
	}
	pressing := serialData[0] == 144
	button := int(serialData[1])

	if pressing {
		kb, err := keybd_event.NewKeyBonding()
		if err != nil {
			fmt.Println(err)
			return
		}
		res := setKeyForkB(&kb, button)
		if res == 0 {
			fmt.Println("Pressing", button)
			kb.Launching()
		} else {
			kb.Press()
			boundingsForKeys[button] = kb
			fmt.Println("Holding", button)
		}
	} else {
		if kb, ok := boundingsForKeys[button]; ok {
			kb.Release()
			fmt.Println("Releasing", ok)
		}
	}
}

func main() {
	ports, err := serial.GetPortsList()
	if err != nil {
		fmt.Println(err)
	}
	for i, port := range ports {
		fmt.Println(i, " => ", port)
	}
	fmt.Println("Enter the port number: ")

	//read console for input
	reader := bufio.NewReader(os.Stdin)
	portNumberString, _ := reader.ReadString('\n')

	portNumber, err := strconv.Atoi(portNumberString[:1])
	if err != nil {
		fmt.Println(err)
		return
	}

	//check if the port number is valid
	if portNumber < 0 || portNumber >= len(ports) {
		fmt.Println("Invalid port number")
		return
	}

	portName := ports[portNumber]

	fmt.Println("Using port: ", portName, " at 9600bps N81")

	// Open the first serial port detected at 9600bps N81
	mode := &serial.Mode{
		BaudRate: 9600,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}
	s, err := serial.Open(portName, mode)

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Listening....")
	for {
		serialReader := bufio.NewReader(s)
		reply, err := serialReader.ReadBytes(0xff)
		if err != nil {
			panic(err)
		}
		convertSerialToKeystrokes(reply)
	}

	// s.Close()
}
