package cmd

import (
	"context"
	"fmt"

	"github.com/charmbracelet/lipgloss/table"
	"github.com/nuki-io/nuki-cli/pkg/bleflows"
	"github.com/spf13/cobra"
)

var keypadCmd = &cobra.Command{
	Use:   "keypad",
	Short: "Manage keypad codes on the device",
}

var keypadListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all keypad codes",
	PreRunE: mustDeviceId,
	RunE: func(cmd *cobra.Command, args []string) error {
		return withAuthenticatedFlow(func(ctx context.Context, flow *bleflows.Flow) error {
			codes, err := flow.GetKeypadCodes(ctx, 0, 100)
			if err != nil {
				return fmt.Errorf("failed to list keypad codes: %w", err)
			}
			if outputFormat == "json" {
				return printJSON(codes)
			}
			t := table.New().Headers("ID", "Code", "Name", "Enabled", "Lock Count", "Created")
			for _, c := range codes {
				t = t.Row(
					fmt.Sprintf("%d", c.CodeID),
					fmt.Sprintf("%d", c.Code),
					c.Name,
					fmt.Sprintf("%t", c.Enabled),
					fmt.Sprintf("%d", c.LockCount),
					c.DateCreated.Local().Format("2006-01-02 15:04"),
				)
			}
			fmt.Println(t)
			return nil
		})
	},
}

var (
	kpCode uint32
	kpName string
)

var keypadAddCmd = &cobra.Command{
	Use:     "add",
	Short:   "Add a keypad code",
	PreRunE: mustDeviceId,
	RunE: func(cmd *cobra.Command, args []string) error {
		return withAuthenticatedFlow(func(ctx context.Context, flow *bleflows.Flow) error {
			id, err := flow.AddKeypadCode(ctx, kpCode, kpName)
			if err != nil {
				return fmt.Errorf("failed to add keypad code: %w", err)
			}
			fmt.Printf("Keypad code added with ID %d.\n", id)
			return nil
		})
	},
}

var kpRemoveID uint16

var keypadRemoveCmd = &cobra.Command{
	Use:     "remove",
	Short:   "Remove a keypad code",
	PreRunE: mustDeviceId,
	RunE: func(cmd *cobra.Command, args []string) error {
		return withAuthenticatedFlow(func(ctx context.Context, flow *bleflows.Flow) error {
			if err := flow.RemoveKeypadCode(ctx, kpRemoveID); err != nil {
				return fmt.Errorf("failed to remove keypad code: %w", err)
			}
			fmt.Printf("Keypad code %d removed.\n", kpRemoveID)
			return nil
		})
	},
}

var (
	kpUpdateID      uint16
	kpUpdateEnabled bool
)

var keypadUpdateCmd = &cobra.Command{
	Use:     "update",
	Short:   "Update a keypad code",
	PreRunE: mustDeviceId,
	RunE: func(cmd *cobra.Command, args []string) error {
		return withAuthenticatedFlow(func(ctx context.Context, flow *bleflows.Flow) error {
			if err := flow.UpdateKeypadCode(ctx, kpUpdateID, kpCode, kpName, kpUpdateEnabled); err != nil {
				return fmt.Errorf("failed to update keypad code: %w", err)
			}
			fmt.Printf("Keypad code %d updated.\n", kpUpdateID)
			return nil
		})
	},
}

func init() {
	bleCmd.AddCommand(keypadCmd)
	keypadCmd.AddCommand(keypadListCmd)
	keypadCmd.AddCommand(keypadAddCmd)
	keypadCmd.AddCommand(keypadRemoveCmd)
	keypadCmd.AddCommand(keypadUpdateCmd)

	keypadAddCmd.Flags().Uint32Var(&kpCode, "code", 0, "Numeric keypad code")
	keypadAddCmd.Flags().StringVar(&kpName, "name", "", "Name for this code (max 20 chars)")
	_ = keypadAddCmd.MarkFlagRequired("code")
	_ = keypadAddCmd.MarkFlagRequired("name")

	keypadRemoveCmd.Flags().Uint16Var(&kpRemoveID, "code-id", 0, "Code ID to remove")
	_ = keypadRemoveCmd.MarkFlagRequired("code-id")

	keypadUpdateCmd.Flags().Uint16Var(&kpUpdateID, "code-id", 0, "Code ID to update")
	keypadUpdateCmd.Flags().Uint32Var(&kpCode, "code", 0, "New numeric keypad code")
	keypadUpdateCmd.Flags().StringVar(&kpName, "name", "", "New name (max 20 chars)")
	keypadUpdateCmd.Flags().BoolVar(&kpUpdateEnabled, "enabled", true, "Enable or disable this code")
	_ = keypadUpdateCmd.MarkFlagRequired("code-id")
	_ = keypadUpdateCmd.MarkFlagRequired("code")
	_ = keypadUpdateCmd.MarkFlagRequired("name")
}
