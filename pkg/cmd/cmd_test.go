package cmd_test

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jenkins-x-labs/cli-doc-gen/pkg/cmd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCmd(t *testing.T) {
	//t.Parallel()

	_, co := cmd.NewCmd()
	outFile, err := ioutil.TempFile("", "")
	require.NoError(t, err, "failed to create tempo file")
	outFileName := outFile.Name()
	co.OutFile = outFileName
	co.TrimPrefix = []string{"retry "}
	co.File = filepath.Join("test_data", "setup_resources.sh")

	err = co.Run()
	require.NoError(t, err, "failed to run")

	expectedFile := filepath.Join("test_data", "expected.txt")
	assert.FileExists(t, expectedFile, "expectedFile results file")
	expectedData, err := ioutil.ReadFile(expectedFile)
	require.NoError(t, err, "failed to load file %s", expectedFile)
	expectedText := strings.TrimSpace(string(expectedData))

	assert.FileExists(t, outFileName, "did not generate the file")
	data, err := ioutil.ReadFile(outFileName)
	require.NoError(t, err, "failed to load file %s", outFileName)
	actualText := strings.TrimSpace(string(data))

	if diff := cmp.Diff(expectedText, actualText); diff != "" {
		t.Errorf("Unexpected generated file")
		t.Log(diff)

		t.Logf("generated file:\n%s\n", actualText)
	}
}
