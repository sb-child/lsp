module main

go 1.16

replace mods/modio => ./pkg/mods/modio

replace mods/mtools => ./pkg/mods/mtools

replace mods/miya => ./pkg/mods/miya

replace mods/yysp => ./pkg/mods/yysp

require (
	github.com/gookit/color v1.4.2
	golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c // indirect
	mods/miya v0.0.0-00010101000000-000000000000
	mods/modio v0.0.0-00010101000000-000000000000
	mods/mtools v0.0.0-00010101000000-000000000000 // indirect
	mods/yysp v0.0.0-00010101000000-000000000000
)
