# i18n-tool

Translation management tool. Framework-agnotics, manage using YAML or Markdown files. This tool is written in Go.

## Install

```sh
go install github.com/willnode/i18n-tool/v2@main
```

## How it works

You have this directory:

```
+ translation
|-+ en
  |-  ui.yaml
```

Run `i18n-tool sync --copy` to add missing files and keys from primary language `en` to the rest of languages.

```
+ translation
|-+ en
| |-  ui.yaml
|-+ es
  |-  ui.yaml
```

Language keys must be in format `en` or `en-US`. The value written is copied verbatim from primary languages.

## Commands

All commands must run in root of translation files.

#### `i18n-tool sync`

Add and copy values from primary languages. Use --copy to copy values too.

#### `i18n-tool detect`

List empty or missing or untranslated keys.

#### `i18n-tool export`

Export to different formats. Options:

+ `--format [format]` Format includes yaml,json
+ `--path [path]` Only export with prefix data path, can be set multiple
+ `--output [path]` Export to specific folder. Default `./build`

