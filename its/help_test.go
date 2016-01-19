package its

import "testing"

func TestHelp(t *testing.T) {
	test := NewIntegrationTest(t)
	defer test.TearDown()

	test.Run("cli: show info", func() {
		test.Cmd("machine").ShouldSucceedWith("Usage:").ShouldSucceedWith("Create and manage machines running Docker")
	})

	test.Run("cli: show active help", func() {
		test.Cmd("machine active -h").ShouldSucceedWith("machine active")
	})

	test.Run("cli: show config help", func() {
		test.Cmd("machine config -h").ShouldSucceedWith("machine config")
	})

	test.Run("cli: show create help", func() {
		test.Cmd("machine create -h").ShouldSucceedWith("machine create")
	})

	test.Run("cli: show env help", func() {
		test.Cmd("machine env -h").ShouldSucceedWith("machine env")
	})

	test.Run("cli: show inspect help", func() {
		test.Cmd("machine inspect -h").ShouldSucceedWith("machine inspect")
	})

	test.Run("cli: show ip help", func() {
		test.Cmd("machine ip -h").ShouldSucceedWith("machine ip")
	})

	test.Run("cli: show kill help", func() {
		test.Cmd("machine kill -h").ShouldSucceedWith("machine kill")
	})

	test.Run("cli: show ls help", func() {
		test.Cmd("machine ls -h").ShouldSucceedWith("machine ls")
	})

	test.Run("cli: show regenerate-certs help", func() {
		test.Cmd("machine regenerate-certs -h").ShouldSucceedWith("machine regenerate-certs")
	})

	test.Run("cli: show restart help", func() {
		test.Cmd("machine restart -h").ShouldSucceedWith("machine restart")
	})

	test.Run("cli: show rm help", func() {
		test.Cmd("machine rm -h").ShouldSucceedWith("machine rm")
	})

	test.Run("cli: show scp help", func() {
		test.Cmd("machine scp -h").ShouldSucceedWith("machine scp")
	})

	test.Run("cli: show ssh help", func() {
		test.Cmd("machine ssh -h").ShouldSucceedWith("machine ssh")
	})

	test.Run("cli: show start help", func() {
		test.Cmd("machine start -h").ShouldSucceedWith("machine start")
	})

	test.Run("cli: show status help", func() {
		test.Cmd("machine status -h").ShouldSucceedWith("machine status")
	})

	test.Run("cli: show stop help", func() {
		test.Cmd("machine stop -h").ShouldSucceedWith("machine stop")
	})

	test.Run("cli: show upgrade help", func() {
		test.Cmd("machine upgrade -h").ShouldSucceedWith("machine upgrade")
	})

	test.Run("cli: show url help", func() {
		test.Cmd("machine url -h").ShouldSucceedWith("machine url")
	})

	test.Run("cli: show version", func() {
		test.Cmd("machine -v").ShouldSucceedWith("version")
	})

	test.Run("cli: show help", func() {
		test.Cmd("machine --help").ShouldSucceedWith("Usage:")
	})
}
