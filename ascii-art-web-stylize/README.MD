# ASCII Art Web

## Description

ASCII Art Web is a Go application that converts user-submitted text into ASCII art. It provides a browser interface and a command-line mode, and supports the `standard`, `shadow`, `thinkertoy`, and `Binary` styles.

## Authors

Ghazal and Mohammad.

## Usage: how to run

From the project root, start the web server:

```sh
go run .
```

Open http://localhost:8080, enter text, choose a style, and use **Generate**. The live preview updates as you type. **See in another page** submits the text and style to `/ascii-art` and opens the plain-text result in a new tab.

The command-line renderer accepts text and an optional style name. The default style is `standard`:

```sh
go run . "Hello"
go run . "Hello" standard
go run . "Hello" Binary
```

## Implementation details: algorithm

The renderer supports printable ASCII characters from code point 32 through 126. It loads the selected font and parses it into glyph rows. The `.txt` banners use eight rows per character with a separator row. `Binary.flf` is parsed using its FIGlet header, including the declared height, comment count, hardblank, endmark, and layout mode.

Literal `\\n` sequences in the input are converted to line breaks. Each input line is rendered by merging the corresponding glyph rows, then a newline is appended. The FIGlet layout value controls whether glyphs are rendered at full width or combined using kerning and smushing. Unsupported characters and invalid form data are returned as HTTP 400 errors by the web endpoint.
