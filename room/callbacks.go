package room

import (
	"encoding/json"

	"github.com/Balhazraell/logger"
)

type callbackStruct struct {
	RoomID  int    `json:"RoomID"`
	UserID  int    `json:"UserID"`
	Status  bool   `json:"Status"`
	Message string `json:"Message"`
}

// CallbackMetods - Перечень методов для получения ответов при запросах.
var CallbackMetods = map[string]func(string){
	"СallbackAPICall":          сallbackAPICall,
	"CallbackRoomConnect":      callbackRoomConnect,
	"CallbackUpdateClientsMap": callbackUpdateClientsMap,
	"CallbackSendErrorMessage": callbackSendErrorMessage,
}

func сallbackAPICall(data string) {
	var callback = callbackStruct{}
	err := json.Unmarshal([]byte(data), &callback)

	if err != nil {
		logger.WarningPrintf("Ошибка распаковки JSON: \nОшибка: %v \nДанные: %v", err, data)
	}

	if !callback.Status {
		logger.WarningPrintf("Ошибка вызова API метода: \n%v", callback.Message)
	}
}

func callbackRoomConnect(data string) {
	var callback = callbackStruct{}
	err := json.Unmarshal([]byte(data), &callback)

	if err != nil {
		logger.WarningPrintf("Ошибка распаковки JSON: \nОшибка: %v \nДанные: %v", err, data)
	}

	if !callback.Status {
		// callbackRoomConnect() - нужно перегенерить id. (не забываем про закрытие канала и открытие нового.)
	}
}

func callbackUpdateClientsMap(data string) {
	var callback = callbackStruct{}
	err := json.Unmarshal([]byte(data), &callback)

	if err != nil {
		logger.WarningPrintf("Ошибка распаковки JSON: \nОшибка: %v \nДанные: %v", err, data)
	}

	if !callback.Status {
		logger.WarningPrintf("Произошла ошибка при обновлении карты пользователя с \nid: %v \nerror: %v", callback.UserID, callback.Message)
		clientDisconnect(callback.UserID)
	}
}

func callbackSendErrorMessage(data string) {
	
}
