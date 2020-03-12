package service

import (
	"log"
	"os/exec"
	"path/filepath"
	"github.com/alphayan/utils/config"
)

var (
	// libreoffice路径
	_libreoffice_path = "C:/Program Files/LibreOffice 5/program/soffice.exe"
	ConvertFileTypes  = []string{".doc", ".docx", ".ppt", ".pptx", ".xls", ".xlsx"}
)

func GetAbsPathFromProj(s string) string {
	_projPath := config.Config.StaticPath
	return filepath.Join(_projPath, s)
}

func GetPathFromProj(s string) string {
	return config.Config.GetAbsProPath(s)
}

// WORD、PPT转PDF
func ConvertToPdf(inPath, outDirPath string) *exec.Cmd {
	libreofficePath := config.Config.GetString("libreoffice_path", "")
	if libreofficePath != "" {
		_libreoffice_path = libreofficePath
	}
	log.Println("_libreoffice_path", _libreoffice_path)
	cmd := exec.Command(_libreoffice_path,
		"--headless", "--convert-to", "pdf")
	if outDirPath != "" {
		outDirPath, _ = filepath.Abs(outDirPath)
		cmd.Args = append(cmd.Args, "--outdir", outDirPath)
	}

	inPath, _ = filepath.Abs(inPath)
	cmd.Args = append(cmd.Args, inPath)
	return cmd
}
