// Port Director - a small TCP proxy for redirect and port forwarding
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

const progname = "pdirector"

func main() {
	var (
		fwdHost = flag.String("fwd-host", "localhost", "An ip or hostname with forwarded port")
		fwdPort = flag.String("fwd-port", "", "An opened port to forward traffic to")
		proxyHost = flag.String("proxy-host", "0.0.0.0", "A local ip or named address for a proxy-port")
		proxyPort = flag.String("proxy-port", "", "A proxy port for a forwarded connection")
	)

	// Parse command line flags and arguments
	flag.Parse()

	// Check required command line args
	if *fwdPort == "" || *proxyPort == "" {
		fmt.Fprintf(os.Stderr, "%s ERROR: Must provide both forward and proxy ports\n", progname)
		flag.PrintDefaults()
		return
	}

	// Start proxy listener
	proxy, err := proxyListen(*proxyHost, *proxyPort)
	if err != nil {
		log.Fatal(err)
	}
	defer proxy.Close()

	// Server loop
	for {
		// Wait for the incomming connections
		c, err := proxy.Accept()
		if err != nil {
			log.Fatal(err)
		}

		// Handle connection
		go handleProxyRequest(c, *fwdHost, *fwdPort)
	}
}

func proxyListen(host, port string) (net.Listener, error) {
	addr := fmt.Sprintf("%s:%s", host, port)
	return net.Listen("tcp", addr)
}

func handleProxyRequest(conn net.Conn, host, port string) {
	defer conn.Close() // Always remember to close connection

	// Open connection to the local service and remember to close it
	lconn, err := fwdConnect(host, port)
	if err != nil {
		// Report error to the remote client first
		reportErrors(conn, err)
		return
	}
	defer lconn.Close()

	// Do the actual proxy and report any erors to the remote client
	err = doProxy(conn, lconn)
	if err != nil {
		reportErrors(conn, err)
	}
}

func fwdConnect(host, port string) (net.Conn, error) {
	return net.Dial("tcp", host+":"+port)
}

func reportErrors(client net.Conn, err error) {
	// Report error to the remote client first
	errMsg := fmt.Sprintf("%s ERROR: %s\n", progname, err)
	// Also, put it in the server log
	fmt.Fprintln(client, errMsg)
	log.Println(errMsg)
}

func doProxy(a, b io.ReadWriter) error {
	errChan := make(chan error, 1)

	// Async copy duplex communications
	go func() {
		_, err := io.Copy(a, b)
		errChan <- err
	}()
	go func() {
		_, err := io.Copy(b, a)
		errChan <- err
	}()

	return <-errChan
}
