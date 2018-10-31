package xtools

import (
	"os"
	"github.com/cihub/seelog"
	"io/ioutil"
	"bytes"
	"io"
	"archive/zip"
)

func Unzip(file, dir string) error {
	bs, err := ioutil.ReadFile(file)
	if err != nil {
		seelog.Info("File Error:", file, err)
		return err
	}
	return UnzipByte(bs, dir)
}

func UnzipByte(bs []byte, dir string) error {
	os.MkdirAll(dir, 0777)
	reader, err := zip.NewReader(bytes.NewReader(bs), int64(len(bs)))
	if err != nil {
		return err
	}
	for _, file := range reader.File {
		if file.FileInfo().IsDir() {
			os.Mkdir(dir+"/"+file.Name, 0777)
		} else {
			w, err := os.Create(dir + "/" + file.Name)
			if err != nil {
				continue
			}
			rc, err := file.Open()
			if err != nil {
				w.Close()
				continue
			}
			io.Copy(w, rc)
			w.Close()
			rc.Close()
		}
	}
	return nil
}
