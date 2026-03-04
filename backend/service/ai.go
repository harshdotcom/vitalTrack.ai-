package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"vita-track-ai/schema"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type AIRequestBody struct {
	Model          string                 `json:"model"`
	Input          []Message              `json:"input"`
	Temperature    float64                `json:"temperature"`
	ResponseFormat map[string]interface{} `json:"response_format"`
}

type AIResponse struct {
	Output []struct {
		Type    string `json:"type"`
		Role    string `json:"role"`
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	} `json:"output"`
}

type MedicalReport struct {
	ReportMetadata struct {
		ReportDate        string `json:"report_date"`
		ReportType        string `json:"report_type"`
		HospitalOrLabName string `json:"hospital_or_lab_name"`
	} `json:"report_metadata"`

	Metrics []struct {
		TestName       string `json:"test_name"`
		Value          string `json:"value"`
		Unit           string `json:"unit"`
		ReferenceRange string `json:"reference_range"`
		Status         string `json:"status"`
	} `json:"metrics"`

	AbnormalFindings  []string `json:"abnormal_findings"`
	SimpleExplanation string   `json:"simple_explanation"`
	OverallRiskLevel  string   `json:"overall_risk_level"`

	Recommendations struct {
		Diet      []string `json:"diet"`
		Lifestyle []string `json:"lifestyle"`
	} `json:"recommendations"`

	FollowUpSuggestions []string `json:"follow_up_suggestions"`
}

func ExtractMedicalReport(aiResp AIResponse) (*MedicalReport, error) {
	for _, item := range aiResp.Output {
		if item.Type == "message" && item.Role == "assistant" {
			for _, content := range item.Content {
				if content.Type == "output_text" {
					var report MedicalReport
					err := json.Unmarshal([]byte(strings.TrimSpace(content.Text)), &report)
					if err != nil {
						return nil, err
					}
					return &report, nil
				}
			}
		}
	}
	return nil, fmt.Errorf("assistant JSON not found")
}

func AnalyzeMedicalReport(fileId string) (*MedicalReport, error) {
	ocrText, err := GenerateOCRText(fileId)
	if err != nil {
		return nil, fmt.Errorf("unable to load AWS config: %w", err)
	}

	systemPrompt := `
You are a medical report analysis assistant.
I Should be able to convert returned JSON String to a struct using the standard library encoding/json
Return ONLY valid JSON.
Do NOT include markdown.
Do NOT include explanations outside JSON.
Do NOT provide diagnosis or medications.
If data is missing, use null.

Follow this exact JSON structure to produce JSON String:

{
  "report_metadata": {
    "report_date": "",
    "report_type": "",
    "hospital_or_lab_name": ""
  },
  "metrics": [
    {
      "test_name": "",
      "value": "",
      "unit": "",
      "reference_range": "",
      "status": ""
    }
  ],
  "abnormal_findings": [],
  "simple_explanation": "",
  "overall_risk_level": "",
  "recommendations": {
    "diet": [],
    "lifestyle": []
  },
  "follow_up_suggestions": []
}
`

	userPrompt := fmt.Sprintf(`
Here is the extracted medical report text:

%s

Generate the structured JSON response following the schema.
`, ocrText)

	requestBody := AIRequestBody{
		Model: "openai.gpt-oss-120b",
		Input: []Message{
			{
				Role:    "system",
				Content: systemPrompt,
			},
			{
				Role:    "user",
				Content: userPrompt,
			},
		},
		Temperature: 0.0,
		ResponseFormat: map[string]interface{}{
			"type":        "json_schema",
			"json_schema": schema.MedicalReportSchema,
		},
	}

	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", os.Getenv("OPENAI_BASE_URL"), bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("OPENAI_API_KEY"))

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		errorBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("status: %d, body: %s", resp.StatusCode, string(errorBody))
	}

	body, _ := io.ReadAll(resp.Body)

	var aiResp AIResponse
	err = json.Unmarshal(body, &aiResp)
	if err != nil {
		return nil, err
	}

	fmt.Print(aiResp)

	return ExtractMedicalReport(aiResp)
}
