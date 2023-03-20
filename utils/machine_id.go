package utils

import (
	"bytes"
	"errors"
	"github.com/yoyo-inc/yoyo/common/logger"
	"os/exec"
	"strings"
)

func GetCpuID() (string, error) {
	var out bytes.Buffer
	var err bytes.Buffer
	cmd := exec.Command("sh", "-c", "dmidecode -t processor | grep ID | head -1 | sed -e 's/[ \\t]*//g' | cut -d ':' -f2")
	cmd.Stdout = &out
	cmd.Stderr = &err
	if err := cmd.Run(); err != nil {
		return "", err
	}
	if err.Len() != 0 {
		return "", errors.New(err.String())
	}

	return out.String(), nil
}

func GetBoardSerialNumber() (string, error) {
	var out bytes.Buffer
	var err bytes.Buffer
	cmd := exec.Command("sh", "-c", "dmidecode -s baseboard-serial-number")
	cmd.Stdout = &out
	cmd.Stderr = &err
	if err := cmd.Run(); err != nil {
		return "", err
	}

	if err.Len() != 0 {
		return "", errors.New(err.String())
	}

	return out.String(), nil
}

func GenMachineID() (string, error) {
	cpuID, err := GetCpuID()
	// throw err when can not get cpu id
	if err != nil {
		return "", err
	}

	bsn, _ := GetBoardSerialNumber()
	rawID := strings.Join([]string{cpuID, bsn}, "-")
	logger.Debugf("raw machineID: %s", rawID)
	return Encrypt(rawID), nil
}
