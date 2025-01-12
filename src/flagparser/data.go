package flagparser

import (
	"flag"
	"fmt"
	resource "github.com/SongZihuan/huan-proxy"
	"github.com/SongZihuan/huan-proxy/src/utils"
	"io"
	"os"
	"reflect"
	"strings"
)

const OptionIdent = "  "
const OptionPrefix = "--"
const UsagePrefixWidth = 10

type flagData struct {
	flagReady  bool
	flagSet    bool
	flagParser bool

	HelpData        bool
	HelpName        string
	HelpUsage       string
	VersionData     bool
	VersionName     string
	VersionUsage    string
	LicenseData     bool
	LicenseName     string
	LicenseUsage    string
	ReportData      bool
	ReportName      string
	ReportUsage     string
	ConfigFileData  string
	ConfigFileName  string
	ConfigFileUsage string
	TermData        bool
	TermName        string
	TermUsage       string

	Usage string
}

func initData() {
	data = flagData{
		flagReady:  false,
		flagSet:    false,
		flagParser: false,

		HelpData:        false,
		HelpName:        "help",
		HelpUsage:       fmt.Sprintf("Show usage of %s. If this option is set, the backend service will not run.", utils.GetArgs0Name()),
		VersionData:     false,
		VersionName:     "version",
		VersionUsage:    fmt.Sprintf("Show version of %s. If this option is set, the backend service will not run.", utils.GetArgs0Name()),
		LicenseData:     false,
		LicenseName:     "license",
		LicenseUsage:    fmt.Sprintf("Show license of %s. If this option is set, the backend service will not run.", utils.GetArgs0Name()),
		ReportData:      false,
		ReportName:      "report",
		ReportUsage:     fmt.Sprintf("Show how to report questions/errors of %s. If this option is set, the backend service will not run.", utils.GetArgs0Name()),
		ConfigFileData:  "config.yaml",
		ConfigFileName:  "config",
		ConfigFileUsage: fmt.Sprintf("%s", "The location of the running configuration file of the backend service. The option is a string, the default value is config.yaml in the running directory."),
		TermData:        false,
		TermName:        "term",
		TermUsage:       "Terminal output mode. If you have this program running in the foreground terminal (standard output is the terminal), you can turn this option on to get formatted logs.",
		Usage:           "",
	}

	data.ready()
}

func (d *flagData) writeUsage() {
	if len(d.Usage) != 0 {
		return
	}

	if d.isFlagSet() || d.isFlagParser() {
		panic("flag is parser")
	}

	var result strings.Builder
	result.WriteString(utils.FormatTextToWidth(fmt.Sprintf("Usage of %s:", utils.GetArgs0Name()), utils.NormalConsoleWidth))
	result.WriteString("\n")

	val := reflect.ValueOf(*d)
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)

		if !strings.HasSuffix(field.Name, "Data") {
			continue
		}

		option := field.Name[:len(field.Name)-4]
		optionName, ok := val.FieldByName(option + "Name").Interface().(string)
		if !ok {
			panic("can not get option name")
		}

		fmt.Printf("OPTION: %s\n", option+"Usage")
		optionUsage, ok := val.FieldByName(option + "Usage").Interface().(string)
		if !ok {
			panic("can not get option Usage")
		}

		var title string
		if field.Type.Kind() == reflect.Bool {
			optionData, ok := val.FieldByName(option + "Data").Interface().(bool)
			if !ok {
				panic("can not get option data")
			}

			if optionData == true {
				panic("bool option can not be true")
			}

			title1 := fmt.Sprintf("%s%s%s", OptionIdent, OptionPrefix, utils.FormatTextToWidth(optionName, utils.NormalConsoleWidth-len(OptionIdent)-len(OptionPrefix)))
			title2 := fmt.Sprintf("%s%s%s", OptionIdent, OptionPrefix, utils.FormatTextToWidth(optionName[0:1], utils.NormalConsoleWidth-len(OptionIdent)-len(OptionPrefix)))
			title = fmt.Sprintf("%s\n%s", title1, title2)
		} else if field.Type.Kind() == reflect.String {
			optionData, ok := val.FieldByName(option + "Data").Interface().(string)
			if !ok {
				panic("can not get option data")
			}

			title1 := fmt.Sprintf("%s%s%s", OptionIdent, OptionPrefix, utils.FormatTextToWidth(fmt.Sprintf("%s string, default: '%s'", optionName, optionData), utils.NormalConsoleWidth-len(OptionIdent)-len(OptionPrefix)))
			title2 := fmt.Sprintf("%s%s%s", OptionIdent, OptionPrefix, utils.FormatTextToWidth(fmt.Sprintf("%s string, default: '%s'", optionName[0:1], optionData), utils.NormalConsoleWidth-len(OptionIdent)-len(OptionPrefix)))
			title = fmt.Sprintf("%s\n%s", title1, title2)
		} else if field.Type.Kind() == reflect.Uint {
			optionData, ok := val.FieldByName(option + "Data").Interface().(uint)
			if !ok {
				panic("can not get option data")
			}

			title1 := fmt.Sprintf("%s%s%s", OptionIdent, OptionPrefix, utils.FormatTextToWidth(fmt.Sprintf("%s number, default: %d", optionName, optionData), utils.NormalConsoleWidth-len(OptionIdent)-len(OptionPrefix)))
			title2 := fmt.Sprintf("%s%s%s", OptionIdent, OptionPrefix, utils.FormatTextToWidth(fmt.Sprintf("%s number, default: %d", optionName[0:1], optionData), utils.NormalConsoleWidth-len(OptionIdent)-len(OptionPrefix)))
			title = fmt.Sprintf("%s\n%s", title1, title2)
		} else {
			panic("error flag type")
		}

		result.WriteString(title)
		result.WriteString("\n")

		usegae := utils.FormatTextToWidthAndPrefix(optionUsage, UsagePrefixWidth, utils.NormalConsoleWidth)
		result.WriteString(usegae)
		result.WriteString("\n\n")
	}

	d.Usage = strings.TrimRight(result.String(), "\n")
}

