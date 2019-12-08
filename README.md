goldmark
==========================================

[![http://godoc.org/github.com/enkogu/goldmark](https://godoc.org/github.com/enkogu/goldmark?status.svg)](http://godoc.org/github.com/enkogu/goldmark)
[![https://github.com/enkogu/goldmark/actions?query=workflow:test](https://github.com/enkogu/goldmark/workflows/test/badge.svg?branch=master&event=push)](https://github.com/enkogu/goldmark/actions?query=workflow:test)
[![https://coveralls.io/github/yuin/goldmark](https://coveralls.io/repos/github/yuin/goldmark/badge.svg?branch=master)](https://coveralls.io/github/yuin/goldmark)
[![https://goreportcard.com/report/github.com/enkogu/goldmark](https://goreportcard.com/badge/github.com/enkogu/goldmark)](https://goreportcard.com/report/github.com/enkogu/goldmark)

> A Markdown parser written in Go. Easy to extend, standard compliant, well structured.

goldmark is compliant with CommonMark 0.29.

TODO
-------

*  [ ] Переделать с пяток основных функций для рендеринга (их 21)
*  [ ] Попробовать запустить, то есть на выходе получить список блоков
*  [ ] Доделать остальные функции рендеринга
*  [ ] То, что писалось рендерером в html атрибуты – надо подумать что и как с этим делать, разобрать все кейсы
*  [ ] Опции парсера/рендерера – надо пощелкать и настроить так, как это должно работать
*  [ ] Собрать и упаковать в один пакет вместе с JohannesKaufmann/html-to-markdown, вырезать всё лишнее
*  [ ] Подключить к мидлу
*  [ ] Тесты


Anytype
---------

1. Так как сейчас скелет мидла пока устаканивается, пишу утилитки для copy/paste, но не только для обычного текстового, но и для html. Когда закончу с этой библиотечкой, смартблоки чуть более подустаканятся и можно будет все разом и подключить.
2. Написал простенький playground вокруг https://github.com/JohannesKaufmann/html-to-markdown и потестил. Работает хорошо. Однако этот конвертер работает на уровне текста и замен, там нет никакого AST, переделать в html -> blocks затруднительно. Но можно его output отдавать в парсер markdown 
3. Парсер markdown. Чекал goldmark: он делает AST, потом по нодам ходит и рендерит их в html строку. Пытался придумать где именно врезаться – брать ли голые ноды и с ними что-то придумывать, или модифицировать соответствующие им рендерные функции.

В goldmark в renderer/html есть 22 рендер-функции, которые принимают ноду и io.writer и пишут туда всякие строки.

Нужно сделать так, чтобы туда передавалась структура состояния, а не `w util.BufWriter`, и переделать эти 22 функции из записи html строк в изменение состояния.

Пока что вложенность будем игнорировать, то есть конвертация будет без row/column.
```js
    {
        textBuffer: "..."
        marksBuffer: []marks
    }

```

Что делать с renderAttributes? Один к одному не переделаешь, надо выписать список возможных атрибутов и подумать что делать с каждым. 

```go
func (r *Renderer) RenderAttributes(w util.BufWriter, node ast.Node) {

    for _, attr := range node.Attributes() {
        _, _ = w.WriteString(" ")
        _, _ = w.Write(attr.Name)
        _, _ = w.WriteString(`="`)
        _, _ = w.Write(util.EscapeHTML(attr.Value.([]byte)))
        _ = w.WriteByte('"')
    }
}
```

Что делать со сложно-вложенными структурами? Например:

```html
    <p>  <!-- state.openedBlock = paragraph -->
        Text  <!-- state.textBuffer += Text -->
        <code> <!-- state.closeCurrentBlock 
         state.blocks.push(currentBlock), 
         state.openedBlock = code -->
            fmt.printLn("Hello world") 
            <!-- state.textBuffer += fmt.printLn("Hello world") -->
        </code> <!-- state.closeCurrentBlock, 
        state.blocks.push(currentBlock)-->
    </p> <!-- IGNORE -->
```

Типичная функция для рендеринга paragraph, попробуем ее переделать.

**BEFORE**

```go
func (r *Renderer) renderParagraph(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
    if entering {
        _, _ = w.WriteString("<p>")
    } else {
        _, _ = w.WriteString("</p>\n")
    }
    return ast.WalkContinue, nil
}
```

**AFTER**

```go
func (r *Renderer) renderParagraph(s RenderState, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
    if entering {
        if s.isCurrentBlock { // если был какой-то блок открыт, то закрываем его, так как мы не поддерживаем вложенность (да ее особо и не может быть после ковертации в markdown)
            _, _ = s.closeCurrentBlock();
            s.pushLastBlockToList();
        }

        _, _ = s.openNewBlock("Paragraph"); // Открываем блок

    } else {
        _, _ = s.closeCurrentBlock();
    }
    return ast.WalkContinue, nil
}
```

Есть `RenderState`, у которого примерно такие интерфейс и структура:

```go
type rState interface {}

type renderState struct {
    isCurrentBlock      bool
    blockBuffer         &model.Block
    textBuffer          string
    marksBuffer         *[]model.Block.Content.Text.Mark
    blocksList          *[]model.Block

    closeCurrentBlock   func()
    openNewBlock        func(blockType string)
    pushLastBlockToList func()
}

```

А маркап может быть вложенным, значит нам нужна очередь 


Motivation
----------------------
I need a Markdown parser for Go that meets following conditions:

- Easy to extend.
    - Markdown is poor in document expressions compared with other light markup languages like reStructuredText.
    - We have extensions to the Markdown syntax, e.g. PHP Markdown Extra, GitHub Flavored Markdown.
- Standard compliant.
    - Markdown has many dialects.
    - GitHub Flavored Markdown is widely used and it is based on CommonMark aside from whether CommonMark is good specification or not.
        - CommonMark is too complicated and hard to implement.
- Well structured.
    - AST based, and preserves source position of nodes.
- Written in pure Go.

[golang-commonmark](https://gitlab.com/golang-commonmark/markdown) may be a good choice, but it seems to be a copy of the [markdown-it](https://github.com/markdown-it).

[blackfriday.v2](https://github.com/russross/blackfriday/tree/v2) is a fast and widely used implementation, but it is not CommonMark compliant and cannot be extended from outside of the package since its AST uses not interfaces but structs.

Furthermore, its behavior differs from other implementations in some cases especially of lists.  ([Deep nested lists don't output correctly #329](https://github.com/russross/blackfriday/issues/329), [List block cannot have a second line #244](https://github.com/russross/blackfriday/issues/244), etc).

This behavior sometimes causes problems. If you migrate your Markdown text to blackfriday-based wikis from GitHub, many lists will immediately be broken.

As mentioned above, CommonMark is too complicated and hard to implement, So Markdown parsers based on CommonMark barely exist.

Features
----------------------

- **Standard compliant.**  goldmark gets full compliance with the latest CommonMark spec.
- **Extensible.**  Do you want to add a `@username` mention syntax to Markdown?
  You can easily do it in goldmark. You can add your AST nodes,
  parsers for block level elements, parsers for inline level elements,
  transformers for paragraphs, transformers for whole AST structure, and
  renderers.
- **Performance.**  goldmark performs pretty much equally to cmark,
  the CommonMark reference implementation written in C.
- **Robust.**  goldmark is tested with [go-fuzz](https://github.com/dvyukov/go-fuzz), a fuzz testing tool.
- **Builtin extensions.**  goldmark ships with common extensions like tables, strikethrough,
  task lists, and definition lists.
- **Depends only on standard libraries.**

Installation
----------------------
```bash
$ go get github.com/enkogu/goldmark
```


Usage
----------------------
Import packages:

```
import (
	"bytes"
	"github.com/enkogu/goldmark"
)
```


Convert Markdown documents with the CommonMark compliant mode:

```go
var buf bytes.Buffer
if err := goldmark.Convert(source, &buf); err != nil {
  panic(err)
}
```

With options
------------------------------

```go
var buf bytes.Buffer
if err := goldmark.Convert(source, &buf, parser.WithContext(ctx)); err != nil {
  panic(err)
}
```

| Functional option | Type | Description |
| ----------------- | ---- | ----------- |
| `parser.WithContext` | A parser.Context | Context for the parsing phase. |

Custom parser and renderer
--------------------------
```go
import (
	"bytes"
	"github.com/enkogu/goldmark"
	"github.com/enkogu/goldmark/extension"
	"github.com/enkogu/goldmark/parser"
	"github.com/enkogu/goldmark/renderer/html"
)

md := goldmark.New(
          goldmark.WithExtensions(extension.GFM),
          goldmark.WithParserOptions(
              parser.WithAutoHeadingID(),
          ),
          goldmark.WithRendererOptions(
              html.WithHardWraps(),
              html.WithXHTML(),
          ),
      )
var buf bytes.Buffer
if err := md.Convert(source, &buf); err != nil {
    panic(err)
}
```

Parser and Renderer options
------------------------------

### Parser options

| Functional option | Type | Description |
| ----------------- | ---- | ----------- |
| `parser.WithBlockParsers` | A `util.PrioritizedSlice` whose elements are `parser.BlockParser` | Parsers for parsing block level elements. |
| `parser.WithInlineParsers` | A `util.PrioritizedSlice` whose elements are `parser.InlineParser` | Parsers for parsing inline level elements. |
| `parser.WithParagraphTransformers` | A `util.PrioritizedSlice` whose elements are `parser.ParagraphTransformer` | Transformers for transforming paragraph nodes. |
| `parser.WithAutoHeadingID` | `-` | Enables auto heading ids. |
| `parser.WithAttribute` | `-` | Enables custom attributes. Currently only headings supports attributes. |

### HTML Renderer options

| Functional option | Type | Description |
| ----------------- | ---- | ----------- |
| `html.WithWriter` | `html.Writer` | `html.Writer` for writing contents to an `io.Writer`. |
| `html.WithHardWraps` | `-` | Render new lines as `<br>`.|
| `html.WithXHTML` | `-` | Render as XHTML. |
| `html.WithUnsafe` | `-` | By default, goldmark does not render raw HTMLs and potentially dangerous links. With this option, goldmark renders these contents as it is. |

### Built-in extensions

- `extension.Table`
  - [GitHub Flavored Markdown: Tables](https://github.github.com/gfm/#tables-extension-)
- `extension.Strikethrough`
  - [GitHub Flavored Markdown: Strikethrough](https://github.github.com/gfm/#strikethrough-extension-)
- `extension.Linkify`
  - [GitHub Flavored Markdown: Autolinks](https://github.github.com/gfm/#autolinks-extension-)
- `extension.TaskList`
  - [GitHub Flavored Markdown: Task list items](https://github.github.com/gfm/#task-list-items-extension-)
- `extension.GFM`
  - This extension enables Table, Strikethrough, Linkify and TaskList.
  - This extension does not filter tags defined in [6.11: Disallowed Raw HTML (extension)](https://github.github.com/gfm/#disallowed-raw-html-extension-).
    If you need to filter HTML tags, see [Security](#security)
- `extension.DefinitionList`
  - [PHP Markdown Extra: Definition lists](https://michelf.ca/projects/php-markdown/extra/#def-list)
- `extension.Footnote`
  - [PHP Markdown Extra: Footnotes](https://michelf.ca/projects/php-markdown/extra/#footnotes)
- `extension.Typographer`
  - This extension substitutes punctuations with typographic entities like [smartypants](https://daringfireball.net/projects/smartypants/).

### Attributes
`parser.WithAttribute` option allows you to define attributes on some elements.

Currently only headings support attributes.

**Attributes are being discussed in the
[CommonMark forum](https://talk.commonmark.org/t/consistent-attribute-syntax/272).
This syntax may possibly change in the future.**


#### Headings

```
## heading ## {#id .className attrName=attrValue class="class1 class2"}

## heading {#id .className attrName=attrValue class="class1 class2"}
```

```
heading {#id .className attrName=attrValue}
============
```

### Typographer extension

Typographer extension translates plain ASCII punctuation characters into typographic punctuation HTML entities.

Default substitutions are:

| Punctuation | Default entity |
| ------------ | ---------- |
| `'`           | `&lsquo;`, `&rsquo;` |
| `"`           | `&ldquo;`, `&rdquo;` |
| `--`       | `&ndash;` |
| `---`      | `&mdash;` |
| `...`      | `&hellip;` |
| `<<`       | `&laquo;` |
| `>>`       | `&raquo;` |

You can overwrite the substitutions by `extensions.WithTypographicSubstitutions`.

```go
markdown := goldmark.New(
	goldmark.WithExtensions(
		extension.NewTypographer(
			extension.WithTypographicSubstitutions(extension.TypographicSubstitutions{
				extension.LeftSingleQuote:  []byte("&sbquo;"),
				extension.RightSingleQuote: nil, // nil disables a substitution
			}),
		),
	),
)
```



Create extensions
--------------------
**TODO**

See `extension` directory for examples of extensions.

Summary:

1. Define AST Node as a struct in which `ast.BaseBlock` or `ast.BaseInline` is embedded.
2. Write a parser that implements `parser.BlockParser` or `parser.InlineParser`.
3. Write a renderer that implements `renderer.NodeRenderer`.
4. Define your goldmark extension that implements `goldmark.Extender`.

Security
--------------------
By default, goldmark does not render raw HTMLs and potentially dangerous URLs.
If you need to gain more control over untrusted contents, it is recommended to
use an HTML sanitizer such as [bluemonday](https://github.com/microcosm-cc/bluemonday).

Benchmark
--------------------
You can run this benchmark in the `_benchmark` directory.

### against other golang libraries

blackfriday v2 seems the fastest, but it is not CommonMark compliant, so the performance of
blackfriday v2 cannot simply be compared with that of the other CommonMark compliant libraries.

Though goldmark builds clean extensible AST structure and get full compliance with
CommonMark, it is reasonably fast and has lower memory consumption.

```
goos: darwin
goarch: amd64
BenchmarkMarkdown/Blackfriday-v2-12                  326           3465240 ns/op         3298861 B/op      20047 allocs/op
BenchmarkMarkdown/GoldMark-12                        303           3927494 ns/op         2574809 B/op      13853 allocs/op
BenchmarkMarkdown/CommonMark-12                      244           4900853 ns/op         2753851 B/op      20527 allocs/op
BenchmarkMarkdown/Lute-12                            130           9195245 ns/op         9175030 B/op     123534 allocs/op
BenchmarkMarkdown/GoMarkdown-12                        9         113541994 ns/op         2187472 B/op      22173 allocs/op
```

### against cmark (CommonMark reference implementation written in C)

```
----------- cmark -----------
file: _data.md
iteration: 50
average: 0.0037760639 sec
go run ./goldmark_benchmark.go
------- goldmark -------
file: _data.md
iteration: 50
average: 0.0040964230 sec
```

As you can see, goldmark performs pretty much equally to cmark.

Extensions
--------------------

- [goldmark-meta](https://github.com/enkogu/goldmark-meta): A YAML metadata
  extension for the goldmark Markdown parser.
- [goldmark-highlighting](https://github.com/enkogu/goldmark-highlighting): A Syntax highlighting extension
  for the goldmark markdown parser.
- [goldmark-mathjax](https://github.com/litao91/goldmark-mathjax): Mathjax support for goldmark markdown parser

Donation
--------------------
BTC: 1NEDSyUmo4SMTDP83JJQSWi1MvQUGGNMZB

License
--------------------
MIT

Author
--------------------
Yusuke Inuzuka
