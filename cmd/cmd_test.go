package cmd

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/matryer/is"
	"github.com/spf13/cobra"
)

func TestRootCmd(t *testing.T) {
	is := is.New(t)

	// root := &cobra.Command{Use: "root"}
	// cmd.RootCmdFlags(root)

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

func TestOtherCmd(t *testing.T) {
	is := is.New(t)

	testcases := []struct {
		args      []string
		out_check bool
		err       error
		out       string
	}{
		{[]string{}, false, nil, ""},
		{[]string{"wrong"}, true, errors.New("unknown command \"wrong\" for \"tks\""), ""},
		// {[]string{"wrong", "cmd"}, true, errors.New("unknown command \"wrong\" for \"tks\""), ""},

		{[]string{"completion"}, false, nil, ""},
		{[]string{"completion", "bash"}, false, nil, ""},
		{[]string{"completion", "fish"}, false, nil, ""},
		{[]string{"completion", "powershell"}, false, nil, ""},
		{[]string{"completion", "zsh"}, false, nil, ""},

		{[]string{"--config"}, true, errors.New("flag needs an argument: --config"), ""},
		{[]string{"-t"}, false, nil, ""},

		{[]string{"endpoint"}, true, nil, ""},
		{[]string{"endpoint", "register"}, true, nil, ""},
		{[]string{"endpoint", "wrong"}, true, nil, ""},

		{[]string{"help"}, false, nil, ""},
		{[]string{"-t"}, false, nil, ""},
		{[]string{"-h"}, false, nil, ""},

		{[]string{"register"}, true, errors.New("unknown command \"register\" for \"tks\""), "Run 'tks --help' for usage."},
		{[]string{"abcd"}, true, errors.New("unknown command \"abcd\" for \"tks\""), ""},
	}

	for _, tc := range testcases {
		out, err := execute(t, rootCmd, tc.args...)

		is.Equal(tc.err, err)

		if tc.err == nil && tc.out_check {
			is.Equal(tc.out, out)
		}
	}
}

func TestClusterCmd(t *testing.T) {
	is := is.New(t)

	testcases := []struct {
		args      []string
		out_check bool
		err       error
		out       string
	}{
		{[]string{"cluster"}, true, nil, ""},
		{[]string{"cluster", "wrong"}, true, nil, ""},
		{[]string{"cluster", "create"}, true, errors.New("required flag(s) \"contract-id\", \"csp-id\" not set"), ""},
		{[]string{"cluster", "create", "--config"}, true, errors.New("flag needs an argument: --config"), ""},
		{[]string{"cluster", "create", "--config", "xx"}, true, errors.New("required flag(s) \"contract-id\", \"csp-id\" not set"), ""},
		{[]string{"cluster", "create", "-h"}, false, nil, "Create a TKS Cluster to AWS."},
		{[]string{"cluster", "create", "--contract-id"}, true, errors.New("flag needs an argument: --contract-id"), ""},
		{[]string{"cluster", "create", "--contract-id", "1234"}, false, nil, ""},
		{[]string{"cluster", "create", "--csp-id"}, false, errors.New("flag needs an argument: --csp-id"), ""},
		{[]string{"cluster", "create", "--contract-id", "--csp-id"}, false, nil, ""},
		{[]string{"cluster", "create", "--contract-id", "abcd-efg", "--csp-id"}, true, errors.New("flag needs an argument: --csp-id"), ""},
		{[]string{"cluster", "create", "--contract-id", "--csp-id", "2345234", "hihi"}, false, nil, "Create a TKS Cluster to AWS"},
		{[]string{"cluster", "create", "--contract-id", "--csp-id", "2345234", "hihi", "--config"}, true, errors.New("flag needs an argument: --config"), ""},
		{[]string{"cluster", "create", "--contract-id", "abcd-efg", "--csp-id", "2345234", "hihi"}, false, nil, "Create a TKS Cluster to AWS"},

		{[]string{"cluster", "list"}, true, nil, ""},
		{[]string{"cluster", "list", "-h"}, false, nil, "A longer description that spans multiple lines and likely contains examples"},
		{[]string{"cluster", "list", "--config", "xx"}, false, nil, "A longer description that spans multiple lines and likely contains examples"},
		{[]string{"cluster", "list", "--config"}, true, errors.New("flag needs an argument: --config"), ""},
	}

	for _, tc := range testcases {
		out, err := execute(t, rootCmd, tc.args...)

		is.Equal(tc.err, err)

		if tc.err == nil && tc.out_check {
			is.Equal(tc.out, out)
		}
	}
}

func TestServiceCmd(t *testing.T) {
	is := is.New(t)

	testcases := []struct {
		args      []string
		out_check bool
		err       error
		out       string
	}{
		{[]string{"service"}, true, nil, ""},
		{[]string{"service", "-h"}, false, nil, ""},
		{[]string{"service", "create"}, true, errors.New("required flag(s) \"cluster-id\", \"service-name\" not set"), ""},
		{[]string{"service", "create", "--cluster-id"}, true, errors.New("flag needs an argument: --cluster-id"), ""},
		// {[]string{"service", "create", "--cluster-id", "--service-name"}, true, errors.New("required flag(s) \"service-name\" not set"), ""},
		// {[]string{"service", "create", "--cluster-id", "aaa", "--service-name"}, true, errors.New("flag needs an argument: --service-name"), ""},
		{[]string{"service", "create", "--cluster-id", "aaa", "--service-name", "LMA"}, true, nil, ""},
	}

	for _, tc := range testcases {
		out, err := execute(t, rootCmd, tc.args...)

		is.Equal(tc.err, err)

		if tc.err == nil && tc.out_check {
			is.Equal(tc.out, out)
		}
	}
}
