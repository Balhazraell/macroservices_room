package room

import (
	"encoding/json"

	"github.com/Balhazraell/logger"
)

// APIMetods - Перечень доступных API методов.
var APIMetods = map[string]func(string){
	"ClientConnect":    apiClientConnect,
	"ClientDisconnect": apiClientDisconnect,
	"SetChunckState":   apiSetChunckState,
	"UpdateClientsMap": apiUpdateClientsMap,
}

// ------------------------------- Incoming Structures -------------------------------
type setChunckStateStruct struct {
	UserID  int `json:"UserID"`
	ChunkID int `json:"ChunkID"`
}

// ------------------------------- Outgoing Structures -------------------------------
type updateMapStruct struct {
	Map        []byte `json:"Map"`
	ClientsIDs []int  `json:"ClientsIDs"`
}

type sendErrorMessageStruct struct {
	UserID       int    `json:"UserID"`
	ErrorMessage string `json:"ErrorMessage"`
}

func apiClientConnect(data string) {
	var userID int
	err := json.Unmarshal([]byte(data), &userID)

	if err != nil {
		logger.ErrorPrintf("Ошибка распаковки JSON: \nОшибка: %v \nДанные: %v", err, data)
	}

	status, message := validateClientConnect(userID)
	callbackMessage := callbackStruct{
		RoomID:  Room.ID,
		UserID:  userID,
		Status:  status,
		Message: message,
	}

	CreateMessage(callbackMessage, "CallbackClientConnect")

	if status {
		clientConnect(userID)
	}
}

func apiClientDisconnect(data string) {
	var userID int
	err := json.Unmarshal([]byte(data), &userID)

	if err != nil {
		logger.ErrorPrintf("Ошибка распаковки JSON: \nОшибка: %v \nДанные: %v", err, data)
	}

	status, message := validateClientDisconnect(userID)
	callbackMessage := callbackStruct{
		RoomID:  Room.ID,
		UserID:  userID,
		Status:  status,
		Message: message,
	}

	CreateMessage(callbackMessage, "CallbackClientDisconnect")

	if status {
		clientDisconnect(userID)
	}
}

func apiSetChunckState(data string) {
	var сhunckState setChunckStateStruct
	err := json.Unmarshal([]byte(data), &сhunckState)

	if err != nil {
		logger.ErrorPrintf("Ошибка распаковки JSON: \nОшибка: %v \nДанные: %v", err, data)
	}

	userID := сhunckState.UserID
	chunkID := сhunckState.ChunkID

	status, message := validateSetChunckState(userID, chunkID)
	callbackMessage := callbackStruct{
		RoomID:  Room.ID,
		UserID:  userID,
		Status:  status,
		Message: message,
	}

	CreateMessage(callbackMessage, "CallbackSetChunckState")

	if status {
		setChunckState(userID, chunkID)
	}
}

func apiUpdateClientsMap(data string) {
	var clientsIDs []int
	var clientsIDsToUpdate []int
	err := json.Unmarshal([]byte(data), &clientsIDs)

	if err != nil {
		logger.ErrorPrintf("Ошибка распаковки JSON: \nОшибка: %v \nДанные: %v", err, data)
	}

	for _, userID := range clientsIDs {
		status, message := validateUpdateClientsMap(userID)
		callbackMessage := callbackStruct{
			RoomID:  Room.ID,
			UserID:  userID,
			Status:  status,
			Message: message,
		}

		CreateMessage(callbackMessage, "CallbackUpdateClientsMap")

		if status {
			clientsIDsToUpdate = append(clientsIDsToUpdate, userID)
		}
	}

	if len(clientsIDsToUpdate) > 0 {
		updateClientsMap(clientsIDsToUpdate)
	}
}
