package model

import (
	"fmt"
	"math"

	swissqrinvoice "github.com/72nd/swiss-qr-invoice"
	"github.com/signintech/gopdf"
	"github.com/toky03/qr-invoice/document"
)

func (v VariableData) ToCalculatedTableData(debtor DebtorData) CalculatedData {

	pachtzins := v.Pachtzins * debtor.Are
	wasserbezug := v.Wasserbezug * debtor.Are
	gfAbonement := v.GfAbonement
	strom := v.Strom
	versicherung := v.Versicherung
	var mitgliederbeitrag float32
	if debtor.IsVorstand {
		mitgliederbeitrag = 0
	} else {
		mitgliederbeitrag = v.Mitgliederbeitrag
	}
	reparaturfonds := v.Reparaturfonds
	verwaltungskosten := v.Verwaltungskosten

	total := pachtzins + wasserbezug + gfAbonement + strom + versicherung + mitgliederbeitrag + reparaturfonds + verwaltungskosten

	return CalculatedData{
		Pachtzins:         roundHalf(pachtzins),
		Wasserbezug:       roundHalf(wasserbezug),
		GfAbonement:       gfAbonement,
		Strom:             strom,
		Versicherung:      versicherung,
		Mitgliederbeitrag: mitgliederbeitrag,
		Reparaturfonds:    reparaturfonds,
		Verwaltungskosten: verwaltungskosten,
		Are:               debtor.Are,
		Total:             roundHalf(total),
	}

}

func (i InvoiceDetails) ToInvoiceDetails(debtor DebtorData, calculatedData CalculatedData) swissqrinvoice.Invoice {
	return swissqrinvoice.Invoice{
		ReceiverIBAN:    i.Creditor.Account,
		IsQrIBAN:        false,
		ReceiverName:    i.Creditor.Name,
		ReceiverStreet:  i.Creditor.Address,
		ReceiverZIPCode: i.Creditor.Zip,
		ReceiverPlace:   i.Creditor.City,
		ReceiverCountry: i.Creditor.Country,
		PayeeName:       debtor.Debtor.Name,
		PayeeStreet:     debtor.Debtor.Address,
		PayeeZIPCode:    debtor.Debtor.Zip,
		PayeePlace:      debtor.Debtor.City,
		PayeeCountry:    debtor.Debtor.Country,
		Currency:        "CHF",
		Amount:          fmt.Sprintf("%.2f", calculatedData.Total),
		AdditionalInfo:  fmt.Sprintf("Parzelle %s", debtor.Parzelle),
		Language:        debtor.Language,
	}
}

func (debtor DebtorData) ToReceiverAdress() document.ReceiverAdress {
	var heading string
	if debtor.Language == "de" {
		heading = fmt.Sprintf("Parzelle %s", debtor.Parzelle)
	} else {
		heading = fmt.Sprintf("Parcelle %s", debtor.Parzelle)
	}
	return document.ReceiverAdress{
		Header: heading,
		Name:   debtor.Debtor.Name,
		Adress: debtor.Debtor.Address,
		City:   fmt.Sprintf("%s %s", debtor.Debtor.Zip, debtor.Debtor.City),
	}
}

func (invoiceDetails InvoiceDetails) ToTitle(language string) document.TitleWithDate {
	var title string
	if language == "de" {
		title = invoiceDetails.Ueberschrift.De
	} else {
		title = invoiceDetails.Ueberschrift.Fr
	}

	return document.TitleWithDate{
		Title: title,
		City:  invoiceDetails.Creditor.City,
	}

}

