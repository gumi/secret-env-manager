package model

type Secret struct {
	Name      string
	CreatedAt string
	Version   string
}

type Secrets struct {
	Secrets []Secret
}
