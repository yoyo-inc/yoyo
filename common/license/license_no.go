//go:build !license

package license

import "github.com/yoyo-inc/yoyo/utils"

var Service utils.License

func Setup() {
	Service = &utils.NoLicense{}
}
