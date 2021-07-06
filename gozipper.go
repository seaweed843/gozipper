/*MIT License Copyright (c) 2021 seaweed843

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
==============================================================================*/

package gozipper

import (
	"archive/zip"
	"bytes"
	"errors"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type SkipType int

const (
	FileToSkip SkipType = iota
	DirToSkip
)

func ZipPath(srcPath string, dstPathOptional ...string /*nullable*/) (errRet error) {

	var dstDir string
	var dstFileName string

	var skips = map[string]SkipType{
		".DS_Store": FileToSkip,
		"thumbs.db": FileToSkip,
	}

	if len(dstPathOptional) > 0 {
		dstDir = dstPathOptional[0]
	}
	if len(dstPathOptional) > 1 {
		dstFileName = dstPathOptional[1]
	}

	if len(dstPathOptional) > 2 {
		skips = make(map[string]SkipType)
	}
	for i := 2; i < len(dstPathOptional); i++ {
		skips[dstPathOptional[i]] = FileToSkip
	}

	srcPathInfo, errRet := os.Stat(srcPath)
	if errRet != nil {
		log.Fatal(errRet)
		return
	}

	if srcPathInfo.IsDir() {
		errRet = zipFolder(srcPath, dstDir, dstFileName, skips)
	} else {
		errRet = zipFile(srcPath, dstDir, dstFileName)
	}

	return
}

func zipFolder(folderToZip string, dstDir string, dstFileName string, skips map[string]SkipType) (errRet error) {

	absFolderToZip, errRet := filepath.Abs(folderToZip)
	if errRet != nil {
		log.Fatal(errRet)
		return
	}

	if len(dstDir) == 0 {
		dstDir = filepath.Dir(absFolderToZip)
	}

	if len(dstFileName) == 0 {
		dstFileName = filepath.Base(absFolderToZip) + ".zip"
	}

	zippedPath, errRet := filepath.Abs(filepath.Join(dstDir, dstFileName))
	zippedFolderName := strings.Split(dstFileName, ".")[0]

	if _, errRet = os.Stat(zippedPath); !errors.Is(errRet, fs.ErrNotExist) {
		errRet = os.Remove(zippedPath)
		if errRet != nil {
			log.Println(errRet)
			return
		}
	}

	bytesBuffer := new(bytes.Buffer)
	zipWriter := zip.NewWriter(bytesBuffer)
	byteSliceBuffer := make([]byte, 4096)

	errRet = filepath.WalkDir(absFolderToZip, func(path string, fsDirEntry fs.DirEntry, err error) error {
		if err != nil {
			log.Fatal(err)
			return err
		}

		var ioWriter io.Writer

		relativePath, err := filepath.Rel(absFolderToZip, path)
		relativeSlashPath := filepath.ToSlash(filepath.Join(zippedFolderName, relativePath))

		log.Println("visiting", relativeSlashPath)

		fsInfo, _ := fsDirEntry.Info()

		if fsDirEntry.IsDir() {
			if v, found := skips[fsInfo.Name()]; found && v == 1 {
				log.Println("skipping", fsInfo.Name())
				return filepath.SkipDir
			}

			ioWriter, err = zipWriter.Create(relativeSlashPath + "/")
			if err != nil {
				log.Fatal(err)
				return err
			}
			ioWriter.Write(nil)
		} else {
			if v, found := skips[fsInfo.Name()]; found && v == 0 {
				log.Println("skipping", fsInfo.Name())
			} else {

				ioWriter, err = zipWriter.Create(relativeSlashPath)
				if err != nil {
					log.Fatal(err)
					return err
				}

				file, err := os.Open(path)
				if err != nil {
					log.Fatal(err)
					return err
				}

				for {
					readTotal, err := file.Read(byteSliceBuffer)
					if err != nil {
						if err != io.EOF {
							log.Fatal(err)
						}
						break
					}
					_, err = ioWriter.Write(byteSliceBuffer[:readTotal])
					if err != nil {
						log.Fatal(err)
					}
				}
				file.Close()
			}
		}

		return nil
	})

	if errRet != nil {
		log.Fatal(folderToZip, errRet)
	}
	//
	errRet = zipWriter.Close()
	if errRet != nil {
		log.Fatal(errRet)
	}
	//
	errRet = os.WriteFile(zippedPath, bytesBuffer.Bytes(), 0666)
	if errRet != nil {
		log.Fatal(errRet)
	}

	return
}

func zipFile(filePathToZip string, dstDir string, dstFileName string) (errRet error) {

	absFilePathToZip, errRet := filepath.Abs(filePathToZip)
	if errRet != nil {
		log.Fatal(errRet)
		return
	}

	if len(dstDir) == 0 {
		dstDir = filepath.Dir(absFilePathToZip)
	}

	if len(dstFileName) == 0 {
		dstFileName = filepath.Base(absFilePathToZip) + ".zip"
	}

	zippedPath, errRet := filepath.Abs(filepath.Join(dstDir, dstFileName))

	if _, errRet = os.Stat(zippedPath); !errors.Is(errRet, fs.ErrNotExist) {
		errRet = os.Remove(zippedPath)
		if errRet != nil {
			log.Println(errRet)
			return
		}
	}

	bytesBuffer := new(bytes.Buffer)
	zipWriter := zip.NewWriter(bytesBuffer)
	byteSliceBuffer := make([]byte, 4096)

	relativePath, errRet := filepath.Rel(filepath.Dir(absFilePathToZip), absFilePathToZip)
	relativeSlashPath := filepath.ToSlash(relativePath)

	ioWriter, errRet := zipWriter.Create(relativeSlashPath)
	if errRet != nil {
		log.Fatal(errRet)
		return
	}

	log.Println("visiting", relativeSlashPath)

	file, errRet := os.Open(absFilePathToZip)
	if errRet != nil {
		log.Fatal(errRet)
		return
	}

	for {
		readTotal, errRet := file.Read(byteSliceBuffer)
		if errRet != nil {
			if errRet != io.EOF {
				log.Fatal(errRet)
			}
			break
		}
		_, errRet = ioWriter.Write(byteSliceBuffer[:readTotal])
		if errRet != nil {
			log.Fatal(errRet)
		}
	}
	file.Close()

	//
	errRet = zipWriter.Close()
	if errRet != nil {
		log.Fatal(errRet)
	}
	//
	errRet = os.WriteFile(zippedPath, bytesBuffer.Bytes(), 0666)
	if errRet != nil {
		log.Fatal(errRet)
	}

	return
}
