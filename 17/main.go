package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	var timeout int
	flag.IntVar(&timeout, "timeout", 10, "Connection timeout in seconds")
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s (--timeout seconds) <host> <port>\n", os.Args[0])
		os.Exit(1)
	}

	host := args[0]
	port := args[1]
	addr := net.JoinHostPort(host, port)

	conn, err := net.DialTimeout("tcp", addr, time.Duration(timeout)*time.Second)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to %s: %v\n", addr, err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Printf("Successfully connected to %s\n", addr)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer cancel()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				conn.SetReadDeadline(time.Now().Add(150 * time.Millisecond))

				buffer := make([]byte, 1024)
				n, err := conn.Read(buffer)

				if err != nil {
					if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
						continue
					}

					if err != io.EOF {
						fmt.Fprintf(os.Stderr, "Error reading from the server: %s\n", err)
					}

					fmt.Println("\nConnection closed by server")
					cancel()
					return
				}

				if n > 0 {
					os.Stdout.Write(buffer[:n])
				}
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer cancel()

		reader := bufio.NewReader(os.Stdin)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				line, err := reader.ReadBytes('\n')
				if err != nil {
					if err == io.EOF {
						fmt.Println("\nConnection closed by client (Ctrl+D)")

						if tcpConn, ok := conn.(*net.TCPConn); ok {
							tcpConn.CloseWrite()
						} else {
							conn.Close()
						}
					} else {
						fmt.Fprintf(os.Stderr, "Error reading from stdin: %s\n", err)
					}

					return
				}

				_, err = conn.Write(line)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error writing to the server: %s\n", err)
					return
				}
			}
		}
	}()

	select {
	case sig := <-sigChan:
		fmt.Printf("\nReceived signal: %s. Shutting down...\n", sig)
		cancel()
	case <-ctx.Done():
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
	}

	fmt.Println("Disconnected")
}