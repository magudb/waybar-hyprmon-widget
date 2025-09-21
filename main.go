package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

type WaybarOutput struct {
	Text    string `json:"text"`
	Tooltip string `json:"tooltip"`
	Class   string `json:"class"`
}

type Monitor struct {
	Name        string
	Width       int
	Height      int
	RefreshRate float64
	Scale       float64
	Active      bool
}

var rootCmd = &cobra.Command{
	Use:   "waybar-hyprmon-widget",
	Short: "A waybar widget for hyprland monitor management",
	Long:  `A waybar widget that displays monitor information and provides a TUI for managing hyprmon profiles.`,
	Run: func(cmd *cobra.Command, args []string) {
		monitors, err := getMonitors()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting monitors: %v\n", err)
			os.Exit(1)
		}

		output := WaybarOutput{
			Text:    getDisplayText(monitors),
			Tooltip: getTooltipText(monitors),
			Class:   "hyprmon-widget",
		}

		jsonOutput, err := json.Marshal(output)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error marshaling JSON: %v\n", err)
			os.Exit(1)
		}

		fmt.Println(string(jsonOutput))
	},
}

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Launch the TUI for profile selection",
	Long:  `Launch an interactive terminal user interface for selecting and applying hyprmon profiles.`,
	Run: func(cmd *cobra.Command, args []string) {
		runTUI()
	},
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func getMonitors() ([]Monitor, error) {
	cmd := exec.Command("hyprctl", "monitors", "-j")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute hyprctl: %v", err)
	}

	var hyprMonitors []map[string]any
	if err := json.Unmarshal(output, &hyprMonitors); err != nil {
		return nil, fmt.Errorf("failed to parse hyprctl output: %v", err)
	}

	var monitors []Monitor
	for _, hm := range hyprMonitors {
		monitor := Monitor{
			Name:   hm["name"].(string),
			Active: true,
		}

		if width, ok := hm["width"].(float64); ok {
			monitor.Width = int(width)
		}
		if height, ok := hm["height"].(float64); ok {
			monitor.Height = int(height)
		}
		if refresh, ok := hm["refreshRate"].(float64); ok {
			monitor.RefreshRate = refresh
		}
		if scale, ok := hm["scale"].(float64); ok {
			monitor.Scale = scale
		}

		monitors = append(monitors, monitor)
	}

	return monitors, nil
}

func getDisplayText(monitors []Monitor) string {
	if len(monitors) == 0 {
		return "No monitors"
	}

	if len(monitors) == 1 {
		return "󰍹"
	}

	return fmt.Sprintf("󰍹 %d", len(monitors))
}

func getTooltipText(monitors []Monitor) string {
	if len(monitors) == 0 {
		return "No monitors detected"
	}

	var lines []string

	// Get current active profile
	activeProfile, err := getActiveProfile()
	if err == nil && activeProfile != "" {
		lines = append(lines, fmt.Sprintf("Active Profile: %s", activeProfile))
		lines = append(lines, "")
	}

	lines = append(lines, "Monitor Setup:")
	lines = append(lines, "")

	for i, monitor := range monitors {
		line := fmt.Sprintf("%d. %s - %dx%d@%.1fHz (%.1fx scale)",
			i+1, monitor.Name, monitor.Width, monitor.Height,
			monitor.RefreshRate, monitor.Scale)
		lines = append(lines, line)
	}

	lines = append(lines, "")
	lines = append(lines, "Click to open profile selector")

	return strings.Join(lines, "\n")
}

