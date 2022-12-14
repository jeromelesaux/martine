package m4

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/jeromelesaux/m4client/m4"
	"github.com/jeromelesaux/martine/config"
)

var ErrorNoHostDefined = errors.New("no host defined")

func ImportInM4(cfg *config.MartineConfig) error {
	if cfg.M4Host == "" {
		return ErrorNoHostDefined
	}
	if cfg.M4RemotePath == "" {
		fmt.Fprintf(os.Stdout, "No M4 remote path defined, will copy on folder root.")
		cfg.M4RemotePath = "/"
	}

	client := m4.M4Client{IPClient: cfg.M4Host}
	err := client.ResetCpc()
	if err != nil {
		return err
	}
	if !cfg.Sna {
		fmt.Fprintf(os.Stdout, "Attempt to create remote directory (%s) to host (%s)\n", cfg.M4RemotePath, client.IPClient)
		if err := client.MakeDirectory(cfg.M4RemotePath); err != nil {
			fmt.Fprintf(os.Stderr, "Cannot create directory on M4 (%s) error %v\n", cfg.M4RemotePath, err)
		}

		for _, v := range cfg.DskFiles {
			fmt.Fprintf(os.Stdout, "Attempt to uploading file (%s) on remote path (%s) to host (%s)\n", v, cfg.M4RemotePath, client.IPClient)
			if err := client.Upload(cfg.M4RemotePath, v); err != nil {
				fmt.Fprintf(os.Stderr, "Something is wrong M4 host (%s) local file (%s) remote path (%s) error :%v\n",
					cfg.M4Host,
					v,
					cfg.M4RemotePath,
					err)
			}
		}
	} else {
		if err := client.Remove(cfg.M4RemotePath + "test.sna"); err != nil {
			fmt.Fprintf(os.Stderr, "Cannot create directory on M4 (%s) error %v\n", cfg.M4RemotePath, err)
		}
	}
	if cfg.Dsk {
		dskFile := cfg.Fullpath(".dsk")
		fmt.Fprintf(os.Stdout, "Attempt to uploading file (%s) on remote path (%s) to host (%s)\n", dskFile, cfg.M4RemotePath, client.IPClient)
		if err := client.Upload(cfg.M4RemotePath, dskFile); err != nil {
			fmt.Fprintf(os.Stderr, "Something is wrong M4 host (%s) local file (%s) remote path (%s) error :%v\n",
				cfg.M4Host,
				dskFile,
				cfg.M4RemotePath,
				err)
		}
	}

	if cfg.Sna {
		if err := client.Upload(cfg.M4RemotePath, cfg.SnaPath); err != nil {
			fmt.Fprintf(os.Stderr, "Something is wrong M4 host (%s) local file (%s) remote path (%s) error :%v\n",
				cfg.M4Host,
				cfg.SnaPath,
				cfg.M4RemotePath,
				err)
		}
	}

	if cfg.M4Autoexec {
		if cfg.Sna {
			err := client.Run(cfg.M4RemotePath + "test.sna")
			return err
		}
		p, err := client.Ls(cfg.M4RemotePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot go to the remote path (%s) error :%v\n", cfg.M4RemotePath, err)
		} else {
			fmt.Fprintf(os.Stdout, "Set the remote path (%s) \n", p)
		}

		var overscanFile, basicFile string
		for _, v := range cfg.DskFiles {
			switch path.Ext(v) {
			case ".BAS":
				basicFile = path.Base(v)
			case ".SCR":
				overscanFile = path.Base(v)
			}
		}
		if cfg.Scr {
			fmt.Fprintf(os.Stdout, "Execute basic file (%s)\n", "/"+cfg.M4RemotePath+"/"+basicFile)
			err := client.Run("/" + cfg.M4RemotePath + "/" + basicFile)
			if err != nil {
				return err
			}
		} else {
			if cfg.Overscan {
				fmt.Fprintf(os.Stdout, "Execute overscan file (%s)\n", "/"+cfg.M4RemotePath+"/"+overscanFile)
				err := client.Run("/" + cfg.M4RemotePath + "/" + overscanFile)
				if err != nil {
					return err
				}
			} else {
				fmt.Fprintf(os.Stdout, "Too many importants files, cannot choice.\n")
			}
		}
	}

	return nil
}
