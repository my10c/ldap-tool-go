//
// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//

package configurator

import (
	"fmt"
	"os"

	"my10c.ldap/vars"
	"github.com/BurntSushi/toml"
	"github.com/akamensky/argparse"
	"github.com/my10c/packages-go/is"
	"github.com/my10c/packages-go/print"
)

type (
	Config struct {
		ConfigFile    string
		Server        string
		Cmd           string
		Debug         bool
		AuthValues    Auth
		DefaultValues Defaults
		SudoValues    Sudo
		LogValues     LogConfig
		EnvValues     Envs
		GroupValues   Groups
		ServerValues  Server
		RedisValues   Redis
		LockPID       int
	}

	// the entries structure in the toml file
	GroupMap struct {
		Name string
		Gid  int
	}

	Auth struct {
		AllowUsers []string
		AllowMods  []string
	}

	Defaults struct {
		LockFile      string
		Shell         string
		ValidShells   []string
		UserSearch    string
		GroupSearch   string
		GroupName     string
		GroupId       int
		ShadowMin     int
		ShadowMax     int
		ShadowAge     int
		ShadowWarning int
		Wait          int
		PassLenght    int
		PassComplex   bool
		UidStart      int
		GidStart      int
	}

	Sudo struct {
		ExcludeSudo []string
		SudoersBase string
	}

	LogConfig struct {
		LogsDir       string
		LogFile       string
		LogMaxSize    int
		LogMaxBackups int
		LogMaxAge     int
	}

	Envs struct {
		ValidEnvs []string
	}

	Groups struct {
		SpecialGroups []string
		Groups        []string
		GroupsMap     []GroupMap
	}

	Server struct {
		Server      string `toml:"server,omitempty"`
		BaseDN      string `toml:"basedn,omitempty"`
		Admin       string `toml:"admin,omitempty"`
		AdminPass   string `toml:"adminpass,omitempty"`
		UserDN      string `toml:"userdn,omitempty"`
		GroupDN     string `toml:"groupdn,omitempty"`
		EmailDomain string `toml:"emaildomain,omitempty"`
		NoTLS       bool   `toml:"notls,omitempty"`
		ReadWrite   bool   `toml:"readwrite,omitempty"`
	}

	Redis struct {
		Server  string
		Port    int
		Enabled bool
		TmpFile string
	}

	tomlConfig struct {
		Auth      Auth              `toml:"auth"`
		Defaults  Defaults          `toml:"defaults"`
		Sudo      Sudo              `toml:"sudo"`
		LogConfig LogConfig         `toml:"logconfig"`
		Envs      Envs              `toml:"envs"`
		Groups    Groups            `toml:"groups"`
		Servers   map[string]Server `toml:"servers,omitempty"`
		Redis     Redis             `toml:"redis"`
	}
)

var (
	Is    = is.New()
	Print = print.New()
)

// function to initialize the configuration
func Configurator() *Config {
	// the rest of the values will be filled from the given configuration file
	return &Config{
		ConfigFile: "",
		Server:     "",
		Cmd:        "",
	}
}

