module github.com/toky03/qr-invoice

go 1.21

require (
	github.com/72nd/swiss-qr-invoice v1.0.2
	github.com/signintech/gopdf v0.20.0
	github.com/xuri/excelize/v2 v2.8.0
)

require (
	github.com/72nd/gopdf-wrapper v0.3.1 // indirect
	github.com/GeertJohan/go.rice v1.0.3 // indirect
	github.com/creasty/defaults v1.7.0 // indirect
	github.com/daaku/go.zipexe v1.0.2 // indirect
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826 // indirect
	github.com/phpdave11/gofpdi v1.0.14-0.20211212211723-1f10f9844311 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/richardlehane/mscfb v1.0.4 // indirect
	github.com/richardlehane/msoleps v1.0.3 // indirect
	github.com/skip2/go-qrcode v0.0.0-20200617195104-da1b6568686e // indirect
	github.com/xuri/efp v0.0.0-20230802181842-ad255f2331ca // indirect
	github.com/xuri/nfp v0.0.0-20230819163627-dc951e3ffe1a // indirect
	golang.org/x/crypto v0.12.0 // indirect
	golang.org/x/net v0.14.0 // indirect
	golang.org/x/text v0.12.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

// replace solange notwendig bis die pullrequests https://github.com/72nd/swiss-qr-invoice/pull/3 und https://github.com/72nd/swiss-qr-invoice/pull/2 gemerged sind
replace github.com/72nd/swiss-qr-invoice v1.0.2 => ./dependencies/swiss-qr-invoice
