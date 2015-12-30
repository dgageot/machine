package commands

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"time"

	"io/ioutil"
	"net"

	"github.com/docker/machine/libmachine"
	"github.com/docker/machine/libmachine/persist"

	"golang.org/x/crypto/ssh"
)

const daemonDefaultPort = 8080

func cmdDaemon(c CommandLine, api libmachine.API) error {
	port := c.Int("port")
	fmt.Println("Running on port", port)

	config := &ssh.ServerConfig{
		NoClientAuth: true,
	}

	// You can generate a keypair with 'ssh-keygen -t rsa'
	privateBytes, err := ioutil.ReadFile("id_rsa")
	if err != nil {
		log.Fatal("Failed to load private key (./id_rsa)")
	}

	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		log.Fatal("Failed to parse private key")
	}

	config.AddHostKey(private)

	listener, err := net.Listen("tcp", "0.0.0.0:2200")
	if err != nil {
		log.Fatalf("Failed to listen on 2200 (%s)", err)
	}

	log.Print("Listening on 2200...")
	for {
		tcpConn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept incoming connection (%s)", err)
			continue
		}

		sshConn, chans, reqs, err := ssh.NewServerConn(tcpConn, config)
		if err != nil {
			log.Printf("Failed to handshake (%s)", err)
			continue
		}

		log.Printf("New SSH connection from %s (%s)", sshConn.RemoteAddr(), sshConn.ClientVersion())
		go ssh.DiscardRequests(reqs)
		go handleChannels(api, chans)
	}
}

func handleChannels(api libmachine.API, chans <-chan ssh.NewChannel) {
	for newChannel := range chans {
		go handleChannel(api, newChannel)
	}
}

func handleChannel(api libmachine.API, newChannel ssh.NewChannel) {
	if t := newChannel.ChannelType(); t != "session" {
		newChannel.Reject(ssh.UnknownChannelType, fmt.Sprintf("unknown channel type: %s", t))
		return
	}

	connection, requests, err := newChannel.Accept()
	if err != nil {
		log.Printf("Could not accept channel (%s)", err)
		return
	}

	go func() {
		for req := range requests {
			switch req.Type {
			case "shell":
				if len(req.Payload) == 0 {
					req.Reply(true, nil)

					bash := exec.Command("bash", "-c", "echo Hello")
					output, err := bash.Output()
					if err != nil {
						log.Println("Bash error", err)
						return
					}

					connection.Write(output)
					connection.Close()
					log.Printf("Session closed")
				}
			case "subsystem":
				command := string(req.Payload[4:])
				req.Reply(true, nil)

				var output []byte
				if command == "machine/ls" {
					output, err = runLs(api)
					if err != nil {
						output = []byte("ERROR: " + err.Error())
					}
				} else {
					output = []byte("UNKNOWN")

				}

				connection.Write(output)
				connection.Close()
				log.Printf("Session closed")
			}
		}
	}()
}

func runLs(api libmachine.API) ([]byte, error) {
	stateTimeoutDuration = 10 * time.Second

	hostList, hostInError, err := persist.LoadAllHosts(api)
	if err != nil {
		return nil, err
	}

	items := getHostListItems(hostList, hostInError)

	bytes, err := json.Marshal(items)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}
