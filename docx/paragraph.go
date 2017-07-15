package docx

import (
	"encoding/xml"
	"errors"
)

// ParagraphItem - параграф
type ParagraphItem struct {
	Params       ParagraphParams `xml:"pPr"`
	Items        []DocItem
	RsidR        string `xml:"rsidR,attr,omitempty"`
	RsidRDefault string `xml:"rsidRDefault,attr,omitempty"`
	RsidP        string `xml:"rsidP,attr,omitempty"`
}

// ParagraphParams - параметры параграфа
type ParagraphParams struct {
	Style   *StringValue  `xml:"pStyle,omitempty"`
	Spacing *SpacingValue `xml:"spacing,omitempty"`
	Jc      *StringValue  `xml:"jc,omitempty"`
	Bidi    *IntValue     `xml:"bidi,omitempty"`
	PBdr    *PBdrValue    `xml"pBdr,omitempty"`
}

// Tag - имя тега элемента
func (item *ParagraphItem) Tag() string {
	return "p"
}

// Type - тип элемента
func (item *ParagraphItem) Type() DocItemType {
	return Paragraph
}

// PlainText - текст
func (item *ParagraphItem) PlainText() string {
	var result string
	for _, i := range item.Items {
		tmp := i.PlainText()
		if len(tmp) > 0 {
			result += tmp
		}
	}
	return result
}

// Clone - клонирование
func (item *ParagraphItem) Clone() DocItem {
	result := new(ParagraphItem)
	result.Items = make([]DocItem, 0)
	for _, i := range item.Items {
		if i != nil {
			result.Items = append(result.Items, i.Clone())
		}
	}
	// Клонирование параметров
	if item.Params.Bidi != nil {
		result.Params.Bidi = new(IntValue)
		result.Params.Bidi.Value = item.Params.Bidi.Value
	}
	if item.Params.Jc != nil {
		result.Params.Jc = new(StringValue)
		result.Params.Jc.Value = item.Params.Jc.Value
	}
	if item.Params.Spacing != nil {
		result.Params.Spacing = new(SpacingValue)
		result.Params.Spacing.After = item.Params.Spacing.After
		result.Params.Spacing.Before = item.Params.Spacing.Before
		result.Params.Spacing.Line = item.Params.Spacing.Line
		result.Params.Spacing.LineRule = item.Params.Spacing.LineRule
	}
	if item.Params.Style != nil {
		result.Params.Style = new(StringValue)
		result.Params.Style.Value = item.Params.Style.Value
	}
	if item.Params.PBdr != nil {
		result.Params.PBdr = new(PBdrValue)
		result.Params.PBdr.From(item.Params.PBdr)
	}
	return result
}

func (item *ParagraphParams) decode(decoder *xml.Decoder) error {
	if decoder != nil {
		var end bool
		for !end {
			token, _ := decoder.Token()
			if token == nil {
				break
			}
			switch element := token.(type) {
			case xml.StartElement:
				{
					if element.Name.Local == "pBdr" {
						decoder.DecodeElement(&item.PBdr, &element)
					} else if element.Name.Local == "pStyle" {
						decoder.DecodeElement(&item.Style, &element)
					} else if element.Name.Local == "spacing" {
						decoder.DecodeElement(&item.Spacing, &element)
					} else if element.Name.Local == "jc" {
						decoder.DecodeElement(&item.Jc, &element)
					} else if element.Name.Local == "bidi" {
						decoder.DecodeElement(&item.Bidi, &element)
					}
				}
			case xml.EndElement:
				{
					if element.Name.Local == "pPr" {
						end = true
					}
				}
			}
		}
		return nil
	}
	return errors.New("Not have decoder")

}

// Декодирование параграфа
func (item *ParagraphItem) decode(decoder *xml.Decoder) error {
	if decoder != nil {
		var end bool
		for !end {
			token, _ := decoder.Token()
			if token == nil {
				break
			}
			switch element := token.(type) {
			case xml.StartElement:
				{
					if element.Name.Local == "pPr" {
						item.Params.decode(decoder)
					} else {
						i := decodeItem(&element, decoder)
						if i != nil {
							item.Items = append(item.Items, i)
						}
					}
				}
			case xml.EndElement:
				{
					if element.Name.Local == "p" {
						end = true
					}
				}
			}
		}
		return nil
	}
	return errors.New("Not have decoder")
}

/* КОДИРОВАНИЕ */

// Кодирование параграфа
func (item *ParagraphItem) encode(encoder *xml.Encoder) error {
	if encoder != nil {
		//
		// RsidR        string `xml:"rsidR,attr,omitempty"`
		// RsidRDefault string `xml:"rsidRDefault,attr,omitempty"`
		// RsidP        string `xml:"rsidP,attr,omitempty"`
		rsidR := xml.Attr{Name: xml.Name{Local: "w:" + "rsidR"}, Value: item.RsidR}
		rsidRDefault := xml.Attr{Name: xml.Name{Local: "w:" + "rsidRDefault"}, Value: item.RsidRDefault}
		rsidP := xml.Attr{Name: xml.Name{Local: "w:" + "rsidP"}, Value: item.RsidP}
		// Начало параграфа

		start := xml.StartElement{Name: xml.Name{Local: "w:" + item.Tag()},
			Attr: []xml.Attr{rsidR, rsidRDefault, rsidP}}
		if err := encoder.EncodeToken(start); err != nil {
			return err
		}
		// Параметры параграфа
		if err := encoder.EncodeElement(&item.Params, xml.StartElement{Name: xml.Name{Local: "w:" + "pPr"}}); err != nil {
			return err
		}
		// Кодируем составные элементы
		for _, i := range item.Items {
			if err := i.encode(encoder); err != nil {
				return err
			}
		}
		// Конец параграфа
		if err := encoder.EncodeToken(start.End()); err != nil {
			return err
		}
		return encoder.Flush()
	}
	return errors.New("Not have encoder")
}