type model struct {
	profiles []string
	cursor   int
	err      error
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.profiles)-1 {
				m.cursor++
			}
		case "enter", " ":
			if len(m.profiles) > 0 {
				selectedProfile := m.profiles[m.cursor]
				if err := applyProfile(selectedProfile); err != nil {
					m.err = fmt.Errorf("failed to apply profile '%s': %v", selectedProfile, err)
					return m, nil
				}
			}
			return m, tea.Quit
		default:
			// Handle number key selection (1-9)
			if len(msg.String()) == 1 && msg.String() >= "1" && msg.String() <= "9" {
				num := int(msg.String()[0] - '0')
				if num > 0 && num <= len(m.profiles) {
					m.cursor = num - 1
					selectedProfile := m.profiles[m.cursor]
					if err := applyProfile(selectedProfile); err != nil {
						m.err = fmt.Errorf("failed to apply profile '%s': %v", selectedProfile, err)
						return m, nil
					}
					return m, tea.Quit
				}
			}
		}
	case tea.MouseMsg:
		switch msg.Action {
		case tea.MouseActionPress:
			if msg.Button == tea.MouseButtonLeft {
				// Calculate which profile was clicked based on Y position
				// Line 0: "Available monitor profiles:"
				// Line 1: empty line
				// Line 2 onwards: profiles (starting from cursor line 2)
				if msg.Y >= 2 && msg.Y < 2+len(m.profiles) {
					profileIndex := msg.Y - 2
					if profileIndex >= 0 && profileIndex < len(m.profiles) {
						m.cursor = profileIndex
						// Apply the selected profile immediately on click
						selectedProfile := m.profiles[m.cursor]
						if err := applyProfile(selectedProfile); err != nil {
							m.err = fmt.Errorf("failed to apply profile '%s': %v", selectedProfile, err)
							return m, nil
						}
						return m, tea.Quit
					}
				}
			}
		case tea.MouseActionMotion:
			if msg.Button == tea.MouseButtonWheelUp {
				if m.cursor > 0 {
					m.cursor--
				}
			} else if msg.Button == tea.MouseButtonWheelDown {
				if m.cursor < len(m.profiles)-1 {
					m.cursor++
				}
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	var s strings.Builder

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(0, 1).
		Bold(true)

	s.WriteString(titleStyle.Render("HyprMon - Monitor Profile Manager"))
	s.WriteString("\n\n")

	if m.err != nil {
		errorStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF5F56")).
			Bold(true)
		s.WriteString(errorStyle.Render(fmt.Sprintf("Error: %v", m.err)))
		s.WriteString("\n\n")

		helpStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			Italic(true)
		s.WriteString(helpStyle.Render("Press q to quit."))
		return s.String()
	}

	if len(m.profiles) == 0 {
		noProfilesStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			Italic(true)
		s.WriteString(noProfilesStyle.Render("No hyprmon profiles found."))
		s.WriteString("\n\n")

		helpStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			Italic(true)
		s.WriteString(helpStyle.Render("Press q to quit."))
		return s.String()
	}

	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#04B575")).
		Bold(true)
	s.WriteString(headerStyle.Render("Available monitor profiles:"))
	s.WriteString("\n\n")

	for i, profile := range m.profiles {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		itemStyle := lipgloss.NewStyle()
		if m.cursor == i {
			itemStyle = itemStyle.Background(lipgloss.Color("#383838"))
		}

		line := fmt.Sprintf("%s %d. %s", cursor, i+1, profile)
		s.WriteString(itemStyle.Render(line))
		s.WriteString("\n")
	}

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#626262")).
		MarginTop(2)

	help := `
Controls:
  ↑/k, ↓/j: Navigate  Enter/Space: Select Profile  1-9: Quick Select
  Mouse Wheel: Navigate  Click: Select  q/Ctrl+c/Esc: Quit`

	s.WriteString(helpStyle.Render(help))

	return s.String()
}

func runTUI() {
	// Set up signal handling to ensure clean exit
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		os.Exit(0)
	}()

	profiles, err := getHyprmonProfiles()

	m := model{
		profiles: profiles,
		cursor:   0,
		err:      err,
	}

	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running TUI: %v\n", err)
		os.Exit(1)
	}
}

func getHyprmonProfiles() ([]string, error) {
	cmd := exec.Command("hyprmon", "--list-profiles")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute hyprmon --list-profiles: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	var profiles []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			// Remove the active profile indicator (*) for display
			profileName := strings.TrimSpace(strings.TrimSuffix(line, "*"))
			profiles = append(profiles, profileName)
		}
	}

	return profiles, nil
}

func getActiveProfile() (string, error) {
	cmd := exec.Command("hyprmon", "--active-profile")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to execute hyprmon --active-profile: %v", err)
	}

	profile := strings.TrimSpace(string(output))
	return profile, nil
}

func applyProfile(profile string) error {
	cmd := exec.Command("hyprmon", "-profile", profile)
	return cmd.Run()
}
