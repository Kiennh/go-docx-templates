package docx

import (
	"encoding/xml"
	"errors"
	"io"
)

// DocItemType - тип элемента
type DocItemType int

// Paragraph - параграф
const (
	Paragraph DocItemType = iota
	Record
	Table
	BookMark
)

// DocItem - интерфейс элемента документа
type DocItem interface {
	Tag() string
	Type() DocItemType
	PlainText() string
	Clone() DocItem
	decode(decoder *xml.Decoder) error
	encode(encoder *xml.Encoder) error
}

// Document - документ разметки DOCX
type Document struct {
	Scheme     map[string]string
	SkipScheme string
	Body       Body `xml:"body"`
}

// Body - тело документа
type Body struct {
	Items  []DocItem
	Params BodyParams `xml:"sectPr"`
}

// BodyParams - параметры тела документа
type BodyParams struct {
	HeaderReference *ReferenceValue `xml:"headerReference,omitempty"`
	FooterReference *ReferenceValue `xml:"footerReference,omitempty"`
	PageSize        SizeValue       `xml:"pgSz"`
	PageMargin      MarginValue     `xml:"pgMar"`
	Bidi            IntValue        `xml:"bidi"`
}

func (b *BodyParams) ToWBodyParams() WBodyParams {
	wp := WBodyParams{PageSize: WSizeValue(b.PageSize),
		PageMargin: WMarginValue(b.PageMargin),
		Bidi:       WIntValue(b.Bidi)}
	if b.HeaderReference != nil {
		wp.HeaderReference = (*WReferenceValue)(b.HeaderReference)
	}
	if b.FooterReference != nil {
		wp.FooterReference = (*WReferenceValue)(b.FooterReference)
	}
	return wp
}

type WBodyParams struct {
	HeaderReference *WReferenceValue `xml:"w:headerReference,omitempty"`
	FooterReference *WReferenceValue `xml:"w:footerReference,omitempty"`
	PageSize        WSizeValue       `xml:"w:pgSz"`
	PageMargin      WMarginValue     `xml:"w:pgMar"`
	Bidi            WIntValue        `xml:"w:bidi"`
}

/* ДЕКОДИРОВАНИЕ */

// Decode (Document) - декодирование документа
func (doc *Document) Decode(reader io.Reader) error {
	decoder := xml.NewDecoder(reader)
	if decoder != nil {
		doc.Scheme = make(map[string]string)
		for {
			token, _ := decoder.Token()
			if token == nil {
				break
			}
			switch element := token.(type) {
			case xml.StartElement:
				{
					if element.Name.Local == "document" {
						for _, attr := range element.Attr {
							if attr.Name.Local == "Ignorable" {
								doc.SkipScheme = attr.Value
							} else {
								doc.Scheme[attr.Name.Local] = attr.Value
							}
						}
					} else if element.Name.Local == "body" {
						err := doc.Body.decode(decoder)
						if err != nil {
							return err
						}
					}
				}
			}
		}
		return nil
	}
	return errors.New("Error create decoder")
}

// Декодирование BODY
func (body *Body) decode(decoder *xml.Decoder) error {
	if decoder != nil {
		if body.Items == nil {
			body.Items = make([]DocItem, 0)
		}
		for {
			token, _ := decoder.Token()
			if token == nil {
				break
			}
			switch element := token.(type) {
			case xml.StartElement:
				{
					if element.Name.Local == "sectPr" {
						decoder.DecodeElement(&body.Params, &element)
					} else {
						// Декодирование элементов
						item := decodeItem(&element, decoder)
						if item != nil {
							body.Items = append(body.Items, item)
						}
					}
				}
			case xml.EndElement:
				{
					if element.Name.Local == "body" {
						break
					}
				}
			}
		}
		return nil
	}
	return errors.New("Not have decoder")
}

func decodeItem(element *xml.StartElement, decoder *xml.Decoder) DocItem {
	if element != nil && decoder != nil {
		var item DocItem
		if element.Name.Local == "p" {
			item = new(ParagraphItem)
			pitem := item.(*ParagraphItem)
			for _, attr := range element.Attr {
				if attr.Name.Local == "rsidR" {
					pitem.RsidR = attr.Value
				}
				if attr.Name.Local == "rsidRDefault" {
					pitem.RsidRDefault = attr.Value
				}
				if attr.Name.Local == "rsidP" {
					pitem.RsidP = attr.Value
				}
				if attr.Name.Local == "rsidRPr" {
					pitem.RsidRPr = attr.Value
				}
			}
		} else if element.Name.Local == "r" {
			item = new(RecordItem)
		} else if element.Name.Local == "tbl" {
			item = new(TableItem)
		}
		if item != nil {
			if item.decode(decoder) == nil {
				return item
			}
		}
	}
	return nil
}

/* КОДИРОВАНИЕ */

// Encode - кодирование
func (doc *Document) Encode(writer io.Writer) error {
	encoder := xml.NewEncoder(writer)
	if encoder != nil {
		// Начало документа
		var attrs = make([]xml.Attr, 0)
		for key, val := range doc.Scheme {
			attrs = append(attrs, xml.Attr{Name: xml.Name{Local: "xmlns:" + key}, Value: val})
		}
		if len(doc.SkipScheme) > 0 {
			attrs = append(attrs, xml.Attr{Name: xml.Name{Local: "mc:Ignorable"}, Value: doc.SkipScheme})
		}
		docStart := xml.StartElement{Name: xml.Name{Local: "w:document"}, Attr: attrs}
		err := encoder.EncodeToken(docStart)
		if err != nil {
			return err
		}
		// Отдаем кодирование глубже - элементам
		err = doc.Body.encode(encoder)
		if err != nil {
			return err
		}
		// Конец документа
		err = encoder.EncodeToken(docStart.End())
		if err != nil {
			return err
		}
		return encoder.Flush()
	}
	return errors.New("Error create encoder")
}

// Кодирование BODY
func (body *Body) encode(encoder *xml.Encoder) error {
	if encoder != nil {
		// Начало BODY
		bodyStart := xml.StartElement{Name: xml.Name{Local: "w:" + "body"}}
		if err := encoder.EncodeToken(xml.StartElement{Name: xml.Name{Local: "w:" + "body"}}); err != nil {
			return err
		}
		// Переходим к элементам
		for _, item := range body.Items {
			if err := item.encode(encoder); err != nil {
				return err
			}
		}
		// Кодируем параметры
		if err := encoder.EncodeElement(body.Params.ToWBodyParams(), xml.StartElement{Name: xml.Name{Local: "w:" + "sectPr"}}); err != nil {
			return err
		}
		// Конец BODY
		if err := encoder.EncodeToken(bodyStart.End()); err != nil {
			return err
		}
		return encoder.Flush()
	}
	return errors.New("Not have encoder")
}
