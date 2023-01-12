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

type FileStat struct {
	Filename string
	Date     string
	Filesize int64
}

func LogFilter(name string) (string, bool) {
	matched := regexp.MustCompile(logFileReg).FindStringSubmatch(name)
	if matched == nil {
		return "", false
	}
	return matched[1], true
}

func ScanLogByRecent(dir string, recent int) ([]FileStat, error) {
	deadline := carbon.Now().SubDays(recent).EndOfDay()

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	stats := make([]FileStat, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			date, ok := LogFilter(entry.Name())
			if !ok {
				break
			}

			if carbon.Parse(date).Lt(deadline) {
				info, err := entry.Info()
				if err != nil {
					logger.Error(err)
					continue
				}
				stats = append(stats, FileStat{
					Filename: entry.Name(),
					Date:     date,
					Filesize: int64(info.Size()),
				})
			}
		}
	}

	return stats, nil
}

func DeleteLogByRecent(dir string, recent int) error {
	stats, err := ScanLogByRecent(dir, recent)
	if err != nil {
		return err
	}

	for _, stat := range stats {
		if err := fileutil.RemoveFile(filepath.Join(dir, stat.Filename)); err != nil {
			logger.Error(err)
		}
	}

	return nil
}

func ArchiveLogByRecent(dir string, recent int) error {
	stats, err := ScanLogByRecent(dir, recent)
	if err != nil {
		return err
	}

	for _, stat := range stats {
		if strings.HasSuffix(stat.Filename, ".zip") {
			continue
		}
		var buf bytes.Buffer
		w := zip.NewWriter(&buf)

		f, err := w.Create(stat.Filename)
		if err != nil {
			logger.Error(err)
			continue
		}

		src, err := os.Open(filepath.Join(dir, stat.Filename))
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

		err = ioutil.WriteFile(filepath.Join(dir, stat.Filename+".zip"), buf.Bytes(), 0644)
		if err != nil {
			logger.Error(err)
			continue
		}

		err = os.Remove(filepath.Join(dir, stat.Filename))
		if err != nil {
			logger.Error(err)
		}
	}

	return nil
}
