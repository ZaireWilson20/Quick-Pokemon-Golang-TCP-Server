package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strings"

	pokeload "example.com/zaire/pokedexrunner"
)

var poke = pokeload.Load("Pokedex.csv")

//Run Server
func main() {

	var addr string
	var network string
	flag.StringVar(&addr, "e", ":32581", "service endpoint [ip addr or socket path]")
	flag.StringVar(&network, "n", "tcp", "network protocol [tcp,unix]")
	flag.Parse()

	// validat supported network protocols
	switch network {
	case "tcp", "tcp4", "tcp6", "unix":
	default:
		log.Fatalln("unsupported network protocol:", network)
	}

	ln, err := net.Listen(network, addr)
	if err != nil {
		log.Fatal("failed to create listener:", err)
	}

	defer ln.Close()
	log.Println("*** Find PoKemon by Type ***")
	log.Printf("Service started: (%s) %s\n", network, addr)

	for {
		log.Println(addr)

		conn, err := ln.Accept()

		if err != nil {

			fmt.Println(err)
			if err := conn.Close(); err != nil {
				log.Println("failed to close listener:", err)

			}
			continue
		} else {

		}
		log.Println("Connect to", conn.RemoteAddr())

		go handleConnection(conn)
	}

}

func handleConnection(conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			log.Println("error closing conneciton:", err)
		}
	}()

	if _, err := conn.Write([]byte("Connecte...\nUsage:GET <pokemon name or type>")); err != nil {
		log.Println("error writing:", err)
		return
	}

	for {
		cmdLine := make([]byte, (1024 * 4))
		n, err := conn.Read(cmdLine)
		if n == 0 || err != nil {
			log.Println("connection read error:", err)
			return
		}
		cmd, param := parseCommand(string(cmdLine[0:n]))
		if cmd == "" {
			if _, err := conn.Write([]byte("Invalid command\n")); err != nil {
				log.Println("failed to write:", err)
				return
			}
			continue
		}

		switch strings.ToUpper(cmd) {
		case "GET":
			result := pokeload.Find(poke, param)
			if len(result) == 0 {
				if _, err := conn.Write([]byte("Nothing found\n")); err != nil {
					log.Println("failed to write:", err)
				}
				continue
			}

			for _, pok := range result {
				_, err := conn.Write([]byte(fmt.Sprintf(
					"%s %s\n",
					pok.Name, pok.Type),
				))

				if err != nil {
					log.Println("failed to write response:", err)
					return
				}
			}
		default:
			if _, err := conn.Write([]byte("Invalid command\n")); err != nil {
				log.Println("failed to write:", err)
				return
			}

		}

	}

}

func parseCommand(cmdLine string) (cmd, param string) {
	parts := strings.Split(cmdLine, " ")
	if len(parts) != 2 {
		return "", ""
	}

	cmd = strings.TrimSpace(parts[0])
	param = strings.TrimSpace(parts[1])
	return
}
