package cli

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/docker/docker/pkg/fileutils"
	"github.com/docker/machine/its"
)

// test.Machinef
// Force shared machine name
// If forced, don't fail if machine already exists
// Copy boot2docker.iso
// $NAME
//
func TestCoreCommands(t *testing.T) {
	//	test := its.NewTest(t)
	test := its.NewTestWithStorage(t, "/Users/dgageot/.docker/its")
	//	defer test.TearDown()

	test.SkipDrivers("ci-test", "none")

	//	name := "its-" + test.DriverName() + "-shared-" + time.Now().Format("15-04-05-000")
	name := "shared"

	if false {
		if err := fileutils.CreateIfNotExists(filepath.Join(test.StoragePath(), "cache"), true); err != nil {
			fmt.Println(err)
		}

		if _, err := fileutils.CopyFile("/Users/dgageot/.docker/machine/cache/boot2docker.iso", filepath.Join(test.StoragePath(), "cache", "boot2docker.iso")); err != nil {
			fmt.Println(err)
		}

		test.Run("create shared machine", func() {
			test.Machine("create -d $DRIVER " + name).Should().Succeed()
		})
	}

	test.Run("machine should not exist", func() {
		test.Machine("inspect UNKNOWN").Should().Fail(`Host does not exist: "UNKNOWN"`)
	})

	test.Run("appears with ls", func() {
		test.Machine("ls -q").Should().Succeed(name)
	})

	test.Run("has status 'started' in ls", func() {
		test.Machine("ls -q --filter state=Running").Should().Succeed(name)
	})

//	test.Run("create with same name fails", func() {
//		test.Machine("create -d $DRIVER $NAME").Should().Fail(`Host already exists: "` + name + `"`)
//	})

	//	test.Run("run busybox container", func() {
	//		test.Machine("run docker $(machine config $NAME) run busybox echo hello world").Should().Succeed( + name + `"`)
	//	})

	test.Run("url", func() {
		test.Machine("url $NAME").Should().Succeed()
	})

//	test.Run("ip", func() {
//		test.Machinef("ip %s", name).Should().Succeed()
//	})
//
//	test.Run("ssh", func() {
//		test.Machinef("ssh %s -- ls -lah /", name).Should().Succeed("total")
//	})
//
//	test.Run("version", func() {
//		test.Machinef("version %s", name).Should().Succeed()
//	})
//
//	test.Run("docker commands with the socket should work", func() {
//		test.Machinef("ssh %s -- sudo docker version", name).Should().Succeed()
//	})
}
