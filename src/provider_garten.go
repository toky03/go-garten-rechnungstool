//go:build !goenner

package main

import (
	"log"

	"github.com/toky03/qr-invoice/core"
	gartenexcelprovider "github.com/toky03/qr-invoice/garten_excel_provider"
)

func initializeDebtorProvider(basePath string) (core.DebtorProvider, func()) {
	gartenDebtorProvider := gartenexcelprovider.CreateExcelDebtorProvider(
		basePath,
		"Mitgliederliste Aktuell.xlsx",
		"Rechnugsvariabeln.xlsx",
	)
	if gartenDebtorProvider == nil {
		log.Panic("could not create garten provider")
	}
	return gartenDebtorProvider, gartenDebtorProvider.Close
}
