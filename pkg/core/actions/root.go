package actions

import "server/pkg/db"

type StorageAction struct {
	Storage *db.Storage
}

func NewStorageAction(storage *db.Storage) StorageAction {
	return StorageAction{
		Storage: storage,
	}
}
