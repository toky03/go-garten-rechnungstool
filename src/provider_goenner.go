//go:build goenner

package main

import (
	"log"

	"github.com/toky03/qr-invoice/core"
	goennerprovider "github.com/toky03/qr-invoice/goenner_provider"
)

func initializeDebtorProvider(basePath string) (core.DebtorProvider, func()) {
	goennerDebtorProvider := goennerprovider.CreateGoennerDebtorProvider(
		basePath,
		"Mitgliederliste Aktuell.xlsx",
		"Rechnugsvariabeln.xlsx",
	)
	if goennerDebtorProvider == nil {
		log.Panic("could not create goenner provider")
	}
	return goennerDebtorProvider, goennerDebtorProvider.Close
}
