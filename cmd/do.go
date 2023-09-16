package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewDoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "do",
		Example: `echo '[{"x":1},{"x":2}]' | jo do echo \$x 
  jo do --in example.json echo \$x
`,
		Short: "Do invokes a command for each object in the json input",
		Long: `Do invokes a command for each object in the json input with shell variable for each field in the object.

Grep to show last 10 errors in all container logs:
    ns=kube-system
    kubectl -n $ns get po -o json | jq '[ .items[] | .metadata.name as $p | (.spec.containers + .spec.initContainers)[] | {"p": $p, "c": .name} ]' | jo do "echo '----' \$p \$c; kubectl -n $ns logs \$p -c \$c | grep -i error | tail -n 10"

Write all container logs to files:
Note that the command passed to jo is quoted (without the quotes all logs will go to a single file)
    kubectl -n $ns get po -o json | jq '[ .items[] | .metadata.name as $p | (.spec.containers + .spec.initContainers)[] | {"p": $p, "c": .name} ]' | jo do "kubectl -n $ns logs --ignore-errors=true \$p -c \$c >$out_dir/\$p+\$c.log"

Tips:
Most of the time jo do flags and command flag are separated ok but when flags are ambiguous add -- like in: jo do --in example.json -- echo -n \$x
Flags can be set with environment variables as well (flag foo-bar maps to environment variable JO_FOO_BAR).
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error

			if len(args) == 0 {
				return fmt.Errorf("command is missing")
			}

			var m []string
			for _, k := range []string{ /* add required flags here */ } {
				if !viper.IsSet(k) {
					m = append(m, k)
				}
			}
			if len(m) > 0 {
				return fmt.Errorf("one or more required flags not set: %v", m)
			}

			in := os.Stdin
			inn := "stdin"
			if s := viper.GetString("in"); s != "" && s != "-" {
				var err error
				in, err = os.Open(s)
				if err != nil {
					return fmt.Errorf("open %s: %w", s, err)
				}
				inn = s
			}
			
			b := new(bytes.Buffer)
			_, err = b.ReadFrom(in)
			if err != nil {
				return fmt.Errorf("reading %s: %w", inn, err)
			}

			var objs []map[string]any
			d := json.NewDecoder(b)
			err = d.Decode(&objs)
			if err != nil {
				var e *json.UnmarshalTypeError
				if !(errors.As(err, &e) && e.Offset == 1) {
					return fmt.Errorf("reading %s: %w", inn, err)
				}

				// we're reading jsonl and at this point the first char of 'in' is already consumed.
				
				b.UnreadRune()
				d = json.NewDecoder(b)
				for {
					var obj map[string]any
					err := d.Decode(&obj)
					if err == io.EOF {
						break
					}
					if err != nil {
						return fmt.Errorf("reading %s jsonl: %w", inn, err)
					}
					objs = append(objs, obj)
				}
			}

			for _, obj := range objs {
				c := exec.Command("sh", "-c", strings.Join(args, " "))
				c.Stdout = os.Stdout
				c.Stderr = os.Stderr

				env := os.Environ()
				for k, v := range obj {
					env = append(env, fmt.Sprintf("%s=%v", k, v))
				}
				c.Env = env

				err = c.Run()
				if err != nil {
					return fmt.Errorf("run with %v: %w", obj, err)
				}
			}

			return nil
		},
	}

	cmd.Flags().String("in", "-",
		"Input file")

	viper.SetEnvPrefix("JO")
	_ = viper.BindPFlags(cmd.Flags())

	return cmd
}
