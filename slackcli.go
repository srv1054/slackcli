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

	var myChannel string
	var myMessage string
	var myEmoji string
	var myName string
	var myConfig string
	var fail string
	var attachments slackmod.Attachment
	var opts slackmod.Slackopts

	opts.Version = "1.04"

	version := flag.Bool("v", false, "Show current version number")
	cfg := flag.String("cfg", "", "Path to optional configuration file (default /etc/slackcli.json)")
	hooker := flag.String("hook", "", "Slackhook URL, if no config file")
	token := flag.String("token", "", "Slack Bearer Token, if no config file")
	postchannel := flag.String("c", "", "Channel to send message to (specify # or @)")
	postname := flag.String("n", "", "Name of bot to post to channel as")
	postemoji := flag.String("e", "", "Emoji to use for bot post (no colons)")
	postmessage := flag.String("m", "", "Message to send to slack channel")

	flag.Parse()

	if *version {
		fmt.Println("slackcli v" + opts.Version + " @srv1054")
		fmt.Println("https://github.com/srv1054/slackcli")
		os.Exit(0)
	}

	if *cfg != "" {
		myConfig = *cfg
	} else {
		myConfig = "default"
	}

	opts, fail = slackmod.LoadConfig(myConfig)
	if fail == "err" {
		os.Exit(1)
	}

	if fail == "nodefault" {
		if *hooker == "" || *postchannel == "" || *postemoji == "" || *postname == "" {
			fmt.Println("Missing all the parameters, without a config file you need -c -n -e -m and -hook minimum")
			os.Exit(1)
		}
		opts.SlackHook = *hooker
		opts.SlackHook = *token
	}

	if *postmessage == "" {
		fmt.Println("No message specified, nothing to send")
		os.Exit(0)
	} else {
		myMessage = *postmessage
	}
	if *postchannel == "" {
		myChannel = opts.SlackDefaultChannel
	} else {
		myChannel = *postchannel
	}
	if *postemoji == "" {
		myEmoji = opts.SlackDefaultEmoji
	} else {
		myEmoji = *postemoji
	}
	if *postname == "" {
		myName = opts.SlackDefaultName
	} else {
		myName = *postname
	}

	slackmod.Wrangler(opts.SlackHook, myMessage, myChannel, myEmoji, myName, attachments)

	// start features for BOT DMs and Snippets here

}
