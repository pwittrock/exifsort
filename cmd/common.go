/*
Copyright © 2020 Michael Rubin <mhr@neverthere.org>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"io"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
)

// common routine to select writer from command line.
func ioWriter(quiet bool) io.Writer {
	if quiet {
		return ioutil.Discard
	}

	return os.Stdout
}

type cmdStringFlag struct {
	shorthand string
	name      string
	required  bool
	usage     string
}

func setStringFlags(cmd *cobra.Command, flags []cmdStringFlag) {
	for _, f := range flags {
		cmd.Flags().StringP(f.name, f.shorthand, "", f.usage)

		if f.required {
			_ = cmd.MarkFlagRequired(f.name)
		}
	}
}
