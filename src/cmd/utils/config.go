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
		os.Getenv("COCOMMIT_CONFIG"),
		os.Getenv("HOME") + "/.config/cocommit",
		os.Getenv("HOME") + "/cocommit",
		"/etc/cocommit",
		"/usr/local/etc/cocommit",
	}
	configName = "config"
	configType = "toml"
)

// type Config struct {
// 	Settings struct {
// 		AuthorFile    string `mapstructure:"author_file"`
// 		StartingScope string `mapstructure:"starting_scope"`
// 		Editor        string `mapstructure:"editor"`
// 		DefaultStoreRepo string `mapstructure:"default_user_store_repository"`
// 	} `mapstructure:"settings"`
// 	Style struct {
// 		Help		string 	`mapstructure:"help"`
// 		Item		string	`mapstructure:"item"`
// 		SelectedItemFG string 	`mapstructure:"selected_item_fg"`
// 		HighlightFG 	string	`mapstructure:"highlight_fg"`
// 		SelectedHighlightFG string `mapstructure:"selected_highlight_fg"`
// 		SelectedItemBG string 	`mapstructure:"selected_item_bg"`
// 		HighlightBG 	string	`mapstructure:"highlight_bg"`
// 		SelectedHighlightBG string `mapstructure:"selected_highlight_bg"`
// 		Delete		string	`mapstructure:"delete"`
// 		Sharing		string 	`mapstructure:"sharing"`
// 		Pasting		string 	`mapstructure:"pasting"`
// 		ActivePaginationDot string `mapstructure:"active_pagination_dot"`
// 		GitScope    string 	`mapstructure:"git_scope"`
// 		LocalScope 	string 	`mapstructure:"local_scope"`
// 		MixedScope	string 	`mapstructure:"mixed_scope"`

// 		CommitMessage struct {
// 			// TUI CommitMessageWriter
// 			Base string 	`mapstructure:"base"`
// 			LineNumber string `mapstructure:"line_number"`
// 		} `mapstructure:"commit_message_editor"`

// 		GH struct {
// 		// Tui GH
// 		Error 	string 	`mapstructure:"error"`
// 		Toggle	string 	`mapstructure:"toggle"`
// 		ActiveToggle string 	`mapstructure:"active_toggle"`
// 		} `mapstructure:"github_creation"`

// 		Author struct {
// 			// TUI author
// 			Focused string  `mapstructure:"focused"`
// 			Blurred string  `mapstructure:"blurred"`
// 			CursorModeHelp  string `mapstructure:"cursor_mode_help"`
// 			Cursor string	`mapstructure:"cursor"`
// 		} `mapstructure:"author_creation"`
// 		Groups struct {
// 			// TUI Groups
// 			ModelStyle string `mapstructure:"group_box"`
// 			FocusedModelStyle string `mapstructure:"group_focused"`
// 		} `mapstructure:"group_selection"`

// 		Light struct {
// 			Help		string 	`mapstructure:"help"`
// 			Item		string	`mapstructure:"item"`
// 			SelectedItemFG string 	`mapstructure:"selected_item_fg"`
// 			HighlightFG 	string	`mapstructure:"highlight_fg"`
// 			SelectedHighlightFG string `mapstructure:"selected_highlight_fg"`
// 			SelectedItemBG string 	`mapstructure:"selected_item_bg"`
// 			HighlightBG 	string	`mapstructure:"highlight_bg"`
// 			SelectedHighlightBG string `mapstructure:"selected_highlight_bg"`
// 			Delete		string	`mapstructure:"delete"`
// 			Sharing		string 	`mapstructure:"sharing"`
// 			Pasting		string 	`mapstructure:"pasting"`
// 			ActivePaginationDot string `mapstructure:"active_pagination_dot"`
// 			GitScope    string 	`mapstructure:"git_scope"`
// 			LocalScope 	string 	`mapstructure:"local_scope"`
// 			MixedScope	string 	`mapstructure:"mixed_scope"`

// 			CommitMessage struct {
// 			// TUI CommitMessageWriter
// 			Base string 	`mapstructure:"base"`
// 			LineNumber string `mapstructure:"line_number"`
// 			} `mapstructure:"commit_message_editor"`

// 			GH struct {
// 			// Tui GH
// 			Error 	string 	`mapstructure:"error"`
// 			Toggle	string 	`mapstructure:"toggle"`
// 			ActiveToggle string 	`mapstructure:"active_toggle"`
// 			} `mapstructure:"github_creation"`

