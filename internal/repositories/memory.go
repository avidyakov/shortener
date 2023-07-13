package repositories

type memoryRepo struct {
	links map[string]string
}

func NewMemoryLink() LinkRepo {
	return &memoryRepo{
		links: make(map[string]string),
	}
}

func (m *memoryRepo) GetLink(shortLinkID string) (originLink string, ok bool) {
	originLink, ok = m.links[shortLinkID]
	return
}

func (m *memoryRepo) CreateLink(shortLinkID, originLink string) {
	m.links[shortLinkID] = originLink
}

func (m *memoryRepo) RemoveLink(shortLinkID string) {
	delete(m.links, shortLinkID)
}
