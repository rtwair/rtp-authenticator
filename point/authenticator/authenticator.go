package authenticator

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	//"strconv"
	"strings"
	"time"
)

// Account represents a 2FA account
type Account struct {
	Name   string
	Secret string
	Issuer string
}

// Authenticator manages 2FA Accounts
type Authenticator struct {
	Accounts []Account
	DataFile string
}

// NewAuthenticator creates a new authenticator instance
func NewAuthenticator() *Authenticator {
	// Get user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Printf("Warning: Could not get home directory: %v", err)
		homeDir = "."
	}

	// Create data directory if it doesn't exist
	dataDir := filepath.Join(homeDir, ".config", "2fa")
	os.MkdirAll(dataDir, 0700) // 0700 = rwx------

	DataFile := filepath.Join(dataDir, "accounts.json")

	auth := &Authenticator{
		Accounts: make([]Account, 0),
		DataFile: DataFile,
	}

	// Load existing Accounts
	auth.loadAccounts()
	return auth
}

// loadAccounts loads Accounts from the data file
func (a *Authenticator) loadAccounts() {
	data, err := os.ReadFile(a.DataFile)
	if err != nil {
		if os.IsNotExist(err) {
			// File doesn't exist yet, that's fine
			return
		}
		log.Printf("Warning: Could not read Accounts file: %v", err)
		return
	}

	if len(data) == 0 {
		return
	}

	err = json.Unmarshal(data, &a.Accounts)
	if err != nil {
		log.Printf("Warning: Could not parse Accounts file: %v", err)
		return
	}
}

// saveAccounts saves Accounts to the data file
func (a *Authenticator) saveAccounts() error {
	data, err := json.MarshalIndent(a.Accounts, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal Accounts: %v", err)
	}

	// Write to temporary file first, then rename (atomic operation)
	tempFile := a.DataFile + ".tmp"
	err = os.WriteFile(tempFile, data, 0600) // 0600 = rw-------
	if err != nil {
		return fmt.Errorf("failed to write temp file: %v", err)
	}

	err = os.Rename(tempFile, a.DataFile)
	if err != nil {
		os.Remove(tempFile) // Clean up temp file on error
		return fmt.Errorf("failed to rename temp file: %v", err)
	}

	return nil
}

// AddAccount adds a new 2FA account
func (a *Authenticator) AddAccount(name, secret, issuer string) error {
	// Validate and normalize the secret
	secret = strings.ToUpper(strings.ReplaceAll(secret, " ", ""))

	// Verify the secret is valid base32
	_, err := base32.StdEncoding.DecodeString(secret)
	if err != nil {
		return fmt.Errorf("invalid secret: %v", err)
	}

	// Check if account already exists
	for _, account := range a.Accounts {
		if account.Name == name {
			return fmt.Errorf("account '%s' already exists", name)
		}
	}

	account := Account{
		Name:   name,
		Secret: secret,
		Issuer: issuer,
	}

	a.Accounts = append(a.Accounts, account)

	// Save to file
	err = a.saveAccounts()
	if err != nil {
		// Remove the account we just added if save fails
		a.Accounts = a.Accounts[:len(a.Accounts)-1]
		return fmt.Errorf("failed to save account: %v", err)
	}

	return nil
}

// ListAccounts displays all Accounts with their current codes
func (a *Authenticator) ListAccounts() {
	if len(a.Accounts) == 0 {
		fmt.Println("No Accounts added yet.")
		return
	}

	fmt.Printf("%-20s %-15s %-8s %s\n", "Account", "Issuer", "Code", "Time Left")
	fmt.Println(strings.Repeat("-", 60))

	for _, account := range a.Accounts {
		code, err := GenerateTOTP(account.Secret, 30)
		if err != nil {
			fmt.Printf("%-20s %-15s %-8s %s\n", account.Name, account.Issuer, "ERROR", "N/A")
			continue
		}

		timeLeft := GetTimeRemaining()
		fmt.Printf("%-20s %-15s %-8s %ds\n", account.Name, account.Issuer, code, timeLeft)
	}
}

// RemoveAccount removes an account by name
func (a *Authenticator) RemoveAccount(name string) error {
	for i, account := range a.Accounts {
		if account.Name == name {
			a.Accounts = append(a.Accounts[:i], a.Accounts[i+1:]...)

			// Save to file
			err := a.saveAccounts()
			if err != nil {
				return fmt.Errorf("failed to save after removing account: %v", err)
			}

			return nil
		}
	}
	return fmt.Errorf("account '%s' not found", name)
}

// GenerateTOTP generates a TOTP code for the given secret
func GenerateTOTP(secret string, timeStep int64) (string, error) {
	// Decode the base32 secret
	key, err := base32.StdEncoding.DecodeString(secret)
	if err != nil {
		return "", fmt.Errorf("failed to decode secret: %v", err)
	}

	// Calculate the time counter (30-second intervals)
	counter := time.Now().Unix() / timeStep

	// Convert counter to 8-byte big-endian
	counterBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(counterBytes, uint64(counter))

	// Generate HMAC-SHA1
	h := hmac.New(sha1.New, key)
	h.Write(counterBytes)
	hash := h.Sum(nil)

	// Dynamic truncation
	offset := hash[len(hash)-1] & 0x0F
	truncated := binary.BigEndian.Uint32(hash[offset:offset+4]) & 0x7FFFFFFF

	// Generate 6-digit code
	code := truncated % uint32(math.Pow10(6))
	return fmt.Sprintf("%06d", code), nil
}

