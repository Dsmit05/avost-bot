package btn

import (
	"strconv"
	"strings"
)

func EncodeStr(event string, args ...string) string {
	var buf strings.Builder

	buf.WriteString(event)

	for _, v := range args {
		buf.WriteString(" ")
		buf.WriteString(v)
	}

	return buf.String()
}

func DecodeStr(data string) (event string, args []string) {
	all := strings.Split(data, " ")
	if len(all) < 1 {
		return
	}

	event = all[0]
	args = all[1:]

	return
}

func EncodeInt(event string, args ...int) string {
	var buf strings.Builder

	buf.WriteString(event)

	for _, v := range args {
		buf.WriteString(" ")
		buf.WriteString(strconv.Itoa(v))
	}

	return buf.String()
}

func DecodeInt(data string) (event string, args []int) {
	all := strings.Split(data, " ")
	if len(all) < 1 {
		return
	}

	event = all[0]

	for _, v := range all[1:] {
		intV, err := strconv.Atoi(v)
		if err != nil {
			intV = 0
		}

		args = append(args, intV)
	}

	return
}
