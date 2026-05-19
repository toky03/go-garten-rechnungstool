package individualinvoiceprovider

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

type IndividualInvoiceProvider struct {
	memberExcelFile      *excelize.File
	invoiceDataExcelFile *excelize.File
	invoiceDetails       gartenexcelprovider.InvoiceDetails
	debtorsByParzelle    map[string]gartenexcelprovider.PaechterData
}

type individualEntry struct {
	Parzelle string
	Betrag   float32
	Betreff  string
}

type individualInvoiceDetailsProvider struct {
	Invoice         swissqrinvoice.Invoice
	MultilineText   string
	ReceiverAddress document.ReceiverAdress
	Title           document.TitleWithDate
	TableData       document.TableData
	SavePath        string
	ImageData       document.ImageData
}

func CreateIndividualInvoiceDebtorProvider(
	basePath, memberExcelFileName, invoiceDataExcelFileName string,
) *IndividualInvoiceProvider {
	memberExcelFile, err := excelize.OpenFile(fmt.Sprintf("%s/%s", basePath, memberExcelFileName))
	if err != nil {
		log.Printf("could not read member excel file %s", err)
		return nil
	}

	invoiceDataExcelFile, err := excelize.OpenFile(
		fmt.Sprintf("%s/%s", basePath, invoiceDataExcelFileName),
	)
	if err != nil {
		log.Printf("could not read invoice data excel file %s", err)
		return nil
	}

	invoiceDetails, err := gartenexcelprovider.ReadInvoiceDetails(invoiceDataExcelFile)
	if err != nil {
		log.Printf("could not read invoice details %s", err)
		return nil
	}

	debtorsByParzelle := make(map[string]gartenexcelprovider.PaechterData)
	for _, debtor := range gartenexcelprovider.ReadPaechterData(memberExcelFile) {
		debtorsByParzelle[strings.TrimSpace(debtor.Parzelle)] = debtor
	}

	return &IndividualInvoiceProvider{
		memberExcelFile:      memberExcelFile,
		invoiceDataExcelFile: invoiceDataExcelFile,
		invoiceDetails:       invoiceDetails,
		debtorsByParzelle:    debtorsByParzelle,
	}
}

func (p *IndividualInvoiceProvider) All() iter.Seq[core.InvoiceDetailsProvider] {
	entries := readIndividualEntries(p.invoiceDataExcelFile)

	return func(yield func(core.InvoiceDetailsProvider) bool) {
		for _, entry := range entries {
			debtor, ok := p.debtorsByParzelle[strings.TrimSpace(entry.Parzelle)]
			if !ok {
				continue
			}

			invoice := p.createInvoice(debtor, entry)
			title := p.invoiceDetails.Ueberschrift.De
			if debtor.Language == "fr" {
				title = p.invoiceDetails.Ueberschrift.Fr
			}

			fileName := fmt.Sprintf(
				"rechnung_%03s_%s.pdf",
				debtor.Parzelle,
				strings.ReplaceAll(strings.ReplaceAll(debtor.Debtor.Name, " ", "_"), "/", ""),
			)

			details := individualInvoiceDetailsProvider{
				Invoice:         invoice,
				MultilineText:   p.invoiceDetails.ToZusatz(debtor.Language),
				ReceiverAddress: debtor.ToReceiverAdress(),
				Title: document.TitleWithDate{
					Title: fmt.Sprintf("%s %s", title, debtor.Parzelle),
					City:  p.invoiceDetails.Creditor.City,
				},
				TableData: document.TableData{},
				SavePath:  fmt.Sprintf("%s/%s", debtor.Language, fileName),
				ImageData: document.ImageData{
					Path:   "logo_neu.png",
					Xpos:   10,
					Ypos:   10,
					Width:  100,
					Height: 33,
				},
			}

			if !yield(details) {
				return
			}
		}
	}
}

func (p *IndividualInvoiceProvider) Close() {
	if err := p.memberExcelFile.Close(); err != nil {
		log.Printf("could not close member excel file %s", err)
	}
	if err := p.invoiceDataExcelFile.Close(); err != nil {
		log.Printf("could not close invoice data excel file %s", err)
	}
}

func (p *IndividualInvoiceProvider) createInvoice(
	debtor gartenexcelprovider.PaechterData,
	entry individualEntry,
) swissqrinvoice.Invoice {
	reference := fmt.Sprintf("%03s05%04d", debtor.Parzelle, time.Now().Year())
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
		PayeeName:       debtor.Debtor.Name,
		PayeeStreet:     debtor.Debtor.Address,
		PayeeZIPCode:    debtor.Debtor.Zip,
		PayeePlace:      debtor.Debtor.City,
		PayeeCountry:    debtor.Debtor.Country,
		Currency:        "CHF",
		Amount:          fmt.Sprintf("%.2f", entry.Betrag),
		AdditionalInfo:  fmt.Sprintf("%s %s", entry.Betreff, debtor.Parzelle),
		Reference:       referenceWithPruefZiffer,
		Language:        debtor.Language,
	}
}

func readIndividualEntries(workbook *excelize.File) []individualEntry {
	entries := make([]individualEntry, 0)

	rows, err := workbook.GetRows("Individuelle Rechnungen")
	if err != nil {
		log.Printf("could not read sheet Individuelle Rechnungen %s", err)
		return entries
	}

	for i, row := range rows {
		if i == 0 {
			continue
		}

		if len(row) < 3 {
			continue
		}

		parzelle := strings.TrimSpace(row[0])
		betragString := strings.TrimSpace(row[1])
		betreff := strings.TrimSpace(row[2])
		if parzelle == "" || betragString == "" {
			continue
		}

		betrag, err := strconv.ParseFloat(strings.ReplaceAll(betragString, "'", ""), 32)
		if err != nil {
			cellName, _ := excelize.CoordinatesToCellName(2, i+1)
			log.Printf(
				"could not convert %s from sheet Individuelle Rechnungen on cell %s to number",
				betragString,
				cellName,
			)
			continue
		}

		entries = append(entries, individualEntry{
			Parzelle: parzelle,
			Betrag:   float32(betrag),
			Betreff:  betreff,
		})
	}

	return entries
}

func calculateCheckSum(reference int) (int, error) {
	prefixedReferenceNumber := fmt.Sprintf("%d271500", reference)
	referenceNumber, err := strconv.Atoi(prefixedReferenceNumber)
	if err != nil {
		return 0, fmt.Errorf("could not convert reference number to int %s", err)
	}
	reminder := referenceNumber % 97
	return 98 - reminder, nil
}

func (d individualInvoiceDetailsProvider) GetInvoice() swissqrinvoice.Invoice {
	return d.Invoice
}

func (d individualInvoiceDetailsProvider) GetMultilineText() string {
	return d.MultilineText
}

func (d individualInvoiceDetailsProvider) GetReceiverAddress() document.ReceiverAdress {
	return d.ReceiverAddress
}

func (d individualInvoiceDetailsProvider) GetTitle() document.TitleWithDate {
	return d.Title
}

func (d individualInvoiceDetailsProvider) GetTableData() document.TableData {
	return d.TableData
}

func (d individualInvoiceDetailsProvider) GetSavePath(basePath string) string {
	return fmt.Sprintf("%s/%s", basePath, d.SavePath)
}

func (d individualInvoiceDetailsProvider) Skip() bool {
	return false
}

func (d individualInvoiceDetailsProvider) GetImageData(basePath string) document.ImageData {
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