// GetTimeRemaining returns seconds until next code generation
func GetTimeRemaining() int {
	return 30 - int(time.Now().Unix()%30)
}

// CopyToClipboard copies text to system clipboard
func CopyToClipboard(text string) error {
	// Try different clipboard commands based on what's available
	var cmd *exec.Cmd

	// Check for xclip (Linux)
	if _, err := exec.LookPath("xclip"); err == nil {
		cmd = exec.Command("xclip", "-selection", "clipboard")
	} else if _, err := exec.LookPath("xsel"); err == nil {
		// Check for xsel (Linux alternative)
		cmd = exec.Command("xsel", "--clipboard", "--input")
	} else if _, err := exec.LookPath("pbcopy"); err == nil {
		// Check for pbcopy (macOS)
		cmd = exec.Command("pbcopy")
	} else if _, err := exec.LookPath("wl-copy"); err == nil {
		// Check for wl-copy (Wayland)
		cmd = exec.Command("wl-copy")
	} else {
		return fmt.Errorf("no clipboard utility found (install xclip, xsel, wl-copy, or pbcopy)")
	}

	cmd.Stdin = strings.NewReader(text)
	return cmd.Run()
}

// DmenuSelect shows a dmenu with account options and returns the selected account
func (a *Authenticator) DmenuSelect() (*Account, error) {
	if len(a.Accounts) == 0 {
		return nil, fmt.Errorf("no Accounts available")
	}

	// Check if dmenu is available
	if _, err := exec.LookPath("dmenu"); err != nil {
		return nil, fmt.Errorf("dmenu not found - please install dmenu")
	}

	// Prepare menu options
	var options []string
	for _, account := range a.Accounts {
		// Format: "Account Name (Issuer)"
		option := account.Name
		if account.Issuer != "" && account.Issuer != "Unknown" {
			option += " (" + account.Issuer + ")"
		}
		options = append(options, option)
	}

	// Run dmenu
	cmd := exec.Command("dmenu", "-i", "-p", "Select 2FA Account:")
	cmd.Stdin = strings.NewReader(strings.Join(options, "\n"))

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("dmenu selection cancelled or failed: %v", err)
	}

	selected := strings.TrimSpace(string(output))
	if selected == "" {
		return nil, fmt.Errorf("no account selected")
	}

	// Find the matching account
	for i, option := range options {
		if option == selected {
			return &a.Accounts[i], nil
		}
	}

	return nil, fmt.Errorf("selected account not found")
}

// DmenuGetCode shows dmenu for account selection and copies code to clipboard
func (a *Authenticator) DmenuGetCode() error {
	account, err := a.DmenuSelect()
	if err != nil {
		return err
	}

	// Generate TOTP code
	code, err := GenerateTOTP(account.Secret, 30)
	if err != nil {
		return fmt.Errorf("failed to generate code for %s: %v", account.Name, err)
	}

	// Copy to clipboard
	err = CopyToClipboard(code)
	if err != nil {
		// If clipboard fails, at least print the code
		fmt.Printf("Code for %s: %s (clipboard copy failed: %v)\n", account.Name, code, err)
		return nil
	}

	// Show success message
	timeLeft := GetTimeRemaining()
	fmt.Printf("Code for %s copied to clipboard: %s (valid for %d seconds)\n",
		account.Name, code, timeLeft)

	return nil
}

// ParseOTPAuthURL parses an otpauth:// URL and extracts account information
func ParseOTPAuthURL(url string) (Account, error) {
	// Basic validation
	if !strings.HasPrefix(url, "otpauth://totp/") {
		return Account{}, fmt.Errorf("invalid otpauth URL format")
	}

	// This is a simplified parser - in production you'd want more robust URL parsing
	parts := strings.Split(url, "?")
	if len(parts) != 2 {
		return Account{}, fmt.Errorf("malformed otpauth URL")
	}

	// Extract account name from path
	pathParts := strings.Split(parts[0], "/")
	if len(pathParts) < 3 {
		return Account{}, fmt.Errorf("missing account name in URL")
	}
	accountName := pathParts[2]

	// Parse query parameters
	params := make(map[string]string)
	for _, param := range strings.Split(parts[1], "&") {
		kv := strings.Split(param, "=")
		if len(kv) == 2 {
			params[kv[0]] = kv[1]
		}
	}

	secret, exists := params["secret"]
	if !exists {
		return Account{}, fmt.Errorf("missing secret parameter")
	}

	issuer := params["issuer"]
	if issuer == "" {
		issuer = "Unknown"
	}

	return Account{
		Name:   accountName,
		Secret: secret,
		Issuer: issuer,
	}, nil
}
