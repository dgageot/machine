package its

import "testing"

func TestUrl(t *testing.T) {
	test := NewIntegrationTest(t)
	defer test.TearDown()

	test.Run("url: show error in case of no args", func() {
		test.Cmd("machine url").ShouldFailWith(`Error: No machine name(s) specified and no "default" machine exists.`)
	})
}
