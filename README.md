# godex

godex is a powerful command-line file manager that simplifies file operations with features for searching, compression, cloud backup integration, and file versioning.

## Features

- **File Search**: Fast and flexible file search functionality
- **Compression Tools**: Zip and unzip files with ease
- **Google Drive Backup**: Seamless cloud backup integration
- **File Versioning**: Create, list, compare, restore and remove file versions
- **Shell Completion**: Built-in shell completion script generation

## Build from Source

### **1️⃣ Install Go (1.21 or later)**

#### **Linux (Debian/Ubuntu)**

```sh
sudo apt update
sudo apt install -y golang
```

#### **Linux (Arch Linux)**

```sh
sudo pacman -S go
```

#### **macOS**

```sh
brew install go
```

Verify the installation:

```sh
go version
```

### **2️⃣ Clone the Repository**

```sh
git clone https://github.com/inodinwetrust10/godex
cd godex
```

### **3️⃣ Build the Binary**

```sh
go build -o godex
```

This will generate an executable named `godex` in the same directory.

To install it system-wide, move it to `/usr/local/bin`:

```sh
sudo mv godex /usr/local/bin/
```

Now you can run:

```sh
godex --help
```

### **4️⃣ Cross-Compile for Different Systems**

If you need to build for multiple platforms:

```sh
# Linux (x86_64)
GOOS=linux GOARCH=amd64 go build -o godex-linux

# macOS (x86_64)
GOOS=darwin GOARCH=amd64 go build -o godex-macos

# macOS (Apple Silicon - M1/M2)
GOOS=darwin GOARCH=arm64 go build -o godex-macos-arm
```

### **5️⃣ Install Dependencies (If Any)**

If your project has missing dependencies, run:

```sh
go mod tidy
```

To fetch dependencies:

```sh
go get -u ./...
```

### **6️⃣ Running godex**

Once built, run:

```sh
./godex
```

Or if installed system-wide:

```sh
godex
```

## Usage

### Command Structure

```bash
godex [flags]
godex [command]
```

### Available Commands

- `search`: Search files with various criteria
- `zip`: Zip one or more files into a .zip archive
- `unzip`: Unzip a .zip archive to a destination directory
- `backup`: Backup file to Google Drive
- `version`: File versioning operations
- `completion`: Generate the autocompletion script for the specified shell
- `help`: Help about any command

### Global Flags

```bash
-h, --help      Help for godex
-t, --toggle    Help message for toggle
-v, --verbose   Enable verbose output
```

### Search Command

Search for files in the specified root directory using various criteria including exact name match, file size range, and modification date range.

```bash
godex search [flags]
```

#### Search Flags

```bash
-h, --help                     Help for search
-M, --max-size int            Maximum file size in bytes
-m, --min-size int            Minimum file size in bytes
-a, --modified-after string   Find files modified after this date (YYYY-MM-DD)
-b, --modified-before string  Find files modified before this date (YYYY-MM-DD)
-n, --name string             Search by exact file name
-p, --path string             Root path for the search (default is current directory)
```

#### Search Examples

Search by exact filename:

```bash
godex search --name "document.pdf"
```

Search by file size range:

```bash
godex search --min-size 1000000 --max-size 5000000
```

Search by modification date:

```bash
godex search --modified-after "2024-01-01" --modified-before "2024-01-31"
```

Combined search:

```bash
godex search --path "/documents" --name "report.pdf" --modified-after "2024-01-01"
```

### Zip Command

Zip one or more files into a .zip archive. The command accepts an output zip filename followed by one or more input files.

```bash
godex zip [output.zip] [files...]
```

#### Zip Flags

```bash
-d, --dir    Zipping directory
-h, --help   Help for zip
```

#### Zip Examples

Zip a single file:

```bash
godex zip archive.zip document.pdf
```

Zip multiple files:

```bash
godex zip documents.zip file1.txt file2.pdf file3.docx
```

Zip a directory:

```bash
godex zip project-backup.zip -d ./myproject/
```

### Unzip Command

Unzip a .zip archive to a destination directory. The command requires an input zip file and a destination directory path.

```bash
godex unzip [input.zip] [destination]
```

#### Unzip Flags

```bash
-h, --help   Help for unzip
```

#### Unzip Examples

Unzip to current directory:

```bash
godex unzip archive.zip .
```

Unzip to specific directory:

```bash
godex unzip documents.zip ./extracted-files
```

Unzip to new directory:

```bash
godex unzip project-backup.zip ./project-restored
```

### Version Command

The version command provides file versioning capabilities, allowing you to create, list, compare, restore, and remove versions of your files.

