package datastore

import (
	"errors"
	"strconv"

	sentry "github.com/getsentry/sentry-go"
	"google.golang.org/protobuf/proto"

	"github.com/brave/go-sync/schema/protobuf/sync_pb"
)

const nigoriTypeID int = 47745

var ErrEntityHasPlaintextSpecifics = errors.New("entity must not contain plaintext specifics")

// Data types exempt from the specifics.encrypted requirement, keyed by
// EntitySpecifics field number from entity_specifics.proto.
var encryptionExemptTypes = map[int]struct{}{
	// Carries the key bag that protects every other type.
	nigoriTypeID: {},
	// The server has to read these to purge its own stored entities, and it
	// must not hold the user's passphrase. The payload is visit ids only.
	HistoryDeleteDirectiveTypeID: {},
	// Encrypted by its own scheme into its own specifics fields, so
	// `specifics.encrypted` stays unset. See
	// DataTypeWorker::EncryptPasswordSpecificsData(). 45873 is the
	// PasswordSpecifics field number.
	45873: {},
}

// checkEntityEncrypted returns an error if the entity carries plaintext specifics.
// Exempt types and server-defined entities are skipped since the server
// creates the latter itself.
func checkEntityEncrypted(entity *sync_pb.SyncEntity, dataType int) error {
	if _, ok := encryptionExemptTypes[dataType]; ok {
		return nil
	}
	if entity.GetServerDefinedUniqueTag() != "" {
		return nil
	}
	if hasPlaintextSpecifics(entity.GetSpecifics()) {
		return ErrEntityHasPlaintextSpecifics
	}
	return nil
}

func hasPlaintextSpecifics(specifics *sync_pb.EntitySpecifics) bool {
	m := specifics.ProtoReflect()
	field := m.WhichOneof(m.Descriptor().Oneofs().ByName("specifics_variant"))
	if field == nil {
		return false
	}
	return proto.Size(m.Get(field).Message().Interface()) > 0
}

// maybeRejectUnencryptedEntity checks the entity's encryption and reports
// failures to Sentry.
func maybeRejectUnencryptedEntity(entity *sync_pb.SyncEntity, dataType int) error {
	err := checkEntityEncrypted(entity, dataType)
	if err == nil {
		return nil
	}
	sentry.WithScope(func(scope *sentry.Scope) {
		scope.SetTag("data_type", strconv.Itoa(dataType))
		sentry.CaptureException(err)
	})
	return err
}
