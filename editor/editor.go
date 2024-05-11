/*
* SPDX-License-Identifier: GPL-3.0-only
* Copyright (C) 2024 Kevin Suñer <ksuner@pm.me>
 */

package editor

import (
	"gopa/colors"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/lipgloss"
)

var styles = map[string]lipgloss.Style{
    "cursor": lipgloss.NewStyle().Foreground(colors.Gray),
    "placeholder-focus": lipgloss.NewStyle().Foreground(colors.Gray),
    "placeholder-blur": lipgloss.NewStyle().Foreground(colors.Gray),
    "cursorline-focus": lipgloss.NewStyle().Background(colors.Background).Foreground(colors.Gray),
    "base-focus": lipgloss.NewStyle().Border(lipgloss.HiddenBorder()),
    "base-blur": lipgloss.NewStyle().Border(lipgloss.HiddenBorder()),
    "end-of-buffer": lipgloss.NewStyle().Foreground(colors.Gray),
}

func New() textarea.Model {
    ta := textarea.New()
    ta.Prompt = ""
    ta.Placeholder = "package main..."
    ta.ShowLineNumbers = false

    ta.Cursor.Style = styles["cursor"]
    ta.FocusedStyle.Placeholder = styles["placeholder-focus"]
    ta.BlurredStyle.Placeholder = styles["placeholder-blur"]
    ta.FocusedStyle.CursorLine = styles["cursorline-focus"]
    ta.FocusedStyle.Base = styles["base-focus"]
    ta.BlurredStyle.Base = styles["base-blur"]
    ta.FocusedStyle.EndOfBuffer = styles["end-of-buffer"]
    ta.BlurredStyle.EndOfBuffer = styles["end-of-buffer"]

    ta.KeyMap.LineNext = key.NewBinding(key.WithKeys("down"))
    ta.KeyMap.LinePrevious = key.NewBinding(key.WithKeys("up"))
    ta.KeyMap.WordForward = key.NewBinding(key.WithKeys("ctrl+right"))
    ta.KeyMap.WordBackward = key.NewBinding(key.WithKeys("ctrl+left"))
    ta.KeyMap.DeleteWordForward = key.NewBinding(key.WithKeys("ctrl+shift+right"))
    ta.KeyMap.DeleteWordBackward = key.NewBinding(key.WithKeys("ctrl+shift+left"))

    ta.Blur()
    return ta
}
