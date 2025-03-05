package main

import (
	"fmt"
	"log"
	"strings"
	"sync"

	swissqrinvoice "github.com/72nd/swiss-qr-invoice"
	"github.com/signintech/gopdf"
	"github.com/toky03/qr-invoice/document"
	"github.com/toky03/qr-invoice/model"
	"github.com/xuri/excelize/v2"
)

var waitGroup sync.WaitGroup

func main() {

	excel, err := excelize.OpenFile("data/Mitgliederliste Aktuell.xlsx")

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

		waitGroup.Add(1)
		calculatedData := variableData.ToCalculatedTableData(debtor)
		invoice := invoiceDetails.ToInvoiceDetails(debtor, calculatedData)
		go createDocument(
			debtor.Parzelle,
			invoice,
			invoiceDetails.ToZusatz(debtor.Language),
			debtor.ToReceiverAdress(),
			invoiceDetails.ToTitle(debtor.Language),
			invoiceDetails.ToTableData(debtor.Language, debtor, variableData, calculatedData),
		)
	}

	waitGroup.Wait()

}

func contains(s []string, e string) bool {
	for _, a := range s {
		if strings.TrimSpace(a) == strings.TrimSpace(e) {
			return true
		}
	}
	return false
}

func createDocument(
	parzelle string,
	invoice swissqrinvoice.Invoice,
	zusatz string,
	receiverAdress document.ReceiverAdress,
	title document.TitleWithDate,
	tableData document.TableData,
) {

	defer waitGroup.Done()
	doc := createDocFromInvoice(invoice)

	document.AddAdressData(doc, receiverAdress)
	document.AddTitle(doc, title)
	document.AddText(doc, zusatz)
	document.AddTable(doc, tableData)

	doc.Image("data/logo_neu.png", 10, 10, &gopdf.Rect{W: 100, H: 33})

	fileName := fmt.Sprintf(
		"rechnung_%03s_%s.pdf",
		parzelle,
		strings.ReplaceAll(strings.ReplaceAll(receiverAdress.Name, " ", "_"), "/", ""),
	)

	if err := doc.WritePdf(fmt.Sprintf("bills/%s/%s", invoice.Language, fileName)); err != nil {
		log.Panic(err)
	}
}

func createDocFromInvoice(invoice swissqrinvoice.Invoice) (doc document.PdfDoc) {
	doc, err := invoice.Doc()
	if err != nil {
		log.Panic(err)
	}
	return doc
}
