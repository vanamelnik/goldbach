package files

import (
	"compress/gzip"
	"encoding/binary"
	"io"
	"os"
)

// UnzipFile - read the file and unzip it using gzip to []uint64 slice
func UnzipFile(name string) ([]uint64, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	zr, err := gzip.NewReader(f)
	if err != nil {
		return nil, err
	}
	defer zr.Close()

	buf := make([]byte, 1048576*100)
	var data []uint64
	var count int
	for {
		n, err := zr.Read(buf)
		if err != nil && err != io.EOF {
			return nil, err
		}
		count += n
		for i := 0; i < n; i += 8 {
			num := binary.BigEndian.Uint64(buf[i:])
			data = append(data, num)
		}
		if err == io.EOF {
			break
		}
	}

	return data, nil
}
