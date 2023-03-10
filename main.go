package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/hpcloud/tail"
	"github.com/mempooler/zerolog"
)

func main() {
	f := flag.String("f", "", "file")
	n := flag.String("n", "10", "number of lines")
	debug := flag.Bool("debug", true, "log level debug")
	trace := flag.Bool("trace", true, "log level trace")
	flag.Parse()
	if *f == "" {
		fmt.Println("usage: zerotail -f <filename> [-n <number of lines>] [-debug] [-trace]")
		os.Exit(1)
	}
	t, err := tail.TailFile(*f, tail.Config{
		Location: &tail.SeekInfo{Whence: getWhence(*n)},
		Follow:   true,
	})
	if err != nil {
		panic(err)
	}
	w := zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC3339Nano,
		Level:      getLevel(*trace, *debug),
	}
	for line := range t.Lines {
		w.Write([]byte(line.Text))
	}
}

func getWhence(n string) int {
	if strings.HasPrefix(n, "+") {
		return os.SEEK_CUR
	}
	return os.SEEK_END
}

func getLevel(trace, debug bool) zerolog.Level {
	if trace {
		return zerolog.TraceLevel
	} else if debug {
		return zerolog.DebugLevel
	} else {
		return zerolog.InfoLevel
	}
}
