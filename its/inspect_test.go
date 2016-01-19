package its

import "testing"

func TestInspect(t *testing.T) {
	test := NewIntegrationTest(t)
	defer test.TearDown()

	test.Run("inspect: show error in case of no args", func() {
		test.Cmd("machine inspect").ShouldFailWith(`Error: No machine name(s) specified and no "default" machine exists.`)
	})
}
