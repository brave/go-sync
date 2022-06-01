package account

import (
	"fmt"

	"github.com/brave/go-sync/datastore"
)

func HandleDelete(clientID string, db datastore.Datastore) error {
	err := db.DeleteClientItems(clientID)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
