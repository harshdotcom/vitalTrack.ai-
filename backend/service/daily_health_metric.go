package service

import (
	"vita-track-ai/models"
	"vita-track-ai/repository"
)

func SaveHealtMetric(req models.SaveHealthMetricRequest, userId int64) (*models.DailyHealthMetric, error) {
	metric := req.ToModel()
	metric.UploadedBy = userId

	err := metric.Validate()
	if err != nil {
		return nil, err
	}

	err = repository.SaveHealthMetric(metric)

	return metric, err
}
