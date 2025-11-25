package main

import (
	"flag"
	"fmt"
	"os"

	"RPC-report/pkg/reporter"
)

func main() {
	// Обработка флагов
	outputFile := flag.String("output", "report.json", "Output JSON file for report")
	postmanFlag := flag.Bool("postman", false, "Generate Postman request file")
	curlFlag := flag.Bool("curl", false, "Generate curl request file")
	flag.Parse()

	// Создаем репортер с конфигурацией по умолчанию
	rep := reporter.New(nil)

	// Получаем host_id
	hostID := reporter.GetHostID()
	fmt.Printf("Generating system report for host: %s\n", hostID)

	// Генерируем отчет
	fmt.Println("Generating system report...")
	report, err := rep.GenerateReport()
	if err != nil {
		fmt.Printf("Error generating report: %v\n", err)
		os.Exit(1)
	}

	// Сохраняем в файл
	fmt.Println("Saving report to file...")
	if err := reporter.SaveReportToJSON(report, *outputFile); err != nil {
		fmt.Printf("Error saving report to file: %v\n", err)
	}

	// Конвертируем отчет для API
	reportData, err := reporter.ConvertToMap(report)
	if err != nil {
		fmt.Printf("Error converting report: %v\n", err)
		os.Exit(1)
	}

	// Создаем Postman запрос если указан флаг
	if *postmanFlag {
		fmt.Println("Generating Postman request...")
		if err := createPostmanRequest(reportData, "postman_request.json"); err != nil {
			fmt.Printf("Error creating Postman request: %v\n", err)
		}
	}

	// Создаем curl запрос если указан флаг
	if *curlFlag {
		fmt.Println("Generating curl request...")
		if err := createCurlRequest(reportData, "curl_request.sh"); err != nil {
			fmt.Printf("Error creating curl request: %v\n", err)
		}
	}

	// Отправляем на API
	fmt.Println("Sending report to API...")

	// Вычисляем хеш отчета
	hash, err := reporter.CalculateReportHash(report)
	if err != nil {
		fmt.Printf("Error calculating report hash: %v\n", err)
	} else {
		fmt.Printf("Report hash: %s\n", hash)
	}

	// Отправляем отчет на API
	if err := rep.GenerateAndSend(); err != nil {
		fmt.Printf("Error sending report to API: %v\n", err)
		fmt.Println("Report was saved locally but failed to send to API")
		os.Exit(1)
	}

	fmt.Printf("Report successfully sent to API for host: %s\n", hostID)
	fmt.Println("System report completed successfully!")
	
	// Выводим информацию о созданных файлах
	if *postmanFlag {
		fmt.Printf("Postman request file created: %s\n", "postman_request.json")
	}
	if *curlFlag {
		fmt.Printf("Curl request file created: %s\n", "curl_request.sh")
	}
}

// Функции для создания Postman и curl запросов (можно вынести в отдельный пакет)
// ... (остальной код функций createPostmanRequest и createCurlRequest)
// createPostmanRequest создает файл для Postman
func createPostmanRequest(reportData map[string]interface{}, filename string) error {
        // Создаем JSON для тела запроса
        requestBody := APIReportRequest{
                Agent:  "system-reporter",
                Report: reportData,
        }

        jsonData, err := json.MarshalIndent(requestBody, "", "  ")
        if err != nil {
                return fmt.Errorf("failed to marshal request body: %v", err)
        }

        // Создаем структуру для Postman
        postmanRequest := PostmanRequest{
                Name: "System Report API",
                Request: PostmanRequestData{
                        Method: "PATCH",
                        Header: []PostmanHeader{
                                {
                                        Key:   "Content-Type",
                                        Value: "application/json",
                                        Type:  "text",
                                },
                        },
                        Body: PostmanBody{
                                Mode: "raw",
                                Raw:  string(jsonData),
                        },
                        URL: PostmanURL{
                                Raw:  APIBaseURL + ReportEndpoint,
                                Host: []string{"31.28.27.171:8123"},
                                Path: []string{"api", "report"},
                        },
                },
        }

        // Сохраняем в файл
        postmanData, err := json.MarshalIndent(postmanRequest, "", "  ")
        if err != nil {
                return fmt.Errorf("failed to marshal postman request: %v", err)
        }

        err = os.WriteFile(filename, postmanData, 0644)
        if err != nil {
                return fmt.Errorf("failed to write postman file: %v", err)
        }

        fmt.Printf("Postman request saved to %s (%d bytes)\n", filename, len(postmanData))
        return nil
}

// createCurlRequest создает файл с curl запросом
func createCurlRequest(reportData map[string]interface{}, filename string) error {
        // Создаем JSON для тела запроса
        requestBody := APIReportRequest{
                Agent:  "system-reporter",
                Report: reportData,
        }

        jsonData, err := json.Marshal(requestBody)
        if err != nil {
                return fmt.Errorf("failed to marshal request body: %v", err)
        }

        // Экранируем JSON для использования в curl
        escapedJSON := strings.ReplaceAll(string(jsonData), `"`, `\"`)
        escapedJSON = strings.ReplaceAll(escapedJSON, "`", "\\`")
        escapedJSON = strings.ReplaceAll(escapedJSON, "$", "\\$")

        // Создаем curl команду
        curlCommand := fmt.Sprintf(`curl -X PATCH "%s" \
  -H "Content-Type: application/json" \
  -d "%s"`, APIBaseURL+ReportEndpoint, escapedJSON)

        // Альтернативный вариант с @filename (более надежный для больших JSON)
        curlCommandAlt := fmt.Sprintf(`# Альтернативный вариант с файлом (рекомендуется для больших JSON):
echo '%s' | curl -X PATCH "%s" \
  -H "Content-Type: application/json" \
  -d @-`, string(jsonData), APIBaseURL+ReportEndpoint)

        // Сохраняем в файл
        content := fmt.Sprintf("#!/bin/bash\n\n# Curl command for system report API\n# Host ID: %s\n\n%s\n\n%s\n", getHostID(), curlCommand, curlCommandAlt)

        err = os.WriteFile(filename, []byte(content), 0755)
        if err != nil {
                return fmt.Errorf("failed to write curl file: %v", err)
        }

        fmt.Printf("Curl request saved to %s (%d bytes)\n", filename, len(content))
        return nil
}

