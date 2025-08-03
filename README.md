# 🔐 AUTHIE 2FA Authenticator

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.20+-blue.svg)](https://golang.org/)
[![Platform](https://img.shields.io/badge/Platform-Linux%20%7C%20macOS%20%7C%20FreeBSD%20%7C%20OpenBSD-lightgrey.svg)](https://github.com/rtwair/rtp-2fa-authenticator)

A lightweight, cross-platform 2FA authenticator written in Go with seamless menu integration for both X11 and Wayland. Compatible with the [Google Authenticator specification](https://github.com/google/google-authenticator/wiki/Key-Uri-Format).

## ✨ Features

- 🚀 **Fast & Lightweight** - Written in Go for optimal performance
- 🔧 **Universal Menu Integration** - Works with dmenu (X11), wofi, rofi, and bemenu (Wayland)
- 🌐 **Cross-Platform** - Works on Linux, macOS, FreeBSD, and OpenBSD
- 📱 **Google Authenticator Compatible** - Supports standard TOTP tokens
- 🎯 **Simple CLI Interface** - Easy-to-use command-line operations
- 📋 **Smart Clipboard Integration** - Automatic copying to X11 and Wayland clipboards

## 🚀 Installation

```bash
git clone https://github.com/rtwair/rtp-2fa-authenticator.git
cd rtp-2fa-authenticator
make
```

## 📖 Usage

```bash
2FA Authenticator - Usage:
  authenticator add <account_name> <secret> [issuer]
  authenticator add-url <otpauth_url>
  authenticator list
  authenticator remove <account_name>
  authenticator generate <secret>
  authenticator dmenu
  authenticator watch
  authenticator info
```

### Examples

```bash
# Add a new account
authenticator add "My Account" "JBSWY3DPEHPK3PXP" "MyService"

# Add from OTP URL
authenticator add-url "otpauth://totp/MyService:user@example.com?secret=JBSWY3DPEHPK3PXP&issuer=MyService"

# List all accounts
authenticator list

# Generate code for a specific secret
authenticator generate "JBSWY3DPEHPK3PXP"

# Launch dmenu interface
authenticator dmenu

# Show storage location
authenticator info
```

## ⚡ Menu Integration

### X11 (dmenu/dwm)
For seamless integration with dwm, add the following to your `config.def.h` file and recompile:

```c
// Add to command definitions
static const char *tfacmd[] = { "authenticator", "dmenu", NULL };

// Add to key bindings
{ MODKEY, XK_e, spawn, {.v = tfacmd } },
```

### Wayland (wofi/sway)
For sway integration, add to your sway config:

```bash
# Add to ~/.config/sway/config
bindsym $mod+e exec authenticator dmenu
```

### Universal (rofi)
Rofi works on both X11 and Wayland:

```bash
# Add to your window manager config
bindsym $mod+e exec authenticator dmenu
```

The menu integration automatically detects and uses the best available menu program:
- **dmenu** (X11)
- **wofi** (Wayland)
- **rofi** (X11/Wayland)
- **bemenu** (Wayland)

Simply fuzzy-search your entries and the 2FA code will be automatically copied to your clipboard.

## 🖥️ Compatibility

Tested and verified on:
- **Linux** (x86_64, ARM)
- **macOS** (Intel, Apple Silicon)
- **FreeBSD**
- **OpenBSD**

## 🛣️ Roadmap

- [ ] macOS integration improvements
- [ ] Sketchybar integration
- [ ] GUI interface
- [ ] Import/export functionality
- [ ] Encrypted storage options

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## ❓ Support

If you have any questions or issues:
- 🐛 [Open an issue on GitHub](https://github.com/rtwair/rtp-2fa-authenticator/issues)
- 📧 Email: riyad@rtp.cc

## ☕ Support the Project

If you find this project helpful, let me know!

riyad@rtp.cc

---

**Author:** Riyad Twair
