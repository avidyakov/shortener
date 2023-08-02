package repositories

type LinkRepo interface {
	GetOriginLink(string) (string, bool)
	GetShortLink(string) (string, bool)
	CreateLink(string, string, int) error
	RemoveLink(string)
	CheckConnection() error
	CreateUser() (int, error)
}
