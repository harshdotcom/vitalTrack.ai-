package models

import (
	"errors"
	"fmt"
	"time"
)

type DailyHealthMetric struct {
	ID            string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UploadedBy    int64
	HeartRate     *int           // bpm
	Weight        *Measurement   `gorm:"serializer:json;type:jsonb"` // supports units
	BloodPressure *BloodPressure `gorm:"serializer:json;type:jsonb"`
	BloodSugar    *Measurement   `gorm:"serializer:json;type:jsonb"`
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
	Systolic  int `json:"systolic"`  // mmHg
	Diastolic int `json:"diastolic"` // mmHg
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

	Notes     *string `json:"notes"`
	Timestamp *string `json:"timestamp" example:"2026-04-16T22:05:00+05:30"`
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
	if parsedTimestamp, err := r.ResolveTimestamp(); err == nil {
		timestamp = parsedTimestamp
	}

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

func (d *DailyHealthMetric) MetricType() string {
	switch {
	case d.BloodPressure != nil:
		return "blood_pressure"
	case d.BloodSugar != nil:
		return "blood_sugar"
	case d.Weight != nil:
		return "weight"
	case d.HeartRate != nil:
		return "heart_rate"
	case d.OxygenLevel != nil:
		return "oxygen_level"
	case d.SleepHours != nil:
		return "sleep_hours"
	case d.Steps != nil:
		return "steps"
	case d.Calories != nil:
		return "calories"
	case d.Notes != nil && *d.Notes != "":
		return "notes"
	default:
		return "unknown"
	}
}

func (d *DailyHealthMetric) MetricLabel() string {
	switch d.MetricType() {
	case "blood_pressure":
		return "Blood Pressure"
	case "blood_sugar":
		return "Blood Sugar"
	case "weight":
		return "Weight"
	case "heart_rate":
		return "Heart Rate"
	case "oxygen_level":
		return "Oxygen Level"
	case "sleep_hours":
		return "Sleep Hours"
	case "steps":
		return "Steps"
	case "calories":
		return "Calories"
	case "notes":
		return "Notes"
	default:
		return "Direct Entry"
	}
}

func (d *DailyHealthMetric) MetricSummary() string {
	switch d.MetricType() {
	case "blood_pressure":
		return fmt.Sprintf("%d/%d mmHg", d.BloodPressure.Systolic, d.BloodPressure.Diastolic)
	case "blood_sugar":
		return fmt.Sprintf("%g %s", d.BloodSugar.Value, d.BloodSugar.Unit)
	case "weight":
		return fmt.Sprintf("%g %s", d.Weight.Value, d.Weight.Unit)
	case "heart_rate":
		return fmt.Sprintf("%d bpm", *d.HeartRate)
	case "oxygen_level":
		return fmt.Sprintf("%g%%", *d.OxygenLevel)
	case "sleep_hours":
		return fmt.Sprintf("%g hrs", *d.SleepHours)
	case "steps":
		return fmt.Sprintf("%d steps", *d.Steps)
	case "calories":
		return fmt.Sprintf("%d kcal", *d.Calories)
	case "notes":
		return *d.Notes
	default:
		return "Logged directly in VitaTrack"
	}
}

func parseMetricTimestamp(raw *string) (time.Time, error) {
	if raw == nil || *raw == "" {
		return time.Time{}, errors.New("timestamp not provided")
	}

	layouts := []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02T15:04",
		"2006-01-02",
	}

	for _, layout := range layouts {
		if layout == "2006-01-02T15:04" || layout == "2006-01-02" {
			if ts, err := time.ParseInLocation(layout, *raw, time.Local); err == nil {
				return ts, nil
			}
			continue
		}

		if ts, err := time.Parse(layout, *raw); err == nil {
			return ts, nil
		}
	}

	return time.Time{}, fmt.Errorf("invalid timestamp format")
}

func (r *SaveHealthMetricRequest) ResolveTimestamp() (time.Time, error) {
	return parseMetricTimestamp(r.Timestamp)
}
