package its

import "testing"

func TestDriverHelp(t *testing.T) {
	test := NewIntegrationTest(t)
	defer test.TearDown()

	test.Run("no --help flag or command specified", func() {
		test.Cmd("machine create -d $DRIVER").ShouldFailWith("Error: No machine name specified")
	})

	test.Run("-h flag specified", func() {
		test.Cmd("machine create -d $DRIVER -h").ShouldSucceedWith(test.DriverName())
	})

	test.Run("--help flag specified", func() {
		test.Cmd("machine create -d $DRIVER --help").ShouldSucceedWith(test.DriverName())
	})
}
