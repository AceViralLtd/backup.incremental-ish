package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/aceviralltd/backup.incremental-ish/internal/config"
	"github.com/otiai10/copy"
	"github.com/pierrre/archivefile/zip"
)

var (
	conf      *Config
	runTime   time.Time
	outputDir string
)

type Config struct {
	LastRun time.Time

	Paths []struct {
		Name string
		Path string
	} `yaml:",flow"`
}

func main() {
	runTime = time.Now()

	if 2 > len(os.Args) {
		panic("must give a destination dir for archive to save to")
	}

	if stat, err := os.Stat(os.Args[1]); nil != err {
		panic("output path not found")
	} else if !stat.IsDir() {
		panic("output path is not a directory")
	}

	outputDir = os.Args[1]

	conf = &Config{}
	if err := config.Load("./incremental-ish.yml", conf); nil != err {
		panic("Failed to load config")
	}

	if 0 == len(conf.Paths) {
		panic("no paths to backup")
	}

	fmt.Println(conf)
	for _, entry := range conf.Paths {
		if stat, err := os.Stat(entry.Path); nil != err {
			fmt.Println(entry.Path)
			panic(err)
		} else if !stat.IsDir() {
			panic("not a dir")
		}

		archiveDir(entry.Path, entry.Name)
	}

	conf.LastRun = runTime
	config.Save("./incremental-ish.yml", conf)
}

func archiveDir(rootPath, name string) {
	dir, _ := os.Getwd()
	tmpDir, err := ioutil.TempDir(dir, name)

	if nil != err {
		log.Fatal("tmp", err)
	}

	_ = filepath.Walk(rootPath, func(fullPath string, stat os.FileInfo, err error) error {
		if stat.IsDir() {
			return nil
		}

		if stat.ModTime().Before(conf.LastRun) {
			return nil
		}

		_ = copy.Copy(
			fullPath,
			path.Join(tmpDir, strings.TrimPrefix(fullPath, rootPath)),
		)

		return nil
	})

	archive, err := os.OpenFile(path.Join(outputDir, fmt.Sprintf("%s.zip", name)), os.O_WRONLY|os.O_CREATE, 0755)
	if nil != err {
		fmt.Println("failed to create zip file")
	}

	log.Printf(tmpDir)

	_ = zip.Archive(
		tmpDir,
		archive,
		nil,
	)

	os.RemoveAll(tmpDir)
}
