package config

import (
	// "errors"

	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"ask_terminal/security"

	"gopkg.in/yaml.v2"
)

// Config holds application configuration
type Config struct {
	BaseURL     string `yaml:"base_url"`     // Changed from mapstructure
	APIKey      string `yaml:"api_key"`      // Changed from mapstructure
	ModelName   string `yaml:"model_name"`   // Changed from mapstructure
	PrivateMode bool   `yaml:"private_mode"` // Changed from mapstructure
	SysPrompt   string `yaml:"sys_prompt"`   // Changed from mapstructure
	Provider    string `yaml:"provider"`     // Added yaml tag
}

// LoadConfig loads configuration from the specified path
func LoadConfig(configPath string) (*Config, error) {
	// If config path is not specified, use default
	if configPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		configPath = filepath.Join(homeDir, ".config", "askta", "config.yaml")
	}

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Create directory structure if it doesn't exist
		configDir := filepath.Dir(configPath)
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create config directory: %w", err)
		}

		// Create a default config with comments
		defaultConfigYaml := `# ASK Terminal AI Configuration

# API service configuration
base_url: "https://api.openai.com/v1/"  # API base URL for your provider
api_key: "your-api-key"                 # Your API key (will be encrypted after first run)
model_name: "gpt-4o-mini"               # Default AI model to use

# Feature configuration
private_mode: false                     # Set to true to not send directory structure
sys_prompt: ""                          # System prompt, WARNING: Please understand what you're modifying before making changes

# Provider configuration (currently only openai-compatible is supported)
provider: "openai-compatible"           # AI provider type, no other options available yet
`

		if err := os.WriteFile(configPath, []byte(defaultConfigYaml), 0600); err != nil {
			return nil, fmt.Errorf("failed to write default config: %w", err)
		}

		return nil, fmt.Errorf("created default config at %s, please add your API key", configPath)
	}

	// Read config file
	data, err := ioutil.ReadFile(configPath)
	fmt.Printf("Debug - Read data: %s\n %s", data, err)
	if err != nil {
		return nil, err
	}

	// Parse YAML
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	// Add debug logging
	fmt.Printf("Debug - Read config file: %s\n", configPath)
	fmt.Printf("Debug - Config values: BaseURL=%s, APIKey=%s, Model=%s\n",
		config.BaseURL,
		config.APIKey, // Show only first 3 chars of API key for safety
		config.ModelName)

	// Validate required fields
	if config.APIKey == "" {
		return nil, fmt.Errorf("api_key is required in configuration: %s,%s, %s", config.APIKey, config.BaseURL, config.ModelName)
	}

	if config.ModelName == "" {
		// Set default model
		config.ModelName = "gpt-4o-mini"
	}

	// Check if API key needs decryption
	decryptedKey := "" // Initialize decryptedKey
	if len(config.APIKey) > 6 && config.APIKey[:6] == "encry_" {
		// Decrypt API key
		decryptedKey, err = security.DecryptAPIKey(config.APIKey[6:])
		if err != nil {
			return nil, err
		}
		config.APIKey = decryptedKey
	} else {
		originalKey := config.APIKey
		// Encrypt API key for future use
		encryptedKey, err := security.EncryptAPIKey(config.APIKey)
		if err != nil {
			return nil, err
		}

		// Update config file with encrypted key
		config.APIKey = encryptedKey
		newData, err := yaml.Marshal(&config)
		if err != nil {
			return nil, err
		}

		// Write updated config back to file
		if err := ioutil.WriteFile(configPath, newData, 0600); err != nil {
			return nil, err
		}

		// Restore unencrypted key for current use
		config.APIKey = originalKey
	}

	return &config, nil
}

// MergeWithArgs merges command line arguments into config
func (c *Config) MergeWithArgs(args map[string]string) {
	// Override config with command line arguments
	if model, ok := args["model"]; ok && model != "" {
		c.ModelName = model
	}

	if baseURL, ok := args["url"]; ok && baseURL != "" {
		c.BaseURL = baseURL
	}

	if apiKey, ok := args["key"]; ok && apiKey != "" {
		c.APIKey = apiKey
	}

	if sysPrompt, ok := args["sys_prompt"]; ok && sysPrompt != "" {
		c.SysPrompt = sysPrompt
	}

	if _, ok := args["private_mode"]; ok {
		c.PrivateMode = true
	}
}