func (c *Config) InitializeArgs() {
	HelpMessage := fmt.Sprintf("commands: %s, %s, %s, %s",
		Print.MessageYellow("search"),
		Print.MessageYellow("create"),
		Print.MessageYellow("modify"),
		Print.MessageYellow("delete"),
	)

	errored := 0
	allowedValues := []string{"create", "modify", "delete", "search"}
	parser := argparse.NewParser(vars.MyProgname, vars.MyDescription)
	configFile := parser.String("c", "configFile",
		&argparse.Options{
			Required: false,
			Help:     "Path to the configuration file to be use",
			Default:  "/usr/local/etc/ldap-tool/config.ini",
		})

	server := parser.String("s", "server",
		&argparse.Options{
			Required: false,
			Help:     "Server profile name",
			Default:  "default",
		})

	cmd := parser.Selector("C", "command", allowedValues,
		&argparse.Options{
			Required: false,
			Help:     HelpMessage,
			Default:  "search",
		})

	debug := parser.Flag("d", "debug",
		&argparse.Options{
			Required: false,
			Help:     "Enable debug",
			Default:  false,
		})

	showInfo := parser.Flag("i", "info",
		&argparse.Options{
			Required: false,
			Help:     "Show information",
		})

	showVersion := parser.Flag("v", "version",
		&argparse.Options{
			Required: false,
			Help:     "Show version",
		})

	err := parser.Parse(os.Args)
	if err != nil {
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}

	if *showVersion {
		Print.ClearScreen()
		Print.PrintYellow(vars.MyProgname + " version: " + vars.MyVersion + "\n")
		os.Exit(0)
	}

	if *showInfo {
		Print.ClearScreen()
		Print.PrintYellow(vars.MyDescription + "\n")
		Print.PrintCyan(vars.MyInfo)
		os.Exit(0)
	}

	if len(*configFile) == 0 {
		Print.PrintRed("the flag -c/--config is required\n")
		errored = 1
	}

	if len(*server) == 0 {
		Print.PrintRed("the flag -s/--server is required\n")
		errored = 1
	}

	if len(*cmd) == 0 {
		Print.PrintRed("the flag -C/--command is required\n")
		errored = 1
	}

	if errored == 1 {
		Print.PrintRed("Aborting..\n")
		os.Exit(1)
	}

	if _, ok, _ := Is.IsExist(*configFile, "file"); !ok {
		Print.PrintRed("Configuration file " + *configFile + " does not exist\n")
		os.Exit(1)
	}

	c.ConfigFile = *configFile
	c.Server = *server
	c.Cmd = *cmd
	c.Debug = *debug
}

// function to add the values to the Config object from the configuration file
func (config *Config) InitializeConfigs() {
	var configValues tomlConfig
	if _, err := toml.DecodeFile(config.ConfigFile, &configValues); err != nil {
		Print.PrintRed("Error reading the configuration file\n")
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if len(configValues.Servers[config.Server].Server) == 0 ||
		len(configValues.Servers[config.Server].BaseDN) == 0 ||
		len(configValues.Servers[config.Server].Admin) == 0 ||
		len(configValues.Servers[config.Server].AdminPass) == 0 ||
		len(configValues.Servers[config.Server].UserDN) == 0 ||
		len(configValues.Servers[config.Server].GroupDN) == 0 {
		Print.PrintRed("\tError reading the configuration file, some value are missing\n")
		Print.PrintBlue("\tRequired fields: server, baseDN, admin, adminPass, userDN and GroupDN\n")
		Print.PrintBlue("\tOptional fields: emailDomain, noTLS and readWrite\n")
		Print.PrintBlue(fmt.Sprintf("\tMake sure there is configuration for the server %s%s%s\n",
			vars.Red, config.Server, vars.Off))
		Print.PrintBlue(fmt.Sprintf("\tThere should be %s[server.%s]%s%s under the %s[servers]%s section%s\n",
			vars.Cyan, config.Server, vars.Off, vars.Yellow, vars.Cyan, vars.Yellow, vars.Off))
		Print.PrintBlue("\tAborting...\n")
		os.Exit(1)
	}

	if configValues.Servers[config.Server].Admin != "cn=admin," + configValues.Servers[config.Server].BaseDN {
		// hardcoded that password length (min 16) and force complex
		// and set password age settings to safe values
		// this is the case the script uses an user configuration instead one from an admin
		if configValues.Defaults.PassLenght < 16 {
			configValues.Defaults.PassLenght = 16
		}
		if configValues.Defaults.PassComplex == false {
			configValues.Defaults.PassComplex = true
		}
		if configValues.Defaults.ShadowMin < 30 {
			configValues.Defaults.ShadowMin = 30
		}
		if configValues.Defaults.ShadowAge < 60 {
			configValues.Defaults.ShadowAge = 60
		}
		if configValues.Defaults.ShadowMax < 90 {
			configValues.Defaults.ShadowMax = 90
		}
	}

	// from the configuration file
	config.AuthValues = configValues.Auth
	config.DefaultValues = configValues.Defaults
	config.SudoValues = configValues.Sudo
	config.LogValues = configValues.LogConfig
	config.EnvValues = configValues.Envs
	config.GroupValues = configValues.Groups
	config.ServerValues = configValues.Servers[config.Server]
	config.RedisValues = configValues.Redis
}
