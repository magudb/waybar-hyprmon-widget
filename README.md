# waybar-hyprmon-widget

A Waybar widget for Hyprland that displays monitor information and provides an interactive TUI for managing hyprmon profiles.

## Features

- 🖥️ **Monitor Detection**: Automatically detects and displays connected monitors
- 📊 **Real-time Status**: Shows current monitor configuration in Waybar
- 🎯 **Profile Management**: Interactive TUI for selecting and applying hyprmon profiles
- ⚡ **Fast Switching**: Quick profile switching with keyboard shortcuts or mouse
- 🎨 **Beautiful Interface**: Styled with Lip Gloss for a modern terminal UI
- 🐧 **Linux Native**: Built specifically for Hyprland on Linux

## Installation

### From Release (Recommended)

Download the latest binary from the [releases page](../../releases):

```bash
# For x86_64 systems
wget https://github.com/yourusername/waybar-hyprmon-widget/releases/latest/download/waybar-hyprmon-widget_Linux_x86_64.tar.gz
tar -xzf waybar-hyprmon-widget_Linux_x86_64.tar.gz
sudo mv waybar-hyprmon-widget /usr/local/bin/

# For ARM64 systems
wget https://github.com/yourusername/waybar-hyprmon-widget/releases/latest/download/waybar-hyprmon-widget_Linux_arm64.tar.gz
tar -xzf waybar-hyprmon-widget_Linux_arm64.tar.gz
sudo mv waybar-hyprmon-widget /usr/local/bin/
```

### From Source

```bash
git clone https://github.com/yourusername/waybar-hyprmon-widget.git
cd waybar-hyprmon-widget
go build -o waybar-hyprmon-widget .
sudo mv waybar-hyprmon-widget /usr/local/bin/
```

## Prerequisites

- **Hyprland**: The window manager this widget is designed for
- **hyprctl**: Hyprland's control utility (comes with Hyprland)
- **hyprmon**: Monitor profile management tool for Hyprland
- **Waybar**: The status bar where this widget will be displayed

## Usage

### Waybar Configuration

Add the following module to your Waybar configuration (`~/.config/waybar/config`):

```json
{
    "modules-left": ["hyprland/workspaces"],
    "modules-center": ["clock"],
    "modules-right": ["custom/hyprmon", "pulseaudio", "network", "battery"],

    "custom/hyprmon": {
        "format": "{}",
        "return-type": "json",
        "exec": "waybar-hyprmon-widget",
        "on-click": "waybar-hyprmon-widget tui",
        "interval": 10,
        "tooltip": true
    }
}
```

### Waybar Styling

Add styling to your Waybar CSS (`~/.config/waybar/style.css`):

```css
#custom-hyprmon {
    padding: 0 10px;
    margin: 0 4px;
    background-color: #2e3440;
    color: #d8dee9;
    border-radius: 4px;
}

#custom-hyprmon:hover {
    background-color: #3b4252;
}
```

### Command Line Usage

```bash
# Display current monitor information (JSON output for Waybar)
waybar-hyprmon-widget

# Launch interactive profile selector
waybar-hyprmon-widget tui

# Show help
waybar-hyprmon-widget --help
```

### TUI Controls

When you click the widget or run `waybar-hyprmon-widget tui`, an interactive terminal interface opens:

- **↑/k, ↓/j**: Navigate between profiles
- **Enter/Space**: Select and apply profile
- **1-9**: Quick select profile by number
- **Mouse**: Click to select, scroll to navigate
- **q/Ctrl+c/Esc**: Quit without changes

## Widget Output

The widget displays:

- **Single Monitor**: `󰍹` (monitor icon)
- **Multiple Monitors**: `󰍹 N` (where N is the number of monitors)
- **No Monitors**: `No monitors`

### Tooltip Information

The tooltip provides detailed information:

```
Active Profile: work-setup

Monitor Setup:

1. DP-1 - 2560x1440@144.0Hz (1.0x scale)
2. HDMI-A-1 - 1920x1080@60.0Hz (1.0x scale)

Click to open profile selector
```

## Development

### Building

```bash
git clone https://github.com/yourusername/waybar-hyprmon-widget.git
cd waybar-hyprmon-widget
go mod download
go build -o waybar-hyprmon-widget .
```

### Testing

```bash
# Run tests
go test -v ./...

# Run linting
golangci-lint run

# Check formatting
go fmt ./...

# Static analysis
staticcheck ./...
```

### Release

Releases are automatically created using GoReleaser when a new tag is pushed:

```bash
git tag v1.0.0
git push origin v1.0.0
```

## Dependencies

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Terminal styling
- [Cobra](https://github.com/spf13/cobra) - CLI framework

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Built for the [Hyprland](https://hyprland.org/) community
- Inspired by other Waybar widgets and the need for easy monitor management
- Uses the excellent [Charm](https://charm.sh/) TUI libraries

## Support

If you encounter any issues or have questions:

1. Check the [issues page](../../issues) for existing solutions
2. Create a new issue with detailed information about your setup
3. Include relevant logs and configuration files

---

**Made with ❤️ for the Hyprland community**