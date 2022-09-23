package common

import (
	"errors"
	"fmt"
	"os"
)

var ErrorIsNotDirectory = errors.New("is not a directory, Quiting")

func CheckOutput(out string) error {
	infos, err := os.Stat(out)
	if os.IsNotExist(err) {
		if err = os.MkdirAll(out, os.ModePerm); err != nil {
			fmt.Fprintf(os.Stderr, "Error while creating directory %s error %v \n", out, err)
			return err
		}
		return nil
	}
	if !infos.IsDir() {
		fmt.Fprintf(os.Stderr, "%s is not a directory can not continue\n", out)
		return ErrorIsNotDirectory
	}
	return nil
}
