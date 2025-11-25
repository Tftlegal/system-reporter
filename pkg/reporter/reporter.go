package reporter

// Reporter основной тип для работы с системными отчетами
type Reporter struct {
	config *Config
}

// New создает новый экземпляр Reporter
func New(config *Config) *Reporter {
	if config == nil {
		config = DefaultConfig()
	}
	return &Reporter{config: config}
}

// GenerateAndSend генерирует и отправляет отчет
func (r *Reporter) GenerateAndSend() error {
	// Генерируем отчет
	report, err := GenerateSystemReport()
	if err != nil {
		return fmt.Errorf("error generating report: %v", err)
	}

	// Конвертируем для API
	reportData, err := ConvertToMap(report)
	if err != nil {
		return fmt.Errorf("error converting report: %v", err)
	}

	// Отправляем на API
	if err := SendReportToAPI(r.config, reportData); err != nil {
		return fmt.Errorf("error sending report to API: %v", err)
	}

	return nil
}

// GenerateReport генерирует отчет без отправки
func (r *Reporter) GenerateReport() (*SystemReport, error) {
	return GenerateSystemReport()
}

// GetConfig возвращает конфигурацию репортера
func (r *Reporter) GetConfig() *Config {
	return r.config
}
