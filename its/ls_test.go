package its

import "testing"

func TestLs(t *testing.T) {
	test := NewIntegrationTest(t)
	defer test.TearDown()

	test.Run("setup", func() {
		test.Cmd("machine create -d none --url none --engine-label app=1 testmachine5").ShouldSucceed()
		test.Cmd("machine create -d none --url none --engine-label foo=bar --engine-label app=1 testmachine4").ShouldSucceed()
		test.Cmd("machine create -d none --url none testmachine3").ShouldSucceed()
		test.Cmd("machine create -d none --url none testmachine2").ShouldSucceed()
		test.Cmd("machine create -d none --url none testmachine").ShouldSucceed()
	})

	test.Run("ls: filter on label", func() {
		test.Cmd("machine ls --filter label=foo=bar").
			ShouldSucceed().
			ShouldContainLines(2).
			ShouldContainLine(0, "NAME").
			ShouldContainLine(1, "testmachine4")
	})

	test.Run("ls: mutiple filters on label", func() {
		test.Cmd("machine ls --filter label=foo=bar --filter label=app=1").
			ShouldSucceed().
			ShouldContainLines(3).
			ShouldContainLine(0, "NAME").
			ShouldContainLine(1, "testmachine4").
			ShouldContainLine(2, "testmachine5")
	})

	test.Run("ls: non-existing filter on label", func() {
		test.Cmd("machine ls --filter label=invalid=filter").
			ShouldSucceed().
			ShouldContainLines(1).
			ShouldContainLine(0, "NAME")
	})

	test.Run("ls: filter on driver", func() {
		test.Cmd("machine ls --filter driver=none").
			ShouldSucceed().
			ShouldContainLines(6).
			ShouldContainLine(0, "NAME").
			ShouldContainLine(1, "testmachine").
			ShouldContainLine(2, "testmachine2").
			ShouldContainLine(3, "testmachine3").
			ShouldContainLine(4, "testmachine4").
			ShouldContainLine(5, "testmachine5")
	})

	test.Run("ls: filter on driver", func() {
		test.Cmd("machine ls -q --filter driver=none").
			ShouldSucceed().
			ShouldContainLines(5).
			ShouldEqualLine(0, "testmachine").
			ShouldEqualLine(1, "testmachine2").
			ShouldEqualLine(2, "testmachine3")
	})

	test.Run("ls: filter on state", func() {
		test.Cmd("machine ls --filter state=Running").
			ShouldSucceed().
			ShouldContainLines(6).
			ShouldContainLine(0, "NAME").
			ShouldContainLine(1, "testmachine").
			ShouldContainLine(2, "testmachine2").
			ShouldContainLine(3, "testmachine3")

		test.Cmd("machine ls -q --filter state=Running").
			ShouldSucceed().
			ShouldContainLines(5).
			ShouldEqualLine(0, "testmachine").
			ShouldEqualLine(1, "testmachine2").
			ShouldEqualLine(2, "testmachine3")

		test.Cmd("machine ls --filter state=None").
			ShouldSucceed().
			ShouldContainLines(1).
			ShouldContainLine(0, "NAME")

		test.Cmd("machine ls --filter state=Paused").
			ShouldSucceed().
			ShouldContainLines(1).
			ShouldContainLine(0, "NAME")

		test.Cmd("machine ls --filter state=Saved").
			ShouldSucceed().
			ShouldContainLines(1).
			ShouldContainLine(0, "NAME")

		test.Cmd("machine ls --filter state=Stopped").
			ShouldSucceed().
			ShouldContainLines(1).
			ShouldContainLine(0, "NAME")

		test.Cmd("machine ls --filter state=Stopping").
			ShouldSucceed().
			ShouldContainLines(1).
			ShouldContainLine(0, "NAME")

		test.Cmd("machine ls --filter state=Starting").
			ShouldSucceed().
			ShouldContainLines(1).
			ShouldContainLine(0, "NAME")

		test.Cmd("machine ls --filter state=Error").
			ShouldSucceed().
			ShouldContainLines(1).
			ShouldContainLine(0, "NAME")
	})

	test.Run("ls: filter on name", func() {
		test.Cmd("machine ls --filter name=testmachine2").
			ShouldSucceed().
			ShouldContainLines(2).
			ShouldContainLine(0, "NAME").
			ShouldContainLine(1, "testmachine2")

		test.Cmd("machine ls -q --filter name=testmachine3").
			ShouldSucceed().
			ShouldContainLines(1).
			ShouldEqualLine(0, "testmachine3")
	})

	test.Run("ls: filter on name with regex", func() {
		test.Cmd("machine ls --filter name=^t.*e[3-5]").
			ShouldSucceed().
			ShouldContainLines(4).
			ShouldContainLine(0, "NAME").
			ShouldContainLine(1, "testmachine3").
			ShouldContainLine(2, "testmachine4").
			ShouldContainLine(3, "testmachine5")

		test.Cmd("machine ls -q --filter name=^t.*e[45]").
			ShouldSucceed().
			ShouldContainLines(2).
			ShouldEqualLine(0, "testmachine4").
			ShouldEqualLine(1, "testmachine5")
	})

	test.Run("setup swarm", func() {
		test.Cmd("machine create -d none --url tcp://127.0.0.1:2375 --swarm --swarm-master --swarm-discovery token://deadbeef testswarm").ShouldSucceed()
		test.Cmd("machine create -d none --url tcp://127.0.0.1:2375 --swarm --swarm-discovery token://deadbeef testswarm2").ShouldSucceed()
		test.Cmd("machine create -d none --url tcp://127.0.0.1:2375 --swarm --swarm-discovery token://deadbeef testswarm3").ShouldSucceed()
	})

	test.Run("ls: filter on swarm", func() {
		test.Cmd("machine ls --filter swarm=testswarm").
			ShouldSucceed().
			ShouldContainLines(4).
			ShouldContainLine(0, "NAME").
			ShouldContainLine(1, "testswarm").
			ShouldContainLine(2, "testswarm2").
			ShouldContainLine(3, "testswarm3")
	})

	test.Run("ls: multi filter", func() {
		test.Cmd("machine ls -q --filter swarm=testswarm --filter name=^t.*e --filter driver=none --filter state=Running").
			ShouldSucceed().
			ShouldContainLines(3).
			ShouldEqualLine(0, "testswarm").
			ShouldEqualLine(1, "testswarm2").
			ShouldEqualLine(2, "testswarm3")
	})

	test.Run("ls: format on driver", func() {
		test.Cmd("machine ls --format {{.DriverName}}").
			ShouldSucceed().
			ShouldContainLines(8).
			ShouldEqualLine(0, "none").
			ShouldEqualLine(1, "none").
			ShouldEqualLine(2, "none").
			ShouldEqualLine(3, "none").
			ShouldEqualLine(4, "none").
			ShouldEqualLine(5, "none").
			ShouldEqualLine(6, "none").
			ShouldEqualLine(7, "none")
	})

	test.Run("ls: format on name and driver", func() {
		test.Cmd("machine ls --format 'table {{.Name}}: {{.DriverName}}'").
			ShouldSucceed().
			ShouldContainLines(9).
			ShouldContainLine(0, "NAME").
			ShouldEqualLine(1, "testmachine: none").
			ShouldEqualLine(2, "testmachine2: none").
			ShouldEqualLine(3, "testmachine3: none").
			ShouldEqualLine(4, "testmachine4: none").
			ShouldEqualLine(5, "testmachine5: none").
			ShouldEqualLine(6, "testswarm: none").
			ShouldEqualLine(7, "testswarm2: none").
			ShouldEqualLine(8, "testswarm3: none")
	})
}
