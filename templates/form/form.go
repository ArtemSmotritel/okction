package form

import (
	"github.com/a-h/templ"
	"github.com/shopspring/decimal"
)

type fieldAutocomplete string

const DecimalPrecision = 2

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
	NumberInputType   fieldInputType = "number"
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
	Min             string
	Max             string
	Step            string
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

	if field.Min != "" {
		attr["min"] = field.Min
	}

	if field.Max != "" {
		attr["max"] = field.Max
	}

	if field.Step != "" {
		attr["step"] = field.Step
	}

	switch field.Value.(type) {
	case string:
		if field.Value != "" {
			attr["value"] = field.Value
		}
	case bool:
		attr["value"] = field.Value
	case decimal.Decimal:
		val, _ := field.Value.(decimal.Decimal)
		attr["value"] = val.StringFixedBank(DecimalPrecision)
	default:
	}

	if value != nil {
		switch value.(type) {
		case string:
			if value != "" {
				attr["value"] = value
			}
		case bool:
			attr["value"] = value
		case decimal.Decimal:
			val, _ := value.(decimal.Decimal)
			attr["value"] = val.StringFixedBank(DecimalPrecision)
		default:
		}
	}

	return attr
}
