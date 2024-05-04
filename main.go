package main

import (
	"io"
	"log"
	"fmt"
	"bufio"
	"bytes"
	"encoding/binary"
	"go.bug.st/serial"
)

type EISData struct {
	Header [3]byte
	Tach   uint16
	CHT    [6]uint16
	EGT    [6]uint16
	AUX5   uint16
	AUX6   uint16
	ASPD   uint16
	ALT    int16
	VOLT   uint16
	FUELF  uint16
	UNIT   uint8
	CARB   int8
	ROCSGN int8
	OATH   int8
	OILT   uint16
	OILP   uint8
	AUX1   uint16
	AUX2   uint16
	AUX3   uint16
	AUX4   uint16
	COOL   uint16
	ETI    uint16
	QTY    uint16
	HRS    uint8
	MIN    uint8
	SEC    uint8
	ENDHRS uint8
	ENDMIN uint8
	BARO   uint16
	MAGHD  uint16
	SPARE  uint8
}

func (e EISData) String() string {
	var buf bytes.Buffer

	fmt.Fprintf(&buf, "Header: %v\n", e.Header)
	fmt.Fprintf(&buf, "Tach: %d\n", e.Tach)
	fmt.Fprintf(&buf, "CHT: %v\n", e.CHT)
	fmt.Fprintf(&buf, "EGT: %v\n", e.EGT)
	fmt.Fprintf(&buf, "AUX5: %d\n", e.AUX5)
	fmt.Fprintf(&buf, "AUX6: %d\n", e.AUX6)
	fmt.Fprintf(&buf, "ASPD: %d\n", e.ASPD)
	fmt.Fprintf(&buf, "ALT: %d\n", e.ALT)
	fmt.Fprintf(&buf, "VOLT: %d\n", e.VOLT)
	fmt.Fprintf(&buf, "FUELF: %d\n", e.FUELF)
	fmt.Fprintf(&buf, "UNIT: %d\n", e.UNIT)
	fmt.Fprintf(&buf, "CARB: %d\n", e.CARB)
	fmt.Fprintf(&buf, "ROCSGN: %d\n", e.ROCSGN)
	fmt.Fprintf(&buf, "OATH: %d\n", e.OATH)
	fmt.Fprintf(&buf, "OILT: %d\n", e.OILT)
	fmt.Fprintf(&buf, "OILP: %d\n", e.OILP)
	fmt.Fprintf(&buf, "AUX1: %d\n", e.AUX1)
	fmt.Fprintf(&buf, "AUX2: %d\n", e.AUX2)
	fmt.Fprintf(&buf, "AUX3: %d\n", e.AUX3)
	fmt.Fprintf(&buf, "AUX4: %d\n", e.AUX4)
	fmt.Fprintf(&buf, "COOL: %d\n", e.COOL)
	fmt.Fprintf(&buf, "ETI: %d\n", e.ETI)
	fmt.Fprintf(&buf, "QTY: %d\n", e.QTY)
	fmt.Fprintf(&buf, "HRS: %d\n", e.HRS)
	fmt.Fprintf(&buf, "MIN: %d\n", e.MIN)
	fmt.Fprintf(&buf, "SEC: %d\n", e.SEC)
	fmt.Fprintf(&buf, "ENDHRS: %d\n", e.ENDHRS)
	fmt.Fprintf(&buf, "ENDMIN: %d\n", e.ENDMIN)
	fmt.Fprintf(&buf, "BARO: %d\n", e.BARO)
	fmt.Fprintf(&buf, "MAGHD: %d\n", e.MAGHD)
	fmt.Fprintf(&buf, "SPARE: %d\n", e.SPARE)

	return buf.String()
}


func main() {
	desiredPort := "/dev/ttyUSB0"

	// lets scan for ports
	ports, err := serial.GetPortsList()
	if err != nil {
		log.Fatal(err)
	}
	if len(ports) == 0 {
		log.Fatal("No serial ports found!")
	}
	for _, port := range ports {
		fmt.Printf("Found port: %v\n", port)
	}

	mode := &serial.Mode{
		BaudRate: 9600,
		DataBits: 8,
		Parity: serial.NoParity,
		StopBits: serial.OneStopBit,
	}

	port, err := serial.Open(desiredPort, mode)
	if err != nil {
		log.Fatal(err)
	}
	defer port.Close()

	// Create a reader for the serial data
	reader := bufio.NewReader(port)

	fmt.Printf("Synchronizing...")
	for {
		// Read bytes until we find the frame header
		for {
			b, err := reader.ReadByte()
			if err != nil {
				log.Fatal(err)
			}
			if b == 0xFE {
				fmt.Printf(" FE")
				b, err = reader.ReadByte()
				if err != nil {
					log.Fatal(err)
				}
				if b == 0xFF {
					fmt.Printf(" FF")
					b, err = reader.ReadByte()
					if err != nil {
						log.Fatal(err)
					}
					if b == 0xFE {
						fmt.Printf(" FE... synchronized!\n")
						break
					}
				} else {
					fmt.Printf(".")
				}
			} else {
				fmt.Printf(".")
			}
		}

		// Read the rest of the frame
		frame := make([]byte, 63)
		_, err = io.ReadFull(reader, frame)
		if err != nil {
			log.Fatal(err)
		}

		// Validate the checksum
		var checksum byte
		for _, b := range frame[:len(frame)-1] {
			checksum += b
		}
		checksum = ^checksum
		if checksum != frame[len(frame)-1] {
			log.Println("invalid checksum")
			continue
		}

		// Parse the frame into an EISData struct
		var data EISData
		err = binary.Read(bytes.NewReader(frame), binary.BigEndian, &data)
		if err != nil {
			log.Fatal(err)
		}

		// Do something with the data
		log.Printf("%s\n", data)
	}
}

