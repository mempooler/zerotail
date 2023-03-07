package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/hpcloud/tail"
	"github.com/mempooler/zerolog"
)

func main() {
	filename := flag.String("f", "", "file to tail")
	lines := flag.Int("n", 10, "number of lines to tail")
	flag.Parse()

	if *filename == "" {
		fmt.Println("usage: tail -f <filename>")
		os.Exit(1)
	}

	t, err := tail.TailFile(*filename, tail.Config{
		Location: &tail.SeekInfo{
			Offset: -n,
			Whence: os.SEEK_END,
		},
		Follow: true,
	})
	if err != nil {
		panic(err)
	}

	w := zerolog.ConsoleWriter{Out: os.Stderr}
	for line := range t.Lines {
		w.Write([]byte(line.Text))
	}
}
