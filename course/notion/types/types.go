package types

// Use struct to define the request structure for better type safety and readability
type TextContent struct {
	Content string      `json:"content"`
	Link    interface{} `json:"link,omitempty"`
}

type TextAnnotation struct {
	Bold          bool   `json:"bold"`
	Italic        bool   `json:"italic"`
	Strikethrough bool   `json:"strikethrough"`
	Underline     bool   `json:"underline"`
	Code          bool   `json:"code"`
	Color         string `json:"color"`
}

type RichText struct {
	Type        string          `json:"type,omitempty"`
	Text        *TextContent    `json:"text"`
	Annotations *TextAnnotation `json:"annotations"`
	PlainText   string          `json:"plain_text,omitempty"`
	Href        interface{}     `json:"href"`
}

type DateObject struct {
	Start    string      `json:"start,omitempty"`
	End      interface{} `json:"end,omitempty"`
	TimeZone interface{} `json:"time_zone,omitempty"`
}

type StatusObject struct {
	ID    string `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Color string `json:"color,omitempty"`
}
type Title struct {
	ID    string     `json:"id,omitempty"`
	Type  string     `json:"type,omitempty"`
	Title []RichText `json:"title"`
}

type Text struct {
	ID       string     `json:"id,omitempty"`
	Type     string     `json:"type,omitempty"`
	RichText []RichText `json:"rich_text"`
}

type Deadline struct {
	ID   string      `json:"id,omitempty"`
	Type string      `json:"type,omitempty"`
	Date *DateObject `json:"date,omitempty"`
}

type Courses struct {
	ID       string              `json:"id,omitempty"`
	Type     string              `json:"type,omitempty"`
	Relation []map[string]string `json:"relation,omitempty"`
	HasMore  bool                `json:"has_more,omitempty"`
}

type Type struct {
	ID     string            `json:"id,omitempty"`
	Type   string            `json:"type,omitempty"`
	Select map[string]string `json:"select,omitempty"`
}

type Status struct {
	ID     string        `json:"id,omitempty"`
	Type   string        `json:"type,omitempty"`
	Status *StatusObject `json:"status,omitempty"`
}

type TODO struct {
	ID       string     `json:"id,omitempty"`
	Type     string     `json:"type,omitempty"`
	RichText []RichText `json:"rich_text,omitempty"`
}

type Properties struct {
	Name       Title `json:"Course Name,omitempty"`
	Code       Text  `json:"Course Code,omitempty"`
	RoomNumber Text  `json:"Room Number,omitempty"`
	Duration   Text  `json:"Course Duration,omitempty"`
}

type PageRequest struct {
	Cover  *interface{} `json:"cover,omitempty"`
	Icon   *interface{} `json:"icon,omitempty"`
	Parent struct {
		Type       string `json:"type,omitempty"`
		DatabaseID string `json:"database_id,omitempty"`
	} `json:"parent,omitempty"`
	Archived   bool        `json:"archived,omitempty"`
	InTrash    bool        `json:"in_trash,omitempty"`
	Properties *Properties `json:"properties,omitempty"`
}
