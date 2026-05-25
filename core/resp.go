package core

import (
	"errors"
)

func Decode(data []byte) (interface{}, error) {
	if len(data) == 0 {
		return nil, errors.New("no data")
	}
	value, _, err := DecodeOne(data)
	return value, err

}

func DecodeOne(data []byte) (interface{}, int, error) {
	if len(data) == 0 {
		return nil, 0, errors.New("no data")
	}
	switch data[0] {
	case '+':
		return readSimpeString(data)
	case '-':
		return readError(data)
	case ':':
		return readInt64(data)
	case '$':
		return readBulkString(data)
	case '*':
		return readArray(data)
	}
	return nil, 0, nil
}

func readError(data []byte) (string, int, error) {
	return readSimpeString(data)
}

func readSimpeString(data []byte) (string, int, error) {
	pos := 1
	for ; data[pos] != '\r'; pos++ {

	}
	return string(data[1:pos]), pos + 2, nil
}

func readInt64(data []byte) (int64, int, error) {
	var pos int = 0
	negative := false
	if data[1] == '-' {
		negative = true
		pos = 2
	} else {
		pos = 1
	}
	var val int64 = 0
	for ; data[pos] != '\r'; pos++ {
		val = val*10 + int64(data[pos]-'0')
	}
	if negative {
		val = -val
	}
	return val, pos + 2, nil
}

func readBulkString(data []byte) (string, int, error) {
	pos := 1
	len, delta := readLength(data[pos:])
	pos += delta
	return string(data[pos:(pos + len)]), pos + len + 2, nil
}

func readLength(data []byte) (int, int) {
	// $5\r\nhello\r\n
	pos, len := 0, 0
	for pos = range data {
		b := data[pos]
		if !(b >= '0' && b <= '9') {
			return len, pos + 2
		}
		len = len*10 + int(b-'0')
	}
	return 0, 0
}

func readArray(data []byte) (interface{}, int, error) {
	pos := 1
	len, delta := readLength(data[pos:])
	pos += delta
	var elems []interface{} = make([]interface{}, len)
	for i := range elems {
		elem, delta, err := DecodeOne(data[pos:])
		if err != nil {
			return nil, 0, err
		}
		elems[i] = elem
		pos += delta
	}
	return elems, pos, nil
}
