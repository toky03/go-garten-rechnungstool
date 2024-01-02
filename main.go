package main

import (
	"fmt"
	"log"

	"github.com/toky03/qr-invoice/document"
	"github.com/toky03/qr-invoice/model"
	"github.com/xuri/excelize/v2"
)

func main() {

	excel, err := excelize.OpenFile("data/mitgliederliste.xlsx")

	if err != nil {
		log.Printf("could not read excel file %s", err)
	}

	defer func() {
		if err := excel.Close(); err != nil {
			log.Printf("could not close excel file %s", err)
		}
	}()

	debtors := model.ReadDebtorData(excel)
	invoiceDetails, err := model.ReadInvoiceDetails(excel)
	variableData, err := model.ReadVariableData(excel)

	for _, debtor := range debtors {
		calculatedData := variableData.ToCalculatedTableData(debtor)
		invoice := invoiceDetails.ToInvoiceDetails(debtor, calculatedData)

		doc, err := invoice.Doc()
		if err != nil {
			log.Panic(err)
		}

		document.AddAdressData(doc, debtor)
		document.AddTitle(doc, debtor.Language, invoiceDetails)
		document.AddTableHeader(doc, debtor.Language, invoiceDetails)
		document.AddTableData(doc, debtor.Language, debtor, invoiceDetails, variableData, calculatedData)

		doc.Image("data/logo.png", 10, 10, nil)

		if err := doc.WritePdf(fmt.Sprintf("rechnung_%s.pdf", debtor.Parzelle)); err != nil {
			log.Panic(err)
		}

	}

}
