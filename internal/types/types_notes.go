package types

// NotionPageFetchRequest represents a request to fetch a Notion page by its ID.
type NotionPageFetchRequest struct {
	PageID string `json:"page_id"`
}

// NotionPageProperties represents the properties of a Notion notes page.
type NotionPageProperties struct {
	Doc        Doc        `json:"Doc"`
	Transcript Transcript `json:"Transcript"`
	Keywords   Keywords   `json:"Keywords"`
	Title      Title      `json:"Title"`
	Date       Date       `json:"Date"`
	Courses    Courses    `json:"Courses"`
}

// NotionPageFetchResponse represents the response from fetching a Notion notes page.
type NotionPageFetchResponse struct {
	Title      string               `json:"title"`
	Properties NotionPageProperties `json:"properties"`
	Content    string               `json:"content"`
	URL        string               `json:"url"`
}

// Date represents the Notion date property.
type Date struct {
	ID   string     `json:"id"`
	Type string     `json:"type"`
	Date *DateValue `json:"date"`
}

type DateValue struct {
	Start    *string `json:"start,omitempty"`
	End      *string `json:"end,omitempty"`
	TimeZone *string `json:"time_zone,omitempty"`
}

// Keywords represents the Notion keywords property.
type Keywords struct {
	ID       string         `json:"id"`
	Type     string         `json:"type"`
	RichText []RichTextItem `json:"rich_text"`
}

// Transcript represents the Notion transcript property.
type Transcript struct {
	ID       string         `json:"id"`
	Type     string         `json:"type"`
	RichText []RichTextItem `json:"rich_text"`
}

// RichTextItem represents a Notion rich_text item.
type RichTextItem struct {
	Type        string       `json:"type"`
	Text        *TitleText   `json:"text,omitempty"`
	Annotations *Annotations `json:"annotations,omitempty"`
	PlainText   string       `json:"plain_text"`
	Href        *string      `json:"href"`
}

// Doc represents the Notion files property.
type Doc struct {
	ID    string    `json:"id"`
	Type  string    `json:"type"`
	Files []DocFile `json:"files"`
}

type DocFile struct {
	Name       string      `json:"name"`
	Type       string      `json:"type"`
	File       *FileObject `json:"file,omitempty"`
	ExpiryTime *string     `json:"expiry_time,omitempty"`
}

type FileObject struct {
	URL string `json:"url"`
}

type TitleText struct {
	Content string  `json:"content"`
	Link    *string `json:"link"`
}

type Annotations struct {
	Bold          bool   `json:"bold"`
	Italic        bool   `json:"italic"`
	Strikethrough bool   `json:"strikethrough"`
	Underline     bool   `json:"underline"`
	Code          bool   `json:"code"`
	Color         string `json:"color"`
}

var NOTES_COLUMNS = map[string]string{
	"LiYj":   "doc",
	"QaGq":   "courses",
	"QgLP":   "date",
	"_S%7B_": "keywords",
	"v%3ESl": "transcript",
	"title":  "title",
}
