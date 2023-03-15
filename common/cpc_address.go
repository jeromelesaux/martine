package common

import (
	"errors"
	"strconv"
	"strings"

	"github.com/jeromelesaux/martine/log"
)

var (
	ErrorCannotBeParsed = errors.New("can not be parsed")
)

func ParseHexadecimal16(address string) (uint16, error) {

	switch address[0] {
	case '#':
		value := strings.Replace(address, "#", "", -1)
		v, err := strconv.ParseUint(value, 16, 16)
		if err != nil {
			log.GetLogger().Error("cannot get the hexadecimal value fom %s, error : %v\n", address, err)
			return 0, ErrorCannotBeParsed
		} else {
			return uint16(v), nil
		}
	case '0':
		value := strings.Replace(address, "0x", "", -1)
		v, err := strconv.ParseUint(value, 16, 16)
		if err != nil {
			log.GetLogger().Error("cannot get the hexadecimal value fom %s, error : %v\n", address, err)
			return 0, ErrorCannotBeParsed
		} else {
			return uint16(v), nil
		}
	default:
		v, err := strconv.ParseUint(address, 10, 16)
		if err != nil {
			log.GetLogger().Error("cannot get the hexadecimal value fom %s, error : %v\n", address, err)
			return 0, ErrorCannotBeParsed
		} else {
			return uint16(v), nil
		}
	}
}

func ParseHexadecimal8(address string) (uint8, error) {

	switch address[0] {
	case '#':
		value := strings.Replace(address, "#", "", -1)
		v, err := strconv.ParseUint(value, 16, 8)
		if err != nil {
			log.GetLogger().Error("cannot get the hexadecimal value fom %s, error : %v\n", address, err)
			return 0, ErrorCannotBeParsed
		} else {
			return uint8(v), nil
		}
	case '0':
		value := strings.Replace(address, "0x", "", -1)
		v, err := strconv.ParseUint(value, 16, 8)
		if err != nil {
			log.GetLogger().Error("cannot get the hexadecimal value fom %s, error : %v\n", address, err)
			return 0, ErrorCannotBeParsed
		} else {
			return uint8(v), nil
		}
	default:
		v, err := strconv.ParseUint(address, 10, 8)
		if err != nil {
			log.GetLogger().Error("cannot get the hexadecimal value fom %s, error : %v\n", address, err)
			return 0, ErrorCannotBeParsed
		} else {
			return uint8(v), nil
		}
	}
}
