package common

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)

func WilcardedFiles(filespath []string) ([]string, error) {
	fullfilespath := make([]string, 0)

	for _, v := range filespath {
		dir := filepath.Dir(v)
		fis, err := os.ReadDir(dir)
		if err != nil {
			return fullfilespath, err
		}
		reg := filepath.Base(v)
		// fmt.Fprintf(os.Stdout, "Regular to match (%s)\n", reg)
		for _, f := range fis {
			if !f.IsDir() {
				check := filepath.Join(dir, f.Name())
				//	fmt.Fprintf(os.Stdout, "Checking regex for (%s) matches (%s)\n", f.Name(), reg)
				matchs, err := filepath.Match(reg, f.Name())
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error while checking match with error %v\n", err)
					return fullfilespath, err
				}
				// fmt.Fprintf(os.Stdout, "Returns %v\n", matchs)
				if matchs {
					//	fmt.Fprintf(os.Stdout, "Ok (%s) matches (%s)\n", reg, check)
					if !ContainsFilepath(fullfilespath, check) {
						fullfilespath = append(fullfilespath, check)
					}
				}
			}
		}
	}
	sort.Slice(
		fullfilespath,
		func(i, j int) bool {
			return sortName(fullfilespath[i]) < sortName(fullfilespath[j])
		},
	)
	return fullfilespath, nil
}

func ContainsFilepath(filespath []string, filePath string) bool {
	for _, v := range filespath {
		if v == filePath {
			return true
		}
	}
	return false
}

func sortName(filename string) string {
	ext := filepath.Ext(filename)
	name := filename[:len(filename)-len(ext)]
	// split numeric suffix
	i := len(name) - 1
	for ; i >= 0; i-- {
		if '0' > name[i] || name[i] > '9' {
			break
		}
	}
	i++
	// string numeric suffix to uint64 bytes
	// empty string is zero, so integers are plus one
	b64 := make([]byte, 64/8)
	s64 := name[i:]
	if len(s64) > 0 {
		u64, err := strconv.ParseUint(s64, 10, 64)
		if err == nil {
			binary.BigEndian.PutUint64(b64, u64+1)
		}
	}
	// prefix + numeric-suffix + ext
	return name[:i] + string(b64) + ext
}

func StructToBytes(v interface{}) ([]byte, error) {
	var b bytes.Buffer
	err := binary.Write(&b, binary.LittleEndian, v)
	return b.Bytes(), err
}
