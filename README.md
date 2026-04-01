# clipcut

A simple command line tool that performs cut operations directly on the content of your system clipboard.

It reads text from the clipboard and extracts fields or characters based on the options, similar to the standard Unix cut command.

## Features

- Select fields using delimiters (-f)
- Select by character positions (-c)
- Custom delimiter support (-d)
- Suppress lines that do not contain the delimiter (-s)
- Handles large clipboard content efficiently
- Cross platform support (Linux, macOS, Windows)

Usage:
clipcut [options]

Options
```
-d string   Delimiter (default is TAB)
-f string   Fields to select (example: 1,3-5,7-)
-c string   Characters to select (example: 1-20,30-)
-s          Suppress lines that do not contain the delimiter
```
Examples
```
ccut -f 1,3
ccut -d, -f 2-4
ccut -d' ' -f 1-
ccut -c 1-50
ccut -d: -s -f 3
ccut -d',' -f 1,4-
```
