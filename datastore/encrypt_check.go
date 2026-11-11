package datastore

import (
	"errors"
	"strconv"

	sentry "github.com/getsentry/sentry-go"
	"google.golang.org/protobuf/proto"

	"github.com/brave/go-sync/schema/protobuf/sync_pb"
)

// Data type IDs are EntitySpecifics field numbers from entity_specifics.proto.
const (
	nigoriTypeID   int = 47745
	passwordTypeID int = 45873
)

var ErrEntityHasPlaintextSpecifics = errors.New("entity must not contain plaintext specifics")

// checkEntityEncrypted returns an error if the entity carries plaintext
// specifics, reporting it to Sentry.
func checkEntityEncrypted(entity *sync_pb.SyncEntity, dataType int) error {
	if isEncryptionExempt(entity.GetSpecifics(), dataType) {
		return nil
	}
	if !hasPlaintextSpecifics(entity.GetSpecifics()) {
		return nil
	}
	sentry.WithScope(func(scope *sentry.Scope) {
		scope.SetTag("data_type", strconv.Itoa(dataType))
		sentry.CaptureException(ErrEntityHasPlaintextSpecifics)
	})
	return ErrEntityHasPlaintextSpecifics
}

func isEncryptionExempt(specifics *sync_pb.EntitySpecifics, dataType int) bool {
	switch dataType {
	case nigoriTypeID:
		// Carries the key bag that protects every other type.
		return true
	case HistoryDeleteDirectiveTypeID:
		// The server may read these in the future to delete its own stored
		// entities. URL directives would expose browsing history, so they are
		// not exempt.
		return specifics.GetHistoryDeleteDirective().GetUrlDirective().GetUrl() == ""
	case passwordTypeID:
		// Encrypted by its own scheme into its own specifics fields, so
		// `specifics.encrypted` stays unset. See
		// DataTypeWorker::EncryptPasswordSpecificsData().
		return specifics.GetPassword().GetUnencryptedMetadata() == nil
	}
	return false
}

func hasPlaintextSpecifics(specifics *sync_pb.EntitySpecifics) bool {
	m := specifics.ProtoReflect()
	field := m.WhichOneof(m.Descriptor().Oneofs().ByName("specifics_variant"))
	if field == nil {
		return false
	}
	return proto.Size(m.Get(field).Message().Interface()) > 0
}
