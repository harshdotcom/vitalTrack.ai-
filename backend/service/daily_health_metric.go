package service

import (
	"vita-track-ai/models"
	"vita-track-ai/repository"
)

func SaveHealtMetric(req models.SaveHealthMetricRequest, userId int64) (*models.DailyHealthMetric, error) {
	if req.Timestamp != nil && *req.Timestamp != "" {
		if _, err := req.ResolveTimestamp(); err != nil {
			return nil, err
		}
	}

	metric := req.ToModel()
	metric.UploadedBy = userId

	err := metric.Validate()
	if err != nil {
		return nil, err
	}

	err = repository.SaveHealthMetric(metric)

	return metric, err
}

func DeleteHealthMetric(id string, userID int64) error {
	return repository.DeleteHealthMetric(id, userID)
}
