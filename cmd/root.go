/*
Package cmd

Copyright © 2024 Pete Wall <pete@petewall.net>
*/
package cmd

import (
	"fmt"
	"github.com/go-logfmt/logfmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

const TimeKey = "time"
const MessageKey = "msg"

func dateParts() (string, string) {
	return TimeKey, time.Now().Format(time.RFC3339)
}

func labelParts(label string) (string, string, error) {
	parts := strings.Split(label, "=")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("label must be in the format <key>=<value>")
	}
	return parts[0], parts[1], nil
}

var rootCmd = &cobra.Command{
	Use:   "logfmtwrite",
	Short: "Write a statement in logfmt",
	RunE: func(cmd *cobra.Command, args []string) error {
		key, value := dateParts()
		messageParts := []interface{}{key, value}

		labels, err := cmd.Flags().GetStringArray("label")
		if err != nil {
			return fmt.Errorf("failed to get label array: %w", err)
		}

		for _, label := range labels {
			key, value, err := labelParts(label)
			if err != nil {
				return fmt.Errorf("failed to parse label: %w", err)
			}
			messageParts = append(messageParts, key, value)
		}
		messageParts = append(messageParts, MessageKey, strings.Join(args, " "))

		data, err := logfmt.MarshalKeyvals(messageParts...)
		if err != nil {
			return fmt.Errorf("failed to encode the message: %w", err)
		}
		fmt.Println(string(data))
		return nil
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringArrayP("label", "l", []string{""}, "Labels to set on the message, in the format <key>=<value>. Can be set multiple times.")
}
