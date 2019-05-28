// Copyright 2019 The Hugo Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package commands

import (
	"os"

	"github.com/gohugoio/hugo/config"

	"github.com/gohugoio/hugo/hugofs"

	"github.com/gohugoio/hugo/mods"
	"github.com/spf13/cobra"
)

var _ cmder = (*modCmd)(nil)

type modCmd struct {
	*baseBuilderCmd
}

func (b *commandsBuilder) newModCmd() *modCmd {
	c := &modCmd{}

	cmd := &cobra.Command{
		Use:   "mod",
		Short: "Various Hugo Modules helpers.",
		RunE:  nil,
	}

	cmd.AddCommand(
		&cobra.Command{
			// go get [-d] [-m] [-u] [-v] [-insecure] [build flags] [packages]
			Use:   "get",
			Short: "TODO(bep)",
			RunE: func(cmd *cobra.Command, args []string) error {
				if len(args) >= 1 {
					c, err := c.newModsClient(nil)
					if err != nil {
						return err
					}
					return c.Get(args[0])
				}

				// Collect any modules defined in config.toml
				_, err := c.initConfig()
				return err

			},
		},
		&cobra.Command{
			Use:   "graph",
			Short: "TODO(bep)",
			RunE: func(cmd *cobra.Command, args []string) error {
				return c.withModsClient(func(c *mods.Client) error {
					return c.Graph(os.Stdout)
				})
			},
		},
		&cobra.Command{
			Use:   "init",
			Short: "TODO(bep) ",
			RunE: func(cmd *cobra.Command, args []string) error {
				var path string
				if len(args) >= 1 {
					path = args[0]
				}
				return c.withModsClient(func(c *mods.Client) error {
					return c.Init(path)
				})
			},
		},
		&cobra.Command{
			Use:   "vendor",
			Short: "TODO(bep)",
			RunE: func(cmd *cobra.Command, args []string) error {
				return c.withModsClient(func(c *mods.Client) error {
					return c.Vendor()
				})
			},
		},
		&cobra.Command{
			Use:   "tidy",
			Short: "TODO(bep)",
			RunE: func(cmd *cobra.Command, args []string) error {
				return c.withModsClient(func(c *mods.Client) error {
					return c.Tidy()
				})
			},
		},
	)

	c.baseBuilderCmd = b.newBuilderCmd(cmd)

	return c

}

func (c *modCmd) withModsClient(f func(*mods.Client) error) error {
	com, err := c.initConfig()
	if err != nil {
		return err
	}
	client, err := c.newModsClient(com.Cfg)
	if err != nil {
		return err
	}
	return f(client)
}

func (c *modCmd) initConfig() (*commandeer, error) {
	com, err := initializeConfig(true, false, &c.hugoBuilderCommon, c, nil)
	if err != nil {
		return nil, err
	}
	return com, nil
}

func (c *modCmd) newModsClient(cfg config.Provider) (*mods.Client, error) {
	var (
		workingDir   string
		themesDir    string
		themes       []string
		ignoreVendor bool
	)

	if c.source != "" {
		workingDir = c.source
	} else {
		var err error
		workingDir, err = os.Getwd()
		if err != nil {
			return nil, err
		}
	}

	if cfg != nil {
		// TODO(bep) mod remember this if we change
		themesDir = cfg.GetString("themesDir")
		themes = cfg.GetStringSlice("theme")
		ignoreVendor = cfg.GetBool("ignoreVendor")
	}

	return mods.NewClient(hugofs.Os, ignoreVendor, workingDir, themesDir, themes), nil
}
