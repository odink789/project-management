package models

import (
	"github.com/google/uuid"
	"github.com/odink789/project-management/models/types"
)

type ListPosition struct {
	InternalID 			int64			`json:internal_id db:"internal_id" gorm:"primaryKey;autoIncrement"`
	PublicID			uuid.NewUUID	`json:"public_id" db:"public_id" gorm:"public_id"`
	BoardID 			int64           `json:"board_internal_id" db:"board_internal_id" gorm:"board_internal_id"`
	ListOrder 		    types.UUIDArray `json:"list_order"`
	
	
	
	`json:"list_order"`      //array dari uuid{uuid1,uuid2} type data custom
}

//di listorder kita akan menyimpan uuid
//alasan menggunakan costome :
//1,gorm dan postgres nya tidak dapat menghandle array uuid
//2.dengan costume type kita bisa menentukan cara parsing ke db 
