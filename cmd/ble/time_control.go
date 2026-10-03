package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss/table"
	"github.com/nuki-io/nuki-cli/pkg/blecommands"
	"github.com/nuki-io/nuki-cli/pkg/bleflows"
	"github.com/spf13/cobra"
)

var timeControlCmd = &cobra.Command{
	Use:   "time-control",
	Short: "Manage time control entries on the device",
}

var timeControlListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all time control entries",
	PreRunE: mustDeviceId,
	RunE: func(cmd *cobra.Command, args []string) error {
		return withAuthenticatedFlow(func(ctx context.Context, flow *bleflows.Flow) error {
			entries, err := flow.GetTimeControlEntries(ctx)
			if err != nil {
				return fmt.Errorf("failed to list time control entries: %w", err)
			}
			if outputFormat == "json" {
				return printJSON(entries)
			}
			t := table.New().Headers("ID", "Enabled", "Weekdays", "Time", "Action")
			for _, e := range entries {
				t = t.Row(
					fmt.Sprintf("%d", e.EntryID),
					fmt.Sprintf("%t", e.Enabled),
					weekdayNames(e.Weekdays),
					fmt.Sprintf("%02d:%02d", e.TimeHour, e.TimeMinute),
					e.LockAction.String(),
				)
			}
			fmt.Println(t)
			return nil
		})
	},
}

var (
	tcWeekdays string
	tcTime     string
	tcAction   string
)

var timeControlAddCmd = &cobra.Command{
	Use:     "add",
	Short:   "Add a time control entry",
	PreRunE: mustDeviceId,
	RunE: func(cmd *cobra.Command, args []string) error {
		weekdays, err := parseWeekdays(tcWeekdays)
		if err != nil {
			return err
		}
		hour, minute, err := parseTime(tcTime)
		if err != nil {
			return err
		}
		action, err := parseLockAction(tcAction)
		if err != nil {
			return err
		}
		return withAuthenticatedFlow(func(ctx context.Context, flow *bleflows.Flow) error {
			id, err := flow.AddTimeControlEntry(ctx, weekdays, hour, minute, action)
			if err != nil {
				return fmt.Errorf("failed to add time control entry: %w", err)
			}
			fmt.Printf("Time control entry added with ID %d.\n", id)
			return nil
		})
	},
}

var tcRemoveID uint8

var timeControlRemoveCmd = &cobra.Command{
	Use:     "remove",
	Short:   "Remove a time control entry",
	PreRunE: mustDeviceId,
	RunE: func(cmd *cobra.Command, args []string) error {
		return withAuthenticatedFlow(func(ctx context.Context, flow *bleflows.Flow) error {
			if err := flow.RemoveTimeControlEntry(ctx, tcRemoveID); err != nil {
				return fmt.Errorf("failed to remove time control entry: %w", err)
			}
			fmt.Printf("Time control entry %d removed.\n", tcRemoveID)
			return nil
		})
	},
}

var (
	tcUpdateID      uint8
	tcUpdateEnabled bool
)

var timeControlUpdateCmd = &cobra.Command{
	Use:     "update",
	Short:   "Update a time control entry",
	PreRunE: mustDeviceId,
	RunE: func(cmd *cobra.Command, args []string) error {
		weekdays, err := parseWeekdays(tcWeekdays)
		if err != nil {
			return err
		}
		hour, minute, err := parseTime(tcTime)
		if err != nil {
			return err
		}
		action, err := parseLockAction(tcAction)
		if err != nil {
			return err
		}
		return withAuthenticatedFlow(func(ctx context.Context, flow *bleflows.Flow) error {
			if err := flow.UpdateTimeControlEntry(ctx, tcUpdateID, tcUpdateEnabled, weekdays, hour, minute, action); err != nil {
				return fmt.Errorf("failed to update time control entry: %w", err)
			}
			fmt.Printf("Time control entry %d updated.\n", tcUpdateID)
			return nil
		})
	},
}

// weekdayNames converts a weekday bitmask to human-readable string.
// Bit layout (LSB to MSB): SU SA FR TH WE TU MO (bit 0 = Sunday, bit 6 = Monday)
func weekdayNames(mask byte) string {
	names := []string{"Su", "Sa", "Fr", "Th", "We", "Tu", "Mo"}
	var days []string
	for i, name := range names {
		if mask&(1<<i) != 0 {
			days = append(days, name)
		}
	}
	if len(days) == 0 {
		return "all"
	}
	return strings.Join(days, ",")
}

