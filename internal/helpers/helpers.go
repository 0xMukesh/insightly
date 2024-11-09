package helpers

import (
	"encoding/json"
	"os/exec"
	"strings"
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

type Pa11yOutputErr struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Context string `json:"context"`
}

func RunPa11yReport(website string) ([]Pa11yOutputErr, error) {
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
