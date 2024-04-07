package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/zenlor/blackfriday-steam"
)

func isStdinReady() bool {
	stat, _ := os.Stdin.Stat()
	return (stat.Mode() & os.ModeCharDevice) == 0
}

func main() {
	inputFilename := flag.String("i", "-", "input filename (default stdin)")
	outputFilename := flag.String("o", "-", "output filename (default stdout)")
	flag.Parse()

	if !isStdinReady() && *inputFilename == "-" {
		fmt.Println("Not enough arguments.")
		flag.Usage()
		return
	}

	var err error
	var inputData []byte
	if *inputFilename == "-" {
		buf := make([]byte, 0, 4*1024)
		r := bufio.NewReader(os.Stdin)
		for {
			n, err := r.Read(buf[:cap(buf)])
			buf = buf[:n]
			if n == 0 {
				if err == nil {
					continue
				}

				if err == io.EOF {
					break
				}

				slog.Error("Error reading stdin", slog.Any("error", err))
				return
			}

			inputData = append(inputData, buf...)
		}
	} else {
		inputData, err = os.ReadFile(*inputFilename)
		if err != nil {
			slog.Error("Error reading input file", slog.Any("error", err))
			return
		}
	}

	var output *bufio.Writer
	if *outputFilename == "-" {
		output = bufio.NewWriter(os.Stdout)
	} else {
		fd, err := os.Open(*inputFilename)
		if err != nil {
			slog.Error("Error opening output file", slog.Any("error", err))
			return
		}
		defer fd.Close()

		output = bufio.NewWriter(fd)
	}

	data := steam.Run(inputData)
	if _, err = output.Write(data); err != nil {
		slog.Error("error writing chunk", slog.Any("error", err))
		return
	}
	if err = output.Flush(); err != nil {
		slog.Error("error writing chunk", slog.Any("error", err))
		return
	}
}
