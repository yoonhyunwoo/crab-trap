package worker

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/your-username/moltbook-prompt-injector/internal/env"
)

type Generator struct {
	serverURL string
	patterns  map[string]map[string]string
}

func NewGenerator(serverURL string) *Generator {
	return &Generator{
		serverURL: serverURL,
		patterns:  make(map[string]map[string]string),
	}
}

func (g *Generator) LoadPatterns(patternsPath string) error {
	data, err := os.ReadFile(patternsPath)
	if err != nil {
		return fmt.Errorf("failed to read patterns: %w", err)
	}

	if err := json.Unmarshal(data, &g.patterns); err != nil {
		return fmt.Errorf("failed to parse patterns: %w", err)
	}

	return nil
}

func (g *Generator) GenerateAll() []string {
	envInfo := env.Detect()

	var prompts []string

	if envVars, ok := g.patterns[string(envInfo.OS)]; ok {
		for _, templatePath := range envVars {
			if prompt := g.generateFromTemplate(templatePath); prompt != "" {
				prompts = append(prompts, prompt)
			}
		}
	}

	return prompts
}

func (g *Generator) Generate(templatePath string) string {
	return g.generateFromTemplate(templatePath)
}

func (g *Generator) generateFromTemplate(templatePath string) string {
	content, err := os.ReadFile(templatePath)
	if err != nil {
		return ""
	}

	template := string(content)

	template = strings.ReplaceAll(template, "{{SERVER_URL}}", g.serverURL)

	envInfo := env.Detect()
	envValues := env.GetAllValues(envInfo.OS)

	for varName, varValue := range envValues {
		pattern := envInfo.Pattern + varName
		template = strings.ReplaceAll(template, pattern, varValue)
	}

	return template
}

func (g *Generator) GenerateWithOS() map[string]string {
	result := make(map[string]string)

	unixPrompts := g.generateForOS(env.OSUnix)
	windowsPrompts := g.generateForOS(env.OSWindows)

	for k, v := range unixPrompts {
		result["unix_"+k] = v
	}

	for k, v := range windowsPrompts {
		result["windows_"+k] = v
	}

	return result
}

func (g *Generator) generateForOS(osType env.OSType) map[string]string {
	result := make(map[string]string)

	if osVars, ok := g.patterns[string(osType)]; ok {
		for name, templatePath := range osVars {
			if prompt := g.generateFromTemplate(templatePath); prompt != "" {
				result[name] = prompt
			}
		}
	}

	return result
}
