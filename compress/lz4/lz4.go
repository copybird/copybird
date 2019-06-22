package lz4

import (
	"bufio"
	"errors"
	"io"
	"os"

	"github.com/pierrec/lz4"
)

var (
	errCompLevel = errors.New("compression level must be between -1 and 9")
)

// Compressor represents ...
type Compressor struct {
	reader io.Reader
	writer io.Writer
	level  int
}

// Init initializez comressor struct with reader, write
func (c *Compressor) Init(w io.Writer, r io.Reader) error {
	c.reader = r
	c.writer = w
	return nil
}

// InitLevel  initializez compressor with level
func (c *Compressor) InitLevel(level int) error {
	if level < -1 || level > 9 {
		return errCompLevel
	}
	c.level = level
	return nil
}

func (c *Compressor) CompressFile(inputFile, outputFile string) error {

	// open input file
	fin, err := os.Open(inputFile)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := fin.Close(); err != nil {
			panic(err)
		}
	}()
	// make a read buffer
	c.reader = bufio.NewReader(fin)

	// open output file
	fout, err := os.Create(outputFile)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := fout.Close(); err != nil {
			panic(err)
		}
	}()

	if err = c.Compress(); err != nil {
		return err
	}

	return nil
}

// Compress compresses file from input, saves to output
func (c *Compressor) Compress() error {

	// make a buffer to keep chunks that are read
	buf := make([]byte, 4096)

	// make an lz4 write buffer
	w := lz4.NewWriter(c.writer)

	for {
		// read a chunk
		n, err := c.reader.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if n == 0 {
			break
		}

		// write a chunk
		if _, err := w.Write(buf[:n]); err != nil {
			panic(err)
		}
	}

	if err := w.Flush(); err != nil {
		panic(err)
	}

	return nil
}

// Close closes compressor
func (c *Compressor) Close() error {
	return nil
}
