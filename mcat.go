package main

import (
	"fmt"
	"os"

	"github.com/Comcast/gots/packet"
)

func main() {
	filename := "./scenario1.ts"
	file, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open %s: %s\n", filename, err.Error())
		os.Exit(1)
	}
	er2 := mcat(file, os.Stdout)
	if er2 != nil {
		fmt.Fprintf(os.Stderr, "mcat: error %s\n", err.Error())
	}
}

func mcat(input *os.File, output *os.File) error {
	pkt := make([]byte, packet.PacketSize)

	for read, err := input.Read(pkt); read > 0 && err == nil; read, err = input.Read(pkt) {
		if err != nil {
			return fmt.Errorf("error reading from %s: %s", input.Name(), err.Error())
		}
		_, err2 := os.Stdout.Write(pkt)
		if err2 != nil {
			return fmt.Errorf("error writing output to %s: %s", output.Name(), err2.Error())
		}
	}
	return nil
}
