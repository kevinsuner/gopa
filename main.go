package main

import (
	"os"
	"os/exec"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Think of switching UIs as an array of UIs ui[0]

type playground struct {
    editor *tview.TextArea
    console *tview.TextView
    menu *tview.Box
    flex *tview.Flex
}

func newPlayground() playground {
    console := newConsole()
    editor := newEditor(console)
    menu := newMenu()

    return playground{
        editor: editor,
        console: console,
        menu: menu,
        flex: tview.NewFlex().
            AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
                AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
                    AddItem(editor, 0, 1, false).
                    AddItem(console, 0, 1, false), 0, 1, false).
                AddItem(menu, 5, 1, false), 0, 1, false),
    }
}

func newConsole() *tview.TextView {
    console := tview.NewTextView().SetWordWrap(true)
    console.SetTitle("Console").SetBorder(true)
    return console
}

func newEditor(console *tview.TextView) *tview.TextArea {
    editor := tview.NewTextArea().
        SetWrap(false).
        SetPlaceholder("Type some code here...")
    
    editor.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
        if event.Key() == tcell.KeyCtrlR {
            out, err := runCode(editor.GetText())
            if err != nil {
                panic(err)
            }
            
            console.SetText(out)
            return nil
        }

        return event
    })

    editor.SetTitle("Editor").SetBorder(true)
    return editor
}

func runCode(input string) (string, error) {
    file, err := createTempFile()
    if err != nil {
        return "", err
    }

    defer file.Close()
    defer os.Remove(file.Name())

    _, err = file.WriteString(input)
    if err != nil {
        return "", err
    }

    cmd := exec.Command("go", "run", file.Name())
    out, _ := cmd.CombinedOutput()
    return string(out), nil
}

func createTempFile() (*os.File, error) {
    return os.CreateTemp("./", "tempfile-*.go")
}

func newMenu() *tview.Box {
    return tview.NewBox().SetBorder(true)
}

func main() {
    app := tview.NewApplication()
    playground := newPlayground()
    if err := app.SetRoot(playground.flex, true).EnableMouse(true).EnablePaste(true).Run(); err != nil {
        panic(err)
    }
}
