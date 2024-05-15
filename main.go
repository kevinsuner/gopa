package main

import (
	"github.com/rivo/tview"
)

func main() {
    app := tview.NewApplication()
    
    textarea := tview.NewTextArea().
                SetWrap(false).
                SetPlaceholder("Type some code here...")
    textarea.SetTitle("Editor").SetBorder(true)
    
    flex := tview.NewFlex().
            AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
                    AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
                            AddItem(textarea, 0, 1, false).
                            AddItem(tview.NewBox().SetTitle("Console").SetBorder(true), 0, 1, false), 0, 1, false).
                    AddItem(tview.NewBox().SetBorder(true), 5, 1, false), 0, 1, false)
    if err := app.SetRoot(flex, true).EnableMouse(true).EnablePaste(true).Run(); err != nil {
        panic(err)
    }
}
