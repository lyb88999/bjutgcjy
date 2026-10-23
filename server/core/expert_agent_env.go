package core

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// loadExpertAgentEnv loads only the assistant settings from the ignored file beside config.yaml.
// Explicit process environment values always take priority. Credentials never enter GVA_CONFIG.
func loadExpertAgentEnv(configPath string) error {
	file, err := os.Open(filepath.Join(filepath.Dir(configPath), ".env.local"))
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("无法读取本地助手配置")
	}
	defer file.Close()
	allowed := map[string]bool{"EXPERT_AGENT_API_KEY": true, "EXPERT_AGENT_BASE_URL": true, "EXPERT_AGENT_MODEL": true}
	values := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for lineNumber := 1; scanner.Scan(); lineNumber++ {
		line := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(scanner.Text()), "export "))
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, found := strings.Cut(line, "=")
		key = strings.TrimSpace(key)
		if !allowed[key] {
			continue
		}
		if !found {
			return fmt.Errorf("本地助手配置第 %d 行格式无效", lineNumber)
		}
		value = strings.TrimSpace(value)
		if strings.HasPrefix(value, "\"") {
			value, err = strconv.Unquote(value)
			if err != nil {
				return fmt.Errorf("本地助手配置第 %d 行引号无效", lineNumber)
			}
		} else if strings.HasPrefix(value, "'") {
			if len(value) < 2 || !strings.HasSuffix(value, "'") {
				return fmt.Errorf("本地助手配置第 %d 行引号无效", lineNumber)
			}
			value = value[1 : len(value)-1]
		}
		values[key] = value
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("无法解析本地助手配置")
	}
	for key, value := range values {
		if _, exists := os.LookupEnv(key); !exists {
			if err := os.Setenv(key, value); err != nil {
				return fmt.Errorf("无法加载本地助手配置")
			}
		}
	}
	return nil
}
