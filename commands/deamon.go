package commands

import (
	"net/http"

	"fmt"

	"encoding/json"
	"time"

	"github.com/docker/machine/libmachine"
	"github.com/docker/machine/libmachine/log"
	"github.com/docker/machine/libmachine/persist"
	"github.com/gorilla/mux"
)

const (
	daemonDefaultPort = 8080
)

func cmdDaemon(c CommandLine, api libmachine.API) error {
	port := c.Int("port")

	log.Infof("Running on port %d", port)

	r := mux.NewRouter()

	r.HandleFunc("/ls", wrap(api, lsHandler))

	http.ListenAndServe(fmt.Sprintf(":%d", port), r)

	return nil
}

func lsHandler(api libmachine.API, response http.ResponseWriter, request *http.Request) error {
	stateTimeoutDuration = 10 * time.Second

	hostList, hostInError, err := persist.LoadAllHosts(api)
	if err != nil {
		return err
	}

	items := getHostListItems(hostList, hostInError)

	bytes, err := json.Marshal(items)
	if err != nil {
		return err
	}

	response.WriteHeader(200)
	response.Write(bytes)

	return nil
}

func wrap(api libmachine.API, handler func(api libmachine.API, response http.ResponseWriter, request *http.Request) error) func(http.ResponseWriter, *http.Request) {
	return func(response http.ResponseWriter, request *http.Request) {
		err := handler(api, response, request)
		if err != nil {
			fmt.Printf("Error: %s", err)
		}
	}
}
