package room

import (
	"encoding/json"

	"github.com/Balhazraell/logger"
)

type callbackStruct struct{
	RoomID int `json:"RoomID"`
	UserID int `json:"UserID"`
	Status bool `json:"Status"`
	Message string `json:"Message"`
}

// CallbackMetods - Перечень методов для получения ответов при запросах.
var CallbackMetods = map[string]func(string){
	"СallbackAPICall": сallbackAPICall,
}

func сallbackAPICall(data string){
	var callback = callbackStruct{}
	err := json.Unmarshal([]byte(data), &callback)

	if err != nil {
		logger.ErrorPrintf("Ошибка распаковки JSON: \nОшибка: %v \nДанные: %v", err, data)
	}

	if !callback.Status {
		logger.ErrorPrintf("Ошибка вызова API метода: \n%v", callback.Message)
	}
}