package commands

import "github.com/docker/machine/libmachine/persist"

func cmdIP(c CommandLine, store persist.Store) error {
	return runActionOnHosts("ip", c.Args(), store)
}
