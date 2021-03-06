package cmd

import (
	"kool-dev/kool/cmd/builder"
	"strconv"

	"github.com/spf13/cobra"
)

// KoolLogsFlags holds the flags for the logs command
type KoolLogsFlags struct {
	Tail   int
	Follow bool
}

// KoolLogs holds handlers and functions to implement the logs command logic
type KoolLogs struct {
	DefaultKoolService
	Flags *KoolLogsFlags

	logs builder.Command
}

func init() {
	var (
		logs    = NewKoolLogs()
		logsCmd = NewLogsCommand(logs)
	)

	rootCmd.AddCommand(logsCmd)

	logsCmd.Flags().IntVarP(&logs.Flags.Tail, "tail", "t", 25, "Number of lines to show from the end of the logs for each container. For value equal to 0, all lines will be shown.")
	logsCmd.Flags().BoolVarP(&logs.Flags.Follow, "follow", "f", false, "Follow log output.")
}

// NewKoolLogs creates a new handler for logs logic
func NewKoolLogs() *KoolLogs {
	return &KoolLogs{
		*newDefaultKoolService(),
		&KoolLogsFlags{25, false},
		builder.NewCommand("docker-compose", "logs"),
	}
}

// Execute runs the logs logic with incoming arguments.
func (l *KoolLogs) Execute(args []string) (err error) {
	if l.Flags.Tail == 0 {
		l.logs.AppendArgs("--tail", "all")
	} else {
		l.logs.AppendArgs("--tail", strconv.Itoa(l.Flags.Tail))
	}

	if l.Flags.Follow {
		l.logs.AppendArgs("--follow")
	}

	err = l.logs.Interactive(args...)
	return
}

// NewLogsCommand initializes new kool logs command
func NewLogsCommand(logs *KoolLogs) *cobra.Command {
	return &cobra.Command{
		Use:   "logs [options] [service...]",
		Short: "Displays log output from services.",
		Run: func(cmd *cobra.Command, args []string) {
			logs.SetWriter(cmd.OutOrStdout())

			if err := logs.Execute(args); err != nil {
				logs.Error(err)
				logs.Exit(1)
			}
		},
	}
}