// 			Author struct {
// 				// TUI author
// 				Focused string  `mapstructure:"focused"`
// 				Blurred string  `mapstructure:"blurred"`
// 				CursorModeHelp  string `mapstructure:"cursor_mode_help"`
// 				Cursor string	`mapstructure:"cursor"`
// 			} `mapstructure:"author_creation"`
// 			Groups struct {
// 				// TUI Groups
// 				ModelStyle string `mapstructure:"group"`
// 				FocusedModelStyle string `mapstructure:"group_focused"`
// 			} `mapstructure:"group_selection"`
// 		} `mapstructure:"light"`
// 	} `mapstructure:"style"`
// }

type SettingsConfig struct {
	AuthorFile       string `mapstructure:"author_file"`
	StartingScope    string `mapstructure:"starting_scope"`
	Editor           string `mapstructure:"editor"`
	DefaultStoreRepo string `mapstructure:"default_store_repo"`
}

// PaletteConfig holds one full set of style colors — used once for the
// default (dark) palette and once for the light-mode override.
type PaletteConfig struct {
	Help                string `mapstructure:"help"`
	Item                string `mapstructure:"item"`
	SelectedItemFG      string `mapstructure:"selected_item_fg"`
	HighlightFG         string `mapstructure:"highlight_fg"`
	SelectedHighlightFG string `mapstructure:"selected_highlight_fg"`
	SelectedItemBG      string `mapstructure:"selected_item_bg"`
	HighlightBG         string `mapstructure:"highlight_bg"`
	SelectedHighlightBG string `mapstructure:"selected_highlight_bg"`
	Delete              string `mapstructure:"delete"`
	Sharing             string `mapstructure:"sharing"`
	Pasting             string `mapstructure:"pasting"`
	ActivePaginationDot string `mapstructure:"active_pagination_dot"`
	GitScope            string `mapstructure:"git_scope"`
	LocalScope          string `mapstructure:"local_scope"`
	MixedScope          string `mapstructure:"mixed_scope"`

	CommitMessage CommitMessageStyleConfig `mapstructure:"commit_message_editor"`
	GH            GHStyleConfig            `mapstructure:"github_creation"`
	Author        AuthorStyleConfig        `mapstructure:"author_creation"`
	Groups        GroupsStyleConfig        `mapstructure:"group_selection"`
}

// TUI CommitMessageWriter
type CommitMessageStyleConfig struct {
	Base       string `mapstructure:"base"`
	LineNumber string `mapstructure:"line_number"`
}

// TUI GH
type GHStyleConfig struct {
	Error        string `mapstructure:"error"`
	Toggle       string `mapstructure:"toggle"`
	ActiveToggle string `mapstructure:"active_toggle"`
}

// TUI author
type AuthorStyleConfig struct {
	Focused        string `mapstructure:"focused"`
	Blurred        string `mapstructure:"blurred"`
	CursorModeHelp string `mapstructure:"cursor_mode_help"`
	Cursor         string `mapstructure:"cursor"`
}

// TUI Groups
type GroupsStyleConfig struct {
	ModelStyle        string `mapstructure:"group_box"`
	FocusedModelStyle string `mapstructure:"group_focused"`
}

// StyleConfig embeds PaletteConfig for the default/dark palette (so
// cfg.Style.Help, cfg.Style.CommitMessage.Base, etc. keep working exactly
// as before via promoted fields) and adds Light as the override palette.
type StyleConfig struct {
	PaletteConfig `mapstructure:",squash"`
	Light         PaletteConfig `mapstructure:"light"`
}

type Config struct {
	Settings SettingsConfig `mapstructure:"settings"`
	Style    StyleConfig    `mapstructure:"style"`
}

func (c *Config) String() string {
	return fmt.Sprintf("Author File: %s\nStarting Scope: %s\nEditor: %s",
		c.Settings.AuthorFile,
		c.Settings.StartingScope,
		c.Settings.Editor)
}

func init() {
	configDir, err := os.UserConfigDir()
	if err == nil {
		defaultConfigLocations[0] = filepath.Join(configDir, "cocommit")
	}
}

var v *viper.Viper

