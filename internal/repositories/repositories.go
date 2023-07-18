package repositories

type LinkRepo interface {
	GetOriginLink(string) (string, bool)
	GetShortLink(string) (string, bool)
	CreateLink(string, string) error
	RemoveLink(string)
}