```bash
godex version [command]
```

#### Available Version Commands

- `create`: Create a new version of a file
- `list`: List all versions of a file
- `restore`: Restore your file to a specific version
- `diff`: Check differences between two files or between a file and its version
- `remove`: Remove a specific version or all versions of a file

#### Create Command

Create a new version of a file with an optional commit message.

```bash
godex version create [filepath]
```

##### Create Flags

```bash
-m, --message string   Add a commit message (default "commit")
-h, --help            Help for create
```

##### Create Examples

Create a version with default commit message:

```bash
godex version create document.txt
```

Create a version with custom commit message:

```bash
godex version create document.txt -m "Added section 3"
```

#### List Command

List all versions of a file with their version IDs and commit messages.

```bash
godex version list [filepath]
```

##### List Flags

```bash
-h, --help   Help for list
```

##### List Examples

```bash
godex version list document.txt
```

#### Restore Command

Restore a file to a specific version using the version ID.

```bash
godex version restore [filepath] [versionID]
```

##### Restore Flags

```bash
-h, --help   Help for restore
```

##### Restore Examples

```bash
godex version restore document.txt v2
```

#### Diff Command

Check differences between two files or between a file and its last version.

```bash
godex version diff [filepath1] [filepath2]
```

##### Diff Flags

```bash
-d, --default   Compare with the last version
-h, --help      Help for diff
```

##### Diff Examples

Compare two specific files:

```bash
godex version diff document-v1.txt document-v2.txt
```

Compare a file with its last version:

```bash
godex version diff document.txt -d
```

#### Remove Command

Remove a specific version or all versions of a file.

```bash
godex version remove [filepath]
```

##### Remove Flags

```bash
-v, --version string   Remove a specific version
-h, --help            Help for remove
```

##### Remove Examples

Remove a specific version:

```bash
godex version remove document.txt -v v2
```

Remove all versions:

```bash
godex version remove document.txt
```

### Backup Command

Backup a file to Google Drive. The command requires a file path to backup.

```bash
godex backup [file]
```

#### Backup Flags

```bash
-h, --help   Help for backup
```

#### Backup Examples

Backup a single file:

```bash
godex backup important-document.pdf
```

#### Google Drive Setup

Before using the backup command:

1. Set up Google Cloud Project:

   - Create a project in Google Cloud Console
   - Enable Google Drive API
   - Create credentials (OAuth 2.0 Client ID (Desktop app))
   - Download the client configuration file and rename it credentials.json and place it in ~/.config/godex

2. First-time configuration:
   - Run any backup command
   - Follow the authentication flow in your browser
   - Grant necessary permissions to godex
   - It will show cannot connect to the browser
   - Copy the the url -- http://localhost/?state=state-token&code=4/0IudJceGNktoKZlk-0K-\_X_aCsib7868786pJzH71tR-mjyYEJy\_\_MFw&scope=https://www.googleapis.com/auth/drive.file
   - Copy the code in between &code=xxxxxxxxxxxx&scope the xxxx will be your code
   - Paste it in the terminal

## Configuration

1. For Google Drive integration:
   - Create a Google Cloud project
   - Enable Google Drive API
   - Download credentials file
   - Configure your Google Drive settings

## Dependencies

- Go 1.16+
- Google Drive API client library

# godex Autocompletion

generate the autocompletion script for `godex` for the specified shell.

## Usage

```sh
godex completion [command]
```

## Available Commands

- **bash** Generate the autocompletion script for Bash
- **fish** Generate the autocompletion script for Fish
- **powershell** Generate the autocompletion script for PowerShell
- **zsh** Generate the autocompletion script for Zsh

## Flags

```
-h, --help   help for completion
```

## Installation

To enable autocompletion for your shell, run the appropriate command below:

### Bash

```sh
echo 'source <(godex completion bash)' >> ~/.bashrc
source ~/.bashrc
```

### Zsh

```sh
echo 'source <(godex completion zsh)' >> ~/.zshrc
source ~/.zshrc
```

### Fish

```sh
godex completion fish | source
```

To make it persistent:

```sh
godex completion fish > ~/.config/fish/completions/godex.fish
```

### PowerShell

```powershell
godex completion powershell | Out-String | Invoke-Expression
```

To make it persistent, add it to your PowerShell profile:

```powershell
godex completion powershell > $PROFILE
```

For more details, use:

```sh
godex completion [command] --help
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Support

If you encounter any issues or have questions:

- Open an issue in the GitHub repository
- Contact: [adi4gbsingh@gmail.com]

## Acknowledgments

- Google Drive API team
- Go community
- All contributors

---
