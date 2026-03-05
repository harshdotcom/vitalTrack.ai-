package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"vita-track-ai/models"
	"vita-track-ai/repository"
)

func AnalyzeMedicalReport(fileId string) (*models.MedicalReport, error) {
	reportDb, err := repository.GetMedicalReportByID(fileId)

	if reportDb != nil {
		report, err := models.GetMedicalReportApiFormat(reportDb)
		return report, err
	}

	ocrText, err := GenerateOCRText(fileId)
	if err != nil {
		return nil, fmt.Errorf("unable to get OCR text: %w", err)
	}

	userPrompt := models.GetUserPrompt(ocrText)

	requestBody := models.AIRequestBody{
		Model: "openai.gpt-oss-120b",
		Input: []models.Message{
			{
				Role:    "system",
				Content: models.SystemPrompt,
			},
			{
				Role:    "user",
				Content: userPrompt,
			},
		},
		Temperature: 0.0,
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

	var aiResp models.AIResponse
	err = json.Unmarshal(body, &aiResp)
	if err != nil {
		return nil, err
	}

	report, err := extractMedicalReportFromAiResponse(aiResp)
	if err != nil {
		return nil, err
	}

	reportDb = models.GetMedicalReportDBFormat(report, fileId)
	err = repository.CreateMedicalReport(reportDb)

	if err != nil {
		return nil, err
	}

	return report, err
}

func extractMedicalReportFromAiResponse(aiResp models.AIResponse) (*models.MedicalReport, error) {
	for _, item := range aiResp.Output {
		if item.Type == "message" && item.Role == "assistant" {
			for _, content := range item.Content {
				if content.Type == "output_text" {
					var report models.MedicalReport
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
