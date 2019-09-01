# slackcli
CLI for dumping slack webhooks

## Webhook
Webhook URL is required to use `slackcli`.  You can generate one for your slack server really easily.
Optional for certain functions such as posting *snippets* or *bot DM messages* you will need a slack API token.

You can feed your webhook into `slackcli` in one of three ways:
* via CLI parameter `-hook` for webhook URL and `-token` for API Token
* via ENV variable `slackhook` or `slacktoken`
* via optional `slackcli` config file specified by `-cfg` parameter

If all or any combination exist `slackcli` will take the first one in this order:
* config file
* CLI parameter
* ENV variable

## Optional Config file
`slackcli` will accept everything it needs via the command line parameters interface.  However if its simpler and more convient to have a single config file that contains items that are always the same when you need to execute the application (such as webhook URLs) you can put them into a config file and specify its location.  If the `-cfg` paramter is not present when running `slackcli` it is assume there isn't one and all information is coming from the command line interface parameters.

An example of the config file and its options are available in the source repository: https://github.com/srv1054/slackcli/config.example.json
