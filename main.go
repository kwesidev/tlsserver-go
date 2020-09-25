package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

var (
	serverTime        string  = time.Now().String()
	port              *int64  = flag.Int64("port", 8090, "Port number for the server")
	ipAddress         *string = flag.String("ipaddress", "127.0.0.1", "Ip Address to bind to")
	serverCertificate *string = flag.String("cert", "server.pem", "Location to the server certificate")
	serverPrivateKey  *string = flag.String("privatekey", "server_key.pem", "Location to server key")
)

// workerHandlerClients
// Process client connection
func workerHandleClients(conn net.Conn) {
	defer conn.Close()
	buffer := make([]byte, 1024)
	clientIP := conn.RemoteAddr().String()
	conn.Write([]byte(fmt.Sprintf("Welcome Kwez TCP Server and time is %s \n", serverTime)))

	log.Printf("Client %s has connected.", clientIP)

	for {
		// gets the connection and process
		count, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				log.Println("Failed to read client", err)
				break
			}
		}
		log.Printf(" Client %s says %s", clientIP, buffer[:count])
		// Closes connection if this command is written
		if strings.TrimSpace(string((buffer[:count]))) == "quit" {
			break
		}
	}

	conn.Write([]byte("\n Thanks for connecting to My Tcp Server bye !!!!!!!!!!"))
}

func main() {
	// Loads the server certificate and private key to establish a secure connection
	var connections []net.Conn
	flag.Parse()
	cert, err := tls.LoadX509KeyPair(*serverCertificate, *serverPrivateKey)
	if err != nil {

		panic(err)
	}
	cfg := &tls.Config{Certificates: []tls.Certificate{cert}}
        
	ln, err := tls.Listen("tcp", *ipAddress+":"+strconv.FormatInt(*port, 10), cfg)
	if err != nil {

		log.Printf("Failed to establish tcp server %s",serverTime)
		return
	}
	defer ln.Close()

	fmt.Println("TCP server started waiting for clients to connect ...")
	for {

		conn, err := ln.Accept()
		//skip if fails to read client
		if err != nil {
			log.Printf("Failed to connect to client")
			continue

		}
		connections = append(connections, conn)
		//handle client connections in it own goroutine
		go workerHandleClients(conn)
	}

}
