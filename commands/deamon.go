package commands

import (
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"
	"time"

	"io/ioutil"
	"net"

	"github.com/docker/machine/commands/mcndirs"
	"github.com/docker/machine/libmachine"
	"github.com/docker/machine/libmachine/auth"
	"github.com/docker/machine/libmachine/engine"
	"github.com/docker/machine/libmachine/persist"
	"github.com/docker/machine/libmachine/swarm"

	"strings"

	"github.com/docker/machine/libmachine/drivers"
	"github.com/docker/machine/libmachine/drivers/rpc"
	"github.com/docker/machine/libmachine/host"
	"github.com/docker/machine/libmachine/mcnerror"
	"golang.org/x/crypto/ssh"
)

const daemonDefaultPort = 2200

func cmdDaemon(c CommandLine, api libmachine.API) error {
	port := c.Int("port")

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

	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		log.Fatalf("Failed to listen on %d (%s)", port, err)
	}

	log.Printf("Listening on %d...", port)
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
			case "subsystem":
				command := string(req.Payload[4:])
				req.Reply(true, nil)

				var output []byte
				var err error
				if command == "machine/ls" {
					output, err = runLs(api)
				} else if strings.HasPrefix(command, "machine/create") {
					parts := strings.Split(command, " ")
					output, err = []byte("DONE"), nil
					go runCreate(api, parts[1])
				} else {
					fmt.Println(command)
					output = []byte("UNKNOWN")
				}

				if err != nil {
					fmt.Println(err)
					output = []byte("ERROR: " + err.Error())
				}

				connection.Write(output)
				connection.Close()
				log.Printf("Session closed")
			}
		}
	}()
}

func runLs(api libmachine.API) ([]byte, error) {
	// Not safe
	defer rpcdriver.CloseDrivers()

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

func runCreate(api libmachine.API, name string) ([]byte, error) {
	fmt.Println("CREATE", name)
	validName := host.ValidateHostName(name)
	if !validName {
		return nil, fmt.Errorf("Error creating machine: %s", mcnerror.ErrInvalidHostname)
	}

	// TODO: Fix hacky JSON solution
	rawDriver, err := json.Marshal(&drivers.BaseDriver{
		MachineName: name,
		StorePath:   mcndirs.GetBaseDir(),
	})
	if err != nil {
		return nil, fmt.Errorf("Error attempting to marshal bare driver data: %s", err)
	}

	driverName := "virtualbox"
	driver, err := api.NewPluginDriver(driverName, rawDriver)
	if err != nil {
		return nil, fmt.Errorf("Error loading driver %q: %s", driverName, err)
	}

	h, err := api.NewHost(driver)
	if err != nil {
		return nil, fmt.Errorf("Error getting new host: %s", err)
	}

	h.HostOptions = &host.Options{
		AuthOptions: &auth.Options{
			CertDir:          mcndirs.GetMachineCertDir(),
			CaCertPath:       filepath.Join(mcndirs.GetMachineCertDir(), "ca.pem"),
			CaPrivateKeyPath: filepath.Join(mcndirs.GetMachineCertDir(), "ca-key.pem"),
			ClientCertPath:   filepath.Join(mcndirs.GetMachineCertDir(), "cert.pem"),
			ClientKeyPath:    filepath.Join(mcndirs.GetMachineCertDir(), "key.pem"),
			ServerCertPath:   filepath.Join(mcndirs.GetMachineDir(), name, "server.pem"),
			ServerKeyPath:    filepath.Join(mcndirs.GetMachineDir(), name, "server-key.pem"),
			StorePath:        filepath.Join(mcndirs.GetMachineDir(), name),
		},
		EngineOptions: &engine.Options{
			TLSVerify: true,
		},
		SwarmOptions: &swarm.Options{},
	}

	exists, err := api.Exists(h.Name)
	if err != nil {
		return nil, fmt.Errorf("Error checking if host exists: %s", err)
	}
	if exists {
		return nil, mcnerror.ErrHostAlreadyExists{
			Name: h.Name,
		}
	}

	driverOpts := rpcdriver.RPCFlags{
		Values: map[string]interface{}{
			"virtualbox-memory":                1024,
			"virtualbox-cpu-count":             1,
			"virtualbox-disk-size":             20000,
			"virtualbox-boot2docker-url":       "",
			"virtualbox-import-boot2docker-vm": "",
			"virtualbox-host-dns-resolver":     false,
			"virtualbox-hostonly-cidr":         "192.168.99.1/24",
			"virtualbox-hostonly-nictype":      "82540EM",
			"virtualbox-hostonly-nicpromisc":   "deny",
			"virtualbox-no-share":              false,
			"virtualbox-dns-proxy":             false,
			"virtualbox-no-vtx-check":          false,
			"swarm-master":                     false,
			"swarm-host":                       "",
			"swarm-discovery":                  "",
		},
	}

	if err := h.Driver.SetConfigFromFlags(driverOpts); err != nil {
		return nil, fmt.Errorf("Error setting machine configuration from flags provided: %s", err)
	}

	if err := api.Create(h); err != nil {
		return nil, fmt.Errorf("Error creating machine: %s", err)
	}

	if err := api.Save(h); err != nil {
		return nil, fmt.Errorf("Error attempting to save store: %s", err)
	}

	return []byte("OK"), nil
}
