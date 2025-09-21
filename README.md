# Waybar Hyprmon Widget

A TUI application for Waybar that displays monitor setup information and allows profile selection through hyprmon.

## Features

- Shows current monitor count in Waybar
- Displays detailed monitor information in tooltip
- Interactive TUI for selecting hyprmon profiles
- Integrates seamlessly with Waybar

## Requirements

- Go 1.19+
- Hyprland with hyprctl
- hyprmon for profile management
- Waybar

## Installation

1. Build the application:
   ```bash
   go build -o waybar-hyprmon-widget
   ```

2. Copy to your bin directory:
   ```bash
   sudo cp waybar-hyprmon-widget /usr/local/bin/
   ```

## Waybar Configuration

Add this to your Waybar config:

```json
{
    "custom/hyprmon": {
        "format": "{}",
        "return-type": "json",
        "exec": "/usr/local/bin/waybar-hyprmon-widget",
        "interval": 30,
        "on-click": "/usr/local/bin/waybar-hyprmon-widget tui",
        "tooltip": true
    }
}
```

## Usage

- **Normal mode**: Outputs JSON for Waybar with monitor count and tooltip
- **TUI mode**: Run with `tui` argument to open profile selector

## Commands

- `waybar-hyprmon-widget` - Display monitor info for Waybar
- `waybar-hyprmon-widget tui` - Open interactive profile selector