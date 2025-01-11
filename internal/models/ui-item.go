package models

import "github.com/charmbracelet/glamour"

type UiItem interface {
	GetID() int
	Title() string
	Description() string
	FilterValue() string
	GetURL() string
	GetPreview(g *glamour.TermRenderer) string
}
