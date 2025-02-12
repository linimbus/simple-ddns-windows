package main

import (
	"github.com/lxn/walk"
)

func ErrorBoxAction(form walk.Form, message string) {
	walk.MsgBox(form, "Error", message, walk.MsgBoxOK)
}

func InfoBoxAction(form walk.Form, message string) {
	walk.MsgBox(form, "Info", message, walk.MsgBoxOK)
}
