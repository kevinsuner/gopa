/*
* SPDX-License-Identifier: GPL-3.0-only
* Copyright (C) 2024 Kevin Suñer <ksuner@pm.me>
 */

package main

import (
	"log/slog"
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var menuText string = `
[:#458588][#ebdbb2][::b] GOPA [-:-:-:-] Run code [:#ebdbb2][#282828][::b] ^R [-:-:-:-] Go version [:#ebdbb2][#282828][::b] ^L [-:-:-:-]`

type playground struct {
    app             *tview.Application
    layout          *tview.Flex
    editor          *tview.TextArea
    console         *tview.TextView
    menu            *tview.TextView
    versionsList    *tview.List
}

func newPlayground(app *tview.Application) playground {
    return playground{
        app:            app,
        layout:         tview.NewFlex(),
        editor:         tview.NewTextArea(),
        console:        tview.NewTextView(),
        menu:           tview.NewTextView(),
        versionsList:   tview.NewList(),
    }
}

func (p playground) newLayout() *tview.Flex {
    p.layout.AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
        AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
            AddItem(p.editor, 0, 1, false).
            AddItem(p.console, 0, 1, false), 0, 1, false).
        AddItem(p.menu, 3, 1, false), 0, 1, false)

    return p.layout
}

func (p playground) extendEditor() *tview.TextArea {
    // Colors
    p.editor.SetBackgroundColor(tcell.ColorDefault)
    p.editor.SetBorderColor(tcell.Color245)
    p.editor.SetTitleColor(tcell.Color245)
    p.editor.SetTextStyle(p.editor.GetTextStyle().Background(tcell.ColorDefault))
    p.editor.SetTextStyle(p.editor.GetTextStyle().Foreground(tcell.Color223))
    p.editor.SetFocusFunc(func() {
        p.editor.SetBorderColor(tcell.Color223)
        p.editor.SetTextStyle(p.editor.GetTextStyle().Foreground(tcell.Color223))
        p.editor.SetTitleColor(tcell.Color223)})
    p.editor.SetBlurFunc(func() {
        p.editor.SetBorderColor(tcell.Color245)
        p.editor.SetTextStyle(p.editor.GetTextStyle().Foreground(tcell.Color245))
        p.editor.SetTitleColor(tcell.Color245)})

    // Defaults
    p.editor.SetWrap(false)
    p.editor.SetBorder(true)
    p.editor.SetTitle("Editor")
    p.editor.SetTitleAlign(tview.AlignRight)

    // Commands
    p.editor.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
        if event.Key() == tcell.KeyCtrlR {
            out, err := runCode(p.editor.GetText())
            if err != nil {
                slog.Error("runCode", "error", err.Error())
                os.Exit(1)
            }

            p.console.SetText(out)
            return nil
        }

        return event
    })

    return p.editor
}

func (p playground) extendConsole() *tview.TextView {
    // Colors
    p.console.SetBackgroundColor(tcell.ColorDefault)
    p.console.SetBorderColor(tcell.Color245)
    p.console.SetTitleColor(tcell.Color245)
    p.console.SetTextColor(tcell.Color245)
    p.console.SetFocusFunc(func() {
        p.console.SetBorderColor(tcell.Color223)
        p.console.SetTextColor(tcell.Color223)
        p.console.SetTitleColor(tcell.Color223)})
    p.console.SetBlurFunc(func() {
        p.console.SetBorderColor(tcell.Color245)
        p.console.SetTextColor(tcell.Color245)
        p.console.SetTitleColor(tcell.Color245)})

    // Defaults
    p.console.SetBorder(true)
    p.console.SetWordWrap(true)
    p.console.SetTitle("Console")
    p.console.SetTitleAlign(tview.AlignRight)

    return p.console
}

func (p playground) extendMenu() *tview.TextView {
    // Colors
    p.menu.SetBackgroundColor(tcell.ColorDefault)
    p.menu.SetTextColor(tcell.Color223)

    // Defaults
    p.menu.SetBorder(false)
    p.menu.SetDynamicColors(true)
    p.menu.SetText(menuText)
    p.menu.SetTextAlign(tview.AlignCenter)

    return p.menu
}

func (p playground) newVersionsList() *tview.List {
    p.versionsList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
        if event.Key() == tcell.KeyEsc {
            p.app.SetRoot(p.layout, true)
            return nil
        }

        return event
    })

    return p.versionsList
}

