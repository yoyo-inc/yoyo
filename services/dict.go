package services

import (
	"strconv"

	"github.com/yoyo-inc/yoyo/common/db"
	"github.com/yoyo-inc/yoyo/models"
)

type Entry struct {
	Label string `json:"label"`
	Value any    `json:"value"`
}

func GetEntriesByType(t string) ([]Entry, error) {
	var dicts []models.Dict
	if res := db.Client.Model(&models.Dict{}).Where("type = ?", t).Find(&dicts); res.Error != nil {
		return nil, res.Error
	}

	entries := make([]Entry, 0, len(dicts))
	for _, dict := range dicts {
		v, err := format(dict)
		if err == nil {
			entries = append(entries, Entry{
				Label: dict.Label,
				Value: v,
			})
		}
	}

	return entries, nil
}

func GetLabelByValue(t string, v string) string {
	var dict models.Dict
	if res := db.Client.Where("type = ?", t).Where("value = ?", v).First(&dict); res.Error != nil {
		// return value when error raise
		return v
	}

	return dict.Label
}

func format(dict models.Dict) (any, error) {
	switch dict.ValueType {
	case "integer":
		return strconv.Atoi(dict.Value)
	case "string":
		fallthrough
	default:
		return dict.Value, nil
	}
}
