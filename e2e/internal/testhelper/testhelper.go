package testhelper

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/bradleyjkemp/cupaloy"
)

var TestTimeout = 1 * time.Minute

type TestCase struct {
	name          string
	arguments     []string
	ShouldSucceed bool
	options       TestCaseOptions
	displayStdErr bool
	ignoreForce   bool
}

type TestCaseOptions struct {
	DisplayStdErr bool
	IgnoreForce   bool
}

func NewTestCase(name string, arguments []string, options TestCaseOptions) TestCase {
	return TestCase{
		name:          name,
		arguments:     arguments,
		ShouldSucceed: true,
		options:       options,
		displayStdErr: options.DisplayStdErr,
		ignoreForce:   options.IgnoreForce,
	}
}

func executeApp(t *testing.T, arguments []string) (string, error) {
	cmd, cancel := CreateCommand(arguments)

	buffOut := bytes.NewBuffer(nil)
	buffErr := bytes.NewBuffer(nil)
	cmd.Stdout = buffOut
	cmd.Stderr = buffErr

	var err error

	timer := time.NewTimer(TestTimeout)
	commandFinished := make(chan struct{}, 1)
	combinedOutput := func() string {
		errStr := strings.TrimSuffix(buffErr.String(), "exit status 1\n")
		return buffOut.String() + "\n--\n" + errStr
	}

	go func() {
		err = cmd.Start()

		if err != nil {
			commandFinished <- struct{}{}
			return
		}

		err = cmd.Wait()
		commandFinished <- struct{}{}
	}()

	select {
	case <-timer.C:
		cancel()
		t.Fatalf("command failed to complete on time 'bearer %s':\n%s", strings.Join(arguments, " "), combinedOutput())
	case <-commandFinished:
		cancel()
	}

	return combinedOutput(), err
}

func CreateCommand(arguments []string) (*exec.Cmd, context.CancelFunc) {
	var cmd *exec.Cmd

	ctx, cancel := context.WithCancel(context.Background())

	if os.Getenv("USE_BINARY") != "" {
		cmd = exec.CommandContext(ctx, executablePath(), arguments...)
	} else {
		arguments = append([]string{"run", GetCWD() + "/cmd/bearer/main.go"}, arguments...)
		cmd = exec.CommandContext(ctx, "go", arguments...)
	}

	cmd.Dir = GetCWD()

	return cmd, cancel
}

func executablePath() string {
	if value, ok := os.LookupEnv("BEARER_EXECUTABLE_PATH"); ok {
		return value
	}

	return "./bearer"
}

func GetCWD() string {
	return os.Getenv("GITHUB_WORKSPACE")
}

func RunTestsWithSnapshotSubdirectory(t *testing.T, tests []TestCase, snapshotSubdirectory string) {
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			combinedOutput, err := executeTest(test, t)

			if test.ShouldSucceed && err != nil {
				t.Fatalf("command completed with error %s %s", err, combinedOutput)
			}

			if !test.ShouldSucceed && err == nil {
				t.Fatal("expected command to fail but it succeded instead")
			}

			cupaloyCopy := cupaloy.NewDefaultConfig().WithOptions(cupaloy.SnapshotSubdirectory(snapshotSubdirectory))
			cupaloyCopy.SnapshotT(t, combinedOutput)
		})
	}
}

func RunTests(t *testing.T, tests []TestCase) {
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			combinedOutput, err := executeTest(test, t)

			if test.ShouldSucceed && err != nil {
				t.Fatalf("command completed with error %s %s", err, combinedOutput)
			}

			if !test.ShouldSucceed && err == nil {
				t.Fatal("expected command to fail but it succeded instead")
			}

			cupaloy.SnapshotT(t, combinedOutput)
		})
	}
}

func executeTest(test TestCase, t *testing.T) (string, error) {
	arguments := test.arguments

	if !test.displayStdErr {
		arguments = append(arguments, "--quiet")
	}

	if !test.ignoreForce {
		arguments = append(arguments, "--force")
	}

	return executeApp(t, arguments)
}
