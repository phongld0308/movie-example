package model

type Metadata struct {
	ID          string `json:"id" yaml:"id"`
	Title       string `json:"title" yaml:"title"`
	Description string `json:"description" yaml:"description"`
	Director    string `json:"director" yaml:"director"`
}
