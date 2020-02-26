package common

import (
	"fmt"
	"os"

	"github.com/jeromelesaux/martine/export"
)

func CheckOutput(exportType *export.ExportType) error {
	_, err := os.Stat(exportType.OutputPath)
	if os.IsNotExist(err) {
		if err = os.MkdirAll(exportType.OutputPath, os.ModePerm); err != nil {
			fmt.Fprintf(os.Stderr, "Error while creating directory %s error %v \n", exportType.OutputPath, err)
			return err
		}
	}
	return nil
}
