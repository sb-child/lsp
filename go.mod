module main

go 1.16

replace mods/modio => ./pkg/mods/modio

replace mods/mtools => ./pkg/mods/mtools

replace mods/miya => ./pkg/mods/miya

replace mods/yysp => ./pkg/mods/yysp

require (
	github.com/Luzifer/go-openssl/v4 v4.1.0 // indirect
	github.com/antchfx/xpath v1.2.0 // indirect
	github.com/gojek/valkyrie v0.0.0-20190210220504-8f62c1e7ba45 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/gookit/color v1.4.2
	github.com/stretchr/testify v1.7.0 // indirect
	golang.org/x/crypto v0.0.0-20210817164053-32db794688a5 // indirect
	golang.org/x/net v0.0.0-20210716203947-853a461950ff // indirect
	golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c // indirect
	google.golang.org/protobuf v1.27.1 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
	mods/miya v0.0.0-00010101000000-000000000000
	mods/modio v0.0.0-00010101000000-000000000000
	mods/mtools v0.0.0-00010101000000-000000000000
	mods/yysp v0.0.0-00010101000000-000000000000
)
