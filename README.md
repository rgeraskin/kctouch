# kctouch

A command-line tool for managing macOS Keychain items with TouchID authentication.

## Overview

**kctouch** is a Go-based CLI tool that provides a simple and secure way to store, retrieve, and manage passwords in the macOS Keychain. All operations require TouchID authentication, ensuring that your credentials remain protected even if someone gains access to your terminal.

It works with generic-password items only.

## Features

- 🔐 **TouchID Authentication** - All operations require biometric authentication
- 🔑 **Keychain Integration** - Seamlessly works with macOS Keychain
- 📝 **Multiple Input Methods** - Support for interactive input, stdin, and command-line arguments
- 🛡️ **Secure by Default** - No plaintext password storage in command history
- ⚡ **Fast & Lightweight** - Single binary with minimal dependencies
- 🎯 **Simple Interface** - Intuitive commands with helpful aliases

## Requirements

- macOS with TouchID support
- Go 1.24.4 or later (for building from source)

## Installation

### From Source

```bash
git clone https://github.com/rgeraskin/kctouch.git
cd kctouch
go build -o kctouch
sudo mv kctouch /usr/local/bin/
```

### Using Go Install

```bash
go install github.com/rgeraskin/kctouch@latest
```

## Usage

All commands support the following flags:
- `-s, --service` - Service name (required for add/rm operations)
- `-a, --account` - Account name (optional)
- `-l, --label` - Label for the keychain item (optional, defaults to service name)
- `-v, --verbose` - Enable verbose logging
- `--cache-for` - Cache authentication for specified time duration (e.g. 1h, 10m, 10s)
- `--cache-n` - Cache authentication for N subsequent operations

### Adding Passwords

#### Interactive password input (recommended for security):
```bash
kctouch add -s "MyService"
# You'll be prompted to enter the password securely
```

#### From stdin:
```bash
echo "mypassword" | kctouch add -s "MyService" -p -
```

#### Direct password (not recommended for sensitive data):
```bash
kctouch add -s "MyService" -p "mypassword"
```

#### With custom label and account:
```bash
kctouch add -s "ghpat" -a "myuser" -l "GitHub Personal Access Token"
```

### Retrieving Passwords

```bash
# Get password by service name
kctouch get -s "GitHub"

# Get password by service and account
kctouch get -s "GitHub" -a "myusername"

# Get password with specific label
kctouch get -l "GitHub Personal Access Token"

# Using with scripts
API_KEY=$(kctouch get -s "myapi" -a "production")
curl -H "Authorization: Bearer $API_KEY" https://api.example.com/data

# Using in pipes
kctouch get -l keepass | keepassxc-cli show -a password -s passwords.kdbx my/password
```

### Removing Passwords

```bash
# Remove by service name
kctouch rm -s "GitHub"

# Remove by service and account
kctouch rm -s "GitHub" -a "myusername"
```

### Authentication Caching

To reduce the number of TouchID prompts during multiple operations, you can cache authentication:

```bash
# Get secret and cache authentication for 10 minutes
kctouch get -s /my/secret --cache-for 10m

# Add secret and cache authentication for 5 operations
kctouch add -s /my/secret --cache-n 5

# Remove secret and invalidate authentication cache for duration
kctouch rm -s /my/secret --cache-for 0

# Get secret and invalidate authentication cache for number of operations
kctouch get -s /my/secret --cache-n 0
```

If you set both `--cache-for` and `--cache-n`, the attempts will begin to decrease from `--cache-n` after `--cache-for` expires.

### Command Aliases

For faster typing, kctouch supports several aliases:

- **add**: `a`, `put`, `set`
- **get**: `g`, `find`
- **rm**: `d`, `del`, `delete`, `remove`

Examples:
```bash
kctouch a -s "service" -p "password"
kctouch g -s "service"
kctouch d -s "service"
```

### Debugging and Verbose Output

Use the verbose flag to see detailed logging information:

```bash
# Enable verbose logging for any command
kctouch add -s "service" -v
kctouch get -s "service" --verbose
kctouch noop --cache-for 5m -v
```

## Dependencies

- [cobra](https://github.com/spf13/cobra) - CLI framework
- [go-keychain](https://github.com/keybase/go-keychain) - macOS Keychain bindings
- [go-touchid](https://github.com/ansxuman/go-touchid) - TouchID authentication

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

If you encounter any issues or have questions, please open an issue on GitHub.
