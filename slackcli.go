package main

/* Written by srv1054 (https://github.com/srv1054)
   See LICENSE file for usage and "Stealability"

   Compiles to single binary that can be used to dump info to slack quickly and easily via webhooks
*/

type struct slackopts {
	Version string
	Config string
	SlackHook string
}

func main() {

	var opts slackopts
	opts.Version = "0.1"

	version := flag.Bool("v", false, "Show current version number")
	cfg := flag.String("cfg", "", "Path to optional configuration file")
	hooker := flag.String("hook", "", "Slackhook URL")

	flag.Parse()

	opts.Config = *cfg

	if (*version) {
		fmt.Println("slackcli v" + opts.Version)
	}

	// if cfg is true load it

	// else check ENV for slackhook

	// else take param
	//tiktokOpts.Config.Tkey = os.Getenv("tkey")
}

