package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	glspserver "github.com/tliron/glsp/server"
	"gopkg.in/ini.v1"
)

const version = "0.0.1"

func main() {
	cfgFile := flag.String("config", "~/.config/execlsp.ini", "Path to main config file")
	presets := flag.String("presets", "", "List of presets to load")
	flag.Parse()

	var err error
	allCommands := make(map[string]map[string]string)
	if allCommands, err = loadConfigIfExists(*cfgFile, allCommands); err != nil {
		log.Fatalf("Error loading main config file: %v", err)
	}
	_ = cfgFile

	if allCommands, err = loadConfigIfExists("./.execlsp.ini", allCommands); err != nil {
		log.Fatalf("Error loading local config file: %v", err)
	}

	commands := make(map[string]string)
	if _, ok := allCommands["DEFAULT"]; ok {
		for k, v := range allCommands["DEFAULT"] {
			commands[k] = v
		}
	}

	if len(*presets) > 0 {
		for _, preset := range strings.Split(*presets, ",") {
			if _, ok := allCommands[preset]; ok {
				for k, v := range allCommands[preset] {
					commands[fmt.Sprintf("%s.%s", preset, k)] = v
				}
			} else {
				fmt.Printf("WARNING: No commands defined for preset '%s'", preset)
			}
		}
	}

	cmdNames := make([]string, 0, len(commands))
	for name := range commands {
		cmdNames = append(cmdNames, name)
	}

	handler := protocol.Handler{}

	srv := glspserver.NewServer(&handler, "execlsp", false)

	handler.Initialize = func(context *glsp.Context, params *protocol.InitializeParams) (any, error) {
		capabilities := handler.CreateServerCapabilities()
		capabilities.ExecuteCommandProvider = &protocol.ExecuteCommandOptions{
			Commands: cmdNames,
		}

		v := version
		return protocol.InitializeResult{
			Capabilities: capabilities,
			ServerInfo: &protocol.InitializeResultServerInfo{
				Name:    "Exec LSP",
				Version: &v,
			},
		}, nil
	}

	handler.Initialized = func(context *glsp.Context, params *protocol.InitializedParams) error { return nil }

	handler.Shutdown = func(context *glsp.Context) error { return nil }

	handler.TextDocumentDidOpen = func(context *glsp.Context, params *protocol.DidOpenTextDocumentParams) error { return nil }

	handler.WorkspaceExecuteCommand = func(context *glsp.Context, params *protocol.ExecuteCommandParams) (any, error) {
		if _, ok := commands[params.Command]; !ok {
			return nil, fmt.Errorf("command '%s' is not defined", params.Command)
		}

		return execCommand(commands[params.Command])
	}

	if err = srv.RunStdio(); err != nil {
		fmt.Printf("Error running: %v", err)
	}
}

func execCommand(cmd string) (string, error) {
	c := exec.Command(os.Getenv("SHELL"), "-c", cmd)
	o, err := c.CombinedOutput()
	if err != nil {
		return string("OUTPUUUT"), fmt.Errorf("Error executing command: %w", err)
	}

	return string(o), nil
}

func loadConfigIfExists(path string, cfg map[string]map[string]string) (map[string]map[string]string, error) {
	if strings.HasPrefix(path, "~/") {
		homedir, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}

		path = filepath.Join(homedir, path[2:])
	}

	var err error
	path, err = filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	if _, err = os.Stat(path); os.IsNotExist(err) {
		return cfg, nil
	}

	iniData, err := ini.Load(path)
	if err != nil {
		return nil, err
	}

	for _, sectionName := range iniData.SectionStrings() {
		if _, ok := cfg[sectionName]; !ok {
			cfg[sectionName] = make(map[string]string)
		}
		section := iniData.Section(sectionName)
		for _, key := range section.KeyStrings() {
			cfg[sectionName][key] = section.Key(key).String()
		}
	}

	return cfg, nil
}
