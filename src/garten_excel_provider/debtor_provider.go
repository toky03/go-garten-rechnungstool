package gartenexcelprovider

import (
	"fmt"
	"iter"
	"log"
	"strings"

	swissqrinvoice "github.com/72nd/swiss-qr-invoice"
	"github.com/toky03/qr-invoice/core"
	"github.com/toky03/qr-invoice/document"
	"github.com/xuri/excelize/v2"
)

type gartenDebtorProviderImpl struct {
	excelFile      *excelize.File
	variableData   VariableData
	invoiceDetails InvoiceDetails
}

type debtorDataImpl struct {
	Invoice         swissqrinvoice.Invoice
	MultilineText   string
	ReceiverAddress document.ReceiverAdress
	Title           document.TitleWithDate
	TableData       document.TableData
	SavePath        string
	ImageData       document.ImageData
}

func CreateExcelDebtorProvider(filePath string) *gartenDebtorProviderImpl {
	excelFile, err := excelize.OpenFile(filePath)
	if err != nil {
		log.Printf("could not read excel file %s", err)
		return nil
	}

	variableData, err := ReadVariableData(excelFile)
	if err != nil {
		log.Printf("could not read variable data %s", err)
		return nil
	}
	invoiceDetails, err := ReadInvoiceDetails(excelFile)
	if err != nil {
		log.Printf("could not read invoice details %s", err)
		return nil
	}

	return &gartenDebtorProviderImpl{
		excelFile:      excelFile,
		variableData:   variableData,
		invoiceDetails: invoiceDetails,
	}
}

func (p *gartenDebtorProviderImpl) All() iter.Seq[core.InvoiceDetailsProvider] {
	paechter := ReadPaechterData(p.excelFile)

	imageData := document.ImageData{
		Path:   "../data/logo_neu.png",
		Xpos:   10,
		Ypos:   10,
		Width:  100,
		Height: 33,
	}

	return func(yield func(core.InvoiceDetailsProvider) bool) {
		for _, debtor := range paechter {
			fileName := fmt.Sprintf(
				"rechnung_%03s_%s.pdf",
				debtor.Parzelle,
				strings.ReplaceAll(strings.ReplaceAll(debtor.Debtor.Name, " ", "_"), "/", ""),
			)

			filePath := fmt.Sprintf("../bills/%s/%s", debtor.Language, fileName)

			debtorProvider := debtorDataImpl{
				Invoice:         p.invoiceDetails.ToInvoiceDetails(debtor, p.variableData.ToCalculatedTableData(debtor)),
				MultilineText:   p.invoiceDetails.ToZusatz(debtor.Language),
				ReceiverAddress: debtor.ToReceiverAdress(),
				Title:           p.invoiceDetails.ToTitle(debtor.Language),
				TableData:       p.invoiceDetails.ToTableData(debtor.Language, debtor, p.variableData, p.variableData.ToCalculatedTableData(debtor)),
				SavePath:        filePath,
				ImageData:       imageData,
			}
			if !yield(debtorProvider) {
				return
			}
		}
	}

}

func (d debtorDataImpl) GetInvoice() swissqrinvoice.Invoice {
	return d.Invoice
}
func (d debtorDataImpl) GetMultilineText() string {
	return d.MultilineText
}
func (d debtorDataImpl) GetReceiverAddress() document.ReceiverAdress {
	return d.ReceiverAddress
}
func (d debtorDataImpl) GetTitle() document.TitleWithDate {
	return d.Title
}
func (d debtorDataImpl) GetTableData() document.TableData {
	return d.TableData
}
func (d debtorDataImpl) GetSavePath() string {
	return d.SavePath
}
func (d debtorDataImpl) Skip() bool {
	return false
}
func (d debtorDataImpl) GetImageData() document.ImageData {
	return d.ImageData
}
