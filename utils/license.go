//go:build license

package utils

import (
	"github.com/golang-module/carbon/v2"
	jsoniter "github.com/json-iterator/go"
)

//#cgo LDFLAGS: -llics
//extern char* EchoVerify(_GoString_ code);
import "C"
import (
	"github.com/yoyo-inc/yoyo/common/logger"
)

type License interface {
	GenMachineID() string
	Activate(code string) bool
}

type SNLicenseContent struct {
	MachineCode string            `json:"machinecode"`
	CreateTime  int64             `json:"createtime"`
	Deadline    int               `json:"deadline"`
	Echo        int               `json:"echo"`
	Info        map[string]string `json:"info"`
}

type SNLicense struct {
}

func (*SNLicense) GenMachineID() string {
	machineID, err := GenMachineID()
	if err != nil {
		logger.Error(err)
		return ""
	}
	return machineID
}

func (snl *SNLicense) Activate(code string) bool {
	rawContent := C.GoString(C.EchoVerify(code))
	if rawContent == "" {
		return false
	}

	var content SNLicenseContent
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	if err := json.Unmarshal([]byte(rawContent), &content); err != nil {
		logger.Error(err)
		return false
	}

	if content.Echo != 1 {
		logger.Errorf("This activation code can not echo: %s", code)
		logger.Debugf("Activation code content: %s", rawContent)
		return false
	}

	machineID := snl.GenMachineID()
	if machineID == "" {
		logger.Error("Failed to gen machine id")
		return false
	}

	if machineID != content.MachineCode {
		logger.Error("MachineID not equal")
		return false
	}

	expiredTime := carbon.CreateFromTimestamp(content.CreateTime).AddMonths(content.Deadline)
	if carbon.Now().Gt(expiredTime) {
		return false
	}

	return true
}
