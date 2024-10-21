package datastore

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/brave/go-sync/schema/protobuf/sync_pb"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"
)

const (
	HistoryTypeID                int = 963985
	HistoryDeleteDirectiveTypeID int = 150251
	// Expiration time for history and history delete directive
	// entities in seconds
	HistoryExpirationIntervalSecs = 14 * 24 * 60 * 60 // 14 days
)

// SyncEntity is used to marshal and unmarshal sync items in dynamoDB.
type SyncEntity struct {
	ClientID string
	// ChainID is a synthetic key that is connected to the client id in the SQL db.
	ChainID                *int64 `dynamodbav:"-" db:"chain_id"`
	ID                     string
	ParentID               *string `dynamodbav:",omitempty" db:"parent_id"`
	Version                *int64
	Mtime                  *int64
	Ctime                  *int64
	Name                   *string `dynamodbav:",omitempty"`
	NonUniqueName          *string `dynamodbav:",omitempty" db:"non_unique_name"`
	ServerDefinedUniqueTag *string `dynamodbav:",omitempty" db:"server_defined_unique_tag"`
	Deleted                *bool
	OriginatorCacheGUID    *string `dynamodbav:",omitempty" db:"originator_cache_guid"`
	OriginatorClientItemID *string `dynamodbav:",omitempty" db:"originator_client_item_id"`
	Specifics              []byte
	DataType               *int `db:"data_type"`
	Folder                 *bool
	ClientDefinedUniqueTag *string `dynamodbav:",omitempty" db:"client_defined_unique_tag"`
	UniquePosition         []byte  `dynamodbav:",omitempty" db:"unique_position"`
	DataTypeMtime          *string
	ExpirationTime         *int64
	OldVersion             *int64 `dynamodbav:"-" db:"old_version"`
}

func validatePBEntity(entity *sync_pb.SyncEntity) error {
	if entity == nil {
		return fmt.Errorf("validate SyncEntity error: empty SyncEntity")
	}

	if entity.IdString == nil {
		return fmt.Errorf("validate SyncEntity error: empty IdString")
	}

	if entity.Version == nil {
		return fmt.Errorf("validate SyncEntity error: empty Version")
	}

	if entity.Specifics == nil {
		return fmt.Errorf("validate SyncEntity error: nil Specifics")
	}

	return nil
}

