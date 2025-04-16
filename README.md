# Ask Terminal AI

Ask Terminal AI is a command-line tool that allows users to quickly get and execute shell commands using AI or interact directly with AI in the terminal. It provides both a **virtual terminal mode** for command suggestions and a **conversation mode** for general queries.

> **⚡ DISCLAIMER: 99% of this project's code was generated by AI ⚡**


## Quickly download
1.**Download file**
```bash
wget -O ask "https://github.com/keizman/Ask-Terminal-AI/releases/download/main/ask_linux" 
```
2.**Generate config.yaml and edit**
`If u have not config.yaml at conf dir, it will auto generate and hint you file location`
```bash
ask 
```
`Modify model and apikey`
```bash
vim /root/.config/askta/config.yaml 
```

3.**start to use**


### Quickly transfer
1.**Copy from another machine**
```bash
scp -r root@{ip}:/root/.config/askta /root/.config/
```

2.**Make the binary executable:**
   ```bash
   chmod +x ask
   ```

3.**Optionally install to the system path:**
```bash
   sudo mv ask /usr/local/bin/ask
```

## Building from Source

1.**Clone the repository:**
   ```bash
   git clone https://github.com/keizman/Ask-Terminal-AI.git
   # Or if you already have the code:
   cd Ask-Terminal-AI
   ```

2.**Build the binary:**
   ```bash
   go build -o ask main.go
   ```

3.**Make the binary executable:**
   ```bash
   chmod +x ask
   ```

4.**Optionally install to the system path:**
   ```bash
   sudo mv ask /usr/local/bin/ask
   ```

---

## Configuration

Before using Ask Terminal AI, you need to create a configuration file.

1. **Create a configuration directory:**
   ```bash
   mkdir -p ~/.config/askta
   ```

2. **Create a `config.yaml` file:**
   ```bash
   nano ~/.config/askta/config.yaml
   ```

3. **Add the following content (replace with your actual values):**
   ```yaml
   # AI service configuration
   base_url: "https://api.openai.com/v1/"  # API base URL
   api_key: "your-api-key"                 # API key
   model_name: "gpt-4o-mini"               # Model to use

   # Feature configuration
   private_mode: false                     # Privacy mode
   sys_prompt: "I'm using Linux"           # Custom system prompt

   # Provider configuration
   provider: "openai-compatible"           # Currently only supports openai-compatible
   ```

4. **Save the file:** Use `Ctrl+O`, press `Enter`, then `Ctrl+X` to exit nano.

---

## Usage

### Virtual Terminal Mode (Command Suggestions)

- To get AI-suggested commands, type:
  ```bash
  ask 
  ```
  Then, enter a query like:
  ```bash
  how to find the largest files on my system
  ```

- You'll get a list of suggested commands. Here are the key bindings:
  - **Arrow keys (↑/↓):** Navigate suggestions
  - **Enter:** Execute the selected command
  - **`Ctrl+q` or `Ctrl+C`:** Exit

---

### Conversation Mode

- To have a conversation with the AI, use:
  ```bash
  ask -i "explain the difference between grep and awk"
  ```

---

### Options

| Option               | Description                                                                 |
|-----------------------|-----------------------------------------------------------------------------|
| `-c, --config FILE`   | Specify configuration file location                                        |
| `-m, --model NAME`    | Temporarily specify model to use                                           |
| `-p, --provider NAME` | Temporarily specify AI provider (currently only openai-compatible)         |
| `-u, --url URL`       | Temporarily specify API base URL                                           |
| `-k, --key KEY`       | Temporarily specify API key                                                |
| `-s, --sys-prompt TEXT` | Temporarily specify system prompt                                       |
| `--private-mode`      | Enable privacy mode                                                       |
| `-v, --version`       | Show version information                                                  |
| `-h, --help`          | Show help information                                                     |
| `-show`               | Show command history                                                      |

---

### Example with Options

```bash
ask --model gpt-4 --sys-prompt "I'm using Ubuntu 22.04" "how to install Docker"
```

---

## Logs

- **Command history:** `C:\Users\{user}\AppData\Local\Temp\askta_Chistory.log`  
- **Application logs:** `/tmp/askta_run.log`

---

## Security Notes

- **API keys** are stored encrypted on disk.
- Use `--private-mode` to avoid sending directory structure in queries.

---


## Develop

### generate go mod and sum
`go mod init ask_terminal && go mod tidy`

### Build with stripping and optimization flags 
`go build -ldflags="-s -w" -o ask.exe  main.go`


### Use UPX to compress the binary (install UPX first)
`upx --best ask.exe`
(`choco install upx`)



## License

[Apache 2.0](https://www.apache.org/licenses/LICENSE-2.0)