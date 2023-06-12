package repositories

type LinkRepo interface {
	GetLink(string) (string, bool)
	CreateLink(string, string)
}

type memoryLink struct {
	links map[string]string
}

func NewMemoryLink() LinkRepo {
	return &memoryLink{
		links: make(map[string]string),
	}
}

func (m *memoryLink) GetLink(shortLinkID string) (originLink string, ok bool) {
	originLink, ok = m.links[shortLinkID]
	return
}

func (m *memoryLink) CreateLink(shortLinkID, originLink string) {
	m.links[shortLinkID] = originLink
}
