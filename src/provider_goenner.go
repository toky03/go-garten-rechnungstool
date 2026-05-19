//go:build goenner && !individual_invoice_provider

package main

import (
	"github.com/toky03/qr-invoice/core"
	goennerprovider "github.com/toky03/qr-invoice/goenner_provider"
)

func init() {
	debtorProviderFactory = createGoennerDebtorProvider
}

func createGoennerDebtorProvider(basePath string) (core.DebtorProvider, func()) {
	provider := goennerprovider.CreateGoennerDebtorProvider(
		basePath,
		"Mitgliederliste Aktuell.xlsx",
		"Rechnungsvariabeln.xlsx",
	)
	if provider == nil {
		return nil, nil
	}
	return provider, provider.Close
}
