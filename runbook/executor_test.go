package runbook

import (
	"testing"
	"os"
	"github.com/stretchr/testify/assert"
	"runtime"
)

var testScriptFilePathNonWindows = os.TempDir() + string(os.PathSeparator) + "executorTestScript.sh"
var testScriptFilePathWindows = os.TempDir() + string(os.PathSeparator) + "executorTestScript.bat"

func TestExecuteSuccess(t *testing.T) {
	testEnvironmentVariables := map[string]interface{}{
		"TESTENVVAR": "test env var",
	}

	if runtime.GOOS != "windows" {
		content := "echo \"Test output\"\necho \"Given Environment Variable: TESTENVVAR: $TESTENVVAR\"\n"
		err := createTestScriptFile(content, testScriptFilePathNonWindows)
		defer os.Remove(testScriptFilePathNonWindows)

		if err != nil {
			t.Error(err.Error())
		}

		cmdOutput, cmdErr, err := execute(testScriptFilePathNonWindows, nil, testEnvironmentVariables)

		assert.NoError(t, err, "Error from execute operation was not empty.")
		assert.Equal(t, "", cmdErr, "Error stream from executed file was not empty.")
		assert.Equal(t, "Test output\nGiven Environment Variable: TESTENVVAR: test env var\n", cmdOutput,
			"Output stream was not equal to expected.")
	} else {
		content := "@echo off\r\necho Test output\r\necho Given Environment Variable: TESTENVVAR: %TESTENVVAR%\n"
		err := createTestScriptFile(content, testScriptFilePathWindows)
		defer os.Remove(testScriptFilePathWindows)

		if err != nil {
			t.Error(err.Error())
		}

		cmdOutput, cmdErr, err := execute(testScriptFilePathWindows, nil, testEnvironmentVariables)

		assert.NoError(t, err, "Error from execute operation was not empty.")
		assert.Equal(t, "", cmdErr, "Error stream from executed file was not empty.")
		assert.Equal(t, "Test output\r\nGiven Environment Variable: TESTENVVAR: test env var\r\n", cmdOutput,
			"Output stream was not equal to expected.")
	}
}

func TestExecuteWithErrorStream(t *testing.T) {
	if runtime.GOOS != "windows" {
		content := ">&2 echo \"test error\"\n"
		err := createTestScriptFile(content, testScriptFilePathNonWindows)
		defer os.Remove(testScriptFilePathNonWindows)

		if err != nil {
			t.Error(err.Error())
		}

		cmdOutput, cmdErr, err := execute(testScriptFilePathNonWindows, nil, nil)

		assert.NoError(t, err, "Error from execute operation was not empty.")
		assert.Equal(t, "", cmdOutput, "Output stream from executed file was not empty.")
		assert.Equal(t, "test error\n", cmdErr, "Error stream was not equal to expected.")
	} else {
		content := "@echo off\r\necho test error>&2\r\n"
		err := createTestScriptFile(content, testScriptFilePathWindows)
		defer os.Remove(testScriptFilePathWindows)

		if err != nil {
			t.Error(err.Error())
		}

		cmdOutput, cmdErr, err := execute(testScriptFilePathWindows, nil, nil)

		assert.NoError(t, err, "Error from execute operation was not empty.")
		assert.Equal(t, "", cmdOutput, "Output stream from executed file was not empty.")
		assert.Equal(t, "test error\r\n", cmdErr, "Error stream was not equal to expected.")
	}
}

func TestExecuteWithError(t *testing.T) {
	if runtime.GOOS != "windows" {
		content := "sacmasapan"
		err := createTestScriptFile(content, testScriptFilePathNonWindows)
		defer os.Remove(testScriptFilePathNonWindows)

		if err != nil {
			t.Error(err.Error())
		}

		cmdOutput, cmdErr, err := execute(testScriptFilePathNonWindows, nil, nil)

		assert.Error(t, err, "Error from execute operation was empty.")
		assert.Equal(t, err.Error(), "exit status 127", "Error message was not equal to expected.")
		assert.Equal(t, "", cmdOutput, "Output stream from executed file was not empty.")
		assert.Equal(t, "", cmdErr, "Error stream from executed file was not empty.")
	} else {
		content := "sacmasapan"
		err := createTestScriptFile(content, testScriptFilePathWindows)
		defer os.Remove(testScriptFilePathWindows)

		if err != nil {
			t.Error(err.Error())
		}

		cmdOutput, cmdErr, err := execute(testScriptFilePathWindows, nil, nil)

		assert.Error(t, err, "Error from execute operation was empty.")
		assert.Equal(t, err.Error(), "exit status 1", "Error message was not equal to expected.")
		assert.Equal(t, "", cmdOutput, "Output stream from executed file was not empty.")
		assert.Equal(t, "", cmdErr, "Error stream from executed file was not empty.")
	}
}

func TestDetermineExecutable(t *testing.T) {
	result := determineExecutable("test.bat")
	assert.Equal(t, "cmd", result)
	result = determineExecutable("test.cmd")
	assert.Equal(t, "cmd", result)
	result = determineExecutable("test.ps1")
	assert.Equal(t, "powershell", result)
	result = determineExecutable("test.sh")
	assert.Equal(t, "sh", result)
	result = determineExecutable("/some/path/to/some/executable.bin")
	assert.Equal(t, "", result)
}
