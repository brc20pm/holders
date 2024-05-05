package tools

import (
	"crypto/md5"
	"encoding/hex"
	"strconv"
)

func Float2String(scientificNum float64) string {
	readableStr := strconv.FormatFloat(scientificNum, 'f', -1, 64)
	return readableStr
}

func String2Float(numberStr string) (float64, error) {
	f, err := strconv.ParseFloat(numberStr, 64)
	if err != nil {
		return 0, err
	}
	return f, nil
}

func TokenIdHash(tokenId string) string {
	sum := md5.Sum([]byte(tokenId))
	return hex.EncodeToString(sum[:])
}
