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
    editor          *tview.TextArea
    console         *tview.TextView
    menu            *tview.Box
    layout          *tview.Flex
    versionsList    *tview.List
}

func newPlayground(app *tview.Application) playground {
    return playground{
        app:            app,
        editor:         tview.NewTextArea(),
        console:        tview.NewTextView(),
        menu:           tview.NewBox(),
        layout:         tview.NewFlex(),
        versionsList:   tview.NewList(),
    }
}

func (p playground) newLayout() *tview.Flex {
    p.layout.AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
        AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
            AddItem(p.editor, 0, 1, false).
            AddItem(p.console, 0, 1, false), 0, 1, false).
        AddItem(p.menu, 5, 1, false), 0, 1, false)

    return p.layout
}

func (p playground) newEditor() *tview.TextArea {
    p.editor.SetWrap(false)
    p.editor.SetBorder(true)
    p.editor.SetTitle("Editor")
    p.editor.SetPlaceholder("Type some code here...")
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
    p.console.SetBorder(true)
    p.console.SetWordWrap(true)
    p.console.SetTitle("Console")

    return p.console
}

func (p playground) newMenu() *tview.Box {
    p.menu.SetBorder(true)

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

