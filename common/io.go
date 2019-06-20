package common

import (
	"encoding/binary"
	"strconv"
	"sort"
	"io/ioutil"
	"path/filepath"
)


func WilcardedFiles(filespath []string) ([]string,error) {
	fullfilespath := make([]string,0)

	for _,v := range filespath {
		dir := filepath.Dir(v)
		fis,err := ioutil.ReadDir(dir)
		if err != nil {
			return fullfilespath,err
		}
		for _,f := range fis {
			if ! f.IsDir() {
				check := dir + string(filepath.Separator) + f.Name()
				matchs, err := filepath.Match(v,check)
				if err != nil {
					return fullfilespath,err
				}
				if matchs {
					if !ContainsFilepath(fullfilespath,check) {
						fullfilespath = append(fullfilespath,check)
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
	return fullfilespath,nil
}



func ContainsFilepath(filespath []string, filePath string) bool {
	for _,v:=range filespath {
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
