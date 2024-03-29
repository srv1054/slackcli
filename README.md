# slackcli
CLI for dumping slack messages

# How to use

* `slackcli -v` - get version
* `slackcli -h` - get CLI parameter help

### All the things on the CLI
---
You can specify everything on the CLI if you want.  SlackwebHook, Emoji, Name, Channel, etc...
By doing this `slackcli` will run with no config files necessary

**Example**:   
* `slackcli -hook "https://hooks.slack.com/services/MY-HOOK-TGSABSP4/BS59MA3H/Ooc2gaGwtFcclZ0HXuv" -e "man" -n "RadBot" -c "#general" -m "Hi I'm radbot bot"`

### Configuration File
---
If you have a configuration file the `slackcli` app will default to its settings when they aren't over-ridden by the command line parameters.  This is nice for quick calls or if you are using `slackcli` as the same "bot" all the time.

Be sure to secure your slackcli.json permissions as it contains your Slack Webhook, don't want the evil people stealing that.

**Example**:  
* `slackcli -cfg /etc/myparams.json -m "Hi I'm a bot using my config file for parameters"`

Since you can specify the config file on the CLI you can have as many as you want and reference them to create "different bots" inside slack without needing massive command lines.
The configuration file can be specified using the `-cfg` parameter.  Specify a path and the filename.

**Example**:  
* `slackcli -cfg /etc/filmBot.json -m "Hi I'm film bot talking about films using my config file for parameters"`
* `slackcli -cfg /etc/weatherBot.json -m "Hi I'm a weather bot, sunny as usual, using my config file for parameters"` 


If you do not specify the `-cfg` parameter, `slackcli` does (every time) check a default path to see if one exists.   By using this default path you can simplify the command line for a single use installation.  Or at a minimum to have defaults if you forget a paramemter.

* Windows Default path:  `C:/programdata/slackcli.json`
* Unix Default path:  `/etc/slackcli.json`

**Example**: 
* `slackcli -m "Hi this bot relies on default paths and configs for simplicity"`

CLI parameters individually will over-ride any parameter in the configuration file so you can mix and match

**Example**: 
* `slackcli -m "Hi this bot uses cfg for everything except emoji" -e "bender_dance"`

## Webhook
Webhook URL is required to use `slackcli`.  You can generate one for your slack server really easily -> https://slack.com/help/articles/115005265063-Incoming-WebHooks-for-Slack

An example of the config file and its options are available in the source repository: https://github.com/srv1054/slackcli/blob/master/config.example.json

## Snippets
To send a snippet message into slack, you will need a BOT app token (not a web hook) either in your -cfg or on the CLI directly.

The `-s` parameter will tell slackcli you want to create a snippet and are feeding it in via STDIN (pipe).  You can pipe any text information into slackcli to create the snippet.  You may also specify `-msg` which will place a markdown capable message prior to your snippet (in slack its referred to as the comment) and you may also specify `-title` which will give the snippet a title inside the post.  If omitted slack will labled it `untitled`

**NOTE**: 
The `-s` parameter is need for *piping* files into slackcli as a snippet.  If you wish to just upload a file by referencing a filename, use the `-file` parameter instead.  

**Example Piping**:
* `cat myfile.txt | slackcli -s -title "My File" -msg "<!here> is my file!" -cfg /etc/mycfg.json`

**Example Uploading**:
* `slackcli -file "myfile.txt (or jpg etc)" -cfg /etc/mycfg.json`

# Examples Galore!

Upload a jpg to a slack channel
* `slackcli -file "/my/path/picture.jpg" -cfg /etc/mycfg.json`
* `slackcli -file "/my/path/picture.png" -c "#mychannel" -token "xoxb-asdasdf-asdfasdfasdf-asdfasdfasdf-2134234234"`

Pipe data into a slack snippet
* `cat /my/path/data.txt | slackcli -s -msg "this is the data we talked about" -title "data.txt" -cfg /etc/mycfg.json`
* `cat /my/path/data.txt | slackcli -s -msg "this is the data we talked about" -title "data.txt" -token "xoxb-asdasdf-asdfasdfasdf-asdfasdfasdf-2134234234"`
* `cat /my/path/data.txt | sort -r | unique -c | slackcli -s -msg "this is the data we talked about" -title "data.txt" -cfg /etc/mycfg.json`

Send a basic markdown message to a channel
* `slackcli -m "This is a robot *bleep* _blork_, bloop" -c "#mychannel" -e "robot" -n "Some Robot" -hook "https://slack.com/webhook/asdf90a870982038as098df0a98f70398f2"`
* `slackcli -m "This is a robot *bleep* _blork_, bloop" -cfg /etc/mycfg.json`

# TODO
- [x] Configure to accept piped data and send to Slack Snippet
- [x] Configure to allow file uploads
- [ ] Configure to send DM's as a real bot using Token, instead of showing up via default webhook "Slackbot" DM

