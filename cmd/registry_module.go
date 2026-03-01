// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package cmd

import (
	"math"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/straubt1/tfx/client"
	"github.com/straubt1/tfx/cmd/flags"
	view "github.com/straubt1/tfx/cmd/views"
	"github.com/straubt1/tfx/data"
)

var (
	// `tfx registry module` commands
	registryModuleCmd = &cobra.Command{
		Use:   "module",
		Short: "Modules in Private Registry Commands",
		Long:  "Work with Private Module Registry of a TFx Organization.",
	}

	// `tfx registry module list` command
	registryModuleListCmd = &cobra.Command{
		Use:   "list",
		Short: "List modules",
		Long:  "List modules in the Private Module Registry of a TFx Organization.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdConfig, err := flags.ParseRegistryModuleListFlags(cmd)
			if err != nil {
				return err
			}
			if cmdConfig.All {
				cmdConfig.MaxItems = math.MaxInt
			}
			return registryModuleList(cmdConfig)
		},
	}

	// `tfx registry module create` command
	registryModuleCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a module",
		Long:  "Create a module in the Private Module Registry for a TFx Organization.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdConfig, err := flags.ParseRegistryModuleCreateFlags(cmd)
			if err != nil {
				return err
			}
			return registryModuleCreate(cmdConfig)
		},
	}

	// `tfx registry module show` command
	registryModuleShowCmd = &cobra.Command{
		Use:   "show",
		Short: "Show a module",
		Long:  "Show a module details of a module in the Private Module Registry for a TFx Organization.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdConfig, err := flags.ParseRegistryModuleShowFlags(cmd)
			if err != nil {
				return err
			}
			return registryModuleShow(cmdConfig)
		},
	}

	// `tfx registry module delete` command
	registryModuleDeleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete a module",
		Long:  "Delete a module in the Private Module Registry for a TFx Organization.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdConfig, err := flags.ParseRegistryModuleDeleteFlags(cmd)
			if err != nil {
				return err
			}
			return registryModuleDelete(cmdConfig)
		},
	}
)

func init() {
	// `tfx registry module list` arguments
	registryModuleListCmd.Flags().IntP("max-items", "m", 10, "Max number of results (optional)")
	registryModuleListCmd.Flags().BoolP("all", "a", false, "Retrieve all results regardless of maxItems flag (optional)")

	// `tfx registry module create` arguments
	registryModuleCreateCmd.Flags().StringP("name", "n", "", "Name of the Module (no spaces)")
	registryModuleCreateCmd.Flags().StringP("provider", "p", "", "Name of the provider (no spaces) (i.e. aws, azure, google)")
	registryModuleCreateCmd.MarkFlagRequired("name")
	registryModuleCreateCmd.MarkFlagRequired("provider")

	// `tfx registry module show` arguments
	registryModuleShowCmd.Flags().StringP("name", "n", "", "Name of the Module (no spaces)")
	registryModuleShowCmd.Flags().StringP("provider", "p", "", "Name of the provider (no spaces) (i.e. aws, azure, google)")
	registryModuleShowCmd.MarkFlagRequired("name")
	registryModuleShowCmd.MarkFlagRequired("provider")

	// `tfx registry module delete` arguments
	registryModuleDeleteCmd.Flags().StringP("name", "n", "", "Name of the Module (no spaces)")
	registryModuleDeleteCmd.Flags().StringP("provider", "p", "", "Name of the provider (no spaces) (i.e. aws, azure, google)")
	registryModuleDeleteCmd.MarkFlagRequired("name")
	registryModuleDeleteCmd.MarkFlagRequired("provider")

	registryCmd.AddCommand(registryModuleCmd)
	registryModuleCmd.AddCommand(registryModuleListCmd)
	registryModuleCmd.AddCommand(registryModuleCreateCmd)
	registryModuleCmd.AddCommand(registryModuleShowCmd)
	registryModuleCmd.AddCommand(registryModuleDeleteCmd)
}

func registryModuleList(cmdConfig *flags.RegistryModuleListFlags) error {
	v := view.NewRegistryModuleListView()
	c, err := client.NewFromViper()
	if err != nil {
		return v.RenderError(err)
	}
	v.PrintCommandHeader("List Modules for Organization: %s", c.OrganizationName)
	items, err := data.ListRegistryModules(c, c.OrganizationName, cmdConfig.MaxItems)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to list modules"))
	}
	return v.Render(items)
}

func registryModuleCreate(cmdConfig *flags.RegistryModuleCreateFlags) error {
	v := view.NewRegistryModuleCreateView()
	c, err := client.NewFromViper()
	if err != nil {
		return v.RenderError(err)
	}
	v.PrintCommandHeader("Create Module for Organization: %s", c.OrganizationName)
	module, err := data.CreateRegistryModule(c, c.OrganizationName, cmdConfig.Name, cmdConfig.Provider)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to create module"))
	}
	return v.Render(module)
}

func registryModuleShow(cmdConfig *flags.RegistryModuleShowFlags) error {
	v := view.NewRegistryModuleShowView()
	c, err := client.NewFromViper()
	if err != nil {
		return v.RenderError(err)
	}
	v.PrintCommandHeader("Show Module for Organization: %s", c.OrganizationName)
	module, err := data.ReadRegistryModule(c, c.OrganizationName, cmdConfig.Name, cmdConfig.Provider)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to show module"))
	}
	return v.Render(module)
}

func registryModuleDelete(cmdConfig *flags.RegistryModuleDeleteFlags) error {
	v := view.NewRegistryModuleDeleteView()
	c, err := client.NewFromViper()
	if err != nil {
		return v.RenderError(err)
	}
	v.PrintCommandHeader("Delete Module for Organization: %s", c.OrganizationName)
	err = data.DeleteRegistryModule(c, c.OrganizationName, cmdConfig.Name, cmdConfig.Provider)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to delete module"))
	}
	return v.Render(cmdConfig.Name)
}
