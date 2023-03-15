package common

import (
	"errors"
	"os"

	"github.com/jeromelesaux/martine/log"
)

var ErrorIsNotDirectory = errors.New("is not a directory, Quiting")

func CheckOutput(out string) error {
	infos, err := os.Stat(out)
	if os.IsNotExist(err) {
		if err = os.MkdirAll(out, os.ModePerm); err != nil {
			log.GetLogger().Error("Error while creating directory %s error %v \n", out, err)
			return err
		}
		return nil
	}
	if !infos.IsDir() {
		log.GetLogger().Error("%s is not a directory can not continue\n", out)
		return ErrorIsNotDirectory
	}
	return nil
}
