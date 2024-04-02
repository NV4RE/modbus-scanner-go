package main

import (
	"encoding/json"
	"fmt"
	"github.com/goburrow/modbus"
	"github.com/manifoldco/promptui"
	"os"
	"strconv"
	"strings"
	"time"
)

type Reader struct {
	Name string
	Read func(addr uint16, len uint16) (results []byte, err error)
}

func haltOnError(err error) {
	if err != nil {
		fmt.Printf("Error: %v\n", err)

		for {
			time.Sleep(5 * time.Second)
		}
	}
}

func main() {
	prompt := promptui.Prompt{
		Label: "Port: COM3, /dev/ttyUSB0, /dev/ttyS0",
	}

	port, err := prompt.Run()
	haltOnError(err)

	prompt = promptui.Prompt{
		Label: "Baud Rate: 9600, 115200",
	}

	baudRate, err := prompt.Run()
	haltOnError(err)

	br, err := strconv.Atoi(baudRate)
	haltOnError(err)

	prompt = promptui.Prompt{
		Label: "Client Address: 1",
	}

	slaveId, err := prompt.Run()
	haltOnError(err)

	si, err := strconv.Atoi(slaveId)
	haltOnError(err)

	handler := modbus.NewRTUClientHandler(port)
	handler.BaudRate = br
	handler.DataBits = 8
	handler.Parity = "N"
	handler.StopBits = 1
	handler.SlaveId = uint8(si)
	handler.Timeout = 1 * time.Second

	err = handler.Connect()
	haltOnError(err)

	defer handler.Close()

	client := modbus.NewClient(handler)

	readers := []Reader{
		{
			Name: "ReadCoils",
			Read: client.ReadCoils,
		},
		{
			Name: "ReadDiscreteInputs",
			Read: client.ReadDiscreteInputs,
		},
		{
			Name: "ReadHoldingRegisters",
			Read: client.ReadHoldingRegisters,
		},
		{
			Name: "ReadInputRegisters",
			Read: client.ReadInputRegisters,
		},
	}

	output := map[string]map[uint16][]byte{}
	errors := map[string]string{}

	fmt.Println("Scanning ")
	var addr uint16
	for _, reader := range readers {
		fmt.Printf("Scanning %s\n", reader.Name)
		output[reader.Name] = map[uint16][]byte{}

		for addr = 0; addr < 65535; addr++ {
			fmt.Print(".")
			results, err := reader.Read(addr, 1)
			if err != nil {
				// if error have message "illegal" skip
				if strings.Contains(err.Error(), "illegal") {
					continue
				}

				fmt.Printf("\nRead Error:[%s] 0x%04x | %v\n", reader.Name, addr, err)
				errors[fmt.Sprintf("%s:0x%04x", reader.Name, addr)] = err.Error()
			}

			fmt.Printf("\nFound:[%s] 0x%04x: %d\n", reader.Name, addr, results)
			output[reader.Name][addr] = results
		}
	}

	b, _ := json.Marshal(output)
	err = os.WriteFile("./output.json", b, 0644)
	haltOnError(err)

	b, _ = json.Marshal(errors)
	err = os.WriteFile("./errors.json", b, 0644)
	haltOnError(err)

	fmt.Println("Done, results written to output.json")
	for {
		time.Sleep(5 * time.Second)
	}
}
