package flag

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Context string

const (
	Health Context = "health"
	Empty  Context = ""

	ScannerSAST    = "sast"
	ScannerSecrets = "secrets"

	ErrorLogLevel = "error"
	InfoLogLevel  = "info"
	DebugLogLevel = "debug"
	TraceLogLevel = "trace"
)

var (
	ErrInvalidContext = errors.New("invalid context argument; supported values: health")
	ErrInvalidScanner = errors.New("invalid scanner argument; supported values: sast, secrets")
)

var (
	SkipPathFlag = Flag{
		Name:       "skip-path",
		ConfigName: "scan.skip-path",
		Value:      []string{},
		Usage:      "Specify the comma separated files and directories to skip. Supports * syntax, e.g. --skip-path users/*.go,users/admin.sql",
	}
	DebugFlag = Flag{
		Name:            "debug",
		ConfigName:      "scan.debug",
		Value:           false,
		Usage:           "Enable debug logs. Equivalent to --log-level=debug",
		DisableInConfig: true,
	}
	LogLevelFlag = Flag{
		Name:       "log-level",
		ConfigName: "scan.log-level",
		Value:      "info",
		Usage:      "Set log level (error, info, debug, trace)",
	}
	DisableDomainResolutionFlag = Flag{
		Name:       "disable-domain-resolution",
		ConfigName: "scan.disable-domain-resolution",
		Value:      true,
		Usage:      "Do not attempt to resolve detected domains during classification",
	}
	DomainResolutionTimeoutFlag = Flag{
		Name:       "domain-resolution-timeout",
		ConfigName: "scan.domain-resolution-timeout",
		Value:      3 * time.Second,
		Usage:      "Set timeout when attempting to resolve detected domains during classification, e.g. --domain-resolution-timeout=3s",
	}
	InternalDomainsFlag = Flag{
		Name:       "internal-domains",
		ConfigName: "scan.internal-domains",
		Value:      []string{},
		Usage:      "Define regular expressions for better classification of private or unreachable domains e.g. --internal-domains=\".*.my-company.com,private.sh\"",
	}
	ContextFlag = Flag{
		Name:       "context",
		ConfigName: "scan.context",
		Value:      "",
		Usage:      "Expand context of schema classification e.g., --context=health, to include data types particular to health",
	}
	DataSubjectMappingFlag = Flag{
		Name:       "data-subject-mapping",
		ConfigName: "scan.data_subject_mapping",
		Value:      "",
		Usage:      "Override default data subject mapping by providing a path to a custom mapping JSON file",
	}
	QuietFlag = Flag{
		Name:       "quiet",
		ConfigName: "scan.quiet",
		Value:      false,
		Usage:      "Suppress non-essential messages",
	}
	ForceFlag = Flag{
		Name:       "force",
		ConfigName: "scan.force",
		Value:      false,
		Usage:      "Disable the cache and runs the detections again",
	}
	ExternalRuleDirFlag = Flag{
		Name:       "external-rule-dir",
		ConfigName: "scan.external-rule-dir",
		Value:      []string{},
		Usage:      "Specify directories paths that contain .yaml files with external rules configuration",
	}
	ScannerFlag = Flag{
		Name:       "scanner",
		ConfigName: "scan.scanner",
		Value:      []string{ScannerSAST},
		Usage:      "Specify which scanner to use e.g. --scanner=secrets, --scanner=secrets,sast",
	}
	ParallelFlag = Flag{
		Name:       "parallel",
		ConfigName: "scan.parallel",
		Value:      0,
		Usage:      "Specify the amount of parallelism to use during the scan",
	}
	ExitCodeFlag = Flag{
		Name:       "exit-code",
		ConfigName: "scan.exit-code",
		Value:      -1,
		Usage:      "Force a given exit code for the scan command. Set this to 0 (success) to always return a success exit code despite any findings from the scan.",
	}
)

type ScanFlagGroup struct {
	ScannerFlag                 *Flag
	SkipPathFlag                *Flag
	DebugFlag                   *Flag
	LogLevelFlag                *Flag
	DisableDomainResolutionFlag *Flag
	DomainResolutionTimeoutFlag *Flag
	InternalDomainsFlag         *Flag
	ContextFlag                 *Flag
	DataSubjectMappingFlag      *Flag
	QuietFlag                   *Flag
	ForceFlag                   *Flag
	ExternalRuleDirFlag         *Flag
	ParallelFlag                *Flag
	ExitCodeFlag                *Flag
}

