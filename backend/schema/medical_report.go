package schema

var MedicalReportSchema = map[string]interface{}{
	"name": "medical_report",
	"schema": map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"report_metadata": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"report_date": map[string]interface{}{
						"type":   "string",
						"format": "date-time",
					},
					"report_type": map[string]interface{}{
						"type": "string",
					},
					"hospital_or_lab_name": map[string]interface{}{
						"type": "string",
					},
				},
				"required": []string{
					"report_date",
					"report_type",
					"hospital_or_lab_name",
				},
			},
			"metrics": map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"test_name":       map[string]interface{}{"type": "string"},
						"value":           map[string]interface{}{"type": "string"},
						"unit":            map[string]interface{}{"type": "string"},
						"reference_range": map[string]interface{}{"type": "string"},
						"status":          map[string]interface{}{"type": "string"},
					},
					"required": []string{"test_name", "value", "unit", "reference_range", "status"},
				},
			},
			"abnormal_findings": map[string]interface{}{
				"type":  "array",
				"items": map[string]interface{}{"type": "string"},
			},
			"simple_explanation": map[string]interface{}{"type": "string"},
			"overall_risk_level": map[string]interface{}{"type": "string"},
			"recommendations": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"diet": map[string]interface{}{
						"type":  "array",
						"items": map[string]interface{}{"type": "string"},
					},
					"lifestyle": map[string]interface{}{
						"type":  "array",
						"items": map[string]interface{}{"type": "string"},
					},
				},
				"required": []string{"diet", "lifestyle"},
			},
			"follow_up_suggestions": map[string]interface{}{
				"type":  "array",
				"items": map[string]interface{}{"type": "string"},
			},
		},
		"required": []string{
			"report_metadata",
			"metrics",
			"abnormal_findings",
			"simple_explanation",
			"overall_risk_level",
			"recommendations",
			"follow_up_suggestions",
		},
	},
}
