/*
* SPDX-License-Identifier: GPL-3.0-only
* Copyright (C) 2024 Kevin Suñer <ksuner@pm.me>
 */

package main

import (
	"gopa/editor"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type playground struct {
    width int
    height int
    err error
    editor textarea.Model
    keymap map[string]key.Binding
}

func newPlayground() playground {
    return playground{
        editor: editor.New(),
        keymap: map[string]key.Binding{
            "tab": key.NewBinding(key.WithKeys("tab")),
            "quit": key.NewBinding(key.WithKeys("ctrl+q")),
        },
    }
}

func (p playground) Init() tea.Cmd {
    return textarea.Blink
}

func (p playground) View() string {
    var views []string
    views = append(views, p.editor.View())

    return lipgloss.JoinHorizontal(lipgloss.Top, views...)
}

func (p playground) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmds []tea.Cmd
    var cmd tea.Cmd

    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch {
        case key.Matches(msg, p.keymap["quit"]):
            return p, tea.Quit
        case key.Matches(msg, p.keymap["tab"]):
            p.editor.InsertRune('\t')
        default:
            if !p.editor.Focused() {
                cmd = p.editor.Focus()
                cmds = append(cmds, cmd)
            }
        }
    case tea.WindowSizeMsg:
        p.width = msg.Width
        p.height = msg.Height
    case error:
        p.err = msg
        return p, nil
    }

    p.resizeEditor()

    p.editor, cmd = p.editor.Update(msg)
    cmds = append(cmds, cmd)
    return p, tea.Batch(cmds...)
}

func (p *playground) resizeEditor() {
    p.editor.SetWidth(p.width / 2)
    p.editor.SetHeight(p.height - 5)
}

