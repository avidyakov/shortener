package repositories

type memoryRepo struct {
	links map[string]string
}

func (m *memoryRepo) GetUrlsByUserId(_ int) ([]map[string]string, error) {
	panic("implement me")
}

func (m *memoryRepo) CreateUser() (int, error) {
	panic("implement me")
}

func (m *memoryRepo) CheckConnection() error {
	return nil
}

func NewMemoryRepo() LinkRepo {
	return &memoryRepo{
		links: make(map[string]string),
	}
}

func (m *memoryRepo) GetOriginLink(shortLinkID string) (originLink string, ok bool) {
	originLink, ok = m.links[shortLinkID]
	return
}

func (m *memoryRepo) GetShortLink(originLink string) (shortLinkID string, ok bool) {
	for k, v := range m.links {
		if v == originLink {
			return k, true
		}
	}
	return "", false
}

func (m *memoryRepo) CreateLink(shortLinkID, originLink string, _ int) error {
	m.links[shortLinkID] = originLink
	return nil
}

func (m *memoryRepo) RemoveLink(shortLinkID string) {
	delete(m.links, shortLinkID)
}
