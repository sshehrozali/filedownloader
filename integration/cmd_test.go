package integration

import (
	"filedownloader/cmd"
	"filedownloader/shared/tests"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Run_ShouldDownloadFileFromServerAndSaveLocally(t *testing.T) {
	mockWebServer := tests.StartMockWebServer("/download/file", "sample response data")
	defer mockWebServer.Close()

	testDownloadUrl := mockWebServer.URL + "/download/file"
	tests.SetMockCliInput(testDownloadUrl)

	err := cmd.Run()

	assert.NoError(t, err, "No error occurred in CMD")

	expectedFilePath := "downloaded_file.txt"

	_, fileErr := os.Stat(expectedFilePath)
	assert.NoError(t, fileErr, "File is downloaded successfully")

	defer os.Remove(expectedFilePath)
}
