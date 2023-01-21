package file

type GenericFile struct {
	FileName string            `json:"file_name,omitempty"`
	URL      string            `json:"url,omitempty"`
	Metadata map[string]string `json:"metadata"`
	Private  bool              `json:"private"`
	Data     []byte            `json:"data"`
}
