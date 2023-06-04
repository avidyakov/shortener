package repositories

type LinkRepo interface {
	GetLink(string) (string, bool)
	CreateLink(string, string)
}

type MemoryLink struct {
	links map[string]string
}

func NewMemoryLink() *MemoryLink {
	return &MemoryLink{
		links: make(map[string]string),
	}
}

func (m *MemoryLink) GetLink(shortLinkID string) (originLink string, ok bool) {
	originLink, ok = m.links[shortLinkID]
	return
}

func (m *MemoryLink) CreateLink(shortLinkID, originLink string) {
	m.links[shortLinkID] = originLink
}
