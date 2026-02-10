package services

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MongoChatRepository struct {
	DB *mongo.Database
}

func NewMongoChatRepository(db *mongo.Database) *MongoChatRepository {
	return &MongoChatRepository{DB: db}
}

func (r *MongoChatRepository) FindChatsByPhones(ctx context.Context, phones []string, accountID int, companyID int) ([]Chat, error) {
	coll := r.DB.Collection("chats")

	filter := bson.M{
		"contact.phone": bson.M{"$in": phones},
		"accountId":     accountID,
		"companyId":     companyID,
	}

	cursor, err := coll.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find chats: %w", err)
	}
	defer cursor.Close(ctx)

	var results []struct {
		ID      bson.ObjectID `bson:"_id"`
		Contact struct {
			Phone string `bson:"phone"`
		} `bson:"contact"`
		LeadID *int `bson:"leadId"`
	}

	if err := cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("failed to decode chats: %w", err)
	}

	var chats []Chat
	for _, r := range results {
		chats = append(chats, Chat{
			ID:        r.ID.Hex(),
			Phone:     r.Contact.Phone,
			AccountID: accountID,
			CompanyID: companyID,
			LeadID:    r.LeadID,
		})
	}

	return chats, nil
}

func (r *MongoChatRepository) CreateChat(ctx context.Context, phone string, dialCode string, countryCode string, accountID int, companyID int) (string, error) {
	coll := r.DB.Collection("chats")

	doc := bson.M{
		"contact": bson.M{
			"phone":       phone,
			"dialCode":    dialCode,
			"countryCode": countryCode,
		},
		"accountId": accountID,
		"companyId": companyID,
		"createdAt": time.Now(),
		"updatedAt": time.Now(),
	}

	result, err := coll.InsertOne(ctx, doc)
	if err != nil {
		return "", fmt.Errorf("failed to create chat: %w", err)
	}

	oid, ok := result.InsertedID.(bson.ObjectID)
	if !ok {
		return "", fmt.Errorf("unexpected inserted ID type")
	}

	return oid.Hex(), nil
}

func (r *MongoChatRepository) UpdateChatLeadID(ctx context.Context, chatID string, leadID int) error {
	coll := r.DB.Collection("chats")

	oid, err := bson.ObjectIDFromHex(chatID)
	if err != nil {
		return fmt.Errorf("invalid chat ID: %w", err)
	}

	_, err = coll.UpdateOne(ctx, bson.M{"_id": oid}, bson.M{
		"$set": bson.M{
			"leadId":    leadID,
			"updatedAt": time.Now(),
		},
	})
	if err != nil {
		return fmt.Errorf("failed to update chat: %w", err)
	}

	return nil
}
