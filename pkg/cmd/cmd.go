package cmd

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/jenkins-x-labs/cli-doc-gen/pkg/common"
	"github.com/jenkins-x/jx/pkg/cmd/helper"
	"github.com/jenkins-x/jx/pkg/log"
	"github.com/jenkins-x/jx/pkg/util"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type Options struct {
	File        string
	OutFile     string
	StartPrefix string
	EndPrefix   string
	TrimPrefix  []string
}

// NewCmd creates the new command
func NewCmd() (*cobra.Command, *Options) {
	o := &Options{}
	cmd := &cobra.Command{
		Use:   "cli-doc-gen",
		Short: "a tool to generate website documentation files from shell scripts or other source code",
		Run: func(cmd *cobra.Command, args []string) {
			common.SetLoggingLevel(cmd)
			err := o.Run()
			helper.CheckErr(err)
		},
	}
	cmd.Flags().StringVarP(&o.File, "file", "f", "", "the file name to use to extract a set of commands")
	cmd.Flags().StringVarP(&o.OutFile, "out", "o", "", "the output file to create")
	cmd.Flags().StringVarP(&o.StartPrefix, "start", "s", "# CLI-DOC-GEN-START", "the prefix of the line starting the extraction")
	cmd.Flags().StringVarP(&o.EndPrefix, "end", "e", "# CLI-DOC-GEN-END", "the prefix of the line ending the extraction")
	cmd.Flags().StringArrayVarP(&o.TrimPrefix, "trim-prefix", "", nil, "prefix strings to trim from the start of each line")
	return cmd, o
}

func (o *Options) Run() error {
	if o.File == "" {
		return util.MissingOption("file")
	}
	startPrefix := o.StartPrefix
	if startPrefix == "" {
		return util.MissingOption("start")
	}
	endPrefix := o.EndPrefix
	if startPrefix == "" {
		return util.MissingOption("end")
	}

	data, err := ioutil.ReadFile(o.File)
	if err != nil {
		return errors.Wrapf(err, "failed to read input file %s", o.File)
	}

	lines := strings.Split(string(data), "\n")
	started := false
	ended := false
	buffer := strings.Builder{}
	for i, line := range lines {
		trim := strings.TrimSpace(line)
		if strings.HasPrefix(trim, endPrefix) {
			if started {
				ended = true
				break
			}
			return errors.Errorf("line %d: found end prefix '%s' before a start prefix '%s'", i, endPrefix, startPrefix)
		}
		if strings.HasPrefix(trim, startPrefix) {
			started = true
		} else if started {
			text := line
			for _, tp := range o.TrimPrefix {
				text = strings.TrimPrefix(text, tp)
			}
			buffer.WriteString(text)
			buffer.WriteString("\n")
		}
	}
	if !started {
		return errors.Errorf("never found a start prefix '%s' in the file %s", startPrefix, o.File)
	}
	if !ended {
		return errors.Errorf("never found the end prefix '%s' in the file %s", endPrefix, o.File)
	}

	text := buffer.String()

	if o.OutFile == "" {
		log.Logger().Info(text)
		return nil
	}

	// lets check parent dir exists
	dir := filepath.Dir(o.OutFile)
	err = os.MkdirAll(dir, util.DefaultWritePermissions)
	if err != nil {
		return errors.Wrapf(err, "failed to create output directory %s", dir)
	}

	err = ioutil.WriteFile(o.OutFile, []byte(text), util.DefaultFileWritePermissions)
	if err != nil {
		return errors.Wrapf(err, "failed to save file %s", o.OutFile)
	}

	log.Logger().Infof("saved file %s", util.ColorInfo(o.OutFile))
	return nil
}
