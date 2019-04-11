package gfx

import (
	"errors"
	"fmt"
	"github.com/jeromelesaux/m4client/m4"
	"os"
)

var (
	ErrorNoHostDefined = errors.New("No host defined.")
)

func ImportInM4(exportType *ExportType) error {
	if exportType.M4Host == "" {
		return ErrorNoHostDefined
	}
	if exportType.M4RemotePath == "" {
		fmt.Fprintf(os.Stdout, "No M4 remote path defined, will copy on folder root.")
		exportType.M4RemotePath = "/"
	}

	client := m4.M4Client{IPClient: exportType.M4Host}
	for _, v := range exportType.DskFiles {
		if err := client.Upload(exportType.M4RemotePath, v); err != nil {
			fmt.Fprintf(os.Stderr, "Something is wrong M4 host (%s) local file (%s) remote path (%s) error :%v\n",
				exportType.M4Host,
				v,
				exportType.M4RemotePath,
				err)
		}
	}


	return nil
}
