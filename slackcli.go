package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/srv1054/slackcli/slackmod"
)

/* Written by srv1054 (https://github.com/srv1054)
   See LICENSE file for usage and "Stealability"

   Compiles to single binary that can be used to dump info to slack quickly and easily via webhooks
*/

func main() {

	var shook string
	var stoken string
	var opts slackmod.Slackopts
	opts.Version = "0.2"

	version := flag.Bool("v", false, "Show current version number")
	cfg := flag.String("cfg", "", "Path to optional configuration file")
	hooker := flag.String("hook", "", "Slackhook URL")
	token := flag.String("token", "", "Slack Bearer Token")
	snipme := flag.Bool("snip", false, "Post a snippet")
	botme := flag.Bool("botdm", false, "Send DM message as a bot")

	flag.Parse()

	opts.Config = *cfg
	opts.Snippet = *snipme
	opts.BotDM = *botme
	shook = *hooker
	stoken = *token

	if *version {
		fmt.Println("slackcli v" + opts.Version)
	}

	// if cfg is true load it
	if opts.Config != "" {
		fmt.Println("I gotta a config file at " + opts.Config)
	} else {
		if shook != "" {
			opts.SlackHook = shook
			fmt.Println("I found a hook via params")
		} else {
			key := os.Getenv("slackhook")
			if key != "" {
				opts.SlackHook = key
				fmt.Println("I found a hook via ENV")
			} else {
				fmt.Println("ERR: I found no useable Slack webhook URL")
				os.Exit(0)
			}
		}
		// Token is only required if -snip  or -botdm parameter was used
		if opts.Snippet || opts.BotDM {
			if stoken != "" {
				opts.SlackToken = stoken
				fmt.Println("I found a token via params")
			} else {
				key := os.Getenv("slacktoken")
				if key != "" {
					opts.SlackToken = key
					fmt.Println("I found a token via ENV")
				} else {
					fmt.Println("ERR: I found no useable Slack Token")
					os.Exit(0)
				}
			}
		}
	}
}
