# RPC-report

#Использование пакета system-reporter

#Install go-task Ubuntu Debian Linux Mint
[package]

If you Set up the repository by running:


curl -1sLf 'https://dl.cloudsmith.io/public/task/task/setup.deb.sh' | sudo -E bash

Then you can install Task with:

apt install task

#Проверка кода

task install-formatters 

#Теперь вы можете использовать пакет в других приложениях:
```
package main

import (
	"fmt"
	"time"

	"RPC-report/pkg/reporter"
)

func main() {
	// Создаем кастомную конфигурацию
	config := &reporter.Config{
		APIBaseURL:     "http://localhost:8080/api",
		ReportEndpoint: "/reports",
		Timeout:        30 * time.Second,
                //AgentName:      "system-reporter",
		AgentName:      "my-custom-reporter",
	}

	// Создаем репортер
	rep := reporter.New(config)

	// Генерируем и отправляем отчет
	if err := rep.GenerateAndSend(); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println("Report sent successfully!")
}
```
## Эта структура обеспечивает:

Чистое разделение - логика разделена на отдельные файлы
Переиспользуемость - пакет можно использовать в других проектах
Конфигурируемость - гибкая настройка через Config
Расширяемость - легко добавлять новые функции
Совместимость - сохраняется обратная совместимость с существующим кодом


