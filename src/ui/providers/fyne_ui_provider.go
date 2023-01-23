package UI

import (
	entities "p2p-messenger/src/domain/entities"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"golang.design/x/clipboard"
)

const userIdPrefix = "Your user id: "

type FyneUIProvider struct {
	chatHistory             *widget.Label
	userIdPlace             *widget.Label
	outgoingMessagesChannel chan entities.Message
	newContactChannel       chan entities.Contact
	contactsContainer       *fyne.Container
	chosenContact           entities.Contact
	userId                  string
	chosenContactButton     *widget.Button
	sendMessageButton       *widget.Button
	mainWindow            fyne.Window
}

func NewFyneUIProvider() FyneUIProvider {
	myApp := app.New()
	mainWindow := myApp.NewWindow("Messenger")
	mainWindow.Resize(fyne.NewSize(1024, 768))
	mainWindow.SetFixedSize(true)

	p := FyneUIProvider{
		outgoingMessagesChannel: make(chan entities.Message, 100),
		newContactChannel:       make(chan entities.Contact, 100),
	}
	addNewContactWindow := p.buildAddNewContactWindow(myApp)
	contactsContainer := p.buildContactsCotainer(addNewContactWindow)
	chatContainer := p.buildChatContainer()
	content := container.NewBorder(nil, nil, contactsContainer, nil, chatContainer)

	mainWindow.SetContent(content)
	p.mainWindow = mainWindow
	return p
}

func (p *FyneUIProvider) Run() {
	p.mainWindow.ShowAndRun()
}

func (p *FyneUIProvider) buildChatContainer() *fyne.Container {
	p.chatHistory = widget.NewLabel("")
	p.chatHistory.Wrapping = fyne.TextWrapWord
	textScroller := container.NewVScroll(p.chatHistory)

	messageEntry := widget.NewMultiLineEntry()
	p.sendMessageButton = widget.NewButton("Send", func() {
		p.outgoingMessagesChannel <- entities.Message{
			Text:            messageEntry.Text,
			ReceiverContact: p.chosenContact,
			SenderContact:   entities.Contact{UserId: p.userId},
		}
		p.chatHistory.SetText(p.chatHistory.Text + "\n<<< " + messageEntry.Text)
		p.chatHistory.Refresh()
		messageEntry.SetText("")
	})
	p.sendMessageButton.Disable()

	messageEntryContainer := container.NewBorder(nil, nil, nil, p.sendMessageButton, messageEntry)

	p.userIdPlace = widget.NewLabelWithStyle("Connecting to the network", fyne.TextAlignCenter, fyne.TextStyle{Monospace: true})
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}
	copyUserIdButton := widget.NewButton("Copy", func() { clipboard.Write(clipboard.FmtText, []byte(p.userIdPlace.Text[len(userIdPrefix):])) })
	userIdContainer := container.NewBorder(nil, nil, nil, copyUserIdButton, p.userIdPlace)

	chatContainer := container.NewBorder(userIdContainer, messageEntryContainer, nil, nil, textScroller)
	return chatContainer
}

func (p *FyneUIProvider) buildContactsCotainer(addNewContactWindow fyne.Window) *fyne.Container {
	p.contactsContainer = container.NewVBox()
	contactsScroller := container.NewVScroll(p.contactsContainer)
	addNewContactButton := widget.NewButton("Add a new contact", func() { addNewContactWindow.Show() })
	contactsContainer := container.NewBorder(nil, addNewContactButton, nil, nil, contactsScroller)
	return contactsContainer
}

func (p *FyneUIProvider) buildAddNewContactWindow(myApp fyne.App) fyne.Window {
	addNewContactWindow := myApp.NewWindow("Add a new contact")
	addNewContactWindow.Resize(fyne.NewSize(768, 50))
	addNewContactWindow.SetFixedSize(true)
	newUserIdEntry := widget.NewEntry()
	confirmAddNewContactButton := widget.NewButton("Add", func() {
		p.newContactChannel <- entities.Contact{UserId: newUserIdEntry.Text}
		newUserIdEntry.SetText("")
		addNewContactWindow.Hide()
	})
	addNewContactContainer := container.NewBorder(nil, nil, nil, confirmAddNewContactButton, newUserIdEntry)
	addNewContactWindow.SetContent(addNewContactContainer)
	return addNewContactWindow
}

func (p *FyneUIProvider) ShowNemIncomingMessage(message entities.Message) {
	p.chatHistory.SetText(p.chatHistory.Text + "\n>>> " + message.Text)
}

func (p *FyneUIProvider) GetNewOutgoingMessages() <-chan entities.Message {
	return p.outgoingMessagesChannel
}

func (p *FyneUIProvider) ShowUserId(userId string) {
	p.userId = userId
	p.userIdPlace.SetText(userIdPrefix + userId)
}

func (p *FyneUIProvider) GetNewContacts() <-chan entities.Contact {
	return p.newContactChannel
}

func (p *FyneUIProvider) ShowNewContact(contact entities.Contact) {
	var buttonText = ""
	if len(contact.UserId) < 13 {
		buttonText = contact.UserId
	} else {
		buttonText = contact.UserId[:8] + "..." + contact.UserId[len(contact.UserId)-5:]
	}
	var contactLabel *widget.Button
	contactLabel = widget.NewButton(buttonText, func() {
		if p.chosenContactButton != nil {
			p.chosenContactButton.Enable()
		}
		if p.sendMessageButton != nil && p.sendMessageButton.Disabled() {
			p.sendMessageButton.Enable()
		}
		p.chosenContactButton = contactLabel
		p.chosenContactButton.Disable()
		p.chosenContact = contact
		p.chatHistory.SetText("")
	})
	p.contactsContainer.Add(contactLabel)
	p.contactsContainer.Refresh()
}
