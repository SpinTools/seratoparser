package serato_parser

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
	"math"
	"reflect"
	"strconv"
	"strings"
	"unicode/utf16"
	"unicode/utf8"

)

func parseFilePeek (b *bufio.Reader, n int) ([]byte, bool) {
	bPeek, err := b.Peek(n)
	if err != nil {
		return nil, true
	}

	return bPeek, false
}
func parseFileByte (b *bufio.Reader) (byte, bool) {
	bByte, err := b.ReadByte()
	if err != nil {
		return '\000', true
	}

	return bByte, false
}
func parseFileLen (b *bufio.Reader, n int) (string, bool) {
	var buffer bytes.Buffer
	counter := 0
	for {
		char,eof := parseFileByte(b)
		if eof { return "", true }
		buffer.WriteByte(char)
		counter++
		if counter == n {
			break
		}
	}

	return buffer.String(), false
}

func parseCString (b *bufio.Reader) (string, bool) {
	/*
	 * From the passed string, find the nearest \0 byte and return everything before it
	 */
	var buffer bytes.Buffer
	for {
		char, eof := parseFileByte(b)
		if eof { return "", true }
		if char == '\000' {
			break
		}
		buffer.WriteByte(char)
	}

	return buffer.String(), false
}

func parseField (b *bufio.Reader) (string, string, bool) {
	/*
	 *
	 */
	name, eof := parseFileLen(b, 4)
	if eof { return "", "", true }
	rawlen, eof := parseFileLen(b, 4)
	if eof { return "", "", true }
	length := int(hexBin2Float(rawlen))

	data, eof := parseFileLen(b, length)
	if eof { return "", "", true }

	return name, data, false
}

func matchUtf16 (b *bufio.Reader, s string) (bool) {
	/*
	 * Match utf16 string with next len() bytes
	 */
	chars,_ := parseFileLen(b, len(s))
	if chars == s {
		return true
	}

	return false
}

// UTF16BytesToString converts UTF-16 encoded bytes, in big or little endian byte order,
// to a UTF-8 encoded string.
func utf16BytesToString(b []byte) string {
	utf := make([]uint16, (len(b)+(2-1))/2)
	for i := 0; i+(2-1) < len(b); i += 2 {
		utf[i/2] = binary.BigEndian.Uint16(b[i:])
	}
	if len(b)/2 < len(utf) {
		utf[len(utf)-1] = utf8.RuneError
	}
	return string(utf16.Decode(utf))
}

func makeUtf8 (s string) (string) {
	/*
	 * Convert the passed string s to UTF8 format
	 */
	var buffer bytes.Buffer
	for i := 0; i < len(s); i++ {
		if s[i] == '\000' {
			continue
		}
		buffer.WriteByte(s[i])
	}

	return buffer.String()
}

func makeUtf16 (s string) (string) {
	/*
	 * Convert the passed string s to serato UTF16 format
	 */
	var buffer bytes.Buffer
	for i := 0; i < len(s); i++ {
		buffer.WriteByte(0)
		buffer.WriteByte(s[i])
	}

	return buffer.String()
}

func hexBin2Int (raw string) (int) {
	return int(hexBin2Float(raw))
}

func hexBin2Float (raw string) (val float64) {
	for i := 0; i < len(raw); i++ {
		fl1 := math.Pow(2, 8)
		fl2 := float64((len(raw) - 1) - i)
		val += float64(raw[i]) * math.Pow(fl1, fl2)
	}

	return val
}

func parseOtrk (dataBuffer *bufio.Reader, newEntity *MediaEntity) {
	elem := reflect.ValueOf(newEntity).Elem()
	for {
		dataName, dataValue, eof := parseField(dataBuffer)
		if eof { break }

		newEntity.DVOL = seratoVolume

		v := elem.FieldByName(strings.ToUpper(dataName))
		reflectValue(&v, dataValue)
	}
}

func parseAdat (dataValue string, newEntity interface{}) {
	adatBuffer := bufio.NewReader(strings.NewReader(dataValue))
	elem := reflect.ValueOf(newEntity).Elem()
	for {
		adatFieldHex, adatValue, eof := parseField(adatBuffer)
		if eof { break }

		adatFieldId := hexBin2Int(adatFieldHex)
		adatName := SeratoAdatMap[adatFieldId]

		v := elem.FieldByName(strings.ToUpper(adatName))
		reflectValue(&v, adatValue)
		if v.IsValid() && v.Type().String() == "string" {
			tmpStringVal := v.String()
			v.SetString(tmpStringVal[:len(tmpStringVal)-1])
		}
	}
}

func reflectValue (v *reflect.Value, dataValue string) {
	if v.IsValid() {
		t := v.Type().String()
		switch t {
		case "string":
			//rlog.Debug("%s \n", utf16BytesToString([]byte(dataValue)))
			v.SetString(utf16BytesToString([]byte(dataValue)))
		case "float":
			newFloat := hexBin2Float(dataValue) // convert hexbin to int
			//rlog.Debug("%d(%d) \n", newFloat, int64(newFloat))
			v.SetFloat(newFloat)
		case "int":
			newFloat := hexBin2Int(dataValue) // convert hexbin to int
			//rlog.Debug("%d(%d) \n", newFloat, int64(newFloat))
			v.SetInt(int64(newFloat))
		case "bool":
			newBool,_ := strconv.ParseBool(dataValue)
			//rlog.Debug("%t \n", newBool)
			v.SetBool(newBool)
		case "[]uint8":
			newBytes := make([]byte, 4)
			binary.LittleEndian.PutUint32(newBytes[:], uint32(hexBin2Int(dataValue)))
			//rlog.Debug("%t \n", newBytes)
			v.SetBytes(newBytes[:])
		}
	}
}

func unreadBytes(byteReader *bytes.Reader, n int) {
	for n > 0 {
		_ = byteReader.UnreadByte()
		n--
	}
}

func decodeLenInBytes (byteReader * bytes.Reader, n int) int {
	// Read length in bytes
	lenInBytes := make([]byte, n)
	if _, err := io.ReadAtLeast(byteReader, lenInBytes, n); err != nil {
		return 0
	}
	dataLen := int(binary.BigEndian.Uint32(lenInBytes))
	if dataLen == 0 {
		return 0
	}

	return dataLen
}