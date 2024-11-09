package helpers

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

func IsNodeInstalled() bool {
	cmd := exec.Command("node", "-v")
	_, err := cmd.Output()
	return err == nil
}

func IsPa11yInstalled() bool {
	cmd := exec.Command("pa11y", "--help")
	_, err := cmd.Output()
	return err == nil
}

func IsLighthouseInstalled() bool {
	cmd := exec.Command("lighthouse", "--help")
	_, err := cmd.Output()
	return err == nil
}

type Pa11yOutputErr struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Context string `json:"context"`
}

func GeneratePa11yReport(website string) ([]Pa11yOutputErr, error) {
	cmd := exec.Command("pa11y", website, "--reporter", "json")
	output, err := cmd.Output()

	if err != nil && !strings.Contains(err.Error(), "exit status 2") {
		return nil, err
	}

	var report []Pa11yOutputErr
	if err := json.Unmarshal([]byte(string(output)), &report); err != nil {
		return nil, err
	}

	return report, nil
}

type LighthouseAudit struct {
	Id           string   `json:"id"`
	Score        *float64 `json:"score"`
	NumericValue float64  `json:"numericValue"`
	NumericUnit  string   `json:"numericUnit"`
}
type LighthouseReport struct {
	Audits map[string]LighthouseAudit `json:"audits"`
}

func GenerateLighthouseReport(website string) (LighthouseReport, error) {
	homedir, _ := os.UserHomeDir()
	tmpLighthouseReportFilePath := fmt.Sprintf("%s/.something.lighthouse.tmp.json", homedir)

	go func() {
		cmd := exec.Command("lighthouse", website,
			"--quiet",
			"--no-enable-error-reporting",
			"--output", "json",
			"--output-path", tmpLighthouseReportFilePath,
			"--chrome-flags=--headless")

		err := cmd.Run()
		if err != nil {
			log.Fatal(err.Error())
		}
	}()

	for {
		_, err := os.Stat(tmpLighthouseReportFilePath)
		if err == nil {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	lighthouseReportBytes, err := os.ReadFile(tmpLighthouseReportFilePath)
	if err != nil {
		return LighthouseReport{}, err
	}

	if err := os.Remove(tmpLighthouseReportFilePath); err != nil {
		return LighthouseReport{}, err
	}

	var lighthouseReport LighthouseReport

	if err := json.Unmarshal(lighthouseReportBytes, &lighthouseReport); err != nil {
		return LighthouseReport{}, err
	}

	return lighthouseReport, nil
}
