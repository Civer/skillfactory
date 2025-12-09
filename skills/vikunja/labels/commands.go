package labels

import (
	"fmt"

	"github.com/petervogelmann/skillfactory/skills/vikunja/client"
	"github.com/spf13/cobra"
)

// RegisterCommands creates and returns the labels command group
func RegisterCommands(c *client.Client, printJSON func(interface{}) error) *cobra.Command {
	service := NewService(c)

	cmd := &cobra.Command{
		Use:   "labels",
		Short: "Manage labels",
	}

	// list
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all labels",
		RunE: func(cmd *cobra.Command, args []string) error {
			labels, err := service.List()
			if err != nil {
				return err
			}
			return printJSON(ToLeanSlice(labels))
		},
	}

	// get
	getCmd := &cobra.Command{
		Use:   "get [id]",
		Short: "Get a label by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := parseID(args[0])
			if err != nil {
				return err
			}
			label, err := service.Get(id)
			if err != nil {
				return err
			}
			return printJSON(label.ToLean())
		},
	}

	// create
	var createTitle string
	var createColor string
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new label",
		RunE: func(cmd *cobra.Command, args []string) error {
			if createTitle == "" {
				return fmt.Errorf("--title is required")
			}
			req := CreateLabelRequest{
				Title:    createTitle,
				HexColor: createColor,
			}
			label, err := service.Create(req)
			if err != nil {
				return err
			}
			return printJSON(label.ToLean())
		},
	}
	createCmd.Flags().StringVarP(&createTitle, "title", "t", "", "Label title (required)")
	createCmd.Flags().StringVar(&createColor, "color", "", "Hex color (e.g., #ff0000)")

	// update
	var updateTitle string
	var updateColor string
	updateCmd := &cobra.Command{
		Use:   "update [id]",
		Short: "Update a label",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := parseID(args[0])
			if err != nil {
				return err
			}
			req := UpdateLabelRequest{
				Title:    updateTitle,
				HexColor: updateColor,
			}
			label, err := service.Update(id, req)
			if err != nil {
				return err
			}
			return printJSON(label.ToLean())
		},
	}
	updateCmd.Flags().StringVarP(&updateTitle, "title", "t", "", "New title")
	updateCmd.Flags().StringVar(&updateColor, "color", "", "New hex color")

	// delete
	deleteCmd := &cobra.Command{
		Use:   "delete [id]",
		Short: "Delete a label",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := parseID(args[0])
			if err != nil {
				return err
			}
			if err := service.Delete(id); err != nil {
				return err
			}
			return printJSON(map[string]bool{"deleted": true})
		},
	}

	cmd.AddCommand(listCmd, getCmd, createCmd, updateCmd, deleteCmd)
	return cmd
}

func parseID(s string) (int64, error) {
	var id int64
	_, err := fmt.Sscanf(s, "%d", &id)
	if err != nil {
		return 0, fmt.Errorf("invalid ID: %s", s)
	}
	return id, nil
}
