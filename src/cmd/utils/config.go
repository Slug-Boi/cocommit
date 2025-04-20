package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

var ConfigVar *Config

var (
	defaultConfigLocations = []string{
		"",
		os.Getenv("HOME") + "/.config/cocommit",
		os.Getenv("HOME") + "/cocommit",
		"/etc/cocommit",
		"/usr/local/etc/cocommit",
	}
	configName = "config"
	configType = "toml"
)

type Config struct {
	Settings struct {
		AuthorFile    string `mapstructure:"author_file"`
		StartingScope string `mapstructure:"starting_scope"`
		Editor        string `mapstructure:"editor"`
	} `mapstructure:"settings"`
}

func init() {
	configDir, err := os.UserConfigDir()
	if err == nil {
		defaultConfigLocations[0] = filepath.Join(configDir, "cocommit")
	}
}

func LoadConfig() (*Config, error) {
	v := viper.New()
	v.SetConfigName(configName)
	v.SetConfigType(configType)

	// Set default values
	v.SetDefault("settings.author_file", defaultConfigLocations[0]+"/authors.json")
	v.SetDefault("settings.starting_scope", "git")
	v.SetDefault("settings.editor", "built-in")

	// Add search paths
	for _, path := range defaultConfigLocations {
		if path != "" {
			v.AddConfigPath(path)
		}
	}

	// Try to read config
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			if err := handleMissingConfig(v); err != nil {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("config error: %w", err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("config unmarshal error: %w", err)
	}
	if cfg.Settings.AuthorFile == "" {
		cfg.Settings.AuthorFile = defaultConfigLocations[0] + "/authors.json"
	}

	return &cfg, nil
}

func (c *Config) SetGlobalConfig() {
	if ConfigVar == nil {
		ConfigVar = c
		// This doesnt really do much right now but might be useful later
		viper.WatchConfig()
	} 
}

func handleMissingConfig(v *viper.Viper) error {
	fmt.Println("Config file not found. Would you like to create one? (y/n)")
	var response string
	if _, err := fmt.Scanln(&response); err != nil {
		return fmt.Errorf("error reading response: %w", err)
	}

	yesResponses := map[string]bool{"y": true, "Y": true, "yes": true, "Yes": true, "YES": true}
	if !yesResponses[strings.TrimSpace(response)] {
		return fmt.Errorf("config file not found")
	}

	return createConfig(v)
}

func createConfig(v *viper.Viper) error {
	fmt.Println("Where would you like to create the config file?")
	for i, path := range defaultConfigLocations {
		fmt.Printf("%d. %s\n", i, path)
	}
	fmt.Println("Please enter the number of the location or a custom path:")

	var response string
	if _, err := fmt.Scanln(&response); err != nil {
		return fmt.Errorf("error reading response: %w", err)
	}

	var configPath string
	if num, err := strconv.Atoi(response); err == nil && num >= 0 && num < len(defaultConfigLocations) {
		configPath = defaultConfigLocations[num]
	} else {
		configPath = response
	}

	// Ensure directory exists
	if err := os.MkdirAll(configPath, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Set the config file path
	fullPath := filepath.Join(configPath, fmt.Sprintf("%s.%s", configName, configType))
	v.SetConfigFile(fullPath)

	// Write default config
	if err := v.SafeWriteConfig(); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	fmt.Printf("Config file created at: %s\n", fullPath)
	return nil
}

func (c *Config) Save() error {
    v := viper.New()
    
    // Set all configuration values from the struct
    v.Set("settings.author_file", c.Settings.AuthorFile)
    v.Set("settings.starting_scope", c.Settings.StartingScope)
    v.Set("settings.editor", c.Settings.Editor)
    
    v.SetConfigName(configName)
    v.SetConfigType(configType)
    
    // Try to determine the original config file location
    if viper.ConfigFileUsed() != "" {
        v.SetConfigFile(viper.ConfigFileUsed())
    } else {
        // Fall back to first default location if no existing config
        if len(defaultConfigLocations) > 0 && defaultConfigLocations[0] != "" {
            v.SetConfigFile(filepath.Join(defaultConfigLocations[0], fmt.Sprintf("%s.%s", configName, configType)))
        } else {
            return fmt.Errorf("no config file location available")
        }
    }
    
    // Ensure the directory exists
    configDir := filepath.Dir(v.ConfigFileUsed())
    if err := os.MkdirAll(configDir, 0755); err != nil {
        return fmt.Errorf("failed to create config directory: %w", err)
    }
    
    // Write the config file
    if err := v.WriteConfig(); err != nil {
        return fmt.Errorf("failed to save config: %w", err)
    }
    
    return nil
}