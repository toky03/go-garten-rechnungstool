package document

import (
	"fmt"
	"time"

	gopdf_wrapper "github.com/72nd/gopdf-wrapper"
	"github.com/signintech/gopdf"
	"github.com/toky03/qr-invoice/model"
)

func AddAdressData(doc *gopdf_wrapper.Doc, debtor model.DebtorData) {
	var heading string
	if debtor.Language == "de" {
		heading = fmt.Sprintf("Parzelle %s", debtor.Parzelle)
	} else {
		heading = fmt.Sprintf("Parcelle %s", debtor.Parzelle)
	}
	doc.AddText(130, 50, heading)
	doc.AddText(130, 57, debtor.Debtor.Name)
	doc.AddText(130, 64, debtor.Debtor.Address)
	doc.AddText(130, 71, fmt.Sprintf("%s %s", debtor.Debtor.Zip, debtor.Debtor.City))
}

func AddTitle(doc *gopdf_wrapper.Doc, language string, invoiceDetails model.InvoiceDetails) {
	var title string
	if language == "de" {
		title = invoiceDetails.Ueberschrift.De
	} else {
		title = invoiceDetails.Ueberschrift.Fr
	}

	now := time.Now()

	dateFormatted := fmt.Sprintf("%s, %02d. %02d. %04d", invoiceDetails.Creditor.City, now.Day(), now.Month(), now.Year())
	doc.AddText(130, 80, dateFormatted)

	doc.AddFormattedText(20, 90, title, 14, "bold")
}

func AddTableHeader(doc *gopdf_wrapper.Doc, language string, invoiceDetails model.InvoiceDetails) {

	doc.SetFontSize(12)
	doc.SetFontStyle("bold")

	setTextAligned(doc, 0, 100, tranlate(language, invoiceDetails.TabelleAnzahl), gopdf.Right)
	setTextAligned(doc, 35, 100, tranlate(language, invoiceDetails.TabelleEinheit), gopdf.Left)
	setTextAligned(doc, 65, 100, tranlate(language, invoiceDetails.TabelleBezeichnung), gopdf.Left)
	setTextAligned(doc, 110, 100, tranlate(language, invoiceDetails.TabellePreis), gopdf.Right)
	setTextAligned(doc, 140, 100, tranlate(language, invoiceDetails.TabelleBetrag), gopdf.Right)

	doc.DefaultFontSize()
	doc.DefaultFontStyle()
}

