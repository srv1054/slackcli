package main

import (
	"bufio"
	"errors"
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

// ************ FINISH FILE UPLOLAD MULTIPART DATA SECTION
// FIGURE OUT HOW TO PASS emoji/name/channel to PostSNippet ( NOT SUPPORTED by Slack that I can see in API Docs )
// Finish file posting multipart data stuff

func main() {

	var myChannel string
	var myToken string
	var myMessage string
	var myEmoji string
	var myName string
	var myConfig string
	var fail string
	var attachments slackmod.Attachment
	var opts slackmod.Slackopts
	var buf []byte
	var totalBuf string
	var nBytes = int64(0)

	opts.Version = "1.08.00"

	version := flag.Bool("v", false, "Show current version number")
	cfg := flag.String("cfg", "", "Path to optional configuration file (default /etc/slackcli.json)")
	hooker := flag.String("hook", "", "Slackhook URL, if no config file")
	token := flag.String("token", "", "Slack Bearer Token, if no config file")
	postchannel := flag.String("c", "", "Channel to send message to (specify # or @)")
	postname := flag.String("n", "", "Name of bot to post to channel as")
	postemoji := flag.String("e", "", "Emoji to use for bot post (no colons)")
	postmessage := flag.String("m", "", "Message to send to slack channel")
	snippet := flag.Bool("s", false, "Expect Pipe to Standard In to send to slack as a Snippet attachment.  Note slack does not support Emoji or bot name changes on API Snippet posts or file uploads.  -e and -n do not have any effect")
	snipmsg := flag.String("msg", "", "Pre-Snippet Message, posted with any snippt. Requires -s")
	fileis := flag.String("file", "", "Path and name of file to send as a snippet, can not be binary")
	title := flag.String("title", "", "Title of Snippet file.  Requires -s")

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

	if !*snippet && *snipmsg != "" {
		fmt.Println("-msg requires -s.  See -h for more info")
		os.Exit(0)
	}

	opts, fail = slackmod.LoadConfig(myConfig)
	if fail == "err" {
		os.Exit(1)
	}

	if fail == "nodefault" {
		if (*hooker == "" && *token == "") || *postchannel == "" || *postemoji == "" || *postname == "" {
			fmt.Println("Missing all the parameters, without a config file you need -c -n -e (-m or -s) and (-hook or -token) minimum")
			os.Exit(1)
		}
		opts.SlackHook = *hooker
		opts.SlackHook = *token
	}

	// grab cfg file defaults if CLI is blank
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

	// if we get a -s to make a snippet
	if *snippet {

		// read stdin
		r := bufio.NewReader(os.Stdin)
		buf = make([]byte, 0, 4*1024)

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
			totalBuf = totalBuf + string(buf)

			if err != nil && err != io.EOF {
				log.Fatal(err)
			}
		}

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
		myToken = *token
		if opts.SlackToken != "" {
			myToken = opts.SlackToken
		}

		// post the snippet
		err := slackmod.PostSnippet(myToken, "Plain Text", totalBuf, myChannel, *title, *snipmsg)
		if err != nil {
			fmt.Println("Something failed in PostSnippet function -> " + err.Error())
			os.Exit(1)
		}

		os.Exit(0)

	}

	// if we get a -file to upload a file
	if *fileis != "" {

		// validate token exists (on CLI or in cfg)
		if *token == "" && opts.SlackToken == "" {
			fmt.Println("You must provide a slack bot token in your cfg or on the CLI -token, to leverage snippts")
			os.Exit(1)
		}
		myToken = *token
		if opts.SlackToken != "" {
			myToken = opts.SlackToken
		}

		// UPLOAD whatever file was sent

		// validate path and file exist
		if _, err := os.Stat(*fileis); errors.Is(err, os.ErrNotExist) {
			fmt.Println("Could not find specified file included in the -file parameter of " + *fileis)
			os.Exit(1)
		}
		// verify file type

		// call snippet func
		err := slackmod.PostFile(myToken, myChannel, *fileis)
		if err != nil {
			fmt.Println("Something failed in PostSnippet function -> " + err.Error())
			os.Exit(1)
		}

		os.Exit(0)
	}

	if *postmessage == "" {
		fmt.Println("No message specified, nothing to send")
		os.Exit(0)
	} else {
		myMessage = *postmessage
	}

	// assume we are just web hookering it
	slackmod.Wrangler(opts.SlackHook, myMessage, myChannel, myEmoji, myName, attachments)

}
