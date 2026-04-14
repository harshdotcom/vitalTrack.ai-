package service

import (
	"vita-track-ai/models"
	"vita-track-ai/repository"
)

func SaveHealtMetric(req models.SaveHealthMetricRequest) (*models.DailyHealthMetric, error) {
	metric := req.ToModel()

	err := metric.Validate()
	if err != nil {
		return nil, err
	}

	err = repository.SaveHealthMetric(metric)

	return metric, err
}
