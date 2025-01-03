package mynet

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

func DownloadFile(url, filePath string, configCb func(r *http.Request), progressCb func(readLen, totalLen int64)) error {
	// Create HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	if configCb != nil {
		configCb(req)
		//这个回调用于配制 REQ
		//req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file: status code %d", resp.StatusCode)
	}

	// Get content length for progress calculation
	var contentLength int64
	strInt := resp.Header.Get("Content-Length")
	if strInt != "" {
		contentLength, err = strconv.ParseInt(strInt, 10, 64)
		if err != nil {
			return fmt.Errorf("unable to get content length")
		}
	}

	// Create the file
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	// Wrap the response body with TeeReader to monitor progress
	pReader := &progressReader{
		Reader:    resp.Body,
		TotalSize: int64(contentLength),
		fnCb:      progressCb,
	}

	// Copy data to file while tracking progress
	_, err = io.Copy(file, pReader)
	if err != nil {
		return err
	}

	return nil
}

// ProgressReader tracks progress of a download
type progressReader struct {
	Reader    io.Reader
	TotalSize int64
	ReadBytes int64
	fnCb      func(readLen, totalLen int64)
}

func (pr *progressReader) Read(p []byte) (int, error) {
	n, err := pr.Reader.Read(p)
	pr.ReadBytes += int64(n)

	if pr.fnCb != nil {
		pr.fnCb(pr.ReadBytes, pr.TotalSize)
	}

	return n, err
}
