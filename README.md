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
Webhook URL is required to use `slackcli`.  You can generate one for your slack server [https://slack.com/help/articles/115005265063-Incoming-WebHooks-for-Slack|really easily].

An example of the config file and its options are available in the source repository: https://github.com/srv1054/slackcli/config.example.json
