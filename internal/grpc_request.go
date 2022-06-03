package internal

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"time"
)

func NewHttpRequestWithHeader(ctx context.Context, url string, method string, body []byte, header map[string]string) ([]byte, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	for k, v := range header {
		req.Header.Add(k, v)
	}
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	req.WithContext(ctx)
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return readBody(resp.Body)
}

func readBody(body io.ReadCloser) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, body); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
