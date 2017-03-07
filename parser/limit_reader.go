package parser

import (
	"bufio"
	"io"
)

type limitReader struct {
	io.Reader
	remain int
}

func newLimitReader(r io.Reader, limit int) *limitReader {
	return &limitReader{
		Reader: r,
		remain: limit,
	}
}

func (r *limitReader) Read(buf []byte) (int, error) {
	if r.remain == 0 {
		return 0, io.EOF
	}

	reader, ok := r.Reader.(*bufio.Reader)
	if !ok {
		reader = bufio.NewReader(r.Reader)
	}

	limit := r.remain
	if len(buf) < limit {
		limit = len(buf)
	}

	var idx int
	for i := 0; i < limit; i++ {
		c, _, err := reader.ReadRune()
		if err != nil {
			return idx, err
		}
		bytes := []byte(string(c))
		for _, b := range bytes {
			buf[idx] = b
			idx++
		}
		r.remain--
	}
	return idx, nil
}

func (r *limitReader) Close() error {
	if r.remain > 0 {
		b := make([]byte, 10240)
		for {
			_, err := r.Read(b)
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
		}
	}
	return nil
}
