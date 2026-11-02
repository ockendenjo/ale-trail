package data

type FileWithETag struct {
	File *File   `json:"file"`
	ETag *string `json:"etag"`
}

func (f *FileWithETag) GetStashForDate(date string) *Stash {
	for _, stash := range f.File.Stashes {
		if stash.Date == date {
			return stash
		}
	}
	return nil
}

type File struct {
	Schema   *string  `json:"$schema,omitempty"`
	AdminKey string   `json:"adminKey"`
	Stashes  []*Stash `json:"stashes"`
}
