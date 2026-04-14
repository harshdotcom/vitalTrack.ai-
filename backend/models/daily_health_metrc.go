package models

import (
	"errors"
	"time"
)

type DailyHealthMetric struct {
	ID            string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UploadedBy    int64
	HeartRate     *int           // bpm
	Weight        *Measurement   `gorm:"type:jsonb"` // supports units
	BloodPressure *BloodPressure `gorm:"type:jsonb"`
	BloodSugar    *Measurement   `gorm:"type:jsonb"`
	Notes         *string
	SleepHours    *float64
	Steps         *int
	Calories      *int
	OxygenLevel   *float64
	Timestamp     time.Time
	User          User `gorm:"foreignKey:UploadedBy;references:UserId;constraint:OnDelete:CASCADE"`
}

type Measurement struct {
	Value float64 `json:"value"`
	Unit  string  `json:"unit"` // e.g. "kg", "lb", "mg/dL", "mmol/L"
}

type BloodPressure struct {
	Systolic  int // mmHg
	Diastolic int // mmHg
}

func (d *DailyHealthMetric) Validate() error {

	if d.HeartRate != nil {
		if *d.HeartRate < 30 || *d.HeartRate > 220 {
			return errors.New("heart rate out of realistic range")
		}
	}

	if d.BloodPressure != nil {
		if d.BloodPressure.Systolic < 70 || d.BloodPressure.Systolic > 250 {
			return errors.New("invalid systolic pressure")
		}
		if d.BloodPressure.Diastolic < 40 || d.BloodPressure.Diastolic > 150 {
			return errors.New("invalid diastolic pressure")
		}
	}

	if d.Weight != nil {
		if d.Weight.Value <= 0 {
			return errors.New("weight must be positive")
		}
	}

	return nil
}

type SaveHealthMetricRequest struct {
	HeartRate *int         `json:"heart_rate" example:"72"`
	Weight    *Measurement `json:"weight"`

	BloodPressure *struct {
		Systolic  int `json:"systolic" example:"120"`
		Diastolic int `json:"diastolic" example:"80"`
	} `json:"blood_pressure"`

	BloodSugar *Measurement `json:"blood_sugar"`

	SleepHours  *float64 `json:"sleep_hours" example:"7.5"` // hours (e.g., 7.5)
	Steps       *int     `json:"steps" example:"8500"`      // daily steps
	Calories    *int     `json:"calories" example:"2200"`   // kcal burned/consumed
	OxygenLevel *float64 `json:"oxygen_level" example:"98"` // SpO2 (%)

	Notes *string `json:"notes"`
}

func (r *SaveHealthMetricRequest) ToModel() *DailyHealthMetric {
	var bp *BloodPressure

	if r.BloodPressure != nil {
		bp = &BloodPressure{
			Systolic:  r.BloodPressure.Systolic,
			Diastolic: r.BloodPressure.Diastolic,
		}
	}

	timestamp := time.Now()

	return &DailyHealthMetric{
		Timestamp:     timestamp,
		HeartRate:     r.HeartRate,
		Weight:        r.Weight,
		BloodPressure: bp,
		BloodSugar:    r.BloodSugar,
		SleepHours:    r.SleepHours,
		Steps:         r.Steps,
		Calories:      r.Calories,
		OxygenLevel:   r.OxygenLevel,
		Notes:         r.Notes,
	}
}