// parseWeekdays parses a comma-separated list of day abbreviations into a bitmask.
// Accepted values: mo,tu,we,th,fr,sa,su (case-insensitive). Empty string means all days (0x00).
func parseWeekdays(s string) (byte, error) {
	if s == "" || strings.EqualFold(s, "all") {
		return 0x00, nil
	}
	dayBits := map[string]byte{
		"mo": 0x40, "tu": 0x20, "we": 0x10, "th": 0x08, "fr": 0x04, "sa": 0x02, "su": 0x01,
	}
	var mask byte
	for _, part := range strings.Split(s, ",") {
		bit, ok := dayBits[strings.ToLower(strings.TrimSpace(part))]
		if !ok {
			return 0, fmt.Errorf("unknown day %q, use: mo,tu,we,th,fr,sa,su", part)
		}
		mask |= bit
	}
	return mask, nil
}

// parseTime parses "HH:MM" into hour and minute bytes.
func parseTime(s string) (hour, minute byte, err error) {
	var h, m int
	if _, err = fmt.Sscanf(s, "%d:%d", &h, &m); err != nil || h > 23 || m > 59 {
		return 0, 0, fmt.Errorf("invalid time %q, expected HH:MM", s)
	}
	return byte(h), byte(m), nil
}

// parseLockAction converts a string to an Action value.
func parseLockAction(s string) (blecommands.Action, error) {
	switch strings.ToLower(s) {
	case "unlock":
		return blecommands.Unlock, nil
	case "lock":
		return blecommands.Lock, nil
	case "unlatch":
		return blecommands.Unlatch, nil
	case "lock-n-go", "lockandgo":
		return blecommands.LockAndGo, nil
	case "full-lock", "fulllock":
		return blecommands.FullLock, nil
	default:
		return 0, fmt.Errorf("unknown action %q, use: unlock, lock, unlatch, lock-n-go, full-lock", s)
	}
}

func init() {
	bleCmd.AddCommand(timeControlCmd)
	timeControlCmd.AddCommand(timeControlListCmd)
	timeControlCmd.AddCommand(timeControlAddCmd)
	timeControlCmd.AddCommand(timeControlRemoveCmd)
	timeControlCmd.AddCommand(timeControlUpdateCmd)

	timeControlAddCmd.Flags().StringVar(&tcWeekdays, "weekdays", "", "Comma-separated weekdays (mo,tu,we,th,fr,sa,su) or 'all'")
	timeControlAddCmd.Flags().StringVar(&tcTime, "time", "", "Time in HH:MM format")
	timeControlAddCmd.Flags().StringVar(&tcAction, "action", "", "Lock action: unlock, lock, unlatch, lock-n-go, full-lock")
	_ = timeControlAddCmd.MarkFlagRequired("time")
	_ = timeControlAddCmd.MarkFlagRequired("action")

	timeControlRemoveCmd.Flags().Uint8Var(&tcRemoveID, "entry-id", 0, "Entry ID to remove")
	_ = timeControlRemoveCmd.MarkFlagRequired("entry-id")

	timeControlUpdateCmd.Flags().Uint8Var(&tcUpdateID, "entry-id", 0, "Entry ID to update")
	timeControlUpdateCmd.Flags().BoolVar(&tcUpdateEnabled, "enabled", true, "Enable or disable this entry")
	timeControlUpdateCmd.Flags().StringVar(&tcWeekdays, "weekdays", "", "Comma-separated weekdays (mo,tu,we,th,fr,sa,su) or 'all'")
	timeControlUpdateCmd.Flags().StringVar(&tcTime, "time", "", "Time in HH:MM format")
	timeControlUpdateCmd.Flags().StringVar(&tcAction, "action", "", "Lock action: unlock, lock, unlatch, lock-n-go, full-lock")
	_ = timeControlUpdateCmd.MarkFlagRequired("entry-id")
	_ = timeControlUpdateCmd.MarkFlagRequired("time")
	_ = timeControlUpdateCmd.MarkFlagRequired("action")
}
