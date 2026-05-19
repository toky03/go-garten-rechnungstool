//go:build individual_invoice_provider

package main

import (
	"github.com/toky03/qr-invoice/core"
	individualinvoiceprovider "github.com/toky03/qr-invoice/individual_invoice_provider"
)

func init() {
	debtorProviderFactory = createIndividualInvoiceDebtorProvider
}

func createIndividualInvoiceDebtorProvider(basePath string) (core.DebtorProvider, func()) {
	provider := individualinvoiceprovider.CreateIndividualInvoiceDebtorProvider(
		basePath,
		"Mitgliederliste Aktuell.xlsx",
		"Rechnungsvariabeln.xlsx",
	)
	if provider == nil {
		return nil, nil
	}
	return provider, provider.Close
}
