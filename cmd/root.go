// cmd/root.go
package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"ask_terminal/config"
	"ask_terminal/terminal"
	"ask_terminal/utils"

	"github.com/spf13/cobra"
)

var (
	cfgFile     string
	modelName   string
	provider    string
	baseURL     string
	apiKey      string
	sysPrompt   string
	temperature float64
	maxTokens   uint
	privateMode bool
	showHistory bool
	proxyURL    string
	interactive bool // Add this for interactive mode
)

var rootCmd = &cobra.Command{
	Use:   "ask [query]",
	Short: "ASK Terminal AI - AI assistant for your terminal",
	Run: func(cmd *cobra.Command, args []string) {
		// Initialize logger
		logger := utils.NewLogger()

		// If -show flag is present, display command history
		if showHistory {
			displayCommandHistory(logger)
			return
		}

		// Load configuration
		conf, err := config.LoadConfig(cfgFile)
		if err != nil {
			logger.LogApplication(fmt.Sprintf("Error loading config: %v", err))
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			os.Exit(1)
		}

		// Apply command line overrides with flag changed checks
		if modelName != "" {
			conf.ModelName = modelName
		}
		if provider != "" {
			conf.Provider = provider
		}
		if baseURL != "" {
			conf.BaseURL = baseURL
		}
		if apiKey != "" {
			conf.APIKey = apiKey
		}
		if sysPrompt != "" {
			conf.SysPrompt = sysPrompt
		}

		// Only override temperature if the flag was changed
		if cmd.Flags().Changed("temp") {
			conf.Temperature = temperature
		}

		// Only override max_tokens if the flag was changed
		if cmd.Flags().Changed("max-tokens") {
			conf.MaxTokens = maxTokens
		}

		if privateMode {
			conf.PrivateMode = true
		}

		if proxyURL != "" {
			conf.Proxy = proxyURL
		}

		// Check if stdin has data (is being piped)
		stdinInfo, _ := os.Stdin.Stat()
		isPipe := (stdinInfo.Mode() & os.ModeCharDevice) == 0

		var query string
		if isPipe {
			// Read from stdin
			stdinData, err := io.ReadAll(os.Stdin)
			if err != nil {
				logger.LogApplication(fmt.Sprintf("Error reading stdin: %v", err))
				fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
				os.Exit(1)
			}

			// If args are provided, use them as the query and the stdin data as context
			if len(args) > 0 {
				query = strings.Join(args, " ") + "\n\nContent:\n" + string(stdinData)
			} else {
				// If no args, just use the stdin data as query
				query = string(stdinData)
			}

			// Start conversation mode with piped data
			terminal.StartConversationMode(query, conf)
		} else if len(args) > 0 || interactive {
			// Join all args to form the query if any
			query = strings.Join(args, " ")
			// Conversation mode (with args or interactive flag)
			terminal.StartConversationMode(query, conf)
		} else {
			// Virtual terminal mode
			terminal.StartVirtualTerminalMode(conf)
		}
	},
}

func init() {
	// Existing flags
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "Config file path")
	rootCmd.PersistentFlags().StringVarP(&modelName, "model", "m", "", "Model name to use")
	rootCmd.PersistentFlags().StringVarP(&provider, "provider", "p", "", "AI provider (openai-compatible)")
	rootCmd.PersistentFlags().StringVarP(&baseURL, "url", "u", "", "API base URL")
	rootCmd.PersistentFlags().StringVarP(&apiKey, "key", "k", "", "API key")
	rootCmd.PersistentFlags().StringVarP(&sysPrompt, "sys-prompt", "s", "", "System prompt")

	// Add temperature and maxTokens flags (without trying to set Changed callback)
	rootCmd.PersistentFlags().Float64Var(&temperature, "temp", 0, "Temperature (0.0-1.0)")
	rootCmd.PersistentFlags().UintVar(&maxTokens, "max-tokens", 0, "Max tokens (0 for unlimited)")

	// Existing boolean flags
	rootCmd.PersistentFlags().BoolVar(&privateMode, "private-mode", false, "Enable private mode")
	rootCmd.PersistentFlags().BoolVar(&showHistory, "show", false, "Show recent command history")

	// Add proxyURL flag
	rootCmd.PersistentFlags().StringVarP(&proxyURL, "proxy", "x", "", "Proxy URL (e.g., http://user:pass@host:port)")

	// Add interactive mode flag
	rootCmd.PersistentFlags().BoolVarP(&interactive, "interactive", "i", false, "Use interactive conversation mode")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// displayCommandHistory shows the recent command history
func displayCommandHistory(logger *utils.Logger) {
	items, err := logger.GetRecentCommands(1000)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error retrieving command history: %v\n", err)
		os.Exit(1)
	}

	if len(items) == 0 {
		fmt.Println("No command history found.")
		return
	}

	fmt.Printf("Recent commands (showing %d entries):\n\n", len(items))
	for i, item := range items {
		fmt.Printf("%d. [%s] Query: %s\n", i+1, item.Timestamp, item.Query)
		for cmd := range item.Commands {
			fmt.Printf("   - %s\n", cmd)
		}
		fmt.Println()
	}
}
