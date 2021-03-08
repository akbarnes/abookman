# abookman
Bookmarks Exporter/Importer for the Amfora base32 encoded TOML bookmarks format

Currently only export is supported. Supported output formats are

- GemText
- TOML
- YAML
- JSON


# Usage
1. Copy `~/.local/share/amfora/bookmarks.toml` into the same folder as `abookmark.go`
2. To export to yaml, run `go run abookman.go -to yaml >bookmarks.yaml`

Specifiers for other formats are
1. GemText: `gemtext`, `gmi`
2. TOML: `toml`
3. YAML: `yaml`, `yml`
4. JSON: `json`
