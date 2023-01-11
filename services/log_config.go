package services

import (
	"archive/zip"
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/duke-git/lancet/v2/fileutil"
	"github.com/golang-module/carbon/v2"
	"github.com/yoyo-inc/yoyo/common/logger"
)

var logFileReg = `.*(\d{4}-\d{2}-\d{2}).(log|zip)`

func logFilter(name string) (string, bool) {
	matched := regexp.MustCompile(logFileReg).FindStringSubmatch(name)
	if matched == nil {
		return "", false
	}
	return matched[1], true
}

func ScanLogByRecent(dir string, recent int) ([]string, error) {
	deadline := carbon.Now().SubDays(recent).StartOfDay()

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			date, ok := logFilter(entry.Name())
			if !ok {
				break
			}

			if carbon.Parse(date).Lt(deadline) {
				names = append(names, entry.Name())
			}
		}
	}

	return names, nil
}

func DeleteLogByRecent(dir string, recent int) error {
	names, err := ScanLogByRecent(dir, recent)
	if err != nil {
		return err
	}

	for _, name := range names {
		if err := fileutil.RemoveFile(filepath.Join(dir, name)); err != nil {
			logger.Error(err)
		}
	}

	return nil
}

func ArchiveLogByRecent(dir string, recent int) error {
	names, err := ScanLogByRecent(dir, recent)
	if err != nil {
		return err
	}

	for _, name := range names {
		if strings.HasSuffix(name, ".zip") {
			continue
		}
		var buf bytes.Buffer
		w := zip.NewWriter(&buf)

		f, err := w.Create(name)
		if err != nil {
			logger.Error(err)
			continue
		}

		src, err := os.Open(filepath.Join(dir, name))
		if err != nil {
			logger.Error(err)
			continue
		}
		_, err = io.Copy(f, src)
		if err != nil {
			logger.Error(err)
			continue
		}
		if err := w.Close(); err != nil {
			logger.Error(err)
			continue
		}

		err = ioutil.WriteFile(name+".zip", buf.Bytes(), 0644)
		if err != nil {
			logger.Error(err)
		}
	}

	return nil
}
