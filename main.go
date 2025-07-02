package main

import (
	"fmt"
	"os"
	"rtp-cc/authenticator/point/authenticator"
	"time"
)

func printUsage() {
	fmt.Println("2FA Authenticator - Usage:")
	fmt.Println("  authenticator add <account_name> <secret> [issuer]")
	fmt.Println("  authenticator add-url <otpauth_url>")
	fmt.Println("  authenticator list")
	fmt.Println("  authenticator remove <account_name>")
	fmt.Println("  authenticator generate <secret>")
	fmt.Println("  authenticator dmenu")
	fmt.Println("  authenticator watch")
	fmt.Println("  authenticator info")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  authenticator add \"My Account\" \"JBSWY3DPEHPK3PXP\" \"MyService\"")
	fmt.Println("  authenticator add-url \"otpauth://totp/MyService:user@example.com?secret=JBSWY3DPEHPK3PXP&issuer=MyService\"")
	fmt.Println("  authenticator list")
	fmt.Println("  authenticator generate \"JBSWY3DPEHPK3PXP\"")
	fmt.Println("  authenticator dmenu    # Quick dmenu selection")
	fmt.Println("  authenticator info     # Show storage location")
}

func main() {
	auth := authenticator.NewAuthenticator()

	if len(os.Args) < 2 {
		printUsage()
		return
	}

	command := os.Args[1]

	switch command {
	case "info":
		fmt.Printf("2FA Authenticator Storage Info:\n")
		fmt.Printf("Data file: %s\n", auth.DataFile)
		fmt.Printf("Accounts stored: %d\n", len(auth.Accounts))

		// Check if file exists and show its permissions
		if info, err := os.Stat(auth.DataFile); err == nil {
			fmt.Printf("File size: %d bytes\n", info.Size())
			fmt.Printf("File permissions: %s\n", info.Mode().String())
			fmt.Printf("Last modified: %s\n", info.ModTime().Format("2006-01-02 15:04:05"))
		} else {
			fmt.Printf("Data file status: Not created yet\n")
		}

	case "add":
		if len(os.Args) < 4 {
			fmt.Println("Usage: authenticator add <account_name> <secret> [issuer]")
			return
		}

		name := os.Args[2]
		secret := os.Args[3]
		issuer := "Unknown"
		if len(os.Args) > 4 {
			issuer = os.Args[4]
		}

		err := auth.AddAccount(name, secret, issuer)
		if err != nil {
			fmt.Printf("Error adding account: %v\n", err)
			return
		}
		fmt.Printf("Account '%s' added successfully!\n", name)

	case "add-url":
		if len(os.Args) < 3 {
			fmt.Println("Usage: authenticator add-url <otpauth_url>")
			return
		}

		account, err := authenticator.ParseOTPAuthURL(os.Args[2])
		if err != nil {
			fmt.Printf("Error parsing URL: %v\n", err)
			return
		}

		err = auth.AddAccount(account.Name, account.Secret, account.Issuer)
		if err != nil {
			fmt.Printf("Error adding account: %v\n", err)
			return
		}
		fmt.Printf("Account '%s' added successfully from URL!\n", account.Name)

	case "list":
		auth.ListAccounts()

	case "remove":
		if len(os.Args) < 3 {
			fmt.Println("Usage: authenticator remove <account_name>")
			return
		}

		name := os.Args[2]
		if auth.RemoveAccount(name) == nil {
			fmt.Printf("Account '%s' removed successfully!\n", name)
		} else {
			fmt.Printf("Account '%s' not found.\n", name)
		}

	case "generate":
		if len(os.Args) < 3 {
			fmt.Println("Usage: authenticator generate <secret>")
			return
		}

		secret := os.Args[2]
		code, err := authenticator.GenerateTOTP(secret, 30)
		if err != nil {
			fmt.Printf("Error generating code: %v\n", err)
			return
		}

		timeLeft := authenticator.GetTimeRemaining()
		fmt.Printf("Current code: %s (valid for %d seconds)\n", code, timeLeft)

	case "dmenu":
		err := auth.DmenuGetCode()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

	case "watch":
		// Continuous display mode
		for {
			// Clear screen (works on Unix-like systems)
			fmt.Print("\033[2J\033[H")
			fmt.Printf("2FA Codes - %s\n\n", time.Now().Format("15:04:05"))
			auth.ListAccounts()

			time.Sleep(1 * time.Second)
		}

	default:
		fmt.Printf("Unknown command: %s\n\n", command)
		printUsage()
	}
}
