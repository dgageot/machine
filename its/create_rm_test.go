package its

import "testing"

func TestCreateRm(t *testing.T) {
	test := NewIntegrationTest(t)
	defer test.TearDown()

	test.Run("non-existent driver fails", func() {
		test.Cmd("machine create -d bogus bogus").ShouldFailWith(`Driver "bogus" not found. Do you have the plugin binary accessible in your PATH?`)
	})

	test.Run("non-existent driver fails", func() {
		test.Cmd("machine create -d bogus bogus").ShouldFailWith(`Driver "bogus" not found. Do you have the plugin binary accessible in your PATH?`)
	})

	test.Run("create with no name fails", func() {
		test.Cmd("machine create -d none").ShouldFailWith(`Error: No machine name specified`)
	})

	test.Run("create with invalid name fails", func() {
		test.Cmd("machine create -d none --url none ∞").ShouldFailWith(`Error creating machine: Invalid hostname specified. Allowed hostname chars are: 0-9a-zA-Z . -`)
	})

	test.Run("create with invalid name fails", func() {
		test.Cmd("machine create -d none --url none -").ShouldFailWith(`Error creating machine: Invalid hostname specified. Allowed hostname chars are: 0-9a-zA-Z . -`)
	})

	test.Run("create with invalid name fails", func() {
		test.Cmd("machine create -d none --url none .").ShouldFailWith(`Error creating machine: Invalid hostname specified. Allowed hostname chars are: 0-9a-zA-Z . -`)
	})

	test.Run("create with invalid name fails", func() {
		test.Cmd("machine create -d none --url none ..").ShouldFailWith(`Error creating machine: Invalid hostname specified. Allowed hostname chars are: 0-9a-zA-Z . -`)
	})

	test.Run("create with weird but valid name succeeds", func() {
		test.Cmd("machine create -d none --url none a").ShouldSucceed()
	})

	test.Run("fail with extra argument", func() {
		test.Cmd("machine create -d none --url none a extra").ShouldFailWith(`Invalid command line. Found extra arguments [extra]`)
	})

	test.Run("create with weird but valid name", func() {
		test.Cmd("machine create -d none --url none 0").ShouldSucceed()
	})

	test.Run("rm with no name fails", func() {
		test.Cmd("machine rm -y").ShouldFailWith(`Error: Expected to get one or more machine names as arguments`)
	})

	test.Run("rm non existent machine fails", func() {
		test.Cmd("machine rm ∞ -y").ShouldFailWith(`Error removing host "∞": Host does not exist: "∞"`)
	})

	test.Run("rm existing machine", func() {
		test.Cmd("machine rm 0 -y").ShouldSucceed()
	})

	test.Run("rm ask user confirmation when -y is not provided", func() {
		test.Cmd("machine create -d none --url none ba").ShouldSucceed()
		test.Cmd("echo y | machine rm ba").ShouldSucceed()
	})

	test.Run("rm deny user confirmation when -y is not provided", func() {
		test.Cmd("machine create -d none --url none ab").ShouldSucceed()
		test.Cmd("echo n | machine rm ab").ShouldSucceed()
	})

	test.Run("rm never prompt user confirmation when -f is provided", func() {
		test.Cmd("machine create -d none --url none c").ShouldSucceed()
		test.Cmd("machine rm -f c").ShouldSucceedWith("Successfully removed c")
	})
}
