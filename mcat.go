package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Comcast/gots/packet"
)

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		// read from stdin
		err := mcat(os.Stdin, os.Stdout)
		if err != nil {
			handleError(err)
		}
	} else {
		// read each input file and concat
		err := mcatFiles(flag.Args())
		if err != nil {
			handleError(err)
		}
	}
}

func handleError(err error) {
	fmt.Fprintf(os.Stderr, "mcat: error %s\n", err.Error())
	os.Exit(1)
}

func mcatFiles(files []string) error {
	for _, filename := range files {
		file, err := os.Open(filename)
		defer func() {
			if file != nil {
				file.Close()
			}
		}()

		if err != nil {
			return fmt.Errorf("Failed to open %s: %s", filename, err.Error())
		}

		err2 := mcat(file, os.Stdout)
		if err2 != nil {
			return err
		}
	}
	return nil
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