func (d *flagData) setFlag() {
	if d.isFlagSet() {
		return
	}

	flag.BoolVar(&d.HelpData, data.HelpName, data.HelpData, data.HelpUsage)
	flag.BoolVar(&d.HelpData, data.HelpName[0:1], data.HelpData, data.HelpUsage)

	flag.BoolVar(&d.VersionData, data.VersionName, data.VersionData, data.VersionUsage)
	flag.BoolVar(&d.VersionData, data.VersionName[0:1], data.VersionData, data.VersionUsage)

	flag.BoolVar(&d.LicenseData, data.LicenseName, data.LicenseData, data.LicenseUsage)
	flag.BoolVar(&d.LicenseData, data.LicenseName[0:1], data.LicenseData, data.LicenseUsage)

	flag.BoolVar(&d.ReportData, data.ReportName, data.ReportData, data.ReportUsage)
	flag.BoolVar(&d.ReportData, data.ReportName[0:1], data.ReportData, data.ReportUsage)

	flag.StringVar(&d.ConfigFileData, data.ConfigFileName, data.ConfigFileData, data.ConfigFileUsage)
	flag.StringVar(&d.ConfigFileData, data.ConfigFileName[0:1], data.ConfigFileData, data.ConfigFileUsage)

	flag.BoolVar(&d.TermData, data.TermName, data.TermData, data.TermUsage)
	flag.BoolVar(&d.TermData, data.TermName[0:1], data.TermData, data.TermUsage)

	flag.Usage = func() {
		_, _ = d.PrintUsage()
	}
	d.flagSet = true
}

func (d *flagData) parser() {
	if d.flagParser {
		return
	}

	if !d.isFlagSet() {
		panic("flag not set")
	}

	flag.Parse()
	d.flagParser = true
}

func (d *flagData) ready() {
	if d.isReady() {
		return
	}

	d.writeUsage()
	d.setFlag()
	d.parser()
	d.flagReady = true
}

func (d *flagData) isReady() bool {
	return d.isFlagSet() && d.isFlagParser() && d.flagReady
}

func (d *flagData) isFlagSet() bool {
	return d.flagSet
}

func (d *flagData) isFlagParser() bool {
	return d.flagParser
}

func (d *flagData) Help() bool {
	if !d.isReady() {
		panic("flag not ready")
	}

	return d.HelpData
}

func (d *flagData) FprintUsage(writer io.Writer) (int, error) {
	return fmt.Fprintf(writer, "%s\n", d.Usage)
}

func (d *flagData) PrintUsage() (int, error) {
	return d.FprintUsage(flag.CommandLine.Output())
}

func (d *flagData) Version() bool {
	if !d.isReady() {
		panic("flag not ready")
	}

	return d.VersionData
}

func (d *flagData) FprintVersion(writer io.Writer) (int, error) {
	version := utils.FormatTextToWidth(fmt.Sprintf("Version of %s: %s", utils.GetArgs0Name(), resource.Version), utils.NormalConsoleWidth)
	return fmt.Fprintf(writer, "%s\n", version)
}

func (d *flagData) PrintVersion() (int, error) {
	return d.FprintVersion(flag.CommandLine.Output())
}

func (d *flagData) FprintLicense(writer io.Writer) (int, error) {
	title := utils.FormatTextToWidth(fmt.Sprintf("License of %s:", utils.GetArgs0Name()), utils.NormalConsoleWidth)
	license := utils.FormatTextToWidth(resource.License, utils.NormalConsoleWidth)
	return fmt.Fprintf(writer, "%s\n%s\n", title, license)
}

func (d *flagData) PrintLicense() (int, error) {
	return d.FprintLicense(flag.CommandLine.Output())
}

func (d *flagData) FprintReport(writer io.Writer) (int, error) {
	// 不需要title
	report := utils.FormatTextToWidth(resource.Report, utils.NormalConsoleWidth)
	return fmt.Fprintf(os.Stderr, "%s\n", report)
}

func (d *flagData) PrintReport() (int, error) {
	return d.FprintReport(flag.CommandLine.Output())
}

func (d *flagData) FprintLF(writer io.Writer) (int, error) {
	return fmt.Fprintf(os.Stderr, "\n")
}

func (d *flagData) PrintLF() (int, error) {
	return d.FprintLF(flag.CommandLine.Output())
}

func (d *flagData) License() bool {
	if !d.isReady() {
		panic("flag not ready")
	}

	return d.LicenseData
}

func (d *flagData) Report() bool {
	if !d.isReady() {
		panic("flag not ready")
	}

	return d.ReportData
}

func (d *flagData) ConfigFile() string {
	if !d.isReady() {
		panic("flag not ready")
	}

	return d.ConfigFileData
}

func (d *flagData) Term() bool {
	if !d.isReady() {
		panic("flag not ready")
	}

	return d.TermData
}

func (d *flagData) SetOutput(writer io.Writer) {
	flag.CommandLine.SetOutput(writer)
}
