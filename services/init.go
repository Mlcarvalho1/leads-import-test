package services

import (
	"sync"

	"leads-import/database"
)

var (
	importService     *LeadImportService
	importServiceOnce sync.Once
)

func GetImportService() *LeadImportService {
	importServiceOnce.Do(func() {
		db := database.GetDB()

		var chats ChatRepository
		mongoDB := database.GetMongo()
		if mongoDB != nil {
			chats = NewMongoChatRepository(mongoDB)
		} else {
			chats = &NoopChatRepository{}
		}

		importService = &LeadImportService{
			DB:       db,
			Chats:    chats,
			WhatsApp: &NoopWhatsAppValidator{},
			Events:   &NoopEventEmitter{},
			Cache:    &NoopCacheClearer{},
		}
	})
	return importService
}
