package cmd

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
				t.Errorf("Test parameter: %s \nunexpected error message: %s \nreturned error message: %s\n",
					tc.args, tc.err_msg, fmt.Sprint(err))
			}
		}

		if err == nil && tc.out_check {
			if !strings.Contains(fmt.Sprint(out), tc.out) {
				t.Errorf("Test parameter: %s \nunexpected output message: %s \n returned output message: %s\n",
					tc.args, tc.out, fmt.Sprint(out))
			}
		}
	}
}

func TestOtherCmd(t *testing.T) {
	testcases := []Testcase{
		{[]string{}, false, "", ""},
		{[]string{"wrong"}, true, "unknown command \"wrong\" for \"tks\"", ""},
		{[]string{"wrong", "cmd"}, true, "unknown command \"wrong\" for \"tks\"", ""},

		{[]string{"completion"}, false, "", ""},
		{[]string{"completion", "bash"}, false, "", ""},
		{[]string{"completion", "fish"}, false, "", ""},
		{[]string{"completion", "powershell"}, false, "", ""},
		{[]string{"completion", "zsh"}, false, "", ""},

		{[]string{"--config"}, true, "flag needs an argument: --config", ""},
		{[]string{"-t"}, false, "", ""},

		{[]string{"endpoint"}, true, "", ""},
		{[]string{"endpoint", "register"}, true, "", ""},
		{[]string{"endpoint", "wrong"}, true, "", ""},

		{[]string{"help"}, false, "", ""},
		{[]string{"-t"}, false, "", ""},
		{[]string{"-h"}, false, "", ""},

		{[]string{"register"}, true, "unknown command \"register\" for \"tks\"", "Run 'tks --help' for usage."},
		{[]string{"abcd"}, true, "unknown command \"abcd\" for \"tks\"", ""},
	}

	doTest(t, testcases)
}

func TestClusterCmd(t *testing.T) {
	testcases := []Testcase{

		{[]string{"cluster"}, true, "", ""},
		{[]string{"cluster", "wrong"}, true, "", ""},
		{[]string{"cluster", "create"}, true, "Usage: tks cluster create <CLUSTERNAME>", ""},
		{[]string{"cluster", "create", "--config"}, true, "flag needs an argument: --config", ""},
		{[]string{"cluster", "create", "--config", "xx"}, true, "Usage: tks cluster create <CLUSTERNAME>", ""},
		{[]string{"cluster", "create", "-h"}, false, "", "Create a TKS Cluster to AWS."},
		{[]string{"cluster", "create", "--contract-id"}, true, "flag needs an argument: --contract-id", ""},
		{[]string{"cluster", "create", "--contract-id", "1234"}, false, "", ""},
		{[]string{"cluster", "create", "--csp-id"}, false, "flag needs an argument: --csp-id", ""},
		{[]string{"cluster", "create", "--contract-id", "--csp-id"}, false, "", ""},
		{[]string{"cluster", "create", "--contract-id", "abcd-efg", "--csp-id"}, true, "flag needs an argument: --csp-id", ""},
		{[]string{"cluster", "create", "--contract-id", "--csp-id", "2345234", "hihi"}, false, "", "Create a TKS Cluster to AWS"},
		{[]string{"cluster", "create", "--contract-id", "--csp-id", "2345234", "hihi", "--config"}, true, "flag needs an argument: --config", ""},
		{[]string{"cluster", "create", "--contract-id", "abcd-efg", "--csp-id", "2345234", "hihi"}, false, "", "Create a TKS Cluster to AWS"},

		{[]string{"cluster", "list"}, true, "", ""},
		{[]string{"cluster", "list", "-h"}, false, "", "A longer description that spans multiple lines and likely contains examples"},
		{[]string{"cluster", "list", "--config", "xx"}, false, "", "A longer description that spans multiple lines and likely contains examples"},
		{[]string{"cluster", "list", "--config"}, true, "flag needs an argument: --config", ""},
	}

	doTest(t, testcases)
}

func TestServiceCmd(t *testing.T) {
	testcases := []Testcase{

		{[]string{"service"}, true, "", ""},
		{[]string{"service", "-h"}, false, "", ""},
		{[]string{"service", "create"}, true, "required flag(s) \"cluster-id\", \"service-name\" not set", ""},
		{[]string{"service", "create", "--cluster-id"}, true, "flag needs an argument: --cluster-id", ""},
		{[]string{"service", "create", "--cluster-id", "--service-name"}, true, "required flag(s) \"service-name\" not set", ""},
		{[]string{"service", "create", "--cluster-id", "aaa", "--service-name"}, true, "flag needs an argument: --service-name", ""},
		{[]string{"service", "create", "--cluster-id", "aaa", "--service-name", "LMA"}, true, "", ""},
	}

	doTest(t, testcases)
}
