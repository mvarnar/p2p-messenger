package UI

import (
	"image/color"
	entities "p2p-messenger/src/domain/entities"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"golang.design/x/clipboard"
)

type FyneUIProvider struct {
	chatHistory             *widget.Label
	userIdPlace             *widget.Label
	outgoingMessagesChannel chan entities.Message
}

func NewFyneUIProvider() FyneUIProvider {
	return FyneUIProvider{outgoingMessagesChannel: make(chan entities.Message, 100)}
}

func (p *FyneUIProvider) Run() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Border Layout")
	myWindow.Resize(fyne.NewSize(1024, 768))
	myWindow.SetFixedSize(true)

	left := canvas.NewText("friends placeholder", color.White)
	p.chatHistory = widget.NewLabel("")
	p.chatHistory.Wrapping = fyne.TextWrapWord
	messageEntry := widget.NewMultiLineEntry()
	sendMessageButton := widget.NewButton("Send", func() {
		p.outgoingMessagesChannel <- entities.Message{Text: messageEntry.Text}
		p.chatHistory.SetText(p.chatHistory.Text + "\n<<< " + messageEntry.Text)
		p.chatHistory.Refresh()
		messageEntry.SetText("")
	})
	textScroller := container.NewVScroll(p.chatHistory)
	messageEntryContainer := container.NewBorder(nil, nil, nil, sendMessageButton, messageEntry)
	p.userIdPlace = widget.NewLabelWithStyle("Connecting to the network", fyne.TextAlignCenter, fyne.TextStyle{Monospace: true})
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}
	copyUserIdButton := widget.NewButton("Copy", func() { clipboard.Write(clipboard.FmtText, []byte(p.userIdPlace.Text)) })
	userIdContainer := container.NewBorder(nil, nil, nil, copyUserIdButton, p.userIdPlace)
	chatContainer := container.NewBorder(userIdContainer, messageEntryContainer, nil, nil, textScroller)
	content := container.NewBorder(nil, nil, left, nil, chatContainer)

	myWindow.SetContent(content)
	myWindow.ShowAndRun()
}

func (p *FyneUIProvider) ShowNemIncomingMessage(message entities.Message) {
	p.chatHistory.SetText(p.chatHistory.Text + "\n>>> " + message.Text)
}

func (p *FyneUIProvider) GetNewOutgoingMessages() <-chan entities.Message {
	return p.outgoingMessagesChannel
}

func (p *FyneUIProvider) ShowUserId(userId string) {
	p.userIdPlace.SetText("Your user id: " + userId)
}
