// SPDX-FileCopyrightText: 2019 Kent Gibson <warthog618@gmail.com>
//
// SPDX-License-Identifier: MIT

// +build linux

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/warthog618/gpiod"
)

func init() {
	rootCmd.AddCommand(infoCmd)
}

var infoCmd = &cobra.Command{
	Use:                   "info [flags] [chip]...",
	Short:                 "Info about chip lines",
	Long:                  `Print information about all lines of the specified GPIO chip(s) (or all gpiochips if none are specified).`,
	Run:                   info,
	DisableFlagsInUseLine: true,
}

func info(cmd *cobra.Command, args []string) {
	rc := 0
	cc := []string(nil)
	cc = append(cc, args...)
	if len(cc) == 0 {
		cc = gpiod.Chips()
	}
	for _, path := range cc {
		c, err := gpiod.NewChip(path)
		if err != nil {
			logErr(cmd, err)
			rc = 1
			continue
		}
		fmt.Printf("%s - %d lines:\n", c.Name, c.Lines())
		for o := 0; o < c.Lines(); o++ {
			li, err := c.LineInfo(o)
			if err != nil {
				logErr(cmd, err)
				rc = 1
				continue
			}
			printLineInfo(li)
		}
		c.Close()
	}
	os.Exit(rc)
}

func printLineInfo(li gpiod.LineInfo) {
	if len(li.Name) == 0 {
		li.Name = "unnamed"
	}
	if li.Used {
		if len(li.Consumer) == 0 {
			li.Consumer = "kernel"
		}
		if strings.Contains(li.Consumer, " ") {
			li.Consumer = "\"" + li.Consumer + "\""
		}
	} else {
		li.Consumer = "unused"
	}
	dirn := "input"
	if li.Config.Direction == gpiod.LineDirectionOutput {
		dirn = "output"
	}
	active := "active-high"
	if li.Config.ActiveLow {
		active = "active-low"
	}
	flags := []string(nil)
	if li.Used {
		flags = append(flags, "used")
	}
	if li.Config.Drive == gpiod.LineDriveOpenDrain {
		flags = append(flags, "open-drain")
	}
	if li.Config.Drive == gpiod.LineDriveOpenSource {
		flags = append(flags, "open-source")
	}
	if li.Config.Bias == gpiod.LineBiasPullUp {
		flags = append(flags, "pull-up")
	}
	if li.Config.Bias == gpiod.LineBiasPullDown {
		flags = append(flags, "pull-down")
	}
	if li.Config.Bias == gpiod.LineBiasDisabled {
		flags = append(flags, "bias-disabled")
	}
	flstr := ""
	if len(flags) > 0 {
		flstr = "[" + strings.Join(flags, " ") + "]"
	}
	fmt.Printf("\tline %3d:%12s%12s%8s%13s%s\n",
		li.Offset, li.Name, li.Consumer, dirn, active, flstr)
}
