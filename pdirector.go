// Port Director - a small TCP proxy to redirect ports openned on local hosts, making them accessible remotely
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

const usage = `

usage: pdirector <local-port> <proxy-port> [<proxy-address>]

	local-port    - An opened port usually bounded to the localhost
	proxy-port    - A proxy port for a localhost connection, which is remotely available
	proxy-address - A specific ip or named address where proxy-port should be opened.
	                Default - 0.0.0.0
`

var PROGNAME = "pdirector"

func init() {
}

func main() {
	var localPort, proxyPort, proxyHost string

	// Parse command line flags and arguments
	flag.Parse()

	// Check required command line args
	switch len(flag.Args()) {
	case 3:
		// Read proxy address
		proxyHost = flag.Arg(2)
		fallthrough
	case 2:
		// Read local and proxy ports
		localPort, proxyPort = flag.Arg(0), flag.Arg(1)

	default:
		fmt.Fprintf(os.Stderr, "%s ERROR: Must provide both local and proxy ports", PROGNAME)
		fmt.Fprintln(os.Stderr, usage)
		return
	}

	// Start proxy listener
	proxy, err := proxyListen(proxyHost, proxyPort)
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
		go handleProxyRequest(c, localPort)
	}
}

func proxyListen(host, port string) (net.Listener, error) {
	addr := fmt.Sprintf("%s:%s", host, port)
	return net.Listen("tcp", addr)
}

func handleProxyRequest(conn net.Conn, localPort string) {
	defer conn.Close() // Always remember to close connection

	// Open connection to the local service and remember to close it
	lconn, err := localConnect(localPort)
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

func localConnect(port string) (net.Conn, error) {
	return net.Dial("tcp", "localhost:"+port)
}

func reportErrors(client net.Conn, err error) {
	// Report error to the remote client first
	errMsg := fmt.Sprintf("%s ERROR: %s\n", PROGNAME, err)
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
