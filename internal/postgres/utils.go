package postgres

import (
	"encoding/base64"
	"strconv"
)

// DecodeCursor will decode the cursor string into timestamp
func DecodeCursor(csr string) (timestamp int64, err error) {
	target, err := base64.StdEncoding.DecodeString(csr)
	if err != nil {
		return
	}
	timestamp, err = strconv.ParseInt(string(target), 10, 64)
	return
}

// EncodeCursor will encode the given data to Base64 format for cursor
func EncodeCursor(timestamp int64) (res string) {
	return base64.StdEncoding.EncodeToString([]byte(strconv.FormatInt(timestamp, 10)))
}
