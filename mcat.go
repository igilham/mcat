package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/Comcast/gots/packet"
)

var output = flag.String("o", "", "Output file (default: stdout)")
var loop = flag.Bool("l", false, "Loop a single input")

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
	} else if *loop {
		// loop a single input file
		err := mcatLoop(flag.Arg(0), outputFile)
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

type config struct {
	loop bool
}

// mcatLoop reads from input and starts again when it reaches
// the end of the file
func mcatLoop(filename string, output *os.File) error {
	file, err := os.Open(filename)
	defer func() {
		if file != nil {
			file.Close()
		}
	}()

	if err != nil {
		return fmt.Errorf("Failed to open %s: %s", filename, err.Error())
	}

	for {
		err := mcat(file, output)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			_, err2 := file.Seek(0, 0)
			if err2 != nil {
				return fmt.Errorf("error returning to start of file %s: %s", filename, err2.Error())
			}
			continue
		}
		return fmt.Errorf("error reading from %s: %s", filename, err.Error())
	}
}

// mcatFiles is a thin wrapper around mcat that loops through
// a list of files, calling mcat with each as an input
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

// mcat syncs to the first TS packet in input then reads packets
// continuously while writing them to output.
// input must be opened in a readable way. output must be opened
// for writing. The caller has responsibility for closing both.
func mcat(input *os.File, output *os.File) error {
	reader := bufio.NewReader(input)

	// sync to the first packet header
	_, err := packet.Sync(reader)
	if err != nil {
		return fmt.Errorf("error finding sync byte: %s", err.Error())
	}

	var pkt packet.Packet

	for {
		if _, err := io.ReadFull(reader, pkt[:]); err != nil {
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				return err
			}
			return fmt.Errorf("error reading from %s: %s", input.Name(), err.Error())
		}

		_, err2 := output.Write(pkt[:])
		if err2 != nil {
			return fmt.Errorf("error writing output to %s: %s", output.Name(), err2.Error())
		}
	}

	return nil
}
