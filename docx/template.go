package docx

import (
    //"fmt"
    "regexp"
    "errors"
    "strings"  
    "reflect"  
    "github.com/aymerick/raymond"
    "github.com/legion-zver/go-docx-templates/graph"
)

var (
    rxTemplateItem  = regexp.MustCompile(`\{\{\s*([\w|\.]+)\s*\}\}`)
    rxMergeCellV    = regexp.MustCompile(`\[\s?v-merge\s?\]`)
    rxMergeIndex    = regexp.MustCompile(`\[\s?index\s?:\s?[\d|\.|\,]+\s?\]`)
)

// Функционал шаблонизатора
func renderTemplateDocument(document *Document, v interface{}) error {
    if document != nil {
        // Проходимся по структуре документа
        for _, item := range document.Body.Items {
            if err := renderDocItem(item, v); err != nil {
                return err
            }
        }
        return nil
    }
    return errors.New("Not valid template document")
}

// Поиск элементов шаблона и спаивания текстовых элементов
func findTemplatePatternsInParagraph(p *ParagraphItem) {
    if p != nil {
        // Перебор элементов параграфа и поиск начал {{ и конца }}                
        var startItem *RecordItem        
        for index := 0; index < len(p.Items); index++ {
            i := p.Items[index]
            if i.Type() == Record {
                record := i.(*RecordItem)
                if record != nil {
                    if startItem != nil {
                        startItem.Text += record.Text
                        // Удаляем элемент                        
                        p.Items = append(p.Items[:index], p.Items[index+1:]...)
                        // Проверка на конец
                        if strings.Index(startItem.Text, "}}") < 0 {
                            index--
                            continue
                        }
                        //fmt.Println("Merge records = ", startItem.Text)
                    } else {
                        if strings.Index(record.Text, "{{") >= 0 {
                            startItem = record                            
                            continue
                        }
                    }
                } 
            }
            startItem = nil            
        }
    }
}

// Рендер элемента документа
func renderDocItem(item DocItem, v interface{}) error {
    switch elem := item.(type) {
        // Параграф
        case *ParagraphItem: {
            findTemplatePatternsInParagraph(elem)
            for _, i := range elem.Items {
                if err := renderDocItem(i, v); err != nil {
                    return err
                }
            }
        }
        // Запись
        case *RecordItem: {
            if len(elem.Text) > 0 {
                if rxTemplateItem.MatchString(elem.Text) {
                    out, err := raymond.Render(modeTemplateText(elem.Text), v)
                    if err != nil {
                        return err
                    }
                    elem.Text = out
                }
            }
        }
        // Таблица
        case *TableItem: {            
            for rowIndex := 0; rowIndex < len(elem.Rows); rowIndex++ {
                row := elem.Rows[rowIndex]
                if row != nil {
                    // Если массив
                    if obj, ok := haveArrayInRow(row, v); ok {
                        lines       := objToLines(obj)
                        template    := row.Clone()
                        currentRow  := row                        
                        for _, line := range lines {                            
                            if currentRow == nil {
                                currentRow = template.Clone()
                                // Insert Row
                                elem.Rows = append(elem.Rows[:rowIndex], append([]*TableRow{currentRow}, elem.Rows[rowIndex:]...)...)                                
                            }
                            if err := renderRow(currentRow, &line); err != nil {
                                return err
                            }
                            currentRow = nil
                            rowIndex++
                        }                        
                        template = nil
                        // Балансируем индекс                        
                        rowIndex--
                        continue
                    }
                    // Если нет
                    if err := renderRow(row, v); err != nil {
                        return err
                    }
                }
            }
            // После обхода таблицы проходимся по ячейкам и проверяем merge флаги
            // С конца таблицы, проверяем по ячейкам
            for rowIndex := len(elem.Rows)-1; rowIndex >= 0; rowIndex-- {
                // Обходим ячейки
                for cellIndex, cell := range elem.Rows[rowIndex].Cells {
                    if len(cell.Items) > 0 {
                        plainText := plainTextFromTableCell(cell)
                        // Если найден флаг соединения
                        if rxMergeCellV.MatchString(plainText) {
                            if rowIndex > 0 {
                                topCell := elem.Rows[rowIndex-1].Cells[cellIndex]                                
                                if topCell != nil {                                    
                                    if plainText == plainTextFromTableCell(topCell) {
                                        cell.Params.VerticalMerge = new(StringValue)                                                                                
                                        cell.Items = make([]DocItem, 0) // Clear
                                        continue
                                    }
                                }
                            }
                            cell.Params.VerticalMerge = new(StringValue)
                            cell.Params.VerticalMerge.Value = "restart"
                            // Очищаяем контент ячейки от индексов и merge флагов
                            removeMergeIndexFromCell(cell)
                        }
                    }
                }
            } 
        }
    }
    return nil
}

