package upnode

import (
	"strconv"
	"strings"
)

type Options struct {
	port      uint16
	ste       string //string to expect
	keywords  []uint64
	protocol  string
	redirects uint8
	forceWWW  bool
	withBody  bool
	follow    bool
}

func GetDefaultOptions() *Options {
	option := new(Options)

	option.port = 0
	option.ste = "default"
	option.follow = true
	option.forceWWW = false
	option.redirects = 3
	option.withBody = false

	return option
}

func (opt *Options) ParseFromString(optionString string) {
	newOptions := GetDefaultOptions()

	elements := strings.Split(optionString, "#")

	for _, el := range elements {

		temp := strings.Split(el, ":=")

		if len(temp) != 2 {
			continue
		}

		key, val := temp[0], temp[1]

		switch key {
		case "p":
			conv, _ := strconv.ParseUint(val, 10, 16)
			newOptions.port = uint16(conv)
		case "ste":
			newOptions.ste = val
		case "pr":
			newOptions.protocol = val
		case "flw":
			newOptions.follow = ParseOptionBool(val)
		case "rdr":
			conv, _ := strconv.ParseUint(val, 10, 16)
			newOptions.redirects = uint8(conv)
		case "fwww":
			newOptions.forceWWW = ParseOptionBool(val)
		case "wb":
			newOptions.withBody = ParseOptionBool(val)
		case "kwd":
			conv, _ := strconv.ParseUint(val, 10, 16)
			newOptions.keywords = append(newOptions.keywords, conv)
		default:
		}
	}
	*opt = *newOptions
}

func ParseOptionBool(val string) (eval bool) {
	numeric, _ := strconv.ParseUint(val, 10, 16)
	if numeric == 1 {
		eval = true
	} else {
		eval = false
	}

	return
}
