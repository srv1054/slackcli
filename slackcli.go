package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
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

	opts.Version = "1.06.02"

	version := flag.Bool("v", false, "Show current version number")
	cfg := flag.String("cfg", "", "Path to optional configuration file (default /etc/slackcli.json)")
	hooker := flag.String("hook", "", "Slackhook URL, if no config file")
	token := flag.String("token", "", "Slack Bearer Token, if no config file")
	postchannel := flag.String("c", "", "Channel to send message to (specify # or @)")
	postname := flag.String("n", "", "Name of bot to post to channel as")
	postemoji := flag.String("e", "", "Emoji to use for bot post (no colons)")
	postmessage := flag.String("m", "", "Message to send to slack channel")
	snippet := flag.Bool("s", false, "Expect Pipe to Standard In to send to slack as a Snippet attachment")

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

	if *snippet {

		nBytes := int64(0)

		// read stdin
		r := bufio.NewReader(os.Stdin)
		buf := make([]byte, 0, 4*1024)

		for {

			n, err := r.Read(buf[:cap(buf)])
			buf = buf[:n]

			if n == 0 {

				if err == nil {
					continue
				}

				if err == io.EOF {
					break
				}

				log.Fatal(err)
			}

			nBytes += int64(len(buf))

			// fail out with error if no stdin
			if nBytes == 0 {
				fmt.Println("Standard input was blank or 0 bytes, you must pipe valid text in")
				os.Exit(1)
			}

			// validate token exists (on CLI or in cfg)
			if *token == "" && opts.SlackToken == "" {
				fmt.Println("You must provide a slack bot token in your cfg or on the CLI -token, to leverage snippts")
				os.Exit(1)
			}

			err = slackmod.PostSnippet(opts, "text", string(buf), "server-messages", "filename")
			if err != nil {
				fmt.Println("Something failed in PostSnippet function -> " + err.Error())
				os.Exit(1)
			}

			fmt.Println(string(buf))

			if err != nil && err != io.EOF {
				log.Fatal(err)
			}
		}

		os.Exit(0)
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
