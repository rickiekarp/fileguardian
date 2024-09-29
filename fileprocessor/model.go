package fileprocessor

type FileGuardianEventMessage struct {
	Id         *int64 `json:"id,omitempty"`
	Type       string `json:"type,omitempty"`
	Source     string `json:"source,omitempty"`
	Target     string `json:"target,omitempty"`
	Context    string `json:"context,omitempty"`
	Inserttime *int64 `json:"inserttime,omitempty"`
}
