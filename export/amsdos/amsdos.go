package amsdos

import (
	"encoding/binary"
	"os"
	"path/filepath"
	"strings"

	"github.com/jeromelesaux/m4client/cpc"
	"github.com/jeromelesaux/martine/config"
	"github.com/jeromelesaux/martine/log"
)

func AmsdosFilename(inputPath, ext string) string {
	file := strings.ToUpper(filepath.Base(inputPath))
	filename := config.RemoveUnsupportedChar(strings.TrimSuffix(file, filepath.Ext(file)))
	filenameSize := len(filename)
	if filenameSize > 8 {
		filenameSize = 8
	}
	osFile := make([]byte, filenameSize)
	copy(osFile, filename[0:filenameSize])
	return string(osFile) + ext
}

func SaveAmsdosFile(filename, extension string, data []byte, fileType, user byte, loadingAddress, executionAddress uint16) error {
	filesize := len(data)
	header := cpc.CpcHead{
		Type: fileType, User: user, Address: loadingAddress, Exec: executionAddress,
		Size:        uint16(filesize),
		Size2:       uint16(filesize),
		LogicalSize: uint16(filesize),
	}
	cpcFilename := AmsdosFilename(filename, extension)
	copy(header.Filename[:], strings.Replace(cpcFilename, ".", "", -1))
	header.Checksum = uint16(header.ComputedChecksum16())
	log.GetLogger().Error("filesize:%d,#%.2x\n", filesize, filesize)
	log.GetLogger().Error("Data length %d\n", binary.Size(data))
	fw, err := os.Create(filename)
	if err != nil {
		log.GetLogger().Error("Error while creating file (%s) error :%s\n", filename, err)
		return err
	}
	err = binary.Write(fw, binary.LittleEndian, header)
	if err != nil {
		return err
	}
	err = binary.Write(fw, binary.LittleEndian, data)
	if err != nil {
		return err
	}

	return fw.Close()
}

func SaveOSFile(filename string, data []byte) error {
	fw, err := os.Create(filename)
	if err != nil {
		log.GetLogger().Error("Error while creating file (%s) error :%s\n", filename, err)
		return err
	}
	err = binary.Write(fw, binary.LittleEndian, data)
	if err != nil {
		return err
	}

	return fw.Close()
}

func SaveStringOSFile(filename string, data string) error {
	fw, err := os.Create(filename)
	if err != nil {
		log.GetLogger().Error("Error while creating file (%s) error :%s\n", filename, err)
		return err
	}
	_, err = fw.WriteString(data)
	if err != nil {
		return err
	}
	return fw.Close()
}
