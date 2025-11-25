package reporter

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

// DefaultConfig возвращает конфигурацию по умолчанию
func DefaultConfig() *Config {
	return &Config{
		APIBaseURL:     "http://go.sytes.ru:8123/api",
		ReportEndpoint: "/report",
		Timeout:        60 * time.Second,
		AgentName:      "system-reporter",
	}
}

// SendReportToAPI отправляет отчет на API
func SendReportToAPI(config *Config, reportData map[string]interface{}) error {
	request := APIReportRequest{
		Agent:  config.AgentName,
		Report: reportData,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal report: %v", err)
	}

	fmt.Printf("Sending report to API (%d bytes)...\n", len(jsonData))

	client := &http.Client{Timeout: config.Timeout}
	apiURL := config.APIBaseURL + config.ReportEndpoint

	// Сначала пробуем PATCH
	req, err := http.NewRequest("PATCH", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create PATCH request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send PATCH request: %v", err)
	}
	defer resp.Body.Close()

	// Если PATCH не поддерживается, пробуем PUT
	if resp.StatusCode == http.StatusMethodNotAllowed {
		fmt.Println("PATCH not supported, trying PUT...")
		req, err = http.NewRequest("PUT", apiURL, bytes.NewBuffer(jsonData))
		if err != nil {
			return fmt.Errorf("failed to create PUT request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err = client.Do(req)
		if err != nil {
			return fmt.Errorf("failed to send PUT request: %v", err)
		}
		defer resp.Body.Close()
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("API returned status: %d", resp.StatusCode)
	}

	fmt.Println("Report successfully sent to API")
	return nil
}

// CalculateReportHash вычисляет хеш отчета для идентификации
func CalculateReportHash(report *SystemReport) (string, error) {
	jsonData, err := json.Marshal(report)
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256(jsonData)
	return hex.EncodeToString(hash[:]), nil
}

// ConvertToMap конвертирует SystemReport в map для API
func ConvertToMap(report *SystemReport) (map[string]interface{}, error) {
	jsonData, err := json.Marshal(report)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// SaveReportToJSON сохраняет отчет в JSON файл
func SaveReportToJSON(report *SystemReport, filename string) error {
	jsonData, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}

	err = os.WriteFile(filename, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	fmt.Printf("Report saved to %s (%d bytes)\n", filename, len(jsonData))
	return nil
}
