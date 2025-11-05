package command

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/google/uuid"

	"github.com/brave/go-sync/datastore"
	"github.com/brave/go-sync/schema/protobuf/sync_pb"
)

const (
	nigoriName          string = "Nigori"
	nigoriTag           string = "google_chrome_nigori"
	bookmarksName       string = "Bookmarks"
	bookmarksTag        string = "google_chrome_bookmarks"
	otherBookmarksName  string = "Other Bookmarks"
	otherBookmarksTag   string = "other_bookmarks"
	syncedBookmarksName string = "Synced Bookmarks"
	syncedBookmarksTag  string = "synced_bookmarks"
	bookmarkBarName     string = "Bookmark Bar"
	bookmarkBarTag      string = "bookmark_bar"
)

func createServerDefinedUniqueEntity(name string, serverDefinedTag string, clientID string, parentID string, specifics *sync_pb.EntitySpecifics) (*datastore.SyncEntity, error) {
	now := time.Now().UnixMilli()
	deleted := false
	folder := true
	version := int64(1)
	idString := uuid.New().String()

	pbEntity := &sync_pb.SyncEntity{
		Ctime: &now, Mtime: &now, Deleted: &deleted, Folder: &folder,
		Name: aws.String(name), ServerDefinedUniqueTag: aws.String(serverDefinedTag),
		Version: &version, ParentIdString: &parentID,
		IdString: &idString, Specifics: specifics}

	return datastore.CreateDBSyncEntity(pbEntity, nil, clientID)
}

// InsertServerDefinedUniqueEntities inserts the server defined unique tag
// entities if it is not in the DB yet for a specific client.
func InsertServerDefinedUniqueEntities(db datastore.Datastore, clientID string) error {
	var entities []*datastore.SyncEntity
	// Check if they're existed already for this client.
	// If yes, just return directly.
	ready, err := db.HasServerDefinedUniqueTag(clientID, nigoriTag)
	if err != nil {
		return fmt.Errorf("error checking if entity with a server tag existed: %w", err)
	}
	if ready {
		return nil
	}

	// Create nigori top-level folder
	nigoriSpecific := &sync_pb.NigoriSpecifics{}
	nigoriEntitySpecific := &sync_pb.EntitySpecifics_Nigori{Nigori: nigoriSpecific}
	specifics := &sync_pb.EntitySpecifics{SpecificsVariant: nigoriEntitySpecific}
	entity, err := createServerDefinedUniqueEntity(nigoriName, nigoriTag, clientID, "0", specifics)
	if err != nil {
		return fmt.Errorf("error creating entity with a server tag: %w", err)
	}
	entities = append(entities, entity)

	// Create bookmarks top-level folder
	bookmarkSpecific := &sync_pb.BookmarkSpecifics{}
	bookmarkEntitySpecific := &sync_pb.EntitySpecifics_Bookmark{Bookmark: bookmarkSpecific}
	specifics = &sync_pb.EntitySpecifics{SpecificsVariant: bookmarkEntitySpecific}
	entity, err = createServerDefinedUniqueEntity(bookmarksName, bookmarksTag, clientID, "0", specifics)
	if err != nil {
		return fmt.Errorf("error creating entity with a server tag: %w", err)
	}
	entities = append(entities, entity)

	// Create other bookmarks, synced bookmarks, and bookmark bar sub-folders
	bookmarkRootID := entity.ID
	bookmarkSecondLevelFolders := map[string]string{
		otherBookmarksName:  otherBookmarksTag,
		syncedBookmarksName: syncedBookmarksTag,
		bookmarkBarName:     bookmarkBarTag}
	for name, tag := range bookmarkSecondLevelFolders {
		entity, err := createServerDefinedUniqueEntity(
			name, tag, clientID, bookmarkRootID, specifics)
		if err != nil {
			return fmt.Errorf("error creating entity with a server tag: %w", err)
		}
		entities = append(entities, entity)
	}

	// Start a transaction to insert all server defined unique entities
	err = db.InsertSyncEntitiesWithServerTags(entities)
	if err != nil {
		return fmt.Errorf("error inserting entities with server tags: %w", err)
	}
	return nil
}
