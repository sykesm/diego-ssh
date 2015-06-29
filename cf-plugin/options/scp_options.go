package options

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/pborman/getopt"
)

type SCPOptions struct {
	Sources []FileLocation
	Target  FileLocation

	Verbose            bool
	PreserveAttributes bool
	Recurse            bool

	getoptSet                *getopt.Set
	verboseOption            getopt.Option
	preserveAttributesOption getopt.Option
	recurseOption            getopt.Option
}

type FileLocation struct {
	AppName string // empty string is for local path
	Index   uint
	Path    string
}

func NewSCPOptions() *SCPOptions {
	scpOptions := &SCPOptions{}

	opts := getopt.New()

	scpOptions.verboseOption = opts.BoolVar(
		&scpOptions.Verbose,
		'v',
		"enable verbose output",
	)

	scpOptions.preserveAttributesOption = opts.BoolVar(
		&scpOptions.PreserveAttributes,
		'p',
		"preserve file times and permissions",
	)

	scpOptions.recurseOption = opts.BoolVar(
		&scpOptions.Recurse,
		'r',
		"recurse into directories",
	)

	scpOptions.getoptSet = opts

	return scpOptions
}

func (o *SCPOptions) Parse(args []string) error {
	if len(args) < 1 || args[0] != "scp" {
		return UsageError
	}

	opts := o.getoptSet

	err := opts.Getopt(args, nil)
	if err != nil {
		return err
	}

	if len(opts.Args()) < 2 {
		return errors.New("Source and target must be provided")
	}

	locs, err := o.parseLocations(opts.Args())
	if err != nil {
		return err
	}

	o.Sources = locs[:len(locs)-1]
	o.Target = locs[len(locs)-1]

	return nil
}

func (o *SCPOptions) parseLocations(args []string) ([]FileLocation, error) {
	locs := []FileLocation{}

	for _, arg := range args {
		location, err := ParseLocation(arg)
		if err != nil {
			return nil, err
		}
		locs = append(locs, location)
	}

	return locs, nil
}

func ParseLocation(arg string) (FileLocation, error) {
	parts := splitFirstUnescapedColon(arg)

	location := FileLocation{}
	if len(parts) == 2 {
		host, index, err := splitHostIndex(parts[0])
		if err != nil {
			return location, err
		}
		location.AppName = host
		location.Index = index
		location.Path = parts[1]
	} else {
		location.Path = arg
	}

	location.Path = strings.Replace(location.Path, "\\:", ":", -1)

	return location, nil
}

func splitFirstUnescapedColon(arg string) []string {
	for i := 1; i < len(arg); i++ {
		if arg[i] != ':' {
			continue
		}
		if arg[i-1] == '\\' {
			continue
		}

		return []string{arg[:i], arg[i+1:]}
	}

	return []string{arg}
}

func splitHostIndex(arg string) (string, uint, error) {
	parts := strings.Split(arg, "/")
	switch len(parts) {
	case 1:
		return arg, 0, nil
	case 2:
		index, err := strconv.ParseUint(parts[1], 10, 32)
		if err != nil {
			return arg, 0, nil
		}
		return parts[0], uint(index), nil
	default:
		return "", 0, fmt.Errorf("invalid host/index format: %q", arg)
	}
}

func SCPUsage() string {
	b := &bytes.Buffer{}

	o := NewSCPOptions()
	o.getoptSet.SetProgram("scp")
	o.getoptSet.SetParameters("[app[/index]:]file1 ... [app[/index]:]file2")
	o.getoptSet.PrintUsage(b)

	return b.String()
}