func AddTableData(doc *gopdf_wrapper.Doc, language string, debtorData model.DebtorData, invoiceDetails model.InvoiceDetails, variableData model.VariableData, tableData model.CalculatedData) {
	doc.SetFontSize(11)

	column_1 := 0
	column_2 := 35
	column_3 := 65
	column_4 := 110
	column_5 := 140

	setTextAligned(doc, column_1, 106, fmt.Sprintf("%.1f", debtorData.Are), gopdf.Right)
	setTextAligned(doc, column_2, 106, tranlate(language, invoiceDetails.TabelleAaren), gopdf.Left)
	setTextAligned(doc, column_3, 106, tranlate(language, variableData.TextPachtzins), gopdf.Left)
	setTextAligned(doc, column_4, 106, fmt.Sprintf("CHF %.2f", variableData.Pachtzins), gopdf.Right)
	setTextAligned(doc, column_5, 106, fmt.Sprintf("CHF %.2f", tableData.Pachtzins), gopdf.Right)

	setTextAligned(doc, column_1, 112, fmt.Sprintf("%.1f", debtorData.Are), gopdf.Right)
	setTextAligned(doc, column_2, 112, tranlate(language, invoiceDetails.TabelleAaren), gopdf.Left)
	setTextAligned(doc, column_3, 112, tranlate(language, variableData.TextWasserbezug), gopdf.Left)
	setTextAligned(doc, column_4, 112, fmt.Sprintf("CHF %.2f", variableData.Wasserbezug), gopdf.Right)
	setTextAligned(doc, column_5, 112, fmt.Sprintf("CHF %.2f", tableData.Wasserbezug), gopdf.Right)

	setTextAligned(doc, column_1, 118, "1", gopdf.Right)
	setTextAligned(doc, column_2, 118, tranlate(language, invoiceDetails.TabelleJahre), gopdf.Left)
	setTextAligned(doc, column_3, 118, tranlate(language, variableData.TextGfAbonement), gopdf.Left)
	setTextAligned(doc, column_4, 118, fmt.Sprintf("CHF %.2f", variableData.GfAbonement), gopdf.Right)
	setTextAligned(doc, column_5, 118, fmt.Sprintf("CHF %.2f", tableData.GfAbonement), gopdf.Right)

	setTextAligned(doc, column_1, 124, "1", gopdf.Right)
	setTextAligned(doc, column_2, 124, tranlate(language, invoiceDetails.TabelleJahre), gopdf.Left)
	setTextAligned(doc, column_3, 124, tranlate(language, variableData.TextStrom), gopdf.Left)
	setTextAligned(doc, column_4, 124, fmt.Sprintf("CHF %.2f", variableData.Strom), gopdf.Right)
	setTextAligned(doc, column_5, 124, fmt.Sprintf("CHF %.2f", tableData.Strom), gopdf.Right)

	setTextAligned(doc, column_1, 130, "1", gopdf.Right)
	setTextAligned(doc, column_2, 130, tranlate(language, invoiceDetails.TabelleJahre), gopdf.Left)
	setTextAligned(doc, column_3, 130, tranlate(language, variableData.TextVersicherung), gopdf.Left)
	setTextAligned(doc, column_4, 130, fmt.Sprintf("CHF %.2f", variableData.Versicherung), gopdf.Right)
	setTextAligned(doc, column_5, 130, fmt.Sprintf("CHF %.2f", tableData.Versicherung), gopdf.Right)

	var mitgliederBeitrag string

	if tableData.Mitgliederbeitrag == 0 {
		mitgliederBeitrag = "CHF -"
	} else {
		mitgliederBeitrag = fmt.Sprintf("CHF %.2f", tableData.Mitgliederbeitrag)
	}

	setTextAligned(doc, column_1, 136, "1", gopdf.Right)
	setTextAligned(doc, column_2, 136, tranlate(language, invoiceDetails.TabelleJahre), gopdf.Left)
	setTextAligned(doc, column_3, 136, tranlate(language, variableData.TextMitgliederbeitrag), gopdf.Left)
	setTextAligned(doc, column_4, 136, fmt.Sprintf("CHF %.2f", variableData.Mitgliederbeitrag), gopdf.Right)
	setTextAligned(doc, column_5, 136, mitgliederBeitrag, gopdf.Right)

	setTextAligned(doc, column_1, 142, "1", gopdf.Right)
	setTextAligned(doc, column_2, 142, tranlate(language, invoiceDetails.TabelleJahre), gopdf.Left)
	setTextAligned(doc, column_3, 142, tranlate(language, variableData.TextVerwaltungskosten), gopdf.Left)
	setTextAligned(doc, column_4, 142, fmt.Sprintf("CHF %.2f", variableData.Verwaltungskosten), gopdf.Right)
	setTextAligned(doc, column_5, 142, fmt.Sprintf("CHF %.2f", tableData.Verwaltungskosten), gopdf.Right)

	setTextAligned(doc, column_1, 148, "1", gopdf.Right)
	setTextAligned(doc, column_2, 148, tranlate(language, invoiceDetails.TabelleJahre), gopdf.Left)
	setTextAligned(doc, column_3, 148, tranlate(language, variableData.TextReparaturFonds), gopdf.Left)
	setTextAligned(doc, column_4, 148, fmt.Sprintf("CHF %.2f", variableData.Reparaturfonds), gopdf.Right)
	setTextAligned(doc, column_5, 148, fmt.Sprintf("CHF %.2f", tableData.Reparaturfonds), gopdf.Right)

	setTextAligned(doc, column_1, 154, "1", gopdf.Right)
	setTextAligned(doc, column_2, 154, tranlate(language, invoiceDetails.TabelleJahre), gopdf.Left)
	setTextAligned(doc, column_3, 154, tranlate(language, variableData.TextVerwaltungskosten), gopdf.Left)
	setTextAligned(doc, column_4, 154, fmt.Sprintf("CHF %.2f", variableData.Verwaltungskosten), gopdf.Right)
	setTextAligned(doc, column_5, 154, fmt.Sprintf("CHF %.2f", tableData.Verwaltungskosten), gopdf.Right)

	doc.SetFontSize(12)
	doc.SetFontStyle("bold")

	setTextAligned(doc, column_1, 162, "Total", gopdf.Right)
	setTextAligned(doc, column_5, 162, fmt.Sprintf("CHF %.2f", tableData.Total), gopdf.Right)

	doc.DefaultFontSize()
	doc.DefaultFontStyle()
}

func setTextAligned(doc *gopdf_wrapper.Doc, x, y int, text string, alignment int) {

	doc.SetFillColor(50, 50, 50)
	doc.SetPosition(float64(x), float64(y))
	doc.CellWithOption(&gopdf.Rect{W: 30, H: 11}, text, gopdf.CellOption{Align: alignment})

}

func tranlate(language string, text model.TranslatedText) string {
	if language == "fr" {
		return text.Fr
	}
	return text.De
}
