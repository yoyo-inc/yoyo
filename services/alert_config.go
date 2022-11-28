package services

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"path"
	"text/template"

	"github.com/duke-git/lancet/v2/fileutil"
	"github.com/duke-git/lancet/v2/netutil"
	"github.com/duke-git/lancet/v2/slice"
	jsoniter "github.com/json-iterator/go"
	"github.com/yoyo-inc/yoyo/common/logger"
	"github.com/yoyo-inc/yoyo/models"
	"github.com/yoyo-inc/yoyo/resources"
	"github.com/yoyo-inc/yoyo/vo"
)

const (
	prometheusConfigDirPath     = "/etc/prometheus"            // prometheus 配置目录
	prometheusRulesDirPath      = "/etc/prometheus/rules"      // prometheus 规则目录
	alertmanagerConfigDirPath   = "/etc/alertmanager"          // alertmanager 配置目录
	alertmanagerTemplateDirPath = "/etc/alertmanager/template" // alertmanager 模板路径
	configFileMode              = 0644
)

// GeneratePrometheusConfig generates prometheus config
func GeneratePrometheusConfig() (err error) {
	if !fileutil.IsExist(prometheusConfigDirPath) {
		if err = fileutil.CreateDir(prometheusConfigDirPath); err != nil {
			return
		}
	}
	if !fileutil.IsExist(prometheusRulesDirPath) {
		if err = fileutil.CreateDir(prometheusRulesDirPath); err != nil {
			return
		}
	}

	prometheusConfig := path.Join(prometheusConfigDirPath, "prometheus.yml")
	if !fileutil.IsExist(prometheusConfig) {
		tpl, _ := resources.PrometheusDir.ReadFile("prometheus.yml")
		if err = ioutil.WriteFile(prometheusConfig, tpl, configFileMode); err != nil {
			return
		}
	}

	hostRules := path.Join(prometheusRulesDirPath, "host.rules")
	if !fileutil.IsExist(hostRules) {
		tpl, _ := resources.PrometheusDir.ReadFile("rules/host.rules")
		if err = ioutil.WriteFile(hostRules, tpl, configFileMode); err != nil {
			return
		}
	}

	servicesRules := path.Join(prometheusRulesDirPath, "services.rules")
	if !fileutil.IsExist(servicesRules) {
		tpl, _ := resources.PrometheusDir.ReadFile("rules/services.rules")
		if err = ioutil.WriteFile(servicesRules, tpl, configFileMode); err != nil {
			return
		}
	}

	// reload prometheus
	err = reload("localhost:9090/-/reload")
	return err
}

// GenerateAlertManagerConfig generates alertmanager config
func GenerateAlertManagerConfig(alertConfig models.AlertConfig) (err error) {
	logger.Info("Start to generate alertmanager config")
	if !fileutil.IsExist(alertmanagerConfigDirPath) {
		if err = fileutil.CreateDir(alertmanagerConfigDirPath); err != nil {
			return
		}
	}
	if !fileutil.IsExist(alertmanagerTemplateDirPath) {
		if err = fileutil.CreateDir(alertmanagerTemplateDirPath); err != nil {
			return
		}
	}

	alertmanagerConfigFile := path.Join(alertmanagerConfigDirPath, "alertmanager.yml")
	if !fileutil.IsExist(alertmanagerConfigFile) {
		tpl, _ := resources.AlertmanagerDir.ReadFile("alertmanager.yml")
		var buf bytes.Buffer
		t := template.Must(template.New("alertmanager").Parse(string(tpl)))

		var receivers []vo.SmtpReceiver
		json := jsoniter.ConfigCompatibleWithStandardLibrary
		json.Unmarshal(alertConfig.SmtpReceivers, &receivers)

		if err = t.Execute(&buf, map[string]interface{}{
			"EmailEnable":      alertConfig.EmailEnable,
			"SmtpServer":       alertConfig.SmtpServer,
			"SmtpSender":       alertConfig.SmtpSender,
			"SmtpAuthUser":     alertConfig.SmtpAuthUser,
			"SmtpAuthPassword": alertConfig.SmtpAuthPassword,
			"SmtpReceivers": slice.Filter(receivers, func(_ int, receiver vo.SmtpReceiver) bool {
				return receiver.Enable
			}),
		}); err != nil {
			return
		}
		if err = ioutil.WriteFile(alertmanagerConfigFile, tpl, configFileMode); err != nil {
			return
		}
	}

	alertTemplateFile := path.Join(alertmanagerTemplateDirPath, "alert.yml")
	if !fileutil.IsExist(alertTemplateFile) {
		tpl, _ := resources.AlertmanagerDir.ReadFile("template/alert.yml")
		if err = ioutil.WriteFile(alertTemplateFile, tpl, configFileMode); err != nil {
			return
		}
	}
	logger.Info("Success to generate alertmanager config")

	// reload alertmanager
	err = reload("localhost:9093/-/reload")
	return
}

func reload(url string) (err error) {
	// reload alertmanager
	var res *http.Response
	res, err = netutil.NewHttpClientWithConfig(&netutil.HttpClientConfig{ResponseTimeout: 5}).Post(url, "application/json", nil)
	if err != nil {
		return
	}
	if res.StatusCode != 200 {
		var buf bytes.Buffer
		_, err = buf.ReadFrom(res.Body)
		if err != nil {
			return
		}
		return errors.New(buf.String())
	}
	return
}
