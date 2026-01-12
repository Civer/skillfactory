package keys

import (
	"fmt"

	"habitwire/client"

	"github.com/spf13/cobra"
)

// RegisterCommands creates and returns the keys command group
func RegisterCommands(c *client.Client, printJSON func(interface{}) error) *cobra.Command {
	service := NewService(c)

	cmd := &cobra.Command{
		Use:   "keys",
		Short: "Manage API keys",
	}

	// list
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all API keys",
		RunE: func(cmd *cobra.Command, args []string) error {
			keys, err := service.List()
			if err != nil {
				return err
			}
			return printJSON(ToLeanSlice(keys))
		},
	}

	// create
	var createName string
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new API key",
		Long:  "Create a new API key. The key value is only shown once upon creation.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if createName == "" {
				return fmt.Errorf("--name is required")
			}
			key, err := service.Create(CreateKeyRequest{Name: createName})
			if err != nil {
				return err
			}
			return printJSON(key.ToLean())
		},
	}
	createCmd.Flags().StringVarP(&createName, "name", "n", "", "Key name (required)")

	// delete
	deleteCmd := &cobra.Command{
		Use:   "delete [id]",
		Short: "Delete an API key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := service.Delete(args[0]); err != nil {
				return err
			}
			return printJSON(map[string]bool{"deleted": true})
		},
	}

	cmd.AddCommand(listCmd, createCmd, deleteCmd)
	return cmd
}
