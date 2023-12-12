package mycrypto

import (
	"bytes"
	"compress/zlib"
	"io"
)

func CompressZip(inputData []byte) []byte {
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	defer w.Close()
	w.Write(inputData)
	return b.Bytes()
}

func DecompressZip(inputZipData []byte) []byte {
	b := bytes.NewReader(inputZipData)
	r, err := zlib.NewReader(b)
	if err != nil {
		return nil
	}
	defer r.Close()
	var returnByte bytes.Buffer
	io.Copy(&returnByte, r)
	return returnByte.Bytes()
}
