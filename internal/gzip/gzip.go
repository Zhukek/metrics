package gzip

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
)

type GzipWriter struct {
	writer    http.ResponseWriter
	zipWriter *gzip.Writer
}

func (c *GzipWriter) Close() error {
	return c.zipWriter.Close()
}

func (c *GzipWriter) Write(b []byte) (int, error) {
	return c.zipWriter.Write(b)
}

func (c *GzipWriter) Header() http.Header {
	return c.writer.Header()
}

func (c *GzipWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		c.writer.Header().Set("Content-Encoding", "gzip")
	}
	c.writer.WriteHeader(statusCode)
}

func NewGzipWriter(writer http.ResponseWriter) *GzipWriter {
	return &GzipWriter{
		writer:    writer,
		zipWriter: gzip.NewWriter(writer),
	}
}

type GzipReader struct {
	r         io.ReadCloser
	zipReader *gzip.Reader
}

func (c *GzipReader) Read(p []byte) (n int, err error) {
	return c.zipReader.Read(p)
}

func (c *GzipReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zipReader.Close()
}

func NewGzipReader(reader io.ReadCloser) (*GzipReader, error) {
	zipReader, err := gzip.NewReader(reader)

	if err != nil {
		return nil, err
	}

	return &GzipReader{
		r:         reader,
		zipReader: zipReader,
	}, nil
}

func GzipCompress(data []byte) ([]byte, error) {
	var b bytes.Buffer

	writer := gzip.NewWriter(&b)

	_, err := writer.Write(data)

	if err != nil {
		return nil, err
	}

	err = writer.Close()

	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}
