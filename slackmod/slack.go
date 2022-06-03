package slackmod

// handles slack API interface for sending webhooks back with responses

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strings"

	"github.com/parnurzeal/gorequest"
)

const (
	fileUploadURL       string = "https://slack.com/api/files.upload"
	channelCreateURL    string = "https://slack.com/api/channels.create"
	channelArchiveURL   string = "https://slack.com/api/channels.archive"
	channelListURL      string = "https://slack.com/api/conversations.list"
	channelInviteURL    string = "https://slack.com/api/channels.invite"
	channelTopicSetURL  string = "https://slack.com/api/channels.setTopic"
	channelUnArchiveURL string = "https://slack.com/api/channels.unarchive"
)

// Slackopts - slackCLI Options
type Slackopts struct {
	Version             string
	Config              string
	SlackHook           string
	SlackToken          string
	SlackDefaultName    string
	SlackDefaultChannel string
	SlackDefaultEmoji   string
	Snippet             bool
	BotDM               bool
}

// Field - struct
type Field struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

// BasicSlackPayload - returns most basic slack response payload
type BasicSlackPayload struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error"`
}

// BotDMPayload - struct for bot DMs
type BotDMPayload struct {
	Token          string       `json:"token,omitempty"`
	Channel        string       `json:"channel,omitempty"`
	Text           string       `json:"text,omitempty"`
	AsUser         bool         `json:"as_user,omitempty"`
	Attachments    []Attachment `json:"attachments,omitempty"`
	IconEmoji      string       `json:"icon_emoji,omitempty"`
	IconURL        string       `json:"icon_url,omitempty"`
	LinkNames      bool         `json:"link_names,omitempty"`
	Mkrdwn         bool         `json:"mrkdwn,omitempty"`
	Parse          string       `json:"parse,omitempty"`
	ReplyBroadcast bool         `json:"reply_broadcast,omitempty"`
	ThreadTS       string       `json:"thread_ts,omitempty"`
	UnfurlLinks    bool         `json:"unfurl_links,omitempty"`
	UnfurlMedia    bool         `json:"unfurl_media,omitempty"`
	Username       string       `json:"username,omitempty"`
}

// Attachment - struct
type Attachment struct {
	Fallback   string   `json:"fallback,omitempty"`
	Color      string   `json:"color,omitempty"`
	PreText    string   `json:"pretext,omitempty"`
	AuthorName string   `json:"author_name,omitempty"`
	AuthorLink string   `json:"author_link,omitempty"`
	AuthorIcon string   `json:"author_icon,omitempty"`
	Title      string   `json:"title,omitempty"`
	TitleLink  string   `json:"title_link,omitempty"`
	Text       string   `json:"text,omitempty"`
	ImageURL   string   `json:"image_url,omitempty"`
	Fields     []*Field `json:"fields,omitempty"`
	Footer     string   `json:"footer,omitempty"`
	FooterIcon string   `json:"footer_icon,omitempty"`
	Timestamp  int64    `json:"ts,omitempty"`
	MarkdownIn []string `json:"mrkdwn_in,omitempty"`
}

// Payload - struct
type Payload struct {
	Parse       string       `json:"parse,omitempty"`
	Username    string       `json:"username,omitempty"`
	IconURL     string       `json:"icon_url,omitempty"`
	IconEmoji   string       `json:"icon_emoji,omitempty"`
	Channel     string       `json:"channel,omitempty"`
	Text        string       `json:"text,omitempty"`
	LinkNames   string       `json:"link_names,omitempty"`
	Attachments []Attachment `json:"attachments,omitempty"`
	UnfurlLinks bool         `json:"unfurl_links,omitempty"`
	UnfurlMedia bool         `json:"unfurl_media,omitempty"`
}

// AddField - add fields
func (attachment *Attachment) AddField(field Field) *Attachment {
	attachment.Fields = append(attachment.Fields, &field)
	return attachment
}

func redirectPolicyFunc(req gorequest.Request, via []gorequest.Request) error {
	return fmt.Errorf("Incorrect token (redirection)")
}

// PostSnippet - Post a snippet of any type to slack channel
func PostSnippet(token string, fileType string, fileContent string, channel string, title string, comment string) error {

	form := url.Values{}

	form.Set("channels", channel)
	form.Set("content", fileContent)
	form.Set("filetype", fileType)
	form.Set("title", title)
	form.Set("initial_comment", comment)

	s := form.Encode()

	req, err := http.NewRequest("POST", fileUploadURL, strings.NewReader(s))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "Bearer "+token)

	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	//fmt.Println(resp.Status)

	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	//fmt.Println(string(body))

	return nil
}

// PostFile - Post a file of any type to slack channel
func PostFile(token string, channel string, fileName string) error {

	form := url.Values{}

	form.Set("channels", channel)
	//form.Set("filetype", fileType)
	form.Set("file", fileName)
	form.Set("filename", fileName)

	s := form.Encode()

	req, err := http.NewRequest("POST", fileUploadURL, strings.NewReader(s))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "Bearer "+token)

	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	fmt.Println(resp.Status)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(body))

	return nil
}

// Send - send message
func Send(webhookURL string, proxy string, payload Payload) []error {
	request := gorequest.New().Proxy(proxy)
	resp, _, err := request.
		Post(webhookURL).
		RedirectPolicy(redirectPolicyFunc).
		Send(payload).
		End()

	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		return []error{fmt.Errorf("Error sending msg. Status: %v", resp.Status)}
	}

	return nil
}

// WranglerDM - Send chat.Post API DM messages "as the bot"
func WranglerDM(opts Slackopts, payload BotDMPayload) error {
	url := "https://slack.com/api/chat.postMessage"

	payload.Token = opts.SlackToken
	payload.AsUser = true

	jsonStr, err := json.Marshal(&payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+opts.SlackToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return err

}

// Wrangler - wrangle slack calls
func Wrangler(webhookURL string, message string, myChannel string, emojiName string, botName string, attachments Attachment) {

	payload := Payload{
		Text:        message,
		Username:    botName,
		Channel:     myChannel,
		IconEmoji:   emojiName,
		Attachments: []Attachment{attachments},
	}
	err := Send(webhookURL, "", payload)
	if len(err) > 0 {
		fmt.Printf("Slack Messaging Error in Wrangler function in slack.go: %s\n", err)
	}
}

// LoadConfig - Load Main Configuration JSON
func LoadConfig(path string) (opts Slackopts, fail string) {
	var fileName string

	if path == "default" {
		if runtime.GOOS == "windows" {
			fileName = "c:/programdata/slackcli.json"
		} else {
			fileName = "/etc/slackcli.json"
		}
	} else {
		fileName = path
	}

	file, err := os.Open(fileName)
	if err != nil {
		if path == "default" {
			return opts, "nodefault"
		}

		return opts, "err"
	}

	decoded := json.NewDecoder(file)
	err = decoded.Decode(&opts)
	if err != nil {
		fmt.Println("Error reading invalid JSON file: " + fileName + "(" + err.Error() + ")")
		return opts, "err"
	}

	if opts.SlackDefaultEmoji == "" {
		opts.SlackDefaultEmoji = "robot_face"
	}
	if opts.SlackDefaultName == "" {
		opts.SlackDefaultName = "Slack Robot"
	}

	if opts.SlackDefaultChannel == "" {
		opts.SlackDefaultChannel = "#general"
	}

	return opts, "loaded"
}
