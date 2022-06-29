package actions

import db "server/pkg/core/storage"

type StorageAction struct {
	Storage *db.Storage
}

func NewStorageAction(storage *db.Storage) StorageAction {
	return StorageAction{
		Storage: storage,
	}
}
