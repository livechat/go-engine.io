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

func (r *limitReader) Read(b []byte) (int, error) {
	if r.remain == 0 {
		return 0, io.EOF
	}

	reader, ok := r.Reader.(*bufio.Reader)
	if !ok {
		reader = bufio.NewReader(r.Reader)
	}

	max := r.remain
	if len(b) < max {
		max = len(b)
	}

	count := 0
	for i := 0; i < max; i++ {
		myRune, n, err := reader.ReadRune()
		if err != nil {
			return count, err
		}
		myBytes := []byte(string(myRune))
		for j := 0; j < n; j++ {
			b[count] = myBytes[j]
			count++
			r.remain--
		}
	}
	return count, nil
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
