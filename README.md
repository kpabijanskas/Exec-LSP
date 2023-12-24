# Exec-LSP

Exec LSP is a very simple LSP server to execute commands.

## Build
```sh
CGO_ENABLED=0 go build -o execlsp main.go
```

## Usage
To load default commands only
```sh
./execlsp
```


To load default commands AND presets (comma-separated and defined in your config)
```sh
 ./execlsp -presets go,git
```

## Config
This LSP reads and merges configurations from two places:
- ~/.config/execlsp.ini (global configuration that is always loaded - this is also configurable via `-config` flag)
- ./.execlsp.ini (local per-project configuration - duplicates overwrite the global configuration)

Config file example:
```ini
# Commands outside of a secion always get loaded
pull = zellij run -f -- git pull

# Commands inside a section are called 'presets' and only get loaded when requested
[go]
test = zellij run -f -- go test ./...
```

Key is the command name exposed via LSP, and value is the command to be run inside your `$SHELL`.

## How it works
All commands get passed to `$SHELL` executable via `$SHELL -c "COMMAND"`.

It's not meant to run full scripts, but rather trigger short commands. If you want to trigger something more complex, just wrap it in a shell script.

# Motivation
I use [Helix](https://helix-editor.com/) as my text editor, running inside [Zellij](https://zellij.dev/), both for development and for general note-taking. I love Helix's editing model, but I generally have two issues with Helix itself (I am using 23.10 release at the time of writing htis):
- You cannot define per-language keybindings
- LSP implementation does not support sending arguments to workspace commands

This is important, as I generally need to run predefined commands periodically while developing. Say I am developing in Go and need to run `go test` or `go run` on demand. I could just bind it to a keybinding, but then I will have it as a keybinding whenever I am developing in other languages. If I need to run equivalent commands in other languages, I will end up with clashes (or a needlesly deep, per-language keybinding tree with a subtree for each language, and years of vim has taught me to try to minimise the amount of keystrokes I make).

For the second issue, some commands may need to be run with specific arguments. For example, I use [zk](https://github.com/mickael-menu/zk/) for my note taking, but the LSP commands it exposes requires some arguments (usually path to the notes folder as the first argument). Since LSPs expose the list of commands supported, you can't just join the arguments with the command and split it up inside LSP - Helix won't let that command even reach the LSP. So for such commands to work, you need to expose the full command, including all of its arguments, from the LSP directly. Also, not all LSPs even let you define custom commands, at least without forking the code.

Most commands generally have the same idea behind them, in various projects and languages (for example, a `build`, `run test suite`, etc), so I really want to be able to define a general set of keybindings and just have them automaticlly be translated to the appropriate command for the project. At the same time, since I am in Zellij, I really like its floating windows. Usually I want the output to just be displayed in a floating window, and then the floating window closed after I view it.

So I wrote this LSP for myself as a way to get around these Helix limitations. It simply allows me to extend Helix with a per-language set of keybindings (or LSP commands, speaking more generally) in the way it is suitable for me.
