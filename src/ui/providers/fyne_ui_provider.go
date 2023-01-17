package UI

import (
	"image/color"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	entities "p2p-messenger/src/domain/entities"
)

type FyneUIProvider struct{
	chatHistory *widget.Label
}

func (p* FyneUIProvider) Run(){
	myApp := app.New()
	myWindow := myApp.NewWindow("Border Layout")
	myWindow.Resize(fyne.NewSize(1024, 768))
	myWindow.SetFixedSize(true)

	left := canvas.NewText("friends placeholder", color.White)
	p.chatHistory = widget.NewLabel("")
	p.chatHistory.Wrapping = fyne.TextWrapWord
	messageEntry := widget.NewMultiLineEntry()
	sendMessageButton := widget.NewButton("Send", func() {
		p.chatHistory.SetText(p.chatHistory.Text + "\n<<< " + messageEntry.Text)
		p.chatHistory.Refresh()
		messageEntry.SetText("")
	})
	textScroller := container.NewVScroll(p.chatHistory)
	messageEntryContainer := container.NewBorder(nil, nil, nil, sendMessageButton, messageEntry)
	chatContainer := container.NewBorder(nil, messageEntryContainer, nil, nil, textScroller)
	content := container.NewBorder(nil, nil, left, nil, chatContainer)

	myWindow.SetContent(content)
	myWindow.ShowAndRun()
}

func (p* FyneUIProvider) ShowNemIncomingMessage(message entities.Message){
	p.chatHistory.SetText(p.chatHistory.Text + "\n>>> " + message.Text)
}