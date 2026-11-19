# go-reloaded

Reads a text file, applies the modification markers/rules below, and writes
the result to a second file.

## Usage

```
go run . <input file> <output file>
```

## Supported rules

- `(hex)` / `(bin)` — convert the preceding hex/binary word to decimal.
- `(up)` / `(low)` / `(cap)` — upper/lower/capitalize the preceding word.
- `(up,N)` / `(low,N)` / `(cap,N)` — apply the case change to the preceding
  N words instead of just one.
- `. , ! ? : ;` — normalized so there's no space before and exactly one
  space after, with runs like `...` or `!?` kept glued together.
- `' ... '` — the space just inside each quote pair is trimmed.
- `a` → `an` when the next word starts with a vowel or `h` (case of the
  article is preserved).

## Tests

```
go test ./...
```
