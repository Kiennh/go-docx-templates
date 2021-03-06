package docx

import (
	"encoding/xml"
	"errors"
)

type Text struct {
	Value string `xml:",chardata"`
	Space string `xml:"space,attr,omitempty"`
}

type WText struct {
	Value string `xml:",chardata"`
	Space string `xml:"xml:space,attr,omitempty"`
}

func (t Text) ToWText() *WText {
	return &WText{Value: t.Value, Space: t.Space}
}

// RecordItem - record item
type RecordItem struct {
	Params  *RecordParams `xml:"rPr,omitempty"`
	Text    Text          `xml:"t,omitempty"`
	Tab     bool          `xml:"tab,omitempty"`
	Break   bool          `xml:"br,omitempty"`
	Drawing *Drawing      `xml:"drawing,omitempty"`
}

// RecordParams - params record
type RecordParams struct {
	Fonts     *RecordFonts `xml:"rFonts,omitempty"`
	Rtl       *IntValue    `xml:"rtl,omitempty"`
	Size      *IntValue    `xml:"sz,omitempty"`
	SizeCs    *IntValue    `xml:"szCs,omitempty"`
	Lang      *StringValue `xml:"lang,omitempty"`
	Underline *ShadowValue `xml:"u,omitempty"`
	Italic    *EmptyValue  `xml:"i,omitempty"`
	Bold      *EmptyValue  `xml:"b,omitempty"`
	BoldCS    *EmptyValue  `xml:"bCs,omitempty"`
	Color     *StringValue `xml:"color,omitempty"`
	Highlight *StyleValue  `xml:"highlight,omitempty"`
	VertAlign *StyleValue  `xml:"vertAlign,omitempty"`
	Strike    *EmptyValue  `xml:"strike,omitempty"`
	NoProof   *EmptyValue  `xml:"noProof,omitempty"`
}

func (rp *RecordParams) Clone() *RecordParams {

	result := new(RecordParams)
	if rp.Bold != nil {
		result.Bold = new(EmptyValue)
	}
	if rp.BoldCS != nil {
		result.BoldCS = new(EmptyValue)
	}
	if rp.Italic != nil {
		result.Italic = new(EmptyValue)
	}
	if rp.Underline != nil {
		result.Underline = new(ShadowValue)
		result.Underline.Value = rp.Underline.Value
	}
	if rp.Color != nil {
		result.Color = new(StringValue)
		result.Color.Value = rp.Color.Value
	}
	if rp.Lang != nil {
		result.Lang = new(StringValue)
		result.Lang.Value = rp.Lang.Value
	}
	if rp.Rtl != nil {
		result.Rtl = new(IntValue)
		result.Rtl.Value = rp.Rtl.Value
	}
	if rp.Size != nil {
		result.Size = new(IntValue)
		result.Size.Value = rp.Size.Value
	}
	if rp.SizeCs != nil {
		result.SizeCs = new(IntValue)
		result.SizeCs.Value = rp.SizeCs.Value
	}
	if rp.Fonts != nil {
		result.Fonts = new(RecordFonts)
		result.Fonts.ASCII = rp.Fonts.ASCII
		result.Fonts.CS = rp.Fonts.CS
		result.Fonts.EastAsia = rp.Fonts.EastAsia
		result.Fonts.HandleANSI = rp.Fonts.HandleANSI
		result.Fonts.HandleInt = rp.Fonts.HandleInt
	}
	return result
}

func (rp *RecordParams) ToWRecordParams() *WRecordParams {
	wrp := WRecordParams{}
	// Fonts     *WRecordFonts `xml:"w:rFonts,omitempty"`
	if rp.Fonts != nil {
		wrp.Fonts = (*WRecordFonts)(rp.Fonts)
	}
	// Rtl       *WIntValue    `xml:"w:rtl,omitempty"`
	if rp.Rtl != nil {
		wrp.Rtl = (*WIntValue)(rp.Rtl)
	}
	// Size      *WIntValue    `xml:"w:sz,omitempty"`
	if rp.Size != nil {
		wrp.Size = (*WIntValue)(rp.Size)
	}
	// SizeCs    *WIntValue    `xml:"w:szCs,omitempty"`
	if rp.SizeCs != nil {
		wrp.SizeCs = (*WIntValue)(rp.SizeCs)
	}
	// Lang      *WStringValue `xml:"w:lang,omitempty"`
	if rp.Lang != nil {
		wrp.Lang = (*WStringValue)(rp.Lang)
	}
	// Underline *WStringValue `xml:"w:u,omitempty"`
	if rp.Underline != nil {
		wrp.Underline = (*WShadowValue)(rp.Underline)
	}
	// Italic    *WEmptyValue  `xml:"w:i,omitempty"`
	if rp.Italic != nil {
		wrp.Italic = (*WEmptyValue)(rp.Italic)
	}
	// Bold      *WEmptyValue  `xml:"w:b,omitempty"`
	if rp.Bold != nil {
		wrp.Bold = (*WEmptyValue)(rp.Bold)
	}
	// BoldCS    *WEmptyValue  `xml:"w:bCs,omitempty"`
	if rp.BoldCS != nil {
		wrp.BoldCS = (*WEmptyValue)(rp.BoldCS)
	}
	// Color     *WStringValue `xml:"w:color,omitempty"`
	if rp.Color != nil {
		wrp.Color = (*WStringValue)(rp.Color)
	}
	if rp.Highlight != nil {
		wrp.Highlight = (*WStyleValue)(rp.Highlight)
	}
	if rp.VertAlign != nil {
		wrp.VertAlign = (*WStyleValue)(rp.VertAlign)
	}
	if rp.Strike != nil {
		wrp.Strike = (*WEmptyValue)(rp.Strike)
	}
	if rp.NoProof != nil {
		wrp.NoProof = (*WEmptyValue)(rp.NoProof)
	}
	return &wrp
}

