package model

import (
	"fmt"
	"math"

	swissqrinvoice "github.com/72nd/swiss-qr-invoice"
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

func roundHalf(number float32) float32 {
	return float32(math.Round(float64(number)*20) / 20)
}
