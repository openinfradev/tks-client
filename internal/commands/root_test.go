package commands

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/matryer/is"
	"github.com/spf13/cobra"
)

func TestRootCmd(t *testing.T) {
	is := is.New(t)
	err := rootCmd.Execute()
	is.NoErr(err)
}

func execute(t *testing.T, c *cobra.Command, args ...string) (string, error) {
	t.Helper()

	buf := new(bytes.Buffer)
	c.SetOut(buf)
	c.SetErr(buf)
	c.SetArgs(args)

	err := c.Execute()
	return strings.TrimSpace(buf.String()), err
}

type Testcase struct {
	args      []string
	out_check bool
	err_msg   string
	out       string
}

func doTest(t *testing.T, testcases []Testcase) {
	for _, tc := range testcases {
		fmt.Printf("Test with arguments: %s\n", tc.args)
		out, err := execute(t, rootCmd, tc.args...)

		if err != nil {
			if !strings.Contains(err.Error(), tc.err_msg) {
				t.Errorf("Test parameter: %s \nunexpected error message: %s \nreturned error message: %s",
					tc.args, tc.err_msg, fmt.Sprint(err))
			}
		}

		if err == nil && tc.out_check {
			if !strings.Contains(fmt.Sprint(out), tc.out) {
				t.Errorf("Test parameter: %s \nunexpected output message: %s \n returned output message: %s",
					tc.args, tc.out, fmt.Sprint(out))
			}
		}
	}
}

func TestContractCmd(t *testing.T) {
	testcases := []Testcase{
		{[]string{"contract"}, true, "", ""},
		{[]string{"contract", "create"}, true, "Usage: tksadmin contract create <CONTRACT NAME>", ""},
		{[]string{"contract", "create", "cli-unit-test"}, true, "Contract Name:  cli-unit-test\nYou must specify tksContractUrl at config file", ""},
	}

	doTest(t, testcases)
}

func TestOtherCmd(t *testing.T) {
	testcases := []Testcase{
		{[]string{}, false, "", ""},
		{[]string{"wrong"}, true, "unknown command", ""},
		{[]string{"wrong", "cmd"}, true, "unknown command", ""},

		{[]string{"completion"}, false, "", ""},
		{[]string{"completion", "bash"}, false, "", ""},
		{[]string{"completion", "fish"}, false, "", ""},
		{[]string{"completion", "powershell"}, false, "", ""},
		{[]string{"completion", "zsh"}, false, "", ""},

		{[]string{"--config"}, true, "flag needs an argument:", ""},
		{[]string{"-t"}, false, "", ""},
	}

	doTest(t, testcases)
}
