package core

import (
	"iter"

	swissqrinvoice "github.com/72nd/swiss-qr-invoice"
	"github.com/toky03/qr-invoice/document"
)

type DebtorProvider interface {
	All() iter.Seq[InvoiceDetailsProvider]
}

type InvoiceDetailsProvider interface {
	GetInvoice() swissqrinvoice.Invoice
	GetMultilineText() string
	GetReceiverAddress() document.ReceiverAdress
	GetTitle() document.TitleWithDate
	GetTableData() document.TableData
	GetSavePath() string
	Skip() bool
	GetImageData() document.ImageData
}