type ScanOptions struct {
	Target                  string        `mapstructure:"target" json:"target" yaml:"target"`
	SkipPath                []string      `mapstructure:"skip-path" json:"skip-path" yaml:"skip-path"`
	Debug                   bool          `mapstructure:"debug" json:"debug" yaml:"debug"`
	LogLevel                string        `mapstructure:"log-level" json:"log-level" yaml:"log-level"`
	DisableDomainResolution bool          `mapstructure:"disable-domain-resolution" json:"disable-domain-resolution" yaml:"disable-domain-resolution"`
	DomainResolutionTimeout time.Duration `mapstructure:"domain-resolution-timeout" json:"domain-resolution-timeout" yaml:"domain-resolution-timeout"`
	InternalDomains         []string      `mapstructure:"internal-domains" json:"internal-domains" yaml:"internal-domains"`
	Context                 Context       `mapstructure:"context" json:"context" yaml:"context"`
	DataSubjectMapping      string        `mapstructure:"data_subject_mapping" json:"data_subject_mapping" yaml:"data_subject_mapping"`
	Quiet                   bool          `mapstructure:"quiet" json:"quiet" yaml:"quiet"`
	Force                   bool          `mapstructure:"force" json:"force" yaml:"force"`
	ExternalRuleDir         []string      `mapstructure:"external-rule-dir" json:"external-rule-dir" yaml:"external-rule-dir"`
	Scanner                 []string      `mapstructure:"scanner" json:"scanner" yaml:"scanner"`
	Parallel                int           `mapstructure:"parallel" json:"parallel" yaml:"parallel"`
	ExitCode                int           `mapstructure:"exit-code" json:"exit-code" yaml:"exit-code"`
	DiffBaseBranch          string        `mapstructure:"diff_base_branch" json:"diff_base_branch" yaml:"diff_base_branch"`
	GithubToken             string        `mapstructure:"github_token" json:"github_token" yaml:"github_token"`
}

func NewScanFlagGroup() *ScanFlagGroup {
	return &ScanFlagGroup{
		SkipPathFlag:                &SkipPathFlag,
		DebugFlag:                   &DebugFlag,
		LogLevelFlag:                &LogLevelFlag,
		DisableDomainResolutionFlag: &DisableDomainResolutionFlag,
		DomainResolutionTimeoutFlag: &DomainResolutionTimeoutFlag,
		InternalDomainsFlag:         &InternalDomainsFlag,
		ContextFlag:                 &ContextFlag,
		DataSubjectMappingFlag:      &DataSubjectMappingFlag,
		QuietFlag:                   &QuietFlag,
		ForceFlag:                   &ForceFlag,
		ExternalRuleDirFlag:         &ExternalRuleDirFlag,
		ScannerFlag:                 &ScannerFlag,
		ParallelFlag:                &ParallelFlag,
		ExitCodeFlag:                &ExitCodeFlag,
	}
}

func (f *ScanFlagGroup) Name() string {
	return "Scan"
}

func (f *ScanFlagGroup) Flags() []*Flag {
	return []*Flag{
		f.SkipPathFlag,
		f.DebugFlag,
		f.LogLevelFlag,
		f.DisableDomainResolutionFlag,
		f.DomainResolutionTimeoutFlag,
		f.InternalDomainsFlag,
		f.ContextFlag,
		f.DataSubjectMappingFlag,
		f.QuietFlag,
		f.ForceFlag,
		f.ExternalRuleDirFlag,
		f.ScannerFlag,
		f.ParallelFlag,
		f.ExitCodeFlag,
	}
}

func (f *ScanFlagGroup) ToOptions(args []string) (ScanOptions, error) {
	var target string
	if len(args) == 1 {
		target = args[0]
	}

	context := getContext(f.ContextFlag)
	switch context {
	case Empty, Health:
	default:
		return ScanOptions{}, ErrInvalidContext
	}

	scanners := getStringSlice(f.ScannerFlag)
	for _, scanner := range scanners {
		switch scanner {
		case ScannerSAST:
		case ScannerSecrets:
		default:
			return ScanOptions{}, ErrInvalidScanner
		}
	}

	debug := getBool(f.DebugFlag)
	logLevel := getString(f.LogLevelFlag)
	if debug {
		logLevel = DebugLogLevel
	}

	return ScanOptions{
		SkipPath:                getStringSlice(f.SkipPathFlag),
		Debug:                   debug,
		LogLevel:                logLevel,
		DisableDomainResolution: getBool(f.DisableDomainResolutionFlag),
		DomainResolutionTimeout: getDuration(f.DomainResolutionTimeoutFlag),
		InternalDomains:         getStringSlice(f.InternalDomainsFlag),
		Context:                 context,
		DataSubjectMapping:      getString(f.DataSubjectMappingFlag),
		Quiet:                   getBool(f.QuietFlag),
		Force:                   getBool(f.ForceFlag),
		Target:                  target,
		ExternalRuleDir:         getStringSlice(f.ExternalRuleDirFlag),
		Scanner:                 scanners,
		Parallel:                viper.GetInt(f.ParallelFlag.ConfigName),
		ExitCode:                viper.GetInt(f.ExitCodeFlag.ConfigName),
		DiffBaseBranch:          os.Getenv("DIFF_BASE_BRANCH"),
		GithubToken:             os.Getenv("GITHUB_TOKEN"),
	}, nil
}

func getContext(flag *Flag) Context {
	if flag == nil {
		return ""
	}

	flagStr := strings.ToLower(getString(flag))
	return Context(flagStr)
}
