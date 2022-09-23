package net

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/jeromelesaux/m4client/m4"
	x "github.com/jeromelesaux/martine/export"
)

var (
	ErrorNoHostDefined = errors.New("no host defined")
)

func ImportInM4(cont *x.MartineContext) error {
	if cont.M4Host == "" {
		return ErrorNoHostDefined
	}
	if cont.M4RemotePath == "" {
		fmt.Fprintf(os.Stdout, "No M4 remote path defined, will copy on folder root.")
		cont.M4RemotePath = "/"
	}

	client := m4.M4Client{IPClient: cont.M4Host}
	client.ResetCpc()
	if !cont.Sna {
		fmt.Fprintf(os.Stdout, "Attempt to create remote directory (%s) to host (%s)\n", cont.M4RemotePath, client.IPClient)
		if err := client.MakeDirectory(cont.M4RemotePath); err != nil {
			fmt.Fprintf(os.Stderr, "Cannot create directory on M4 (%s) error %v\n", cont.M4RemotePath, err)
		}

		for _, v := range cont.DskFiles {
			fmt.Fprintf(os.Stdout, "Attempt to uploading file (%s) on remote path (%s) to host (%s)\n", v, cont.M4RemotePath, client.IPClient)
			if err := client.Upload(cont.M4RemotePath, v); err != nil {
				fmt.Fprintf(os.Stderr, "Something is wrong M4 host (%s) local file (%s) remote path (%s) error :%v\n",
					cont.M4Host,
					v,
					cont.M4RemotePath,
					err)
			}
		}
	} else {
		if err := client.Remove(cont.M4RemotePath + "test.sna"); err != nil {
			fmt.Fprintf(os.Stderr, "Cannot create directory on M4 (%s) error %v\n", cont.M4RemotePath, err)
		}
	}
	if cont.Dsk {
		dskFile := cont.Fullpath(".dsk")
		fmt.Fprintf(os.Stdout, "Attempt to uploading file (%s) on remote path (%s) to host (%s)\n", dskFile, cont.M4RemotePath, client.IPClient)
		if err := client.Upload(cont.M4RemotePath, dskFile); err != nil {
			fmt.Fprintf(os.Stderr, "Something is wrong M4 host (%s) local file (%s) remote path (%s) error :%v\n",
				cont.M4Host,
				dskFile,
				cont.M4RemotePath,
				err)
		}
	}

	if cont.Sna {
		if err := client.Upload(cont.M4RemotePath, cont.SnaPath); err != nil {
			fmt.Fprintf(os.Stderr, "Something is wrong M4 host (%s) local file (%s) remote path (%s) error :%v\n",
				cont.M4Host,
				cont.SnaPath,
				cont.M4RemotePath,
				err)
		}
	}

	if cont.M4Autoexec {
		if cont.Sna {
			client.Run(cont.M4RemotePath + "test.sna")
			return nil
		}
		p, err := client.Ls(cont.M4RemotePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot go to the remote path (%s) error :%v\n", cont.M4RemotePath, err)
		} else {
			fmt.Fprintf(os.Stdout, "Set the remote path (%s) \n", p)
		}

		var overscanFile, basicFile string
		for _, v := range cont.DskFiles {
			switch path.Ext(v) {
			case ".BAS":
				basicFile = path.Base(v)
			case ".SCR":
				overscanFile = path.Base(v)
			}
		}
		if cont.Scr {
			fmt.Fprintf(os.Stdout, "Execute basic file (%s)\n", "/"+cont.M4RemotePath+"/"+basicFile)
			client.Run("/" + cont.M4RemotePath + "/" + basicFile)
		} else {
			if cont.Overscan {
				fmt.Fprintf(os.Stdout, "Execute overscan file (%s)\n", "/"+cont.M4RemotePath+"/"+overscanFile)
				client.Run("/" + cont.M4RemotePath + "/" + overscanFile)
			} else {
				fmt.Fprintf(os.Stdout, "Too many importants files, cannot choice.\n")
			}
		}
	}

	return nil
}
