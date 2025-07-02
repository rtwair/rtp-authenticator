### RTP-2FA Authenticator

#### About

This is a 2FA Authenticator written in Go with integration functionality in dmenu. It is a simple implementation of the [Google Authenticator](https://github.com/google/google-authenticator/wiki/Key-Uri-Format) specification.

#### Usage

```
2FA Authenticator - Usage:
  go run main.go add <account_name> <secret> [issuer]
  go run main.go add-url <otpauth_url>
  go run main.go list
  go run main.go remove <account_name>
  go run main.go generate <secret>
  go run main.go dmenu
  go run main.go watch
  go run main.go info

Examples:
  go run main.go add "My Account" "JBSWY3DPEHPK3PXP" "MyService"
  go run main.go add-url "otpauth://totp/MyService:user@example.com?secret=JBSWY3DPEHPK3PXP&issuer=MyService"
  go run main.go list
  go run main.go generate "JBSWY3DPEHPK3PXP"
  go run main.go dmenu    # Quick dmenu selection
  go run main.go info     # Show storage location
```

#### Questions

If you have any questions, please open an issue on [Github](https://github.com/rtwair/rtp-2fa-authenticator/issues) or email me at riyad@rtp.cc.

Author: Riyad Twair
