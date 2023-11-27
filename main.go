package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	setupSignalHandler()
	startServer()
}

func setupSignalHandler() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch)

	go func() {
		for {
			s := <-ch
			signalHandler(s)
		}
	}()
}

func signalHandler(signal os.Signal) {
	if signal == syscall.SIGTERM {
		fmt.Println("Got kill signal. ")
		fmt.Println("Program will terminate now.")
		os.Exit(0)
	} else if signal == syscall.SIGINT {
		fmt.Println("Got CTRL+C signal")
		fmt.Println("Closing.")
		os.Exit(0)
	} else if signal == syscall.SIGURG {
		//fmt.Println("SIGURG: ", signal)
	} else {
		fmt.Println("Ignoring signal: ", signal)
	}
}

func startServer() {
	l, err := net.Listen("tcp", "0.0.0.0:7777")
	if err != nil {
		fmt.Println("Can't listen on 0.0.0.0:7777")
		os.Exit(1)
	}

	defer l.Close()
	fmt.Println("Listening on 127.0.0.1:7777 ... ")

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			continue
		}

		fmt.Printf("New connection: %+v\n", conn.RemoteAddr())

		c := Connection{conn: conn}
		if err := c.Init(); err != nil {
			fmt.Println("Error during socket init: ", err.Error())
			continue
		}
		// TODO: monitor active connection number to limit it
		go c.Handle()

	}
}
