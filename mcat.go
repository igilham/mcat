package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/Comcast/gots/packet"
)

var output = flag.String("o", "", "Output file (default: stdout)")

func main() {
	flag.Parse()

	// parse output file argument
	outputFile := os.Stdout
	var errOutput error
	if len(*output) > 0 {
		outputFile, errOutput = os.OpenFile(*output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if errOutput != nil {
			handleError(errOutput)
		}
	}
	defer func() {
		if outputFile != os.Stdout && outputFile != nil {
			outputFile.Close()
		}
	}()

	if flag.NArg() == 0 {
		// read from stdin
		err := mcat(os.Stdin, outputFile)
		if err != nil {
			handleError(err)
		}
	} else {
		// read each input file and concat
		err := mcatFiles(flag.Args(), outputFile)
		if err != nil {
			handleError(err)
		}
	}
}

func handleError(err error) {
	fmt.Fprintf(os.Stderr, "mcat: error %s\n", err.Error())
	os.Exit(1)
}

func mcatFiles(files []string, output *os.File) error {
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

		err2 := mcat(file, output)
		if err2 != nil {
			return err
		}
	}
	return nil
}

func mcat(input *os.File, output *os.File) error {
	reader := bufio.NewReader(input)

	// sync to the first packet header
	_, err := packet.Sync(reader)
	if err != nil {
		return fmt.Errorf("error finding sync byte: %s", err.Error())
	}

	pkt := make([]byte, packet.PacketSize)

	for read, err := reader.Read(pkt); read > 0 && err == nil; read, err = reader.Read(pkt) {
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
