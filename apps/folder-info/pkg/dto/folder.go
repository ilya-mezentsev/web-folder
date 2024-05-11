package dto

type (
	DirInfo struct {
		Path string `json:"path"`

		Files []File `json:"files"`
		Dirs  []Dir  `json:"dirs"`
	}

	File struct {
		Name string `json:"name"`
		Type string `json:"type"`
		Size string `json:"size"`
	}

	Dir struct {
		Name string `json:"name"`
		Size string `json:"size"`
	}
)
