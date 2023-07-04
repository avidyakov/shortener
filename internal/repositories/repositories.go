package repositories

type LinkRepo interface {
	GetLink(string) (string, bool)
	CreateLink(string, string)
	RemoveLink(string)
}
