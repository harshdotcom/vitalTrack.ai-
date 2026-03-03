package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
)

func AnalyzeMedicalReport(fileId string) (string, error) {
	ocrText, err := GenerateOCRText(fileId)
	if err != nil {
		return "", err
	}
	ctx := context.Background()

	// Load AWS config
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return "", fmt.Errorf("unable to load AWS config: %w", err)
	}

	client := bedrockruntime.NewFromConfig(cfg)

	modelID := "anthropic.claude-3-sonnet-20240229-v1:0"

	systemPrompt := `
You are a medical report analysis assistant.

Return ONLY valid JSON.
Do NOT include markdown.
Do NOT include explanations outside JSON.
Do NOT provide diagnosis or medications.
If data is missing, use null.

Follow this exact JSON structure:

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

Generate the structured JSON response.
`, ocrText)

	requestBody := map[string]interface{}{
		"anthropic_version": "bedrock-2023-05-31",
		"max_tokens":        1500,
		"temperature":       0.1,
		"system":            systemPrompt,
		"messages": []map[string]interface{}{
			{
				"role":    "user",
				"content": userPrompt,
			},
		},
	}

	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	resp, err := client.InvokeModel(ctx, &bedrockruntime.InvokeModelInput{
		ModelId:     &modelID,
		ContentType: awsString("application/json"),
		Body:        bodyBytes,
	})
	if err != nil {
		return "", err
	}

	var responseBody map[string]interface{}
	if err := json.Unmarshal(resp.Body, &responseBody); err != nil {
		return "", err
	}

	// Extract model output
	content := responseBody["content"].([]interface{})
	first := content[0].(map[string]interface{})
	jsonText := first["text"].(string)

	return jsonText, nil
}

func awsString(s string) *string {
	return &s
}
