package form

import "github.com/a-h/templ"

type fieldAutocomplete string

const (
	EmailAutocomplete       fieldAutocomplete = "email"
	PasswordAutocomplete    fieldAutocomplete = "new-password"
	OldPasswordAutocomplete fieldAutocomplete = "old-password"
	FullNameAutocomplete    fieldAutocomplete = "name"
	PhoneAutocomplete       fieldAutocomplete = "tel"
	OffAutocomplete         fieldAutocomplete = "off"
)

type fieldInputType = string

const (
	TextInputType     fieldInputType = "text"
	TextAreaInputType fieldInputType = "textarea"
	ButtonInputType   fieldInputType = "button"
	EmailInputType    fieldInputType = "email"
	PhoneInputType    fieldInputType = "tel"
)

type Field struct {
	Name            string
	Autocomplete    fieldAutocomplete
	Required        bool
	ID              string
	Type            fieldInputType
	Placeholder     string
	AriaInvalid     string
	AriaDescribedBy string
	Value           any
	Disabled        bool
}

func (f *Field) Attributes(value any) templ.Attributes {
	return formFieldToTemplAttr(f, value)
}

func (f *Field) WithErrors(errs map[string]string) *Field {
	if errs != nil && len(errs) > 0 {
		name := f.Name
		if _, ok := errs[name]; ok {
			f.AriaInvalid = "true"
			f.AriaDescribedBy = name + "-helper"
		}
	}

	return f
}

func formFieldToTemplAttr(field *Field, value any) templ.Attributes {
	attr := templ.Attributes{}

	if field.Name != "" {
		attr["name"] = field.Name
	}

	if field.Autocomplete != "" {
		attr["autocomplete"] = field.Autocomplete
	}

	if field.Required {
		attr["required"] = true
	}

	if field.ID != "" {
		attr["id"] = field.ID
	}

	if field.Type != "" {
		attr["type"] = field.Type
	}

	if field.Placeholder != "" {
		attr["placeholder"] = field.Placeholder
	}

	if field.AriaInvalid != "" {
		attr["aria-invalid"] = field.AriaInvalid
	}

	if field.AriaDescribedBy != "" {
		attr["aria-describedby"] = field.AriaDescribedBy
	}

	switch field.Value.(type) {
	case string:
		if field.Value != "" {
			attr["value"] = field.Value
		}
	case bool:
		attr["value"] = field.Value
	default:
		// TODO log this case
	}

	if value != nil {
		switch value.(type) {
		case string:
			if value != "" {
				attr["value"] = value
			}
		case bool:
			attr["value"] = value
		default:
			// TODO log this case
		}
	}

	return attr
}
