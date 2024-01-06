package document

import (
	"fmt"
	"time"

	"github.com/signintech/gopdf"
)

var now = time.Now

type PdfDoc interface {
	AddText(x, y float64, content string) error
	AddFormattedText(x, y float64, content string, size int, style string)
	SetFontSize(size int) error
	SetFontStyle(style string) error
	DefaultFontSize()
	DefaultFontStyle()
	SetPosition(x, y float64)
	CellWithOption(rectangle *gopdf.Rect, text string, opt gopdf.CellOption) error
	Image(picPath string, x float64, y float64, rect *gopdf.Rect) error
	WritePdf(pdfPath string) error
	AddMultilineText(x, y float64, content string)
}

func AddAdressData(doc PdfDoc, receiver ReceiverAdress) {

	doc.AddText(130, 50, receiver.Header)
	doc.AddText(130, 57, receiver.Name)
	doc.AddText(130, 64, receiver.Adress)
	doc.AddText(130, 71, receiver.City)
}

func AddTitle(doc PdfDoc, titleWithDate TitleWithDate) {

	now := now()

	dateFormatted := fmt.Sprintf("%s, %02d.%02d.%04d", titleWithDate.City, now.Day(), now.Month(), now.Year())
	doc.AddText(130, 89, dateFormatted)

	doc.AddFormattedText(20, 100, titleWithDate.Title, 14, "bold")
}

func AddText(doc PdfDoc, multilineText string) {
	doc.AddMultilineText(20, 110, multilineText)
}

func AddTable(doc PdfDoc, tableData TableData) {
	doc.SetFontSize(12)
	doc.SetFontStyle("bold")

	initialHeight := 130
	rowHeight := 6

	cursor := 20.0

	for _, column := range tableData.Columns {
		setTextAligned(doc, cursor, float64(initialHeight), column.Header, column.Alignment, column.Width, gopdf.Bottom)
		cursor += column.Width
	}

	doc.DefaultFontSize()
	doc.DefaultFontStyle()

	cursor = 20.0
	doc.SetFontSize(11)

	for _, column := range tableData.Columns {
		for i, columnRow := range column.Rows {

			if i == (len(column.Rows)-1) && tableData.LastRowBold {
				doc.SetFontSize(12)
				doc.SetFontStyle("bold")
			}

			setTextAligned(doc, cursor, float64(initialHeight+(i+1)*rowHeight), columnRow, column.Alignment, column.Width, 0)

		}
		cursor += column.Width
		doc.DefaultFontSize()
		doc.DefaultFontStyle()

	}

}

func setTextAligned(doc PdfDoc, x, y float64, text string, alignment int, width float64, border int) {
	doc.SetPosition(x, y)
	doc.CellWithOption(&gopdf.Rect{W: width, H: 4.5}, text, gopdf.CellOption{Align: alignment, Border: border})

}
