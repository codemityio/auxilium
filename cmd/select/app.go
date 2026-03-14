package slct

import (
	"fmt"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/urfave/cli/v2"
)

type input struct {
	Name  string
	Value string
}

var App = cli.Command{ //nolint:gochecknoglobals,exhaustruct
	Name:         "select",
	Aliases:      nil,
	Usage:        "",
	UsageText:    "",
	Description:  "Configurable shell select",
	ArgsUsage:    "",
	Category:     "",
	BashComplete: nil,
	Before:       nil,
	After:        nil,
	Action: func(ctx *cli.Context) error {
		list := ctx.String("list")
		label := ctx.String("label")
		delimiter := ctx.String("delimiter")
		selectNameLabel := ctx.String("select-name-label")
		selectValueLabel := ctx.String("select-value-label")
		exitName := ctx.String("exit-name")
		exitValue := ctx.String("exit-value")

		if strings.TrimSpace(list) == "" {
			return nil
		}

		lst := strings.Split(list, "\n")

		var options []input

		if exitName != "" {
			options = append(options, input{
				Name:  exitName,
				Value: exitValue,
			})
		}

		for i := range lst {
			ll := strings.Split(lst[i], delimiter)

			options = append(options, input{
				Name:  ll[0],
				Value: strings.Join(ll[1:], delimiter),
			})
		}

		templates := &promptui.SelectTemplates{ //nolint:exhaustruct
			Label:    "{{ . | yellow }}{{ \":\" | yellow }}",
			Active:   "> {{ .Name | cyan | bold }}",
			Inactive: "  {{ .Name | cyan }}",
			Selected: "{{ \"Selected:\" | bold }} {{ .Name }}",
			Details: fmt.Sprintf(`
{{ "Select:" | yellow }}
{{ "%s:" | faint }} {{ .Name }}
{{ if .Value }}{{ "%s:" | faint }} {{ .Value }}{{ end }}`, selectNameLabel, selectValueLabel),
		}

		searcher := func(input string, index int) bool {
			option := options[index]
			name := strings.ReplaceAll(strings.ToLower(option.Name), " ", "")
			input = strings.ReplaceAll(strings.ToLower(input), " ", "")

			return strings.Contains(name, input)
		}

		prompt := promptui.Select{ //nolint:exhaustruct
			Label:     label,
			Items:     options,
			Templates: templates,
			Size:      selectSize,
			Searcher:  searcher,
			Stdout:    os.Stderr,
		}

		index, _, err := prompt.Run()
		if err != nil {
			return nil
		}

		if index == 0 {
			if _, e := os.Stdout.WriteString(options[index].Value); e != nil {
				return fmt.Errorf("%w: unable to write string", e)
			}

			return nil
		}

		var output string

		if options[index].Value == "" {
			output = options[index].Name
		}

		if options[index].Value != "" {
			output = options[index].Name + delimiter + options[index].Value
		}

		if _, e := os.Stdout.WriteString(output); e != nil {
			return fmt.Errorf("%w: unable to write string", e)
		}

		return nil
	},
	OnUsageError: nil,
	Flags: []cli.Flag{
		&cli.StringFlag{ //nolint:exhaustruct
			Name:     "list",
			Usage:    "a list of options (each option as a new line)",
			Value:    "",
			Required: true,
		},
		&cli.StringFlag{ //nolint:exhaustruct
			Name:  "label",
			Usage: "a label to be used in the select",
			Value: "Choose",
		},
		&cli.StringFlag{ //nolint:exhaustruct
			Name:  "delimiter",
			Usage: "an option delimiter used to split option names from the remaining content to be used as description",
			Value: " ",
		},
		&cli.StringFlag{ //nolint:exhaustruct
			Name:  "select-name-label",
			Usage: "a label for the details menu name field label",
			Value: "Name",
		},
		&cli.StringFlag{ //nolint:exhaustruct
			Name:  "select-value-label",
			Usage: "a label for the details menu value field label",
			Value: "Value",
		},
		&cli.StringFlag{ //nolint:exhaustruct
			Name:  "exit-name",
			Usage: "a name to be used for the exit option",
			Value: ".",
		},
		&cli.StringFlag{ //nolint:exhaustruct
			Name:  "exit-value",
			Usage: "a value to be returned for the exit option",
			Value: "",
		},
	},
	Subcommands: nil,
}
