package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
	"filedownloader/network/clients"
	"filedownloader/network"
	"filedownloader/file"
)

// GLOBALs
var exit = os.Exit

func loader(done chan bool) {
	symbols := []string{"|", "/", "-", "\\"}

	i := 0
	for {
		select {
		case <-done:
			return
		default:
			fmt.Printf("\r%s", symbols[i%len(symbols)])
			i++
			time.Sleep(100 * time.Millisecond)

		}
	}
}

func fileExtensionRetriever(mimeType string) string {
	commonMIMEtypes := map[string]string{
		// Text types
		"text/plain":      ".txt",
		"text/html":       ".html",
		"text/css":        ".css",
		"text/javascript": ".js",

		// Image types
		"image/jpeg":    ".jpeg",
		"image/png":     ".png",
		"image/gif":     ".gif",
		"image/svg+xml": ".svg",

		// Audio types
		"audio/mpeg": ".mp3",
		"audio/wav":  ".wav",
		"audio/ogg":  ".ogg",

		// Video types
		"video/mp4":  ".mp4",
		"video/webm": ".webm",
		"video/ogg":  ".ogv",

		// Application types
		"application/json":                  ".json",
		"application/xml":                   ".xml",
		"application/pdf":                   ".pdf",
		"application/zip":                   ".zip",
		"application/octet-stream":          ".bin",
		"application/x-www-form-urlencoded": ".url",

		// Multipart types
		"multipart/form-data": ".multpart",
	}

	fileExtension := commonMIMEtypes[mimeType]
	return fileExtension
}

type DownloadResult struct {
	Data     []byte
	MimeType string
}

func promptInput() string {
	var url string

	fmt.Print("Enter download URL: ")
	fmt.Scanln(&url)

	if url == "" {
		fmt.Println("\nProvided download URL is empty.")
		exit(1)
	}

	return url
}

func downloadFile(networkClient clients.HttpClient, url string) DownloadResult {
	fmt.Println("\nFetching data from the server...")

	// Start loader in a goroutine
	done := make(chan bool)
	go loader(done)

	response, err := networkClient.Get(url)

	// Stop the loader once request is finished, close the channel
	done <- true
	close(done)

	if err != nil {
		fmt.Println("\nError while performing GET request.")
		exit(1)
	}

	if response.StatusCode != http.StatusOK {
		fmt.Println("Request failed with HTTP status: ", response.StatusCode)
		exit(1)
	}

	mimeType := response.Header.Get(("Content-Type"))
	fmt.Println("\nDetected content-type: ", mimeType)

	body, err := io.ReadAll(response.Body) // Parse body into byte[]

	if err != nil {
		fmt.Println("Error while reading response body.")
		os.Exit(1)
	}

	return DownloadResult{
		Data:     body,
		MimeType: mimeType,
	}
}

func saveLocally(downloadResult DownloadResult) {
	fileName := "downloaded_file" + fileExtensionRetriever(downloadResult.MimeType)
	file, fileErr := os.Create(fileName)
	if fileErr != nil {
		fmt.Printf("Error opening %s\n", fileName)
		defer file.Close()
		os.Exit(1)
	}
	defer file.Close() // Close the file automatically once I/O operations are finished

	bytesWritten, err := file.Write(downloadResult.Data)

	if err != nil {
		fmt.Println("Error copying bytes locally.")
		os.Exit(1)
	}

	if bytesWritten <= 0 {
		fmt.Println("0 bytes copied.")
		os.Exit(1)
	}

	fmt.Printf("Downloaded file size: %d bytes\n", bytesWritten)
	fmt.Printf("File downloaded successfully. Check your current local directory %s\n", fileName)
	os.Exit(0)
}

func run() {
	url := promptInput()

	httpClient := &network.DefaultHttpClient{}
	downloadResult := downloadFile(httpClient, url)

	fileWriter := &file.DefaultFileWriter{}
	saveLocally(downloadResult)
}

func main() {
	run()
}
