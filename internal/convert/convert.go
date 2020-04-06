package convert

import "strconv"

func UInt8SliceTostring(source []uint8) string {
	ba := []byte{}
	for _, b := range source {
		ba = append(ba, byte(b))
	}
	return string(ba)
}

func UInt8SliceToInt64(source []uint8) (int64, error) {
	return strconv.ParseInt(UInt8SliceTostring(source), 10, 64)
}