type WRecordParams struct {
	Fonts     *WRecordFonts `xml:"w:rFonts,omitempty"`
	Rtl       *WIntValue    `xml:"w:rtl,omitempty"`
	Size      *WIntValue    `xml:"w:sz,omitempty"`
	SizeCs    *WIntValue    `xml:"w:szCs,omitempty"`
	Lang      *WStringValue `xml:"w:lang,omitempty"`
	Underline *WShadowValue `xml:"w:u,omitempty"`
	Italic    *WEmptyValue  `xml:"w:i,omitempty"`
	Bold      *WEmptyValue  `xml:"w:b,omitempty"`
	BoldCS    *WEmptyValue  `xml:"w:bCs,omitempty"`
	Color     *WStringValue `xml:"w:color,omitempty"`
	Highlight *WStyleValue  `xml:"w:highlight,omitempty"`
	VertAlign *WStyleValue  `xml:"w:vertAlign,omitempty"`
	Strike    *WEmptyValue  `xml:"w:strike,omitempty"`
	NoProof   *WEmptyValue  `xml:"w:noProof,omitempty"`
}

// RecordFonts - fonts in record
type RecordFonts struct {
	ASCII      string `xml:"ascii,attr"`
	CS         string `xml:"cs,attr"`
	HandleANSI string `xml:"hAnsi,attr"`
	EastAsia   string `xml:"eastAsia,attr"`
	HandleInt  string `xml:"hint,attr,omitempty"`
}

type WRecordFonts struct {
	ASCII      string `xml:"w:ascii,attr"`
	CS         string `xml:"w:cs,attr"`
	HandleANSI string `xml:"w:hAnsi,attr"`
	EastAsia   string `xml:"w:eastAsia,attr"`
	HandleInt  string `xml:"w:hint,attr,omitempty"`
}

// Tag - имя тега элемента
func (item *RecordItem) Tag() string {
	return "r"
}

// Type - тип элемента
func (item *RecordItem) Type() DocItemType {
	return Record
}

// PlainText - текст
func (item *RecordItem) PlainText() string {
	return item.Text.Value
}

// Clone - клонирование
func (item *RecordItem) Clone() DocItem {
	result := new(RecordItem)
	result.Text = item.Text
	result.Tab = item.Tab
	result.Break = item.Break
	// Клонируем параметры

	if item.Params == nil {
		return result
	}
	result.Params = item.Params.Clone()
	return result
}

// Декодирование записи
func (item *RecordItem) decode(decoder *xml.Decoder) error {
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
					if element.Name.Local == "rPr" {
						decoder.DecodeElement(&item.Params, &element)
					} else if element.Name.Local == "t" {
						decoder.DecodeElement(&item.Text, &element)
					} else if element.Name.Local == "br" {
						item.Break = true
					} else if element.Name.Local == "tab" {
						item.Tab = true
					} else if element.Name.Local == "drawing" {
						decoder.DecodeElement(&item.Drawing, &element)
					}
				}
			case xml.EndElement:
				{
					if element.Name.Local == "r" {
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

// Кодирование записи
func (item *RecordItem) encode(encoder *xml.Encoder) error {
	if encoder != nil {
		// Начало записи
		start := xml.StartElement{Name: xml.Name{Local: "w:" + item.Tag()}}
		if err := encoder.EncodeToken(start); err != nil {
			return err
		}
		// Параметры записи
		if item.Params != nil {
			if err := encoder.EncodeElement(item.Params.ToWRecordParams(), xml.StartElement{Name: xml.Name{Local: "w:" + "rPr"}}); err != nil {
				return err
			}
		}
		// Текст
		if err := encoder.EncodeElement(item.Text.ToWText(), xml.StartElement{Name: xml.Name{Local: "w:" + "t"}}); err != nil {
			return err
		}
		// todo: Drawing
		if item.Drawing != nil {
			if err := encoder.EncodeElement(item.Drawing.ToWDrawing(), xml.StartElement{Name: xml.Name{Local: "w:" + "drawing"}}); err != nil {
				return err
			}
		}
		// <br />
		if item.Break {
			startBr := xml.StartElement{Name: xml.Name{Local: "w:" + "br"}}
			if err := encoder.EncodeToken(startBr); err != nil {
				return err
			}
			if err := encoder.EncodeToken(startBr.End()); err != nil {
				return err
			}
			if err := encoder.Flush(); err != nil {
				return err
			}
		}
		// Tab
		if item.Tab {
			startTab := xml.StartElement{Name: xml.Name{Local: "w:" + "tab"}}
			if err := encoder.EncodeToken(startTab); err != nil {
				return err
			}
			if err := encoder.EncodeToken(startTab.End()); err != nil {
				return err
			}
			if err := encoder.Flush(); err != nil {
				return err
			}
		}
		// Конец записи
		if err := encoder.EncodeToken(start.End()); err != nil {
			return err
		}
		return encoder.Flush()
	}
	return errors.New("Not have encoder")
}
