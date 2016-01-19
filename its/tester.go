package its

import (
	"os/exec"
	"regexp"
	"strings"
	"testing"

	"io/ioutil"

	"os"

	"fmt"

	"runtime"
)

var (
	regexpCommandLine = regexp.MustCompile("('[^']*')|(\\S+)")
)

type IntegrationTest interface {
	Run(description string, action func())
	Cmd(commandLine string) IntegrationTest
	DriverName() string
	ShouldContainLines(count int) IntegrationTest
	ShouldContainLine(index int, text string) IntegrationTest
	ShouldEqualLine(index int, text string) IntegrationTest
	ShouldSucceed() IntegrationTest
	ShouldSucceedWith(message string) IntegrationTest
	ShouldFail() IntegrationTest
	ShouldFailWith(errorMessage string) IntegrationTest
	TearDown()
}

func NewIntegrationTest(t *testing.T) IntegrationTest {
	storagePath, _ := ioutil.TempDir("", "docker")

	return &dockerMachineTest{
		t:           t,
		storagePath: storagePath,
	}
}

type dockerMachineTest struct {
	t           *testing.T
	storagePath string

	Description string
	RawOutput   string
	Lines       []string
	Err         error
	Failed      bool
}

func (dmt *dockerMachineTest) Run(description string, action func()) {
	dmt.Description = description
	dmt.RawOutput = ""
	dmt.Lines = nil
	dmt.Err = nil
	dmt.Failed = false

	fmt.Print("\033[1;33m[..]\033[0m " + description)
	action()

	if dmt.Failed {
		fmt.Println("\r\033[1;31m[KO]\033[0m " + description)
	} else {
		fmt.Println("\r\033[1;32m[OK]\033[0m " + description)
	}
}

func (dmt *dockerMachineTest) DriverName() string {
	driver := os.Getenv("DRIVER")
	if driver == "" {
		// TEMP
		return "none"
	}

	return driver
}

func (dmt *dockerMachineTest) fail(message string, args ...interface{}) {
	dmt.Failed = true

	allArgs := []interface{}{dmt.Description}
	allArgs = append(allArgs, args...)

	dmt.t.Errorf("%s\nExpected "+message, allArgs...)
}

func (dmt *dockerMachineTest) ShouldContainLines(count int) IntegrationTest {
	if count != len(dmt.Lines) {
		dmt.fail("%d lines but got %d", count, len(dmt.Lines))
	}
	return dmt
}

func (dmt *dockerMachineTest) ShouldContainLine(index int, text string) IntegrationTest {
	if index >= len(dmt.Lines) {
		dmt.Failed = true
		dmt.fail("at least %d lines\nGot %d", index+1, len(dmt.Lines))
	} else if !strings.Contains(dmt.Lines[index], text) {
		dmt.Failed = true
		dmt.fail("line %d to contain '%s'\nGot '%s'", index, text, dmt.Lines[index])
	}
	return dmt
}

func (dmt *dockerMachineTest) ShouldEqualLine(index int, text string) IntegrationTest {
	if index >= len(dmt.Lines) {
		dmt.Failed = true
		dmt.fail("at least %d lines\nGot %d", index+1, len(dmt.Lines))
	} else if text != dmt.Lines[index] {
		dmt.Failed = true
		dmt.fail("line %d to be '%s'\nGot '%s'", index, text, dmt.Lines[index])
	}
	return dmt
}

func (dmt *dockerMachineTest) ShouldSucceed() IntegrationTest {
	if dmt.Err != nil {
		dmt.Failed = true
		dmt.fail("to succeed\nFailed with %s", dmt.Err)
	}
	return dmt
}

func (dmt *dockerMachineTest) ShouldSucceedWith(text string) IntegrationTest {
	if dmt.Err != nil {
		dmt.Failed = true
		dmt.fail("to succeed\nFailed with %s", dmt.Err)
	} else if !strings.Contains(dmt.RawOutput, text) {
		dmt.Failed = true
		dmt.fail("output to contain '%s'\nGot '%s'", text, dmt.RawOutput)
	}
	return dmt
}

func (dmt *dockerMachineTest) ShouldFail() IntegrationTest {
	if dmt.Err == nil {
		dmt.Failed = true
		dmt.fail("to fail\nGot success")
	}
	return dmt
}

func (dmt *dockerMachineTest) ShouldFailWith(text string) IntegrationTest {
	if dmt.Err == nil {
		dmt.Failed = true
		dmt.fail("to fail\nGot success")
	} else if !strings.Contains(dmt.RawOutput, text) {
		dmt.Failed = true
		dmt.fail("output to contain '%s'\nGot '%s'", text, dmt.RawOutput)
	}
	return dmt
}

func (dmt *dockerMachineTest) testedBinary() string {
	var path string
	if runtime.GOOS == "windows" {
		path = "..\\bin\\docker-machine.exe"
	} else {
		path = "../bin/docker-machine"
	}

	_, err := os.Stat(path)
	if err != nil {
		dmt.t.Fatalf("%s binary not found", path)
		return ""
	}

	return path
}

func (dmt *dockerMachineTest) replaceMachinePath(commandLine string) string {
	return strings.Replace(commandLine, "machine", dmt.testedBinary(), -1)
}

func (dmt *dockerMachineTest) replaceDriver(commandLine string) string {
	return strings.Replace(commandLine, "$DRIVER", dmt.DriverName(), -1)
}

func parseFields(commandLine string) []string {
	fields := regexpCommandLine.FindAllString(commandLine, -1)

	for i := range fields {
		if len(fields[i]) > 2 && strings.HasPrefix(fields[i], "'") && strings.HasSuffix(fields[i], "'") {
			fields[i] = fields[i][1 : len(fields[i])-1]
		}
	}

	return fields
}

func (dmt *dockerMachineTest) Cmd(commandLine string) IntegrationTest {
	if strings.HasPrefix(commandLine, "machine ") {
		return dmt.cmd(dmt.testedBinary(), parseFields(dmt.replaceDriver(commandLine[len("machine "):]))...)
	}
	return dmt.cmd("bash", "-c", dmt.replaceMachinePath(dmt.replaceDriver(commandLine)))
}

func (dmt *dockerMachineTest) cmd(command string, args ...string) IntegrationTest {
	cmd := exec.Command(command, args...)
	cmd.Env = []string{"MACHINE_STORAGE_PATH=" + dmt.storagePath}

	combinedOutput, err := cmd.CombinedOutput()

	dmt.RawOutput = string(combinedOutput)
	dmt.Lines = strings.Split(strings.TrimSpace(dmt.RawOutput), "\n")
	dmt.Err = err

	return dmt
}

func (dmt *dockerMachineTest) TearDown() {
	os.RemoveAll(dmt.storagePath)
}
