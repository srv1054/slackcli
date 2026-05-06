package slackmod

// handles slack API interface for sending webhooks back with responses

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/parnurzeal/gorequest"
)

const (
	fileGetUploadURL    string = "https://slack.com/api/files.getUploadURLExternal"
	fileCompleteURL     string = "https://slack.com/api/files.completeUploadExternal"
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

type fileUploadURLPayload struct {
	Ok        bool   `json:"ok"`
	Error     string `json:"error"`
	UploadURL string `json:"upload_url"`
	FileID    string `json:"file_id"`
}

type fileCompletePayload struct {
	Files          []fileCompleteItem `json:"files"`
	Channels       string             `json:"channels,omitempty"`
	InitialComment string             `json:"initial_comment,omitempty"`
}

type fileCompleteItem struct {
	ID    string `json:"id"`
	Title string `json:"title,omitempty"`
}

type conversationListPayload struct {
	Ok               bool `json:"ok"`
	Error            string
	Channels         []conversationItem `json:"channels"`
	ResponseMetadata struct {
		NextCursor string `json:"next_cursor"`
	} `json:"response_metadata"`
}

type conversationItem struct {
	ID   string `json:"id"`
	Name string `json:"name"`
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
	if title == "" {
		title = "snippet.txt"
	}

	return uploadExternalFile(token, []byte(fileContent), title, title, fileType, channel, comment)
}

// PostFile - Post a file of any type to slack channel
func PostFile(token string, channel string, fileName string) error {
	fileContent, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}

	title := filepath.Base(fileName)
	return uploadExternalFile(token, fileContent, title, title, "", channel, "")
}

func uploadExternalFile(token string, fileContent []byte, filename string, title string, snippetType string, channel string, comment string) error {
	if snippetType == "Plain Text" {
		snippetType = ""
	}

	uploadURL, fileID, err := getUploadURLExternal(token, filename, int64(len(fileContent)), snippetType)
	if err != nil {
		return err
	}

	if err := uploadFileContent(uploadURL, fileContent); err != nil {
		return err
	}

	return completeUploadExternal(token, fileID, title, channel, comment)
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

func getUploadURLExternal(token string, filename string, length int64, snippetType string) (string, string, error) {
	form := url.Values{}
	form.Set("filename", filename)
	form.Set("length", strconv.FormatInt(length, 10))
	if snippetType != "" {
		form.Set("snippet_type", snippetType)
	}

	req, err := http.NewRequest("POST", fileGetUploadURL, strings.NewReader(form.Encode()))
	if err != nil {
		return "", "", err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	var payload fileUploadURLPayload
	if err := decodeSlackResponse(resp, &payload); err != nil {
		return "", "", err
	}
	if !payload.Ok {
		return "", "", fmt.Errorf("files.getUploadURLExternal failed: %s", payload.Error)
	}
	if payload.UploadURL == "" || payload.FileID == "" {
		return "", "", fmt.Errorf("files.getUploadURLExternal response was missing upload_url or file_id")
	}

	return payload.UploadURL, payload.FileID, nil
}

func uploadFileContent(uploadURL string, fileContent []byte) error {
	req, err := http.NewRequest("POST", uploadURL, bytes.NewReader(fileContent))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/octet-stream")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("upload to Slack file URL failed: %s: %s", resp.Status, strings.TrimSpace(string(body)))
	}

	return nil
}

func completeUploadExternal(token string, fileID string, title string, channel string, comment string) error {
	channel, err := resolveUploadChannels(token, channel)
	if err != nil {
		return err
	}

	payload := fileCompletePayload{
		Files: []fileCompleteItem{
			{
				ID:    fileID,
				Title: title,
			},
		},
		Channels:       channel,
		InitialComment: comment,
	}

	jsonStr, err := json.Marshal(&payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", fileCompleteURL, bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var slackResp BasicSlackPayload
	if err := decodeSlackResponse(resp, &slackResp); err != nil {
		return err
	}
	if !slackResp.Ok {
		return fmt.Errorf("files.completeUploadExternal failed: %s", slackResp.Error)
	}

	return nil
}

func resolveUploadChannels(token string, channels string) (string, error) {
	if channels == "" {
		return "", nil
	}

	channelList := strings.Split(channels, ",")
	for i, channel := range channelList {
		channel = strings.TrimSpace(channel)
		if strings.HasPrefix(channel, "#") {
			channelID, err := findConversationID(token, strings.TrimPrefix(channel, "#"))
			if err != nil {
				return "", err
			}
			channel = channelID
		}
		channelList[i] = channel
	}

	return strings.Join(channelList, ","), nil
}

func findConversationID(token string, channelName string) (string, error) {
	cursor := ""
	for {
		form := url.Values{}
		form.Set("exclude_archived", "true")
		form.Set("limit", "1000")
		form.Set("types", "public_channel,private_channel")
		if cursor != "" {
			form.Set("cursor", cursor)
		}

		req, err := http.NewRequest("POST", channelListURL, strings.NewReader(form.Encode()))
		if err != nil {
			return "", err
		}
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Add("Authorization", "Bearer "+token)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return "", err
		}

		var payload conversationListPayload
		err = decodeSlackResponse(resp, &payload)
		resp.Body.Close()
		if err != nil {
			return "", err
		}
		if !payload.Ok {
			if payload.Error == "missing_scope" {
				return "", fmt.Errorf("missing_scope while resolving #%s via conversations.list: add Slack scopes `channels:read`, `groups:read`, `im:read`, `mpim:read` (then reinstall app) or pass a channel ID like C123456 instead of #%s", channelName, channelName)
			}
			return "", fmt.Errorf("conversations.list failed while resolving #%s: %s", channelName, payload.Error)
		}

		for _, channel := range payload.Channels {
			if channel.Name == channelName {
				return channel.ID, nil
			}
		}

		cursor = payload.ResponseMetadata.NextCursor
		if cursor == "" {
			break
		}
	}

	return "", fmt.Errorf("could not find Slack channel #%s; pass a channel ID like C123456 instead", channelName)
}

func decodeSlackResponse(resp *http.Response, payload interface{}) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("Slack API request failed: %s: %s", resp.Status, strings.TrimSpace(string(body)))
	}
	if err := json.Unmarshal(body, payload); err != nil {
		return err
	}

	return nil
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
