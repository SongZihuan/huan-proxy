package utils

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

func IsUTF8(b []byte) bool {
	return utf8.Valid(b)
}

func IsUTF8String(s string) bool {
	return utf8.ValidString(s)
}

func HasUTF8BOM(data []byte) bool {
	if len(data) >= 3 && data[0] == 0xEF && data[1] == 0xBB && data[2] == 0xBF {
		return true
	}
	return false
}

func RemoveBOMIfExists(data []byte) []byte {
	if HasUTF8BOM(data) {
		return data[3:]
	}
	return data
}

func HasInvisibleByteSlice(data []byte) bool {
	for i := 0; i < len(data); {
		runeValue, size := utf8.DecodeRune(data[i:])
		if !unicode.IsPrint(runeValue) {
			return true
		}
		i += size
	}
	return false
}

func HasInvisibleString(str string) bool {
	for _, runeValue := range str {
		if !unicode.IsPrint(runeValue) {
			fmt.Printf("%d\n", runeValue)
			return true
		}
	}
	return false
}
