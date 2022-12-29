package services

import (
	"errors"
	"fmt"

	"github.com/golang-module/carbon/v2"
	"github.com/robfig/cron/v3"
	"github.com/yoyo-inc/yoyo/common/db"
	"github.com/yoyo-inc/yoyo/common/logger"
	"github.com/yoyo-inc/yoyo/common/sched"
	"github.com/yoyo-inc/yoyo/models"
	"github.com/yoyo-inc/yoyo/services/audit_log"
)

var (
	schedJobMappings    = make(map[string]cron.EntryID)
	ErrExistSchedJob    = errors.New("定时任务已存在")
	ErrNotExistSchedJob = errors.New("定时任务不存在")
	ErrStopSchedJob     = errors.New("停止定时任务失败")
	ErrCreateSchedJob   = errors.New("创建定时任务失败")
	ErrRemoveSchedJob   = errors.New("删除定时任务失败")
)

// AddSchedJob creates sched job
func AddSchedJob(jobID string, jobType string, description string, spec string, job func() error) error {
	ErrExistSchedJob = fmt.Errorf("定时任务(%s)已存在", jobID)
	// judge whether exists the jobID
	if _, ok := schedJobMappings[jobID]; ok {
		logger.Error(ErrExistSchedJob)
		return ErrExistSchedJob
	}

	id, err := sched.C.AddFunc(spec, func() {
		if err := job(); err != nil {
			audit_log.Fail(nil, "定时任务", "执行", fmt.Sprintf("任务ID(%s),任务描述(%s)", jobID, description))
		} else {
			audit_log.Success(nil, "定时任务", "执行", fmt.Sprintf("任务ID(%s),任务描述(%s)", jobID, description))
			if res := db.Client.Model(&models.SchedJob{}).Where("job_id = ?", jobID).Update("last_run_time", carbon.Now().ToDateTimeString()); res.Error != nil {
				logger.Error(res.Error)
			}
		}
	})
	if err != nil {
		return err
	}
	// create relationship between jobID with cron.EntryID
	schedJobMappings[jobID] = id

	var count int64
	if res := db.Client.Model(&models.SchedJob{}).Where("job_id = ?", jobID).Count(&count); res.Error != nil {
		logger.Error(res.Error)
		return ErrCreateSchedJob
	}

	if count == 0 {
		if res := db.Client.Create(&models.SchedJob{
			Spec:        spec,
			Description: description,
			Status:      1,
			JobID:       jobID,
			Type:        jobType,
		}); res.Error != nil {
			logger.Error(res.Error)
			return ErrCreateSchedJob
		}
	}

	return nil
}

// StartSchedJob starts sched job
func StartSchedJob(jobID string) error {
	ErrNotExistSchedJob = fmt.Errorf("定时任务(%s)不存在", jobID)
	ErrStopSchedJob = fmt.Errorf("开启定时任务(%s)失败", jobID)
	entryID, ok := schedJobMappings[jobID]
	if !ok {
		logger.Error(ErrExistSchedJob)
		return ErrExistSchedJob
	}

	sched.C.Resume(entryID)

	if res := db.Client.Model(&models.SchedJob{}).Where("job_id = ?", jobID).Update("status", 1); res.Error != nil {
		logger.Error(res.Error)
		return ErrStopSchedJob
	}

	return nil
}

// StopSchedJob stops sched job
func StopSchedJob(jobID string) error {
	ErrNotExistSchedJob = fmt.Errorf("定时任务(%s)不存在", jobID)
	ErrStopSchedJob = fmt.Errorf("关闭定时任务(%s)失败", jobID)
	entryID, ok := schedJobMappings[jobID]
	if !ok {
		return ErrExistSchedJob
	}

	sched.C.Pause(entryID)

	if res := db.Client.Model(&models.SchedJob{}).Where("job_id = ?", jobID).Update("status", 0); res.Error != nil {
		return ErrStopSchedJob
	}

	return nil
}

// RemoveSchedJob remove sched job by jobID
func RemoveSchedJob(jobID string) error {
	ErrRemoveSchedJob = fmt.Errorf("删除定时任务(%s)失败", jobID)
	entryID, ok := schedJobMappings[jobID]
	if !ok {
		return nil
	}

	sched.C.Remove(entryID)

	if res := db.Client.Delete(&models.SchedJob{}, "job_id = ?", jobID); res.Error != nil {
		return ErrRemoveSchedJob
	}

	delete(schedJobMappings, jobID)

	return nil
}
