package m4

import (
	"errors"
	"path"

	"github.com/jeromelesaux/m4client/m4"
	"github.com/jeromelesaux/martine/config"
	"github.com/jeromelesaux/martine/log"
)

var ErrorNoHostDefined = errors.New("no host defined")

func ImportInM4(cfg *config.MartineConfig) error {
	if cfg.M4Host == "" {
		return ErrorNoHostDefined
	}
	if cfg.M4RemotePath == "" {
		log.GetLogger().Info("No M4 remote path defined, will copy on folder root.")
		cfg.M4RemotePath = "/"
	}

	client := m4.M4Client{IPClient: cfg.M4Host}
	err := client.ResetCpc()
	if err != nil {
		return err
	}
	if !cfg.Sna {
		log.GetLogger().Info("Attempt to create remote directory (%s) to host (%s)\n", cfg.M4RemotePath, client.IPClient)
		if err := client.MakeDirectory(cfg.M4RemotePath); err != nil {
			log.GetLogger().Error("Cannot create directory on M4 (%s) error %v\n", cfg.M4RemotePath, err)
		}

		for _, v := range cfg.DskFiles {
			log.GetLogger().Info("Attempt to uploading file (%s) on remote path (%s) to host (%s)\n", v, cfg.M4RemotePath, client.IPClient)
			if err := client.Upload(cfg.M4RemotePath, v); err != nil {
				log.GetLogger().Error("Something is wrong M4 host (%s) local file (%s) remote path (%s) error :%v\n",
					cfg.M4Host,
					v,
					cfg.M4RemotePath,
					err)
			}
		}
	} else {
		if err := client.Remove(cfg.M4RemotePath + "test.sna"); err != nil {
			log.GetLogger().Error("Cannot create directory on M4 (%s) error %v\n", cfg.M4RemotePath, err)
		}
	}
	if cfg.Dsk {
		dskFile := cfg.Fullpath(".dsk")
		log.GetLogger().Info("Attempt to uploading file (%s) on remote path (%s) to host (%s)\n", dskFile, cfg.M4RemotePath, client.IPClient)
		if err := client.Upload(cfg.M4RemotePath, dskFile); err != nil {
			log.GetLogger().Error("Something is wrong M4 host (%s) local file (%s) remote path (%s) error :%v\n",
				cfg.M4Host,
				dskFile,
				cfg.M4RemotePath,
				err)
		}
	}

	if cfg.Sna {
		if err := client.Upload(cfg.M4RemotePath, cfg.SnaPath); err != nil {
			log.GetLogger().Error("Something is wrong M4 host (%s) local file (%s) remote path (%s) error :%v\n",
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
			log.GetLogger().Error("Cannot go to the remote path (%s) error :%v\n", cfg.M4RemotePath, err)
		} else {
			log.GetLogger().Info("Set the remote path (%s) \n", p)
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
			log.GetLogger().Info("Execute basic file (%s)\n", "/"+cfg.M4RemotePath+"/"+basicFile)
			err := client.Run("/" + cfg.M4RemotePath + "/" + basicFile)
			if err != nil {
				return err
			}
		} else {
			if cfg.Overscan {
				log.GetLogger().Info("Execute overscan file (%s)\n", "/"+cfg.M4RemotePath+"/"+overscanFile)
				err := client.Run("/" + cfg.M4RemotePath + "/" + overscanFile)
				if err != nil {
					return err
				}
			} else {
				log.GetLogger().Info("Too many importants files, cannot choice.\n")
			}
		}
	}

	return nil
}
