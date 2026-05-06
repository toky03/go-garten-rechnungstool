package goennerprovider

import (
	"fmt"
	"iter"
	"log"
	"strconv"
	"strings"
	"time"

	swissqrinvoice "github.com/72nd/swiss-qr-invoice"
	"github.com/toky03/qr-invoice/core"
	"github.com/toky03/qr-invoice/document"
	gartenexcelprovider "github.com/toky03/qr-invoice/garten_excel_provider"
	"github.com/xuri/excelize/v2"
)

type GoennerProvider struct {
	memberExcelFile      *excelize.File
	invoiceDataExcelFile *excelize.File
	invoiceDetails       gartenexcelprovider.InvoiceDetails
}

type goennerDebtorData struct {
	Name     string
	Address  string
	Zip      string
	City     string
	Language string
}

type goennerInvoiceDetailsProvider struct {
	Invoice         swissqrinvoice.Invoice
	MultilineText   string
	ReceiverAddress document.ReceiverAdress
	Title           document.TitleWithDate
	TableData       document.TableData
	SavePath        string
	ImageData       document.ImageData
}

func CreateGoennerDebtorProvider(
	basePath, memberExcelFileName, invoiceDataExcelFileName string,
) *GoennerProvider {
	memberExcelFile, err := excelize.OpenFile(fmt.Sprintf("%s/%s", basePath, memberExcelFileName))
	if err != nil {
		log.Printf("could not read excel file %s", err)
		return nil
	}

	invoiceDataExcelFile, err := excelize.OpenFile(
		fmt.Sprintf("%s/%s", basePath, invoiceDataExcelFileName),
	)
	if err != nil {
		log.Printf("could not read excel file %s", err)
		return nil
	}

	invoiceDetails, err := gartenexcelprovider.ReadInvoiceDetails(invoiceDataExcelFile)
	if err != nil {
		log.Printf("could not read invoice details %s", err)
		return nil
	}

	return &GoennerProvider{
		memberExcelFile:      memberExcelFile,
		invoiceDataExcelFile: invoiceDataExcelFile,
		invoiceDetails:       invoiceDetails,
	}
}

func (p *GoennerProvider) All() iter.Seq[core.InvoiceDetailsProvider] {
	debtors := readGoennerDebtorData(p.memberExcelFile)

	return func(yield func(core.InvoiceDetailsProvider) bool) {
		for _, debtor := range debtors {
			invoice := p.createInvoice(debtor)

			fileName := fmt.Sprintf(
				"rechnung_goenner_%s.pdf",
				strings.ReplaceAll(strings.ReplaceAll(debtor.Name, " ", "_"), "/", ""),
			)
			filePath := fmt.Sprintf("%s/%s", debtor.Language, fileName)

			details := goennerInvoiceDetailsProvider{
				Invoice:       invoice,
				MultilineText: p.invoiceDetails.ToZusatz(debtor.Language),
				ReceiverAddress: document.ReceiverAdress{
					Name:   debtor.Name,
					Adress: debtor.Address,
					City:   fmt.Sprintf("%s %s", debtor.Zip, debtor.City),
					Header: "",
				},
				Title:     document.TitleWithDate{},
				TableData: document.TableData{},
				SavePath:  filePath,
				ImageData: document.ImageData{
					Path:   "",
					Xpos:   0,
					Ypos:   0,
					Width:  0,
					Height: 0,
				},
			}

			if !yield(details) {
				return
			}
		}
	}
}

func (p *GoennerProvider) Close() {
	if err := p.memberExcelFile.Close(); err != nil {
		log.Printf("could not close member excel file %s", err)
	}
	if err := p.invoiceDataExcelFile.Close(); err != nil {
		log.Printf("could not close invoice data excel file %s", err)
	}
}

