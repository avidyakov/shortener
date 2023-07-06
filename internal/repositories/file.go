package repositories

import (
	"encoding/json"
	"github.com/avidyakov/shortener/internal/config"
	"log"
	"os"
)

type fileRepo struct {
	links map[string]string
}

func NewFileRepo() LinkRepo {
	content, err := os.ReadFile(config.Conf.File)
	links := make(map[string]string)

	if err == nil {
		err = json.Unmarshal(content, &links)
		if err != nil {
			log.Fatal(err)
		}
	}

	return &fileRepo{
		links: links,
	}
}

func (r *fileRepo) GetLink(shortLinkID string) (originLink string, ok bool) {
	originLink, ok = r.links[shortLinkID]
	return
}

func (r *fileRepo) CreateLink(shortLinkID string, originLink string) {
	r.links[shortLinkID] = originLink
	r.writeFile()
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

	err = os.WriteFile(config.Conf.File, jsonData, 0644)
	if err != nil {
		log.Println("Error writing file:", err)
		return
	}
}
