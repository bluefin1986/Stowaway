// +build !windows

package initial

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"

	"Stowaway/admin/printer"
	"Stowaway/utils"

	"github.com/nsf/termbox-go"
)

const (
	NORMAL_ACTIVE = iota
	PROXY_ACTIVE
	NORMAL_PASSIVE
)

type Options struct {
	Mode       uint8
	Secret     string
	Listen     string
	Connect    string
	Proxy      string
	ProxyU     string
	ProxyP     string
	Downstream string
	ConfigFile string
}

var Args *Options

func init() {
	Args = new(Options)

	flag.StringVar(&Args.Secret, "s", "", "Communication secret")
	flag.StringVar(&Args.Listen, "l", "", "Listen port")
	flag.StringVar(&Args.Connect, "c", "", "The node address when you actively connect to it")
	flag.StringVar(&Args.Proxy, "proxy", "", "The socks5 server ip:port you want to use")
	flag.StringVar(&Args.ProxyU, "proxyu", "", "socks5 username")
	flag.StringVar(&Args.ProxyP, "proxyp", "", "socks5 password")
	flag.StringVar(&Args.Downstream, "down", "raw", "")
	flag.StringVar(&Args.ConfigFile, "file", "", "config file, if file not exist, default name [cfg.properties] will be use")
	flag.Usage = newUsage
}

func parseConfigFile(configFile string) *Options {
	content := utils.ReadFile(configFile)
	lines := strings.Split(content, "\n")
	Args = new(Options)
	for _, line := range lines {
		config := strings.Split(line, "=")
		if config[0] == "s" {
			Args.Secret = config[1]
			continue
		} else if config[0] == "l" {
			Args.Listen = config[1]
			continue
		} else if config[0] == "c" {
			Args.Connect = config[1]
			continue
		} else if config[0] == "proxy" {
			Args.Proxy = config[1]
			continue
		} else if config[0] == "proxyu" {
			Args.ProxyU = config[1]
			continue
		} else if config[0] == "proxyp" {
			Args.ProxyP = config[1]
			continue
		} else if config[0] == "down" {
			Args.Downstream = config[1]
			continue
		}
	}
	return Args
}

func newUsage() {
	termbox.Close()

	fmt.Fprintf(os.Stderr, `
Usages:
	>> ./stowaway_admin -l <port> -s [secret]
	>> ./stowaway_admin -c <ip:port> -s [secret] 
	>> ./stowaway_admin -c <ip:port> -s [secret] --proxy <ip:port> --proxyu [username] --proxyp [password]
	>> ./stowaway_admin -c <ip:port> -s [secret] --rhostreuse
	>> ./stowaway_admin -c <ip:port> -s [secret] --proxy <ip:port> --proxyu [username] --proxyp [password] --rhostreuse

Options:
`)
	flag.PrintDefaults()
}

// ParseOptions Parsing user's options
func ParseOptions() *Options {
	flag.Parse()

	//add config file support
	if Args.ConfigFile != "" {
		Args = parseConfigFile(Args.ConfigFile)
	}
	if Args.Listen != "" && Args.Connect == "" && Args.Proxy == "" { // ./stowaway_admin -l <port> -s [secret]
		Args.Mode = NORMAL_PASSIVE
		printer.Warning("[*] Starting admin node on port %s\r\n", Args.Listen)
	} else if Args.Connect != "" && Args.Listen == "" && Args.Proxy == "" { // ./stowaway_admin -c <ip:port> -s [secret]
		Args.Mode = NORMAL_ACTIVE
		printer.Warning("[*] Trying to connect node actively")
	} else if Args.Connect != "" && Args.Listen == "" && Args.Proxy != "" { // ./stowaway_admin -c <ip:port> -s [secret] --proxy <ip:port> --proxyu [username] --proxyp [password]
		Args.Mode = PROXY_ACTIVE
		printer.Warning("[*] Trying to connect node actively with proxy %s\r\n", Args.Proxy)
	} else { // Wrong format
		termbox.Close()
		flag.Usage()
		os.Exit(0)
	}

	if err := checkOptions(Args); err != nil {
		termbox.Close()
		printer.Fail("[*] Options err: %s\r\n", err.Error())
		os.Exit(0)
	}

	return Args
}

func checkOptions(option *Options) error {
	var err error

	if Args.Connect != "" {
		_, err = net.ResolveTCPAddr("", option.Connect)
	}

	if Args.Proxy != "" {
		_, err = net.ResolveTCPAddr("", option.Proxy)
	}

	return err
}
