package butane

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
)

func b64Gzip(ibuf *bytes.Buffer) (*bytes.Buffer, error) {
	var ret bytes.Buffer
	b64Enc := base64.NewEncoder(base64.StdEncoding, &ret)
	defer b64Enc.Close()
	gzEnc := gzip.NewWriter(b64Enc)
	defer gzEnc.Close()
	_, err := gzEnc.Write(ibuf.Bytes())
	if err != nil {
		return nil, err
	}
	return &ret, nil
}