// removeMergeIndexFromCell - очищаяем контент ячейки от индексов и merge флагов
func removeMergeIndexFromCell(cell *TableCell) {
    if cell != nil {
        for _, item := range cell.Items {
            removeMergeIndexFromDocItem(item)
        }
    }
}

// removeMergeIndexFromDocItem - очищаяем контент элемента документа от индексов и merge флагов
func removeMergeIndexFromDocItem(item DocItem) {
    if item != nil {
        switch elem := item.(type) {
            case *ParagraphItem: {
                for _, i := range elem.Items {
                    removeMergeIndexFromDocItem(i)
                }
            }
            case *RecordItem: {                
                elem.Text = rxMergeIndex.ReplaceAllString(rxMergeCellV.ReplaceAllString(elem.Text, ""),"")
            }            
        }
    }
}

// objToLines - раскладываем объект на строки
func objToLines(v interface{}) []map[string]interface{} {
    node := new(graph.Node)
    node.FromObject(v)
    return node.ListMap()
}

// renderRow - вывод строки таблицы
func renderRow(row *TableRow, v interface{}) error {
    if row != nil {
        for _, cell := range row.Cells {
            if cell != nil {
                for _, item := range cell.Items {
                    if err := renderDocItem(item, v); err != nil {
                        return err
                    }
                }
            }
        }
    }
    return nil
}

// Модифицируем текст шаблона
func modeTemplateText(tpl string) string {
    //fmt.Println("Mode: ", tpl)    
    tpl = strings.Replace(tpl, "{{", "{{{", -1)
	tpl = strings.Replace(tpl, "}}", "}}}", -1)
    tpl = strings.Replace(tpl,".","_",-1)
    return strings.Replace(tpl,":length","_length",-1) 
}

// haveArrayInRow - содержится ли массив в строке
func haveArrayInRow(row *TableRow, v interface{}) (interface{}, bool) {
    if row != nil {
        for _, cell := range row.Cells {
            if match := rxTemplateItem.FindStringSubmatch(plainTextFromTableCell(cell)); match != nil && len(match) > 1 {                
                names := strings.Split(match[1], ".")
                if len(names) > 0 {
                    t   := reflect.TypeOf(v)
                    val := reflect.ValueOf(v)
                    var lastVal reflect.Value 
                    for _, name := range names {
                        t      := findType(t, name)
                        val, _ := findValue(val, name)
                        if t != nil {
                            if t.Kind() == reflect.Array || t.Kind() == reflect.Slice {
                                if lastVal.IsValid() {                                    
                                    return lastVal.Interface(), true                                                       
                                }
                                return v, true                                
                            }
                        } else {
                            break
                        }                        
                        lastVal = val                        
                    }
                }
            }
        } 
    }
    return nil, false
}

// Простой текс у ячейки
func plainTextFromTableCell(cell *TableCell) string {
    var result string
    if cell != nil {
        for _, item := range cell.Items {
            result += item.PlainText()            
        }
    }
    return result
}

// findType - получаем тип по имени
func findType(t reflect.Type, name string) reflect.Type {
    kind := t.Kind()
    // Если это ссылка, то получаем истенный тип
    if kind == reflect.Ptr || kind == reflect.Interface {
        t = t.Elem()
    }
    kind = t.Kind()
    if kind == reflect.Struct {
        if field, ok := t.FieldByName(name); ok {
            return field.Type
        }
    } 
    return nil
}

// findValue - получаем значение по имени
func findValue(v reflect.Value, name string) (reflect.Value, bool) {
    kind := v.Type().Kind()
    // Если это ссылка, то получаем истенный тип
    if kind == reflect.Ptr || kind == reflect.Interface {
        v = v.Elem()
    }
    kind = v.Type().Kind()
    if kind == reflect.Struct {
        v := v.FieldByName(name)
        if v.IsValid() {
            return v, true
        }        
    } 
    return v, false
}