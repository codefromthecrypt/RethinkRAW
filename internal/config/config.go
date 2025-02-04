package config

import (
	"mime"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/ncruces/go-exiftool"
	"github.com/ncruces/rethinkraw/internal/dcraw"
	"github.com/ncruces/rethinkraw/pkg/osutil"
)

var (
	ServerMode                bool
	BaseDir, DataDir, TempDir string
	DngConverter              string
)

func init() {
	mime.AddExtensionType(".dng", "image/x-adobe-dng")
}

func SetupPaths() (err error) {
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	if exe, err := filepath.EvalSymlinks(exe); err != nil {
		return err
	} else {
		BaseDir = filepath.Dir(exe)
	}

	DataDir = filepath.Join(BaseDir, "data")
	TempDir = filepath.Join(os.TempDir(), "RethinkRAW")

	TempDir, err = osutil.GetANSIPath(TempDir)
	if err != nil {
		return err
	}

	switch runtime.GOOS {
	case "windows":
		ServerMode = filepath.Base(exe) == "RethinkRAW.com"
		dcraw.Path = BaseDir + `\utils\dcraw.wasm`
		exiftool.Exec = BaseDir + `\utils\exiftool\exiftool.exe`
		exiftool.Arg1 = strings.TrimSuffix(exiftool.Exec, ".exe")
		exiftool.Config = BaseDir + `\utils\exiftool_config.pl`
		DngConverter = os.Getenv("PROGRAMFILES") + `\Adobe\Adobe DNG Converter\Adobe DNG Converter.exe`
	case "darwin":
		ServerMode = filepath.Base(exe) == "rethinkraw-server"
		dcraw.Path = BaseDir + "/utils/dcraw.wasm"
		exiftool.Exec = BaseDir + "/utils/exiftool/exiftool"
		exiftool.Config = BaseDir + "/utils/exiftool_config.pl"
		DngConverter = "/Applications/Adobe DNG Converter.app/Contents/MacOS/Adobe DNG Converter"
	}

	if testDataDir() == nil {
		return nil
	}
	if data, err := os.UserConfigDir(); err != nil {
		return err
	} else {
		DataDir = filepath.Join(data, "RethinkRAW")
	}
	return testDataDir()
}

func testDataDir() error {
	if err := os.MkdirAll(DataDir, 0700); err != nil {
		return err
	}
	if f, err := os.Create(filepath.Join(DataDir, "lastrun")); err != nil {
		return err
	} else {
		return f.Close()
	}
}
