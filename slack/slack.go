package slack

type Field struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

type Action struct {
	Type  string `json:"type"`
	Text  string `json:"text"`
	URL   string `json:"url"`
	Style string `json:"style"`
}

type Attachment struct {
	Fallback   *string   `json:"fallback"`
	Color      *string   `json:"color"`
	PreText    *string   `json:"pretext"`
	AuthorName *string   `json:"author_name"`
	AuthorLink *string   `json:"author_link"`
	AuthorIcon *string   `json:"author_icon"`
	Title      *string   `json:"title"`
	TitleLink  *string   `json:"title_link"`
	Text       *string   `json:"text"`
	ImageURL   *string   `json:"image_url"`
	Fields     []*Field  `json:"fields"`
	Footer     *string   `json:"footer"`
	FooterIcon *string   `json:"footer_icon"`
	Timestamp  *int64    `json:"ts"`
	MarkdownIn *[]string `json:"mrkdwn_in"`
	Actions    []*Action `json:"actions"`
	CallbackID *string   `json:"callback_id"`
}

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
	Markdown    bool         `json:"mrkdwn,omitempty"`
}

func (attachment *Attachment) AddField(field Field) *Attachment {
	attachment.Fields = append(attachment.Fields, &field)
	return attachment
}

func (attachment *Attachment) AddAction(action Action) *Attachment {
	attachment.Actions = append(attachment.Actions, &action)
	return attachment
}
