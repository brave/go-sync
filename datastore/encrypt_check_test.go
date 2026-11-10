package datastore_test

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/assert"

	"github.com/brave/go-sync/datastore"
	"github.com/brave/go-sync/schema/protobuf/sync_pb"
)

func emptyBookmarkSpecifics() *sync_pb.EntitySpecifics {
	return &sync_pb.EntitySpecifics{
		SpecificsVariant: &sync_pb.EntitySpecifics_Bookmark{Bookmark: &sync_pb.BookmarkSpecifics{}},
	}
}

func plaintextBookmarkSpecifics() *sync_pb.EntitySpecifics {
	return &sync_pb.EntitySpecifics{
		SpecificsVariant: &sync_pb.EntitySpecifics_Bookmark{
			Bookmark: &sync_pb.BookmarkSpecifics{Url: aws.String("https://example.com")},
		},
	}
}

func newPBEntity(specifics *sync_pb.EntitySpecifics) *sync_pb.SyncEntity {
	return &sync_pb.SyncEntity{
		IdString:      aws.String("id"),
		Version:       aws.Int64(0),
		ClientTagHash: aws.String("client_tag"),
		Specifics:     specifics,
	}
}

func TestCreateDBSyncEntity_EmptySpecificsVariant(t *testing.T) {
	_, err := datastore.CreateDBSyncEntity(newPBEntity(emptyBookmarkSpecifics()), nil, "client1")
	assert.NoError(t, err)
}

func TestCreateDBSyncEntity_EncryptedOnly(t *testing.T) {
	specifics := emptyBookmarkSpecifics()
	specifics.Encrypted = &sync_pb.EncryptedData{KeyName: aws.String("key"), Blob: aws.String("blob")}
	_, err := datastore.CreateDBSyncEntity(newPBEntity(specifics), nil, "client1")
	assert.NoError(t, err)
}

func TestCreateDBSyncEntity_PlaintextSpecifics(t *testing.T) {
	_, err := datastore.CreateDBSyncEntity(newPBEntity(plaintextBookmarkSpecifics()), nil, "client1")
	assert.ErrorIs(t, err, datastore.ErrEntityHasPlaintextSpecifics)
}

func TestCreateDBSyncEntity_EncryptedWithPlaintextSpecifics(t *testing.T) {
	specifics := plaintextBookmarkSpecifics()
	specifics.Encrypted = &sync_pb.EncryptedData{}
	_, err := datastore.CreateDBSyncEntity(newPBEntity(specifics), nil, "client1")
	assert.ErrorIs(t, err, datastore.ErrEntityHasPlaintextSpecifics)
}

func TestCreateDBSyncEntity_TombstoneWithPlaintextSpecifics(t *testing.T) {
	entity := newPBEntity(plaintextBookmarkSpecifics())
	entity.Deleted = aws.Bool(true)
	_, err := datastore.CreateDBSyncEntity(entity, nil, "client1")
	assert.ErrorIs(t, err, datastore.ErrEntityHasPlaintextSpecifics)
}

func TestCreateDBSyncEntity_EncryptionExemptTypes(t *testing.T) {
	for name, specifics := range map[string]*sync_pb.EntitySpecifics{
		"nigori": {SpecificsVariant: &sync_pb.EntitySpecifics_Nigori{
			Nigori: &sync_pb.NigoriSpecifics{KeystoreMigrationTime: aws.Int64(1)},
		}},
		"history_delete_directive": {SpecificsVariant: &sync_pb.EntitySpecifics_HistoryDeleteDirective{
			HistoryDeleteDirective: &sync_pb.HistoryDeleteDirectiveSpecifics{
				UrlDirective: &sync_pb.UrlDirective{Url: aws.String("https://example.com")},
			},
		}},
		"password": {SpecificsVariant: &sync_pb.EntitySpecifics_Password{
			Password: &sync_pb.PasswordSpecifics{Encrypted: &sync_pb.EncryptedData{}},
		}},
	} {
		t.Run(name, func(t *testing.T) {
			_, err := datastore.CreateDBSyncEntity(newPBEntity(specifics), nil, "client1")
			assert.NoError(t, err)
		})
	}
}

func TestCreateDBSyncEntity_ServerDefinedSkipsEncryptionCheck(t *testing.T) {
	entity := newPBEntity(plaintextBookmarkSpecifics())
	entity.ServerDefinedUniqueTag = aws.String("bookmark_bar")
	_, err := datastore.CreateDBSyncEntity(entity, nil, "client1")
	assert.NoError(t, err)
}
