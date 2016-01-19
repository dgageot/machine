package its

import "testing"

func TestStatus(t *testing.T) {
	test := NewIntegrationTest(t)
	defer test.TearDown()

	test.Run("status: show error in case of no args", func() {
		test.Cmd("machine status").ShouldFailWith(`Error: No machine name(s) specified and no "default" machine exists.`)
	})
}
