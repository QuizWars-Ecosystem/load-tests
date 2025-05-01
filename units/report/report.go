package report

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"os"

	"github.com/QuizWars-Ecosystem/load-tests/units/runner"
	"github.com/QuizWars-Ecosystem/load-tests/units/utils"
)

func Print(name string, r runner.Result) {
	avg := utils.AvgLatency(r.Latencies)
	p95 := utils.Percentile(r.Latencies, 95)
	fmt.Printf("📊 [%s] Total: %d | OK: %d | Fail: %d | Avg: %s | P95: %s\n",
		name, r.Total, r.Success, r.Failures, avg, p95)
}

func SaveJSON(name string, r runner.Result) {
	data, _ := json.MarshalIndent(r, "", "  ")
	if err := os.WriteFile(fmt.Sprintf("report_%s.json", name), data, 0o644); err != nil {
		log.Fatal(err)
	}
}

func SaveHTML(name string, r runner.Result) {
	const htmlTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>Load Test Report</title>
	<style>
		body {
			font-family: Arial, sans-serif;
			margin: 20px;
			color: #333;
		}
		h1 {
			text-align: center;
		}
		table {
			width: 100%;
			border-collapse: collapse;
			margin-top: 20px;
		}
		th, td {
			padding: 10px;
			text-align: left;
			border: 1px solid #ddd;
		}
		th {
			background-color: #f4f4f4;
		}
		.success {
			color: green;
		}
		.fail {
			color: red;
		}
		.avg, .p95 {
			font-weight: bold;
		}
	</style>
</head>
<body>
	<h1>Load Test Report: {{.ScenarioName}}</h1>
	<table>
		<tr><th>Total Requests</th><td>{{.Total}}</td></tr>
		<tr><th>Success</th><td class="success">{{.Success}}</td></tr>
		<tr><th>Failures</th><td class="fail">{{.Failures}}</td></tr>
		<tr><th>Average Latency</th><td class="avg">{{.AvgLatency}}</td></tr>
		<tr><th>P95 Latency</th><td class="p95">{{.P95Latency}}</td></tr>
	</table>
</body>
</html>
`

	// Подготавливаем данные для шаблона
	avg := utils.AvgLatency(r.Latencies)
	p95 := utils.Percentile(r.Latencies, 95)

	data := struct {
		ScenarioName string
		Total        int
		Success      int
		Failures     int
		AvgLatency   string
		P95Latency   string
	}{
		ScenarioName: name,
		Total:        r.Total,
		Success:      r.Success,
		Failures:     r.Failures,
		AvgLatency:   avg.String(),
		P95Latency:   p95.String(),
	}

	tmpl, err := template.New("report").Parse(htmlTemplate)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Create(fmt.Sprintf("../reports/report_%s.html", name))
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		_ = file.Close()
	}()

	err = tmpl.Execute(file, data)
	if err != nil {
		log.Fatal(err)
	}
}
