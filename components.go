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

func (p playground) newEditor() *tview.TextArea {
    p.editor.SetBackgroundColor(tcell.ColorDefault)
    p.editor.SetBorderColor(tcell.Color246)
    p.editor.SetTitleColor(tcell.Color246)
    p.editor.SetPlaceholderStyle(
        p.editor.GetPlaceholderStyle().Background(tcell.ColorDefault))
    p.editor.SetPlaceholderStyle(
        p.editor.GetPlaceholderStyle().Foreground(tcell.Color246))
    p.editor.SetTextStyle(
        p.editor.GetTextStyle().Background(tcell.ColorDefault))
    p.editor.SetWrap(false)
    p.editor.SetBorder(true)
    p.editor.SetTitle("Editor")
    p.editor.SetPlaceholder("Type some code here...")
    p.editor.SetFocusFunc(func() {
        p.editor.SetBorderColor(tcell.Color208)
        p.editor.SetTitleColor(tcell.Color208)
    })
    p.editor.SetBlurFunc(func() {
        p.editor.SetBorderColor(tcell.Color246)
        p.editor.SetTitleColor(tcell.Color246)
    })
    p.editor.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
        if event.Key() == tcell.KeyCtrlR {
            out, err := runCode(p.editor.GetText())
            if err != nil {
                slog.Error("RunCode", "error", err.Error())
                os.Exit(1)
            }

            p.console.SetText(out)
            return nil
        }

        return event
    })

    return p.editor
}

func (p playground) newConsole() *tview.TextView {
    p.console.SetBackgroundColor(tcell.ColorDefault)
    p.console.SetBorderColor(tcell.Color246)
    p.console.SetTitleColor(tcell.Color246)
    p.console.SetBorder(true)
    p.console.SetWordWrap(true)
    p.console.SetTitle("Console")
    p.console.SetFocusFunc(func() {
        p.console.SetBorderColor(tcell.Color208)
        p.console.SetTitleColor(tcell.Color208)
    })
    p.console.SetBlurFunc(func() {
        p.console.SetBorderColor(tcell.Color246)
        p.console.SetTitleColor(tcell.Color246)
    })

    return p.console
}

func (p playground) newMenu() *tview.TextView {
    p.menu.SetBackgroundColor(tcell.ColorDefault)
    p.menu.SetBorderColor(tcell.Color208)
    p.menu.SetTextColor(tcell.Color208)
    p.menu.SetBorder(true)
    p.menu.SetText("Run code (Ctrl-R) | Go version (Ctrl-L)")

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

