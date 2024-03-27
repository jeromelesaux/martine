package m4

import (
	"errors"
	"path"

	"github.com/jeromelesaux/m4client/m4"
	"github.com/jeromelesaux/martine/config"
	"github.com/jeromelesaux/martine/log"
)

var ErrorNoHostDefined = errors.New("no host defined")

// nolint:funlen, gocognit
func ImportInM4(cfg *config.MartineConfig) error {
	if cfg.M4.Host == "" {
		return ErrorNoHostDefined
	}
	if cfg.M4.RemotePath == "" {
		log.GetLogger().Info("No M4 remote path defined, will copy on folder root.")
		cfg.M4.RemotePath = "/"
	}

	client := m4.M4Client{IPClient: cfg.M4.Host}
	if err := client.ResetCpc(); err != nil {
		return err
	}
	if !cfg.Sna.Enabled {
		log.GetLogger().Info("Attempt to create remote directory (%s) to host (%s)\n", cfg.M4.RemotePath, client.IPClient)
		if err := client.MakeDirectory(cfg.M4.RemotePath); err != nil {
			log.GetLogger().Error("Cannot create directory on M4 (%s) error %v\n", cfg.M4.RemotePath, err)
		}

		for _, v := range cfg.DskFiles {
			log.GetLogger().Info("Attempt to uploading file (%s) on remote path (%s) to host (%s)\n", v, cfg.M4.RemotePath, client.IPClient)
			if err := client.Upload(cfg.M4.RemotePath, v); err != nil {
				log.GetLogger().Error("Something is wrong M4 host (%s) local file (%s) remote path (%s) error :%v\n",
					cfg.M4.Host,
					v,
					cfg.M4.RemotePath,
					err)
			}
		}
	} else {
		if err := client.Remove(cfg.M4.RemotePath + "test.sna"); err != nil {
			log.GetLogger().Error("Cannot create directory on M4 (%s) error %v\n", cfg.M4.RemotePath, err)
		}
	}
	if cfg.Dsk {
		dskFile := cfg.Fullpath(".dsk")
		log.GetLogger().Info("Attempt to uploading file (%s) on remote path (%s) to host (%s)\n", dskFile, cfg.M4.RemotePath, client.IPClient)
		if err := client.Upload(cfg.M4.RemotePath, dskFile); err != nil {
			log.GetLogger().Error("Something is wrong M4 host (%s) local file (%s) remote path (%s) error :%v\n",
				cfg.M4.Host,
				dskFile,
				cfg.M4.RemotePath,
				err)
		}
	}

	if cfg.Sna.Enabled {
		if err := client.Upload(cfg.M4.RemotePath, cfg.Sna.SnaPath); err != nil {
			log.GetLogger().Error("Something is wrong M4 host (%s) local file (%s) remote path (%s) error :%v\n",
				cfg.M4.Host,
				cfg.Sna.SnaPath,
				cfg.M4.RemotePath,
				err)
		}
	}

	if cfg.M4.Autoexec {
		if cfg.Sna.Enabled {
			return client.Run(cfg.M4.RemotePath + "test.sna")
		}
		p, err := client.Ls(cfg.M4.RemotePath)
		if err != nil {
			log.GetLogger().Error("Cannot go to the remote path (%s) error :%v\n", cfg.M4.RemotePath, err)
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
			log.GetLogger().Info("Execute basic file (%s)\n", "/"+cfg.M4.RemotePath+"/"+basicFile)
			if err := client.Run("/" + cfg.M4.RemotePath + "/" + basicFile); err != nil {
				return err
			}
		} else {
			if cfg.Overscan {
				log.GetLogger().Info("Execute overscan file (%s)\n", "/"+cfg.M4.RemotePath+"/"+overscanFile)
				if err := client.Run("/" + cfg.M4.RemotePath + "/" + overscanFile); err != nil {
					return err
				}
			} else {
				log.GetLogger().Info("Too many importants files, cannot choice.\n")
			}
		}
	}

	return nil
}
