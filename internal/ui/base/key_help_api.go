package base

import "github.com/charmbracelet/bubbles/key"

var (
	QuitKeyMap     = key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("ctrl+c", "quit"))
	CommandKeyMap  = key.NewBinding(key.WithKeys(":"), key.WithHelp(":", "command"))
	EscKeyMap      = key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "return"))
	EditKeyMap     = key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit"))
	AddKeyMap      = key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "add"))
	FilterKeyMap   = key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "filter"))
	EnterKeyMap    = key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "ok"))
	DeleteKeyMap   = key.NewBinding(key.WithKeys("backspace", "delete"), key.WithHelp("delete", "delete"))
	ShowJsonKeyMap = key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "show json"))
	CopyKeyMap     = key.NewBinding(key.WithKeys("c"), key.WithHelp("c", "copy"))
	SwitchKeyMap   = key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "switch"))
)

type KeyHelpApi interface {
	GetKeys() []key.Binding
}

func NewKeyHelpApi(keyBindings ...key.Binding) KeyHelpApi {
	var ks keys
	ks = append(ks, keyBindings...)
	return &ks
}

type keys []key.Binding

func (k keys) GetKeys() []key.Binding {
	return k
}
