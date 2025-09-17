package main

import (
	"log"
	"os"
	"strings"
	"sync"

	swissqrinvoice "github.com/72nd/swiss-qr-invoice"
	"github.com/joho/godotenv"
	"github.com/signintech/gopdf"
	"github.com/toky03/qr-invoice/core"
	"github.com/toky03/qr-invoice/document"
	gartenexcelprovider "github.com/toky03/qr-invoice/garten_excel_provider"
)

var waitGroup sync.WaitGroup

func main() {

	godotenv.Load()

	basePath := os.Getenv("BASE_PATH")
	if basePath == "" {
		basePath = "."
	}

	savePath := os.Getenv("SAVE_PATH")
	if savePath == "" {
		savePath = "."
	}

	debtorProvider := gartenexcelprovider.CreateExcelDebtorProvider(basePath, "mitgliederliste.xlsx")

	defer func() {
		debtorProvider.Close()
	}()

	debtors := debtorProvider.All()
	for debtor := range debtors {

		if debtor.Skip() {
			continue
		}
		waitGroup.Add(1)
		go createDocument(
			basePath,
			savePath,
			debtor,
		)
	}

	waitGroup.Wait()

}

func createDocument(basePath, savePath string,
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

	imageData := invoiceDetailsProvider.GetImageData(basePath)
	if imageData.Path != "" {
		doc.Image(imageData.Path, imageData.Xpos, imageData.Ypos, &gopdf.Rect{W: imageData.Width, H: imageData.Height})
	}

	if err := doc.WritePdf(invoiceDetailsProvider.GetSavePath(savePath)); err != nil {
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
