package components

import (
	"github.com/chenhg5/go-admin/template/types"
	"html/template"
)

type DescriptionAttribute struct {
	Name    string
	Border  string
	Number  string
	Title   template.HTML
	Arrow   string
	Color   template.HTML
	Percent string
	types.Attribute
}

func (compo *DescriptionAttribute) SetNumber(value string) types.DescriptionAttribute {
	compo.Number = value
	return compo
}

func (compo *DescriptionAttribute) SetTitle(value template.HTML) types.DescriptionAttribute {
	compo.Title = value
	return compo
}

func (compo *DescriptionAttribute) SetArrow(value string) types.DescriptionAttribute {
	compo.Arrow = value
	return compo
}

func (compo *DescriptionAttribute) SetPercent(value string) types.DescriptionAttribute {
	compo.Percent = value
	return compo
}

func (compo *DescriptionAttribute) SetColor(value template.HTML) types.DescriptionAttribute {
	compo.Color = value
	return compo
}

func (compo *DescriptionAttribute) SetBorder(value string) types.DescriptionAttribute {
	compo.Border = value
	return compo
}

func (compo *DescriptionAttribute) GetContent() template.HTML {
	return ComposeHtml(compo.TemplateList, *compo, "description")
}
