package compilers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/eris-ltd/common/go/common"
)

var DefaultUrl = "https://compilers.eris.industries:9090/compile"

// Language configuration struct
// New language capabilities can be added to the server simply by
// providing a LangConfig
// Each element in IncludeReplaces is a pair of strings, between which is placed the filename
// CompileCmd is a list of what would be white-space separated tokens on the
// command line, with a `_` to denote the place of the filename
type LangConfig struct {
	URL             string     `json:"url"`
	Net             bool       `json:"net"`
	Extensions      []string   `json:"extensions"`
	IncludeRegexes  []string   `json:"regexes"`
	IncludeReplaces [][]string `json:"replaces"`
	CompileCmd      []string   `json:"cmd"`
	AbiCmd          []string   `json:"abi"`
}

// Append the language extension to the filename
func (l LangConfig) Ext(h string) string {
	return h + "." + l.Extensions[0]
}

// Fill in the filename and return the command line args
func (l LangConfig) Cmd(file string, includes []string) (args []string) {
	for _, s := range l.CompileCmd {
		if s == "_" {
			args = append(args, file)
		} else if s == "imports" {
			args = append(args, includes...)
		} else {
			args = append(args, s)
		}
	}
	logger.Debugf("Command Compiled =>\t\t%v\n", args)
	return
}

func (l LangConfig) Abi(file string, includes []string) (args []string) {
	if len(l.AbiCmd) < 2 {
		return
	}

	for _, s := range l.AbiCmd {
		if s == "_" {
			args = append(args, file)
		} else if s == "imports" {
			args = append(args, includes...)
		} else {
			args = append(args, s)
		}
	}
	return
}

// Global variable mapping languages to their configs
var Languages = map[string]LangConfig{
	"lll": {
		URL:        DefaultUrl,
		Net:        true,
		Extensions: []string{"lll", "def"},
		IncludeRegexes: []string{
			`\(include "(.+?)"\)`,
		},
		IncludeReplaces: [][]string{
			{`(include "`, `.lll")`},
		},
		CompileCmd: []string{
			"lllc",
			"_",
		},
	},

	"se": {
		URL:        DefaultUrl,
		Net:        true,
		Extensions: []string{"se"},
		IncludeRegexes: []string{
			// because I'm not that good with regex and this
			// demonstrates how to have multiple expressions to match :)
			`create\("(.+?)"\)`,
			`create\('(.+?)'\)`,
		},
		IncludeReplaces: [][]string{
			{`create("`, `.se")`},
			{`create('`, `.se')`},
		},
		CompileCmd: []string{
			"sc",
			"compile",
			"_",
		},
		AbiCmd: []string{
			"sc",
			"mk_full_signature",
			"_",
		},
	},
	"sol": {
		URL:        DefaultUrl,
		Net:        true,
		Extensions: []string{"sol"},
		IncludeRegexes: []string{
			`import "(.+?)";`,
		},
		IncludeReplaces: [][]string{
			{`import "`, `.sol";`},
		},
		CompileCmd: []string{
			"/usr/bin/solc",
			"--bin",
			"_",
			"imports",
			"|",
			"grep",
			"Binary",
			"-A1",
			"|",
			"tail",
			"-n",
			"1",
		},
		AbiCmd: []string{
			"/usr/bin/solc",
			"--abi",
			"_",
			"imports",
			"|",
			"grep",
			"ABI",
			"-A1",
			"|",
			"tail",
			"-n",
			"1",
		},
	},
}

func init() {
	common.InitErisDir()
	common.InitDataDir(ClientCache)
	common.InitDataDir(ServerCache)

	f := path.Join(common.LanguagesPath, "config.json")
	err := checkConfig(f)
	if err != nil {
		logger.Errorln(err)
		logger.Errorln("resorting to default language settings")
		return
	}
}

func checkConfig(f string) error {
	if _, err := os.Stat(f); err != nil {
		err := common.WriteJson(&Languages, f)
		if err != nil {
			return fmt.Errorf("Could not write default config to %s: %s", f, err.Error())
		}
	}

	b, err := ioutil.ReadFile(f)
	if err != nil {
		return err
	}

	c := new(map[string]LangConfig)
	err = json.Unmarshal(b, c)
	if err != nil {
		return err
	}

	Languages = *c
	return nil
}

// Set the languages url
func SetLanguageURL(lang, url string) error {
	l, ok := Languages[lang]
	if !ok {
		return UnknownLang(lang)
	}
	l.URL = url
	Languages[lang] = l
	return nil
}

// Set whether the language should use the remote server or compile locally
func SetLanguageNet(lang string, net bool) error {
	l, ok := Languages[lang]
	if !ok {
		return UnknownLang(lang)
	}
	l.Net = net
	Languages[lang] = l
	return nil

}

// Main client struct to wrap a compiler interface and its configuration data
type CompileClient struct {
	config LangConfig
	lang   string
}

// Return the language name
func (c *CompileClient) Lang() string {
	return c.lang //c.Lang()
}

// Return the language's main extension
func (c *CompileClient) Ext(h string) string {
	return c.config.Ext(h)
}

// Return the regex string to match include statements
func (c *CompileClient) IncludeRegexes() []string {
	return c.config.IncludeRegexes
}

// Return the string to replace matched regex expressions
func (c *CompileClient) IncludeReplace(h string, i int) string {
	s := c.config.IncludeReplaces[i]
	return s[0] + h + s[1]
}

// Unknown language error
func UnknownLang(lang string) error {
	return fmt.Errorf("Unknown language %s", lang)
}

// Create a new compile client
func NewCompileClient(lang string) (*CompileClient, error) {
	l, ok := Languages[lang]
	if !ok {
		return nil, UnknownLang(lang)
	}
	cc := &CompileClient{
		config: l,
		lang:   lang,
	}
	return cc, nil
}