// CreateDBSyncEntity converts a protobuf sync entity into a DB sync item.
func CreateDBSyncEntity(entity *sync_pb.SyncEntity, cacheGUID *string, clientID string, chainID int64) (*SyncEntity, error) {
	err := validatePBEntity(entity)
	if err != nil {
		log.Error().Err(err).Msg("Invalid sync_pb.SyncEntity received")
		return nil, fmt.Errorf("error validating protobuf sync entity to create DB sync entity: %w", err)
	}

	// Specifics are always passed and checked by validatePBEntity above.
	var specifics []byte
	specifics, err = proto.Marshal(entity.Specifics)
	if err != nil {
		log.Error().Err(err).Msg("Marshal specifics failed")
		return nil, fmt.Errorf("error marshalling specifics to create DB sync entity: %w", err)
	}

	// Use reflect to find out data type ID defined in protobuf tag.
	structField := reflect.ValueOf(entity.Specifics.SpecificsVariant).Elem().Type().Field(0)
	tag := structField.Tag.Get("protobuf")
	s := strings.Split(tag, ",")
	dataType, _ := strconv.Atoi(s[1])

	var uniquePosition []byte
	if entity.UniquePosition != nil {
		uniquePosition, err = proto.Marshal(entity.UniquePosition)
		if err != nil {
			log.Error().Err(err).Msg("Marshal UniquePosition failed")
			return nil, fmt.Errorf("error marshalling unique position to create DB sync entity: %w", err)
		}
	}

	id := *entity.IdString
	var originatorCacheGUID, originatorClientItemID *string
	if cacheGUID != nil {
		if *entity.Version == 0 {
			idUUID, err := uuid.NewV7()
			if err != nil {
				return nil, err
			}
			id = idUUID.String()
		}
		originatorCacheGUID = cacheGUID
		originatorClientItemID = entity.IdString
	}

	now := time.Now()

	var expirationTime *int64
	if dataType == HistoryTypeID || dataType == HistoryDeleteDirectiveTypeID {
		expirationTime = aws.Int64(now.Unix() + HistoryExpirationIntervalSecs)
	}

	nowMillis := aws.Int64(now.UnixMilli())
	// ctime is only used when inserting a new entity, here we use client passed
	// ctime if it is passed, otherwise, use current server time as the creation
	// time. When updating, ctime will be ignored later in the query statement.
	cTime := nowMillis
	if entity.Ctime != nil {
		cTime = entity.Ctime
	}

	dataTypeMtime := strconv.Itoa(dataType) + "#" + strconv.FormatInt(*nowMillis, 10)

	// Set default values on Deleted and Folder attributes for new entities, the
	// default values are specified by sync.proto protocol.
	deleted := entity.Deleted
	folder := entity.Folder
	if *entity.Version == 0 {
		if entity.Deleted == nil {
			deleted = aws.Bool(false)
		}
		if entity.Folder == nil {
			folder = aws.Bool(false)
		}
	}

	return &SyncEntity{
		ClientID:               clientID,
		ChainID:                &chainID,
		ID:                     id,
		ParentID:               entity.ParentIdString,
		Version:                entity.Version,
		Ctime:                  cTime,
		Mtime:                  nowMillis,
		Name:                   entity.Name,
		NonUniqueName:          entity.NonUniqueName,
		ServerDefinedUniqueTag: entity.ServerDefinedUniqueTag,
		Deleted:                deleted,
		OriginatorCacheGUID:    originatorCacheGUID,
		OriginatorClientItemID: originatorClientItemID,
		ClientDefinedUniqueTag: entity.ClientTagHash,
		Specifics:              specifics,
		Folder:                 folder,
		UniquePosition:         uniquePosition,
		DataType:               aws.Int(dataType),
		DataTypeMtime:          aws.String(dataTypeMtime),
		ExpirationTime:         expirationTime,
	}, nil
}

// CreatePBSyncEntity converts a DB sync item to a protobuf sync entity.
func CreatePBSyncEntity(entity *SyncEntity) (*sync_pb.SyncEntity, error) {
	id := &entity.ID
	// The client tag hash must be used as the primary key
	// for the history type.
	if *entity.DataType == HistoryTypeID {
		id = entity.ClientDefinedUniqueTag
	}

	pbEntity := &sync_pb.SyncEntity{
		IdString:               id,
		ParentIdString:         entity.ParentID,
		Version:                entity.Version,
		Mtime:                  entity.Mtime,
		Ctime:                  entity.Ctime,
		Name:                   entity.Name,
		NonUniqueName:          entity.NonUniqueName,
		ServerDefinedUniqueTag: entity.ServerDefinedUniqueTag,
		ClientTagHash:          entity.ClientDefinedUniqueTag,
		OriginatorCacheGuid:    entity.OriginatorCacheGUID,
		OriginatorClientItemId: entity.OriginatorClientItemID,
		Deleted:                entity.Deleted,
		Folder:                 entity.Folder,
	}

	if entity.Specifics != nil {
		pbEntity.Specifics = &sync_pb.EntitySpecifics{}
		err := proto.Unmarshal(entity.Specifics, pbEntity.Specifics)
		if err != nil {
			log.Error().Err(err).Msg("Unmarshal specifics failed")
			return nil, fmt.Errorf("error unmarshalling specifics to create protobuf sync entity: %w", err)
		}
	}

	if entity.UniquePosition != nil {
		pbEntity.UniquePosition = &sync_pb.UniquePosition{}
		err := proto.Unmarshal(entity.UniquePosition, pbEntity.UniquePosition)
		if err != nil {
			log.Error().Err(err).Msg("Unmarshal UniquePosition failed")
			return nil, fmt.Errorf("error unmarshalling unique position to create protobuf sync entity: %w", err)
		}
	}

	return pbEntity, nil
}
