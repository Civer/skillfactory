package projects

import (
	"fmt"

	"github.com/petervogelmann/skillfactory/skills/vikunja/client"
	"github.com/spf13/cobra"
)

// RegisterCommands creates and returns the projects command group
func RegisterCommands(c *client.Client, printJSON func(interface{}) error) *cobra.Command {
	service := NewService(c)

	cmd := &cobra.Command{
		Use:   "projects",
		Short: "Manage projects",
	}

	// list
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all projects",
		RunE: func(cmd *cobra.Command, args []string) error {
			projects, err := service.List()
			if err != nil {
				return err
			}
			return printJSON(ToLeanSlice(projects))
		},
	}

	// get
	getCmd := &cobra.Command{
		Use:   "get [id]",
		Short: "Get a project by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := parseID(args[0])
			if err != nil {
				return err
			}
			project, err := service.Get(id)
			if err != nil {
				return err
			}
			return printJSON(project.ToLean())
		},
	}

	cmd.AddCommand(listCmd, getCmd)
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
