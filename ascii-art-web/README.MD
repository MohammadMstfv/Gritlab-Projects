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

The generated result can be copied, downloaded as a `.txt` attachment, or saved as a `.png` image.

The command-line renderer accepts text and an optional style name. The default style is `standard`:

```sh
go run . "Hello"
go run . "Hello" standard
go run . "Hello" Binary
```

## Implementation details: algorithm

The renderer supports printable ASCII characters from code point 32 through 126. The `standard`, `shadow`, and `thinkertoy` styles load text banners with eight rows per character and a separator row. The `Binary` style directly renders each non-space character as its 8-bit binary value and preserves spaces between words.

Literal `\\n` sequences in the input are converted to line breaks. Each input line is rendered by joining the corresponding glyph rows, then a newline is appended. Unsupported characters and invalid form data are returned as HTTP 400 errors by the web endpoint.
