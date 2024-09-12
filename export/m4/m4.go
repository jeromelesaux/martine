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
	if cfg.M4cfg.Enabled {
		return ErrorNoHostDefined
	}
	if cfg.M4cfg.RemotePath == "" {
		log.GetLogger().Info("No M4 remote path defined, will copy on folder root.")
		cfg.M4cfg.RemotePath = "/"
	}

	client := m4.M4Client{IPClient: cfg.M4cfg.Host}
	if err := client.ResetCpc(); err != nil {
		return err
	}
	if !cfg.ExportType(config.SnaContainer) {
		log.GetLogger().Info("Attempt to create remote directory (%s) to host (%s)\n", cfg.M4cfg.RemotePath, client.IPClient)
		if err := client.MakeDirectory(cfg.M4cfg.RemotePath); err != nil {
			log.GetLogger().Error("Cannot create directory on M4 (%s) error %v\n", cfg.M4cfg.RemotePath, err)
		}

		for _, v := range cfg.DskFiles {
			log.GetLogger().Info("Attempt to uploading file (%s) on remote path (%s) to host (%s)\n", v, cfg.M4cfg.RemotePath, client.IPClient)
			if err := client.Upload(cfg.M4cfg.RemotePath, v); err != nil {
				log.GetLogger().Error("Something is wrong M4 host (%s) local file (%s) remote path (%s) error :%v\n",
					cfg.M4cfg.Host,
					v,
					cfg.M4cfg.RemotePath,
					err)
			}
		}
	} else {
		if err := client.Remove(cfg.M4cfg.RemotePath + "test.sna"); err != nil {
			log.GetLogger().Error("Cannot create directory on M4 (%s) error %v\n", cfg.M4cfg.RemotePath, err)
		}
	}
	if cfg.ExportType(config.DskContainer) {
		dskFile := cfg.Fullpath(".dsk")
		log.GetLogger().Info("Attempt to uploading file (%s) on remote path (%s) to host (%s)\n", dskFile, cfg.M4cfg.RemotePath, client.IPClient)
		if err := client.Upload(cfg.M4cfg.RemotePath, dskFile); err != nil {
			log.GetLogger().Error("Something is wrong M4 host (%s) local file (%s) remote path (%s) error :%v\n",
				cfg.M4cfg.Host,
				dskFile,
				cfg.M4cfg.RemotePath,
				err)
		}
	}

	if cfg.ExportType(config.SnaContainer) {
		if err := client.Upload(cfg.M4cfg.RemotePath, cfg.ContainerCfg.Path); err != nil {
			log.GetLogger().Error("Something is wrong M4 host (%s) local file (%s) remote path (%s) error :%v\n",
				cfg.M4cfg.Host,
				cfg.ContainerCfg.Path,
				cfg.M4cfg.RemotePath,
				err)
		}
	}

	if cfg.M4cfg.Autoexec {
		if cfg.ExportType(config.SnaContainer) {
			return client.Run(cfg.M4cfg.RemotePath + "test.sna")
		}
		p, err := client.Ls(cfg.M4cfg.RemotePath)
		if err != nil {
			log.GetLogger().Error("Cannot go to the remote path (%s) error :%v\n", cfg.M4cfg.RemotePath, err)
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
		if cfg.ScreenCfg.Type == config.ScreenOldFormat {
			log.GetLogger().Info("Execute basic file (%s)\n", "/"+cfg.M4cfg.RemotePath+"/"+basicFile)
			if err := client.Run("/" + cfg.M4cfg.RemotePath + "/" + basicFile); err != nil {
				return err
			}
		} else {
			if cfg.ScreenCfg.Type == config.FullscreenFormat {
				log.GetLogger().Info("Execute overscan file (%s)\n", "/"+cfg.M4cfg.RemotePath+"/"+overscanFile)
				if err := client.Run("/" + cfg.M4cfg.RemotePath + "/" + overscanFile); err != nil {
					return err
				}
			} else {
				log.GetLogger().Info("Too many importants files, cannot choice.\n")
			}
		}
	}

	return nil
}