func (invoiceDetails InvoiceDetails) ToTableData(language string, debtorData DebtorData, variableData VariableData, calculatedData CalculatedData) document.TableData {

	anzahlColumn := document.TableColumn{
		Header:    tranlate(language, invoiceDetails.TabelleAnzahl),
		Alignment: gopdf.Left,
		Width:     35,
		Rows: []string{
			fmt.Sprintf("%.1f", debtorData.Are), fmt.Sprintf("%.1f", debtorData.Are), "1", "1", "1", "1", "1", "1", "Total"},
	}

	einheitColumn := document.TableColumn{
		Header:    tranlate(language, invoiceDetails.TabelleEinheit),
		Alignment: gopdf.Left,
		Width:     30,
		Rows: []string{
			tranlate(language, invoiceDetails.TabelleAaren),
			tranlate(language, invoiceDetails.TabelleAaren),
			tranlate(language, invoiceDetails.TabelleJahre),
			tranlate(language, invoiceDetails.TabelleJahre),
			tranlate(language, invoiceDetails.TabelleJahre),
			tranlate(language, invoiceDetails.TabelleJahre),
			tranlate(language, invoiceDetails.TabelleJahre),
			tranlate(language, invoiceDetails.TabelleJahre), ""},
	}

	bezeichungColumn := document.TableColumn{
		Header:    tranlate(language, invoiceDetails.TabelleBezeichnung),
		Alignment: gopdf.Left,
		Width:     45,
		Rows: []string{
			tranlate(language, variableData.TextPachtzins),
			tranlate(language, variableData.TextWasserbezug),
			tranlate(language, variableData.TextGfAbonement),
			tranlate(language, variableData.TextStrom),
			tranlate(language, variableData.TextVersicherung),
			tranlate(language, variableData.TextMitgliederbeitrag),
			tranlate(language, variableData.TextReparaturFonds),
			tranlate(language, variableData.TextVerwaltungskosten), ""},
	}

	preisColumn := document.TableColumn{
		Header:    tranlate(language, invoiceDetails.TabellePreis),
		Alignment: gopdf.Right,
		Width:     30,
		Rows: []string{
			fmt.Sprintf("CHF %.2f", variableData.Pachtzins),
			fmt.Sprintf("CHF %.2f", variableData.Wasserbezug),
			fmt.Sprintf("CHF %.2f", variableData.GfAbonement),
			fmt.Sprintf("CHF %.2f", variableData.Strom),
			fmt.Sprintf("CHF %.2f", variableData.Versicherung),
			fmt.Sprintf("CHF %.2f", variableData.Mitgliederbeitrag),
			fmt.Sprintf("CHF %.2f", variableData.Reparaturfonds),
			fmt.Sprintf("CHF %.2f", variableData.Verwaltungskosten), ""},
	}

	var mitgliederBeitrag string

	if calculatedData.Mitgliederbeitrag == 0 {
		mitgliederBeitrag = "CHF     -"
	} else {
		mitgliederBeitrag = fmt.Sprintf("CHF %.2f", calculatedData.Mitgliederbeitrag)
	}

	betragColumn := document.TableColumn{
		Header:    tranlate(language, invoiceDetails.TabelleBetrag),
		Alignment: gopdf.Right,
		Width:     30,
		Rows: []string{
			fmt.Sprintf("CHF %.2f", calculatedData.Pachtzins),
			fmt.Sprintf("CHF %.2f", calculatedData.Wasserbezug),
			fmt.Sprintf("CHF %.2f", calculatedData.GfAbonement),
			fmt.Sprintf("CHF %.2f", calculatedData.Strom),
			fmt.Sprintf("CHF %.2f", calculatedData.Versicherung),
			mitgliederBeitrag,
			fmt.Sprintf("CHF %.2f", calculatedData.Reparaturfonds),
			fmt.Sprintf("CHF %.2f", calculatedData.Verwaltungskosten),
			fmt.Sprintf("CHF %.2f", calculatedData.Total),
		},
	}

	return document.TableData{
		Columns:     []document.TableColumn{anzahlColumn, einheitColumn, bezeichungColumn, preisColumn, betragColumn},
		LastRowBold: true,
	}

}

func roundHalf(number float32) float32 {
	return float32(math.Round(float64(number)*20) / 20)
}

func tranlate(language string, text TranslatedText) string {
	if language == "fr" {
		return text.Fr
	}
	return text.De
}