func (p *GoennerProvider) createInvoice(debtor goennerDebtorData) swissqrinvoice.Invoice {
	reference := fmt.Sprintf("%04d", time.Now().Year())
	refNumber, err := strconv.Atoi(reference)
	if err != nil {
		return swissqrinvoice.Invoice{}
	}
	checkSum, err := calculateCheckSum(refNumber)
	if err != nil {
		return swissqrinvoice.Invoice{}
	}
	referenceWithPruefZiffer := fmt.Sprintf("RF%02d%021d", checkSum, refNumber)

	return swissqrinvoice.Invoice{
		ReceiverIBAN:    p.invoiceDetails.Creditor.Account,
		IsQrIBAN:        false,
		ReceiverName:    p.invoiceDetails.Creditor.Name,
		ReceiverStreet:  p.invoiceDetails.Creditor.Address,
		ReceiverZIPCode: p.invoiceDetails.Creditor.Zip,
		ReceiverPlace:   p.invoiceDetails.Creditor.City,
		ReceiverCountry: p.invoiceDetails.Creditor.Country,
		PayeeName:       debtor.Name,
		PayeeStreet:     debtor.Address,
		PayeeZIPCode:    debtor.Zip,
		PayeePlace:      debtor.City,
		PayeeCountry:    "CH",
		Currency:        "CHF",
		Amount:          fmt.Sprintf("%.2f", 40.0),
		AdditionalInfo:  fmt.Sprintf("Gönnerbeitrag %s", debtor.Name),
		Reference:       referenceWithPruefZiffer,
		Language:        debtor.Language,
	}
}

func readGoennerDebtorData(workbook *excelize.File) []goennerDebtorData {
	debtors := make([]goennerDebtorData, 0)

	rows, err := workbook.GetRows("Gönner 2026")
	if err != nil {
		log.Printf("could not read Sheet %s", err)
		return debtors
	}

	for _, row := range rows {
		if len(row) < 6 {
			continue
		}

		zip := row[4]
		if strings.TrimSpace(zip) == "" {
			continue
		}

		lastName := strings.TrimSpace(row[1])
		firstName := strings.TrimSpace(row[2])
		address := strings.TrimSpace(row[3])
		city := strings.TrimSpace(row[5])

		language := "de"
		if len(row) > 10 && strings.TrimSpace(row[10]) == "F" {
			language = "fr"
		}

		debtors = append(debtors, goennerDebtorData{
			Name:     fmt.Sprintf("%s %s", firstName, lastName),
			Address:  address,
			Zip:      zip,
			City:     city,
			Language: language,
		})
	}

	return debtors
}

func calculateCheckSum(reference int) (int, error) {
	prefixedReferenceNumber := fmt.Sprintf("%d271500", reference)
	referenceNumber, err := strconv.Atoi(prefixedReferenceNumber)
	if err != nil {
		return 0, fmt.Errorf("Could not convert reference Number to int %s", err)
	}
	reminder := referenceNumber % 97
	return 98 - reminder, nil
}

func (d goennerInvoiceDetailsProvider) GetInvoice() swissqrinvoice.Invoice {
	return d.Invoice
}

func (d goennerInvoiceDetailsProvider) GetMultilineText() string {
	return d.MultilineText
}

func (d goennerInvoiceDetailsProvider) GetReceiverAddress() document.ReceiverAdress {
	return d.ReceiverAddress
}

func (d goennerInvoiceDetailsProvider) GetTitle() document.TitleWithDate {
	return d.Title
}

func (d goennerInvoiceDetailsProvider) GetTableData() document.TableData {
	return d.TableData
}

func (d goennerInvoiceDetailsProvider) GetSavePath(basePath string) string {
	return fmt.Sprintf("%s/%s", basePath, d.SavePath)
}

func (d goennerInvoiceDetailsProvider) Skip() bool {
	return false
}

func (d goennerInvoiceDetailsProvider) GetImageData(basePath string) document.ImageData {
	if d.ImageData.Path == "" {
		return document.ImageData{}
	}

	return document.ImageData{
		Path:   fmt.Sprintf("%s/%s", basePath, d.ImageData.Path),
		Xpos:   d.ImageData.Xpos,
		Ypos:   d.ImageData.Ypos,
		Width:  d.ImageData.Width,
		Height: d.ImageData.Height,
	}
}
