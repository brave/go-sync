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
		"history_delete_directive_global_id": {SpecificsVariant: &sync_pb.EntitySpecifics_HistoryDeleteDirective{
			HistoryDeleteDirective: &sync_pb.HistoryDeleteDirectiveSpecifics{
				GlobalIdDirective: &sync_pb.GlobalIdDirective{GlobalId: []int64{1, 2}},
			},
		}},
		"history_delete_directive_time_range": {SpecificsVariant: &sync_pb.EntitySpecifics_HistoryDeleteDirective{
			HistoryDeleteDirective: &sync_pb.HistoryDeleteDirectiveSpecifics{
				TimeRangeDirective: &sync_pb.TimeRangeDirective{StartTimeUsec: aws.Int64(1)},
			},
		}},
		"history_delete_directive_url_directive_without_url": {SpecificsVariant: &sync_pb.EntitySpecifics_HistoryDeleteDirective{
			HistoryDeleteDirective: &sync_pb.HistoryDeleteDirectiveSpecifics{
				UrlDirective: &sync_pb.UrlDirective{EndTimeUsec: aws.Int64(1)},
			},
		}},
		"password_encrypted": {SpecificsVariant: &sync_pb.EntitySpecifics_Password{
			Password: &sync_pb.PasswordSpecifics{
				Encrypted: &sync_pb.EncryptedData{KeyName: aws.String("key"), Blob: aws.String("blob")},
			},
		}},
		"password_encrypted_with_notes_backup": {SpecificsVariant: &sync_pb.EntitySpecifics_Password{
			Password: &sync_pb.PasswordSpecifics{
				Encrypted:            &sync_pb.EncryptedData{KeyName: aws.String("key"), Blob: aws.String("blob")},
				EncryptedNotesBackup: &sync_pb.EncryptedData{KeyName: aws.String("key"), Blob: aws.String("notes")},
			},
		}},
	} {
		t.Run(name, func(t *testing.T) {
			_, err := datastore.CreateDBSyncEntity(newPBEntity(specifics), nil, "client1")
			assert.NoError(t, err)
		})
	}
}

func TestCreateDBSyncEntity_EncryptionExemptTypesWithPlaintext(t *testing.T) {
	for name, specifics := range map[string]*sync_pb.EntitySpecifics{
		"history_delete_directive_url": {SpecificsVariant: &sync_pb.EntitySpecifics_HistoryDeleteDirective{
			HistoryDeleteDirective: &sync_pb.HistoryDeleteDirectiveSpecifics{
				UrlDirective: &sync_pb.UrlDirective{Url: aws.String("https://example.com")},
			},
		}},
		"password_unencrypted_metadata": {SpecificsVariant: &sync_pb.EntitySpecifics_Password{
			Password: &sync_pb.PasswordSpecifics{
				UnencryptedMetadata: &sync_pb.PasswordSpecificsMetadata{Url: aws.String("https://example.com")},
			},
		}},
		"password_encrypted_with_unencrypted_metadata": {SpecificsVariant: &sync_pb.EntitySpecifics_Password{
			Password: &sync_pb.PasswordSpecifics{
				Encrypted:           &sync_pb.EncryptedData{KeyName: aws.String("key"), Blob: aws.String("blob")},
				UnencryptedMetadata: &sync_pb.PasswordSpecificsMetadata{Url: aws.String("https://example.com")},
			},
		}},
	} {
		t.Run(name, func(t *testing.T) {
			_, err := datastore.CreateDBSyncEntity(newPBEntity(specifics), nil, "client1")
			assert.ErrorIs(t, err, datastore.ErrEntityHasPlaintextSpecifics)
		})
	}
}

func TestCreateDBSyncEntity_ServerDefinedEmptySpecifics(t *testing.T) {
	entity := newPBEntity(emptyBookmarkSpecifics())
	entity.ServerDefinedUniqueTag = aws.String("bookmark_bar")
	_, err := datastore.CreateDBSyncEntity(entity, nil, "client1")
	assert.NoError(t, err)
}

func TestCreateDBSyncEntity_ServerDefinedWithPlaintextSpecifics(t *testing.T) {
	entity := newPBEntity(plaintextBookmarkSpecifics())
	entity.ServerDefinedUniqueTag = aws.String("bookmark_bar")
	_, err := datastore.CreateDBSyncEntity(entity, nil, "client1")
	assert.ErrorIs(t, err, datastore.ErrEntityHasPlaintextSpecifics)
}
