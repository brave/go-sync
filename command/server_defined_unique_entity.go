package command

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/brave/go-sync/datastore"
	"github.com/brave/go-sync/schema/protobuf/sync_pb"
	"github.com/brave/go-sync/utils"
	"github.com/google/uuid"
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

func createServerDefinedUniqueEntity(name string, serverDefinedTag string, clientID string, chainID int64, parentID string, specifics *sync_pb.EntitySpecifics) (*datastore.SyncEntity, error) {
	now := utils.UnixMilli(time.Now())
	deleted := false
	folder := true
	version := int64(1)
	idUUID, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	idString := idUUID.String()

	pbEntity := &sync_pb.SyncEntity{
		Ctime: &now, Mtime: &now, Deleted: &deleted, Folder: &folder,
		Name: aws.String(name), ServerDefinedUniqueTag: aws.String(serverDefinedTag),
		Version: &version, ParentIdString: &parentID,
		IdString: &idString, Specifics: specifics}

	return datastore.CreateDBSyncEntity(pbEntity, nil, clientID, chainID)
}

func CreateServerDefinedUniqueEntities(clientID string, chainID int64) (entities []*datastore.SyncEntity, err error) {
	// Create nigori top-level folder
	nigoriSpecific := &sync_pb.NigoriSpecifics{}
	nigoriEntitySpecific := &sync_pb.EntitySpecifics_Nigori{Nigori: nigoriSpecific}
	specifics := &sync_pb.EntitySpecifics{SpecificsVariant: nigoriEntitySpecific}
	entity, err := createServerDefinedUniqueEntity(nigoriName, nigoriTag, clientID, chainID, "0", specifics)
	if err != nil {
		return nil, fmt.Errorf("error creating entity with a server tag: %w", err)
	}
	entities = append(entities, entity)

	// Create bookmarks top-level folder
	bookmarkSpecific := &sync_pb.BookmarkSpecifics{}
	bookmarkEntitySpecific := &sync_pb.EntitySpecifics_Bookmark{Bookmark: bookmarkSpecific}
	specifics = &sync_pb.EntitySpecifics{SpecificsVariant: bookmarkEntitySpecific}
	entity, err = createServerDefinedUniqueEntity(bookmarksName, bookmarksTag, clientID, chainID, "0", specifics)
	if err != nil {
		return nil, fmt.Errorf("error creating entity with a server tag: %w", err)
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
			name, tag, clientID, chainID, bookmarkRootID, specifics)
		if err != nil {
			return nil, fmt.Errorf("error creating entity with a server tag: %w", err)
		}
		entities = append(entities, entity)
	}
	return entities, nil
}