func LoadConfig() (*Config, error) {
	// TODO: create if and give param as default config location
	v = viper.New()
	v.SetConfigName(configName)
	v.SetConfigType(configType)

	// Set default values
	v.SetDefault("settings.author_file", defaultConfigLocations[0]+"/authors.json")
	v.SetDefault("settings.starting_scope", "local")
	v.SetDefault("settings.editor", "built-in")
	v.SetDefault("settings.default_user_store_repository", "https://github.com/Slug-Boi/cocommit_user_store")

	v.SetDefault("style.item", "170")
	v.SetDefault("style.selected_item_fg", "170")
	v.SetDefault("style.selected_item_bg", "236")
	v.SetDefault("style.highlight_fg", "170")
	v.SetDefault("style.highlight_bg", "206")
	v.SetDefault("style.selected_highlight_fg", "90")
	v.SetDefault("style.selected_highlight_bg", "206")
	v.SetDefault("style.delete", "9")
	v.SetDefault("style.sharing", "49")
	v.SetDefault("style.pasting", "86")
	v.SetDefault("style.active_pagination_dot", "170")
	v.SetDefault("style.help", "15")
	v.SetDefault("style.git_scope", "49")
	v.SetDefault("style.local_scope", "170")
	v.SetDefault("style.mixed_scope", "178")

	v.SetDefault("style.github_creation.error", "9")
	v.SetDefault("style.github_creation.toggle", "99")
	v.SetDefault("style.github_creation.active_toggle", "205")

	v.SetDefault("style.author_creation.focused", "170")
	v.SetDefault("style.author_creation.blurred", "240")
	v.SetDefault("style.author_creation.cursor", "170")
	v.SetDefault("style.author_creation.cursor_mode_help", "244")

	v.SetDefault("style.group_selection.group", "241")
	v.SetDefault("style.group_selection.group_focused", "170")

	v.SetDefault("style.commit_message_editor.base", "170")
	v.SetDefault("style.commit_message_editor.line_number", "90")

	// Lightmode
	v.SetDefault("style.light.item", "170")
	v.SetDefault("style.light.selected_item_fg", "170")
	v.SetDefault("style.light.selected_item_bg", "236")
	v.SetDefault("style.light.highlight_fg", "170")
	v.SetDefault("style.light.highlight_bg", "206")
	v.SetDefault("style.light.selected_highlight_fg", "90")
	v.SetDefault("style.light.selected_highlight_bg", "206")
	v.SetDefault("style.light.delete", "9")
	v.SetDefault("style.light.sharing", "49")
	v.SetDefault("style.light.pasting", "86")
	v.SetDefault("style.light.active_pagination_dot", "170")
	v.SetDefault("style.light.help", "15")
	v.SetDefault("style.light.git_scope", "49")
	v.SetDefault("style.light.local_scope", "170")
	v.SetDefault("style.light.mixed_scope", "178")
	v.SetDefault("style.light.github_creation.error", "9")
	v.SetDefault("style.light.github_creation.toggle", "99")
	v.SetDefault("style.light.github_creation.active_toggle", "205")

	v.SetDefault("style.light.github_creation.error", "9")
	v.SetDefault("style.light.github_creation.toggle", "99")
	v.SetDefault("style.light.github_creation.active_toggle", "205")

	v.SetDefault("style.light.author_creation.focused", "170")
	v.SetDefault("style.light.author_creation.blurred", "240")
	v.SetDefault("style.light.author_creation.cursor", "170")
	v.SetDefault("style.light.author_creation.cursor_mode_help", "244")

	v.SetDefault("style.light.group_selection.group", "241")
	v.SetDefault("style.light.group_selection.group_focused", "170")

	v.SetDefault("style.light.commit_message_editor.base", "170")
	v.SetDefault("style.light.commit_message_editor.line_number", "90")


	// Add search paths
	for _, path := range defaultConfigLocations {
		if path != "" {
			v.AddConfigPath(path)
		}
	}

	// Try to read config
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, nil
		}
		// if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		// 	if err := handleMissingConfig(v); err != nil {
		// 		return nil, err
		// 	}
		// } else {
		// 	return nil, fmt.Errorf("config error: %w", err)
		// }
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

func HandleMissingConfig() error {
	fmt.Println("Config file not found. Would you like to create one? (y/n)")
	var response string
	if _, err := fmt.Scanln(&response); err != nil {
		return fmt.Errorf("error reading response: %w", err)
	}

	yesResponses := map[string]bool{"y": true, "Y": true, "yes": true, "Yes": true, "YES": true}
	if !yesResponses[strings.TrimSpace(response)] {
		return fmt.Errorf("config file not found")
	}

	if v == nil {
		v = viper.New()

		v.SetConfigName(configName)
		v.SetConfigType(configType)
	}

	return CreateConfig()
}

func CheckConfig() bool {
	if v == nil {
		return false
	} else if v.ConfigFileUsed() == "" {
		return false
	} else {
		return true
	}
}

func GetConfigFilePath() string {
	if v == nil || v.ConfigFileUsed() == "" {
		return ""
	}
	return v.ConfigFileUsed()
}

func RemoveConfig() error {
	if v == nil || v.ConfigFileUsed() == "" {
		return fmt.Errorf("no config file to remove")
	}

	configPath := v.ConfigFileUsed()
	if err := os.Remove(configPath); err != nil {
		return fmt.Errorf("failed to remove config file: %w", err)
	}

	fmt.Printf("Config file removed: %s\n", configPath)
	return nil
}

func CreateConfig() error {
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
