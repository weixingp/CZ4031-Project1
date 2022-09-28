module CZ4031-Project1

go 1.19

require internal/fs v1.0.0

require internal/bptree v1.0.0

require (
	github.com/grailbio/base v0.0.10 // indirect
	github.com/mattn/go-runewidth v0.0.14 // indirect
	github.com/mitchellh/colorstring v0.0.0-20190213212951-d06e56a500db // indirect
	github.com/rivo/uniseg v0.4.2 // indirect
	github.com/schollz/progressbar/v3 v3.11.0 // indirect
	golang.org/x/sys v0.0.0-20220928140112-f11e5e49a4ec // indirect
	golang.org/x/term v0.0.0-20220919170432-7a66f970e087 // indirect
)

replace internal/fs => ./internal/fs

replace internal/bptree => ./internal/bptree
