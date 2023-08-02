package repositories

import (
	"encoding/json"
	"log"
	"os"
)

type fileRepo struct {
	links       map[string]string
	storagePath string
}

func (r *fileRepo) CreateUser() (int, error) {
	panic("implement me")
}

func (r *fileRepo) CheckConnection() error {
	return nil
}

func NewFileRepo(storagePath string) LinkRepo {
	content, err := os.ReadFile(storagePath)
	links := make(map[string]string)

	if err == nil {
		err = json.Unmarshal(content, &links)
		if err != nil {
			log.Fatal(err)
		}
	}

	return &fileRepo{
		links:       links,
		storagePath: storagePath,
	}
}

func (r *fileRepo) GetOriginLink(shortLinkID string) (originLink string, ok bool) {
	originLink, ok = r.links[shortLinkID]
	return
}

func (r *fileRepo) GetShortLink(originLink string) (shortLinkID string, ok bool) {
	for k, v := range r.links {
		if v == originLink {
			return k, true
		}
	}
	return "", false
}

func (r *fileRepo) CreateLink(shortLinkID string, originLink string, _ int) error {
	r.links[shortLinkID] = originLink
	r.writeFile()
	return nil
}

func (r *fileRepo) RemoveLink(shortLinkID string) {
	delete(r.links, shortLinkID)
	r.writeFile()
}

func (r *fileRepo) writeFile() {
	jsonData, err := json.MarshalIndent(r.links, "", "  ")
	if err != nil {
		log.Println("Error marshaling JSON:", err)
		return
	}

	err = os.WriteFile(r.storagePath, jsonData, 0644)
	if err != nil {
		log.Println("Error writing file:", err)
		return
	}
}
