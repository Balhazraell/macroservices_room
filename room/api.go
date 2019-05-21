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

type callbackStruct struct{
	RoomID int `json:"RoomID"`
	ClientID int `json:"ClientID"`
	Status bool `json:"Status"`
	Message string `json:"Message"`
}

// ------------------------------- Incoming Structures -------------------------------
type setChunckStateStruct struct {
	ClientID int `json:"ClientID"`
	ChunkID  int `json:"ChunkID"`
}

// ------------------------------- Outgoing Structures -------------------------------
type updateMapStruct struct {
	Map        []byte `json:"Map"`
	ClientsIDs []int  `json:"ClientsIDs"`
}

type sendErrorMessageStruct struct {
	ClientID     int    `json:"ClientID"`
	ErrorMessage string `json:"ErrorMessage"`
}

//--------------------- room struct --------------------//
type SetChunckStateStruct struct {
	ClientID int `json:"ClientID"`
	ChunkID  int `json:"ChunkID"`
}

func apiClientConnect(data string) {
	var clientID int
	err := json.Unmarshal([]byte(data), &clientID)

	if err != nil {
		logger.ErrorPrintf("Ошибка API при подключении нового клиента: %s;\n Ошибка в данных: %s", err, data)
	}

	clientConnect(clientID)
}

func apiClientDisconnect(data string) {
	var clientID int
	err := json.Unmarshal([]byte(data), &clientID)

	if err != nil {
		logger.ErrorPrintf("Ошибка API при отключении клиента: %s;\n Ошибка в данных: %s", err, data)
	}

	clientDisconnect(clientID)
}

func apiSetChunckState(data string) {
	var setChunckStateStruct setChunckStateStruct
	err := json.Unmarshal([]byte(data), &setChunckStateStruct)

	if err != nil {
		logger.ErrorPrintf("Ошибка API задании состояния чанку: %s;\n Ошибка в данных: %s", err, data)
	}

	setChunckState(setChunckStateStruct.ClientID, setChunckStateStruct.ChunkID)
}

func apiUpdateClientsMap(data string) {
	var clientsIDs []int
	err := json.Unmarshal([]byte(data), &clientsIDs)

	if err != nil {
		logger.ErrorPrintf("Ошибка при распаковывании сообщения об обновлении карты: %v", err)
	}

	updateClientsMap(clientsIDs)
}
