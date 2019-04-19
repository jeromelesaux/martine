package gfx

import (
	"errors"
	"fmt"
	"github.com/jeromelesaux/m4client/m4"
	"os"
	"path"
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

	fmt.Fprintf(os.Stdout,"Attempt to create remote directory (%s) to host (%s)\n",exportType.M4RemotePath,client.IPClient)
	if err := client.MakeDirectory(exportType.M4RemotePath); err != nil {
		fmt.Fprintf(os.Stderr,"Cannot create directory on M4 (%s) error %v\n",exportType.M4RemotePath,err)
	}

	for _, v := range exportType.DskFiles {
		fmt.Fprintf(os.Stdout,"Attempt to uploading file (%s) on remote path (%s) to host (%s)\n",v,exportType.M4RemotePath,client.IPClient)
		if err := client.Upload(exportType.M4RemotePath, v); err != nil {
			fmt.Fprintf(os.Stderr, "Something is wrong M4 host (%s) local file (%s) remote path (%s) error :%v\n",
				exportType.M4Host,
				v,
				exportType.M4RemotePath,
				err)
		}
	}
	if exportType.Dsk {
		dskFile := exportType.Fullpath(".dsk")
		fmt.Fprintf(os.Stdout,"Attempt to uploading file (%s) on remote path (%s) to host (%s)\n",dskFile,exportType.M4RemotePath,client.IPClient)
		if err := client.Upload(exportType.M4RemotePath,dskFile); err != nil {
			fmt.Fprintf(os.Stderr, "Something is wrong M4 host (%s) local file (%s) remote path (%s) error :%v\n",
				exportType.M4Host,
				dskFile,
				exportType.M4RemotePath,
				err)
		}
	}

	if exportType.M4Autoexec {
		p, err := client.Ls(exportType.M4RemotePath)
		if err != nil {
			fmt.Fprintf(os.Stderr,"Cannot go to the remote path (%s) error :%v\n", exportType.M4RemotePath, err)
		} else {
			fmt.Fprintf(os.Stdout,"Set the remote path (%s) \n",p)
		}

		var  overscanFile, basicFile string
		for _, v := range exportType.DskFiles {
			switch  path.Ext(v) {
			case ".BAS":
				basicFile = path.Base(v)
			case ".SCR" : 
				overscanFile = path.Base(v)
			}
		}
		if exportType.Scr {	
			client.Run(exportType.M4RemotePath + "/" + basicFile)
		} else {
			if exportType.Overscan {
				client.Run(exportType.M4RemotePath + "/" + overscanFile)
			}  else {
				fmt.Fprintf(os.Stdout,"Too many importants files, cannot choice.\n")
			}
		} 
	}

	return nil
}
