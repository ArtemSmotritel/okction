package templates

import "github.com/a-h/templ"

type RootComponent func(components ...templ.Component) templ.Component

type HTMLPageBuilder struct {
	templates []templ.Component
	container RootComponent
}

func NewHTMLPageBuilder(container RootComponent) *HTMLPageBuilder {
	return &HTMLPageBuilder{
		templates: make([]templ.Component, 0),
		container: container,
	}
}

func (b *HTMLPageBuilder) AppendComponent(component templ.Component) {
	if component != nil {
		b.templates = append(b.templates, component)
	}
}

func (b *HTMLPageBuilder) Build() templ.Component {
	return b.container(b.templates...)
}
