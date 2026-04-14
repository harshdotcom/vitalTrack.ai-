package models

import "fmt"

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type AIRequestBody struct {
	Model       string    `json:"model"`
	Input       []Message `json:"input"`
	Temperature float64   `json:"temperature"`
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

const SystemPrompt string = `
You are a medical report analysis assistant.
I Should be able to convert returned JSON String to a struct using the standard library encoding/json
Return ONLY valid JSON.
Do NOT include markdown.
Do NOT include explanations outside JSON.
Do NOT provide diagnosis or medications.
If data is missing, use null.
If document_date is available, return it in the format DD-MM-YYYY

Follow this exact JSON structure to produce JSON String:

{
  "report_metadata": {
    "document_date": "",
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

func GetUserPrompt(ocrText string) string {
	userPrompt := fmt.Sprintf(`
Here is the extracted medical report text:

%s

Generate the structured JSON response following the schema.
`, ocrText)

	return userPrompt
}
