package main

import (
	"log"
	"strings"
	"sync"

	swissqrinvoice "github.com/72nd/swiss-qr-invoice"
	"github.com/signintech/gopdf"
	"github.com/toky03/qr-invoice/core"
	"github.com/toky03/qr-invoice/document"
	gartenexcelprovider "github.com/toky03/qr-invoice/garten_excel_provider"
	"github.com/xuri/excelize/v2"
)

var waitGroup sync.WaitGroup

func main() {

	var debtorProvider core.DebtorProvider

	debtorProvider = gartenexcelprovider.CreateExcelDebtorProvider("../example_data/mitgliederliste.xlsx")

	excel, err := excelize.OpenFile("../example_data/mitgliederliste.xlsx")

	if err != nil {
		log.Printf("could not read excel file %s", err)
	}

	defer func() {
		if err := excel.Close(); err != nil {
			log.Printf("could not close excel file %s", err)
		}
	}()

	debtors := debtorProvider.All()
	for debtor := range debtors {

		if debtor.Skip() {
			continue
		}
		waitGroup.Add(1)
		go createDocument(
			debtor,
		)
	}

	waitGroup.Wait()

}

func createDocument(
	invoiceDetailsProvider core.InvoiceDetailsProvider,
) {

	defer waitGroup.Done()
	doc := createDocFromInvoice(invoiceDetailsProvider.GetInvoice())

	document.AddAdressData(doc, invoiceDetailsProvider.GetReceiverAddress())

	addIfNotEmpty(invoiceDetailsProvider.GetTitle().Title, func() {
		document.AddTitle(doc, invoiceDetailsProvider.GetTitle())
	})
	addIfNotEmpty(invoiceDetailsProvider.GetMultilineText(), func() {
		document.AddText(doc, invoiceDetailsProvider.GetMultilineText())
	})

	tableData := invoiceDetailsProvider.GetTableData()
	if len(tableData.Columns) > 0 {
		document.AddTable(doc, invoiceDetailsProvider.GetTableData())
	}

	imageData := invoiceDetailsProvider.GetImageData()
	if imageData.Path != "" {
		doc.Image(imageData.Path, imageData.Xpos, imageData.Ypos, &gopdf.Rect{W: imageData.Width, H: imageData.Height})
	}

	if err := doc.WritePdf(invoiceDetailsProvider.GetSavePath()); err != nil {
		log.Panic(err)
	}
}

func addIfNotEmpty(text string, addFunc func()) {
	if strings.TrimSpace(text) != "" {
		addFunc()
	}
}

func createDocFromInvoice(invoice swissqrinvoice.Invoice) (doc document.PdfDoc) {
	doc, err := invoice.Doc()
	if err != nil {
		log.Panic(err)
	}
	return doc
}
