package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bi-zone/go-fileversion"
)

var bases = []string{
	"C:\\Program Files (x86)\\Microsoft\\Edge\\",
	"C:\\Program Files (x86)\\Microsoft\\EdgeWebView\\",
	"C:\\Program Files\\Microsoft\\Edge\\",
	"C:\\Program Files\\Microsoft\\EdgeWebView\\",
}
var files []string
var edgeExe = "msedge.exe"

// test
func AppendBases(path string) {
	bases = append(bases, filepath.Join(path, "AppData", "Local", "Microsoft", "Edge"))
	bases = append(bases, filepath.Join(path, "AppData", "Local", "Microsoft", "EdgeWebView"))

}

func GetBases() {
	baseDir := "C:\\Users\\"
	err := filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && path != baseDir {
			AppendBases(path)
			return filepath.SkipDir
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
}

func VisitFile(path string, info os.FileInfo, err error) error {
	if err != nil {

		return nil
	}
	file_name := info.Name()
	if !info.IsDir() && strings.ToLower(file_name) == edgeExe {

		files = append(files, path)
	}
	return nil
}

func RemoveFile(file string) {
	err := os.Remove(file)
	if err != nil {
		log.Printf("Error on file: %s", file)
		log.Printf("Error: %v", err)
	}
}

func main() {
	fmt.Println("ZacatekProgramu")

	f, err := os.OpenFile("C:\\_install\\logRemoveOldEdge", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	log.SetOutput(f)
	log.Printf("Program started: %s", time.Now())

	var vmajor int
	var vminor int
	var vbuild int
	var vpatch int

	flag.IntVar(&vmajor, "vmajor", 0, "Major version of file")
	flag.IntVar(&vminor, "vminor", 0, "Minor version of file")
	flag.IntVar(&vbuild, "vbuild", 0, "Build version of file")
	flag.IntVar(&vpatch, "vpatch", 0, "Patch version of file")
	flag.Parse()
	if flag.NFlag() < 4 {
		flag.Usage()
		return
	}

	GetBases()

	for _, base := range bases {

		err := filepath.Walk(base, VisitFile)

		if err != nil {
			log.Fatal(err)
		}
	}

	for _, file := range files {
		fmt.Println(file)
		f, err := fileversion.New(file)
		if err != nil {
			fmt.Println(err)
		}
		fmajor := f.FixedInfo().FileVersion.Major
		fminor := f.FixedInfo().FileVersion.Minor
		fbuild := f.FixedInfo().FileVersion.Build
		fpatch := f.FixedInfo().FileVersion.Patch
		if fmajor < uint16(vmajor) {
			RemoveFile(file)
		} else if fmajor == uint16(vmajor) {
			if fminor < uint16(vminor) {
				RemoveFile(file)
			} else if fminor == uint16(vminor) {
				if fbuild < uint16(vbuild) {
					RemoveFile(file)
				} else if fbuild == uint16(vbuild) {
					if fpatch < uint16(vpatch) {
						RemoveFile(file)
					}
				}
			}

		}
	}
}
