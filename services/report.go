package services

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"

	"github.com/duke-git/lancet/v2/condition"
	"github.com/duke-git/lancet/v2/fileutil"
	"github.com/golang-module/carbon/v2"
	"github.com/lampnick/doctron-client-go"
	"github.com/yoyo-inc/yoyo/common/config"
	"github.com/yoyo-inc/yoyo/common/db"
	"github.com/yoyo-inc/yoyo/common/logger"
	"github.com/yoyo-inc/yoyo/core"
	"github.com/yoyo-inc/yoyo/errs"
	"github.com/yoyo-inc/yoyo/models"
	"github.com/yoyo-inc/yoyo/resources"
)

type (
	ReportCallbackData map[string]interface{}
	ReportCallback     func(startTime string, endTime string) ReportCallbackData
)

type GenerateReportOption struct {
	ReportType string
	StartTime  string
	EndTime    string
	ReportName string
}

var reportCallbacks = make(map[string]ReportCallback)

func RegisterReportCallback(tplName string, callback ReportCallback) {
	reportCallbacks[tplName] = callback
}

func RenderReport(option GenerateReportOption) ([]byte, error) {
	fsys, err := fs.Sub(resources.InternalReportTplDir, "report")
	if err != nil {
		return nil, err
	}
	file, err := fsys.Open(filepath.Join(option.ReportType, "template.html"))
	if err != nil {
		return nil, err
	}
	var b bytes.Buffer
	_, err = b.ReadFrom(file)
	if err != nil {
		return nil, err
	}

	tpl, err := template.New(option.ReportType).Parse(b.String())
	if err != nil {
		return nil, err
	}

	data := make(map[string]interface{})
	callback, ok := reportCallbacks[option.ReportType]
	if ok {
		data = callback(option.StartTime, option.EndTime)
	}
	data["dateTimeRange"] = option.StartTime + "至" + option.EndTime

	var buf bytes.Buffer
	err = tpl.Execute(&buf, data)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), err
}

func GenerateReport(option GenerateReportOption) error {
	url := fmt.Sprintf("http://%s:%d/api/report/preview/%s?startTime=%s&endTime=%s", config.GetString("converter.preview_host"), config.GetInt("server.port"), option.ReportType, option.StartTime, option.EndTime)
	filename := fmt.Sprintf("%s%s.pdf", condition.TernaryOperator(option.ReportName != "", option.ReportName, option.ReportType), carbon.Now().ToShortDateTimeString())
	output := filepath.Join(GetReportRootDir(), filename)

	report := models.Report{
		ReportName:   filename,
		ReportType:   option.ReportType,
		ReportStatus: 0,
	}
	if res := db.Client.Omit("resource_id").Create(&report); res.Error != nil {
		return errs.ErrGenerateReport
	}

	var reportStatus int
	err := ConvertHtml2Pdf(url, output)
	if err != nil {
		reportStatus = 2
	} else {
		reportStatus = 1
	}
	if res := db.Client.Model(&models.Report{IModel: core.IModel{ID: report.ID}}).Update("report_status", reportStatus); res.Error != nil {
		return errs.ErrGenerateReport
	}

	if err != nil {
		return err
	}

	stat, err := os.Lstat(output)
	if err != nil {
		return err
	}

	filesize, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(stat.Size())/1024), 64)
	resource := models.Resource{
		ResourceName: filename,
		ResourceType: "report",
		Filename:     filename,
		Filesize:     filesize,
		FileType:     "pdf",
	}
	if res := db.Client.Create(&resource); res.Error != nil {
		return res.Error
	}

	if res := db.Client.Model(&models.Report{IModel: core.IModel{ID: report.ID}}).Update("resource_id", resource.ID); res.Error != nil {
		return res.Error
	}

	return nil
}

// ConvertHtml2Pdf convert html to pdf file
func ConvertHtml2Pdf(url string, output string) error {
	client := doctron.NewClient(context.Background(), config.GetString("converter.url"), "doctron", "lampnick")
	req := doctron.NewDefaultHTML2PdfRequestDTO()
	req.ConvertURL = url
	req.WaitingTime = config.GetInt("converter.waiting_time")
	res, err := client.HTML2Pdf(req)
	if err != nil {
		return err
	}
	err = os.WriteFile(output, res.Data, 0777)
	if err != nil {
		return err
	}
	return nil
}

func GetReportRootDir() string {
	reportDir := GetTypedResourceDir("report")
	if !fileutil.IsExist(reportDir) {
		err := fileutil.CreateDir(reportDir)
		if err != nil {
			logger.Error(err)
		}
	}
	return reportDir
}
