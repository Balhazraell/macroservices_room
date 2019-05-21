package room

import (
	"encoding/json"

	"github.com/Balhazraell/logger"
	"github.com/Balhazraell/tools"
)

// ------------------- Вспомогательные методы -------------------
func createMap() {
	var step = 100
	var y = 0
	var chunckIDCounter = 0

	for i := 0; i < 3; i++ {
		var x = 0
		for j := 0; j < 3; j++ {
			chunc := chunc{
				ID:    chunckIDCounter,
				State: ChuncStateEmpty,
			}
			chunc.Сoordinates = append(
				chunc.Сoordinates,
				[2]int{x, y},
				[2]int{x + step, y},
				[2]int{x + step, y + step},
				[2]int{x, y + step},
				// необходимо указывать 5 элементов, так как последняя точка замыкает фигуру,
				// это нужно для определения пересечения координат мышки и элементов.
				[2]int{x, y},
			)

			x += step

			Room.Map[chunckIDCounter] = &chunc
			chunckIDCounter++
		}
		y += step
	}
}

//---------------------------------------------------------------
func updateClientsMap(clientsIDs []int) {
	logger.InfoPrint("Обновление карт пользователей.")
	gameMap, err := json.Marshal(Room.Map)

	if err != nil {
		logger.WarningPrintf("При формировнии json при упаковке карты, для обновления у клиентов произошла ошибка %v.", err)
		return
	}

	updateMap := updateMapStruct{
		Map:        gameMap,
		ClientsIDs: clientsIDs,
	}

	data, err := json.Marshal(updateMap)

	if err != nil {
		logger.WarningPrintf("При формировнии json при создании сооббщения для обновления карт произошла ошибка %v.", err)
		return
	}

	newMessage := MessageRMQ{
		HandlerName: "UpdateClientsMap",
		Data:        string(data),
	}

	PublishMessage(newMessage)
}

//------------------------- Client Connect -------------------------//

func clientConnect(clientID int){
	logger.InfoPrintf("К комнате %v подключился новый клиент с id=%v.", Room.ID, clientID)

	status, message := validateClientConnect(clientID)
	callbackMessage := callbackStruct{
		RoomID:   Room.ID,
		ClientID: clientID,
		Status:   status,
		Message:  message,
	}

	CreateMessage(callbackMessage, "ClientConnectCallback")

	if status {
		updateClientsMap([]int{clientID})
	}
}

func validateClientConnect(clientID int) (bool, string){
	status := true
	message := ""

	// Проверка идентификатора на существование.
	var elementIndex = tools.FindElementInArray(Room.clients, clientID)
	if elementIndex != -1 {
		logger.WarningPrintf("Попытка подключения пользователя, идентификатор которого уже есть: room_id:%v user_id:%v.", Room.ID, clientID)
		status = false
		message = "Пользователь с таким id уже есть!"
		return status, message
	}

	return status, message
}

//------------------------- Client Disconnect -------------------------//

func clientDisconnect(clientID int){
	status, message := validateClientConnect(clientID)
	callbackMessage := callbackStruct{
		RoomID:   Room.ID,
		ClientID: clientID,
		Status:   status,
		Message:  message,
	}

	CreateMessage(callbackMessage, "ClientDisconnectCallback")

	if status {
		Room.clients = tools.DeleElementFromArraByIndex(Room.clients, clientID)
	}
}

func validateClientDisconnect(clientID int) (bool, string){
	status := true
	message := ""

	// Проверка идентификатора на существование.
	var elementIndex = tools.FindElementInArray(Room.clients, clientID)
	if elementIndex == -1 {
		logger.WarningPrintf("Попытка отключиться пользователя, которого нет в списке пользователей: room_id:%v user_id:%v.", Room.ID, clientID)
		status = false
		message = "Пользователь с таким id не существует!"
		return status, message
	}

	return status, message
}

//------------------------- Set Chunck State -------------------------//
func setChunckState(clientID int, chunkID int) {
	status, message := validatorSetChunckState(clientID, chunkID)
	callbackMessage := callbackStruct{
		RoomID:   Room.ID,
		ClientID: clientID,
		Status:   status,
		Message:  message,
	}

	CreateMessage(callbackMessage, "SetChunckStateCallback")

	if status {
		Room.Map[chunkID].State = Room.GameState

		if Room.GameState == ChuncStateCross {
			Room.GameState = ChuncStateZero
		} else {
			Room.GameState = ChuncStateCross
		}

		updateClientsMap(Room.clients)
	}
}

func validatorSetChunckState(clientID int, chunkID int) (bool, string){
	status := true
	message := ""

	// Проверка на повторное изменение клетки.
	if Room.Map[chunkID].State != ChuncStateEmpty {
		logger.WarningPrintf("Попытка задать значение в поле с заданным значением: room_id:%v user_id:%v.", Room.ID, clientID)
		status = false
		message = "Значение клетки уже задано.\nПовторно изменить нельзя!"
		return status, message
	}

	if Room.lastMovedUser == clientID{
		logger.WarningPrintf("Попытка повторного хода игрока: room_id:%v user_id:%v.", Room.ID, clientID)
		status = false
		message = "Сейчас ход другого игрока."
		return status, message
	}

	return status, message
}
