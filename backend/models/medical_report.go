package models

import (
	"encoding/json"
	"time"

	"gorm.io/datatypes"
)

type MedicalReport struct {
	ReportMetadata      ReportMetadata  `json:"report_metadata"`
	Metrics             []Metric        `json:"metrics"`
	AbnormalFindings    []string        `json:"abnormal_findings"`
	SimpleExplanation   string          `json:"simple_explanation"`
	OverallRiskLevel    string          `json:"overall_risk_level"`
	Recommendations     Recommendations `json:"recommendations"`
	FollowUpSuggestions []string        `json:"follow_up_suggestions"`
}

type MedicalReportDB struct {
	FileID              string         `json:"id" gorm:"column:id;primaryKey"`
	File                File           `gorm:"foreignKey:FileID;references:ID;constraint:OnDelete:CASCADE"`
	ReportMetadata      datatypes.JSON `gorm:"type:jsonb"`
	Metrics             datatypes.JSON `gorm:"type:jsonb"`
	AbnormalFindings    datatypes.JSON `gorm:"type:jsonb"`
	SimpleExplanation   string
	OverallRiskLevel    string
	Recommendations     datatypes.JSON `gorm:"type:jsonb"`
	FollowUpSuggestions datatypes.JSON `gorm:"type:jsonb"`
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

func GetMedicalReportDBFormat(report *MedicalReport, fileId string) *MedicalReportDB {
	metadata, _ := json.Marshal(report.ReportMetadata)
	metrics, _ := json.Marshal(report.Metrics)
	abnormal, _ := json.Marshal(report.AbnormalFindings)
	reco, _ := json.Marshal(report.Recommendations)
	follow, _ := json.Marshal(report.FollowUpSuggestions)

	dbReport := MedicalReportDB{
		FileID:              fileId,
		ReportMetadata:      metadata,
		Metrics:             metrics,
		AbnormalFindings:    abnormal,
		SimpleExplanation:   report.SimpleExplanation,
		OverallRiskLevel:    report.OverallRiskLevel,
		Recommendations:     reco,
		FollowUpSuggestions: follow,
	}

	return &dbReport
}

func GetMedicalReportApiFormat(dbReport *MedicalReportDB) (*MedicalReport, error) {
	var metadata ReportMetadata
	var metrics []Metric
	var abnormal []string
	var reco Recommendations
	var follow []string

	// Unmarshal JSONB fields
	if err := json.Unmarshal(dbReport.ReportMetadata, &metadata); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(dbReport.Metrics, &metrics); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(dbReport.AbnormalFindings, &abnormal); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(dbReport.Recommendations, &reco); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(dbReport.FollowUpSuggestions, &follow); err != nil {
		return nil, err
	}

	report := &MedicalReport{
		ReportMetadata:      metadata,
		Metrics:             metrics,
		AbnormalFindings:    abnormal,
		SimpleExplanation:   dbReport.SimpleExplanation,
		OverallRiskLevel:    dbReport.OverallRiskLevel,
		Recommendations:     reco,
		FollowUpSuggestions: follow,
	}

	return report, nil
}

// ======================================================================================
type Metric struct {
	TestName       string `json:"test_name"`
	Value          string `json:"value"`
	Unit           string `json:"unit"`
	ReferenceRange string `json:"reference_range"`
	Status         string `json:"status"`
}

type Recommendations struct {
	Diet      []string `json:"diet"`
	Lifestyle []string `json:"lifestyle"`
}

type ReportMetadata struct {
	ReportDate    string `json:"report_date"`
	ReportType    string `json:"report_type"`
	HospitalOrLab string `json:"hospital_or_lab_name"`
}
