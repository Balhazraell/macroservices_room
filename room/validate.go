package room

import (
	"fmt"

	"github.com/Balhazraell/logger"
	"github.com/Balhazraell/tools"
)

func validateAPIcall(handlerName string) (bool, string){
	status := true
	message := ""

	_, ok := APIMetods[handlerName]
	if !ok {
		logger.WarningPrintf("Попытка вызвать API которго нет или к которому нет доступа: %v.", handlerName)
		status = false
		message = fmt.Sprintf("API метода %v в комнате room_id:%v не существует.", handlerName, Room.ID)
	}

	return status, message
}

func validateClientConnect(userID int) (bool, string){
	status := true
	message := ""

	// Проверка идентификатора на существование.
	var elementIndex = tools.FindElementInArray(Room.clients, userID)
	if elementIndex != -1 {
		logger.WarningPrintf("Попытка подключения пользователя, идентификатор которого уже есть: room_id:%v user_id:%v.", Room.ID, userID)
		status = false
		message = "Пользователь с таким id уже есть!"
		return status, message
	}

	return status, message
}

func validateClientDisconnect(userID int) (bool, string){
	status := true
	message := ""

	// Проверка идентификатора на существование.
	var elementIndex = tools.FindElementInArray(Room.clients, userID)
	if elementIndex == -1 {
		logger.WarningPrintf("Попытка отключиться пользователя, которого нет в списке пользователей: room_id:%v user_id:%v.", Room.ID, userID)
		status = false
		message = "Пользователь с таким id не существует!"
		return status, message
	}

	return status, message
}

func validateSetChunckState(userID int, chunkID int) (bool, string){
	status := true
	message := ""

	// Проверка на повторное изменение клетки.
	if Room.Map[chunkID].State != ChuncStateEmpty {
		logger.WarningPrintf("Попытка задать значение в поле с заданным значением: room_id:%v user_id:%v.", Room.ID, userID)
		status = false
		message = "Значение клетки уже задано.\nПовторно изменить нельзя!"
		return status, message
	}

	if Room.lastMovedUser == userID{
		logger.WarningPrintf("Попытка повторного хода игрока: room_id:%v user_id:%v.", Room.ID, userID)
		status = false
		message = "Сейчас ход другого игрока."
		return status, message
	}

	return status, message
}

func validateUpdateClientsMap(userID int) (bool, string){
	status := true
	message := ""

	var elementIndex = tools.FindElementInArray(Room.clients, userID)
	if elementIndex == -1 {
		logger.WarningPrintf("Попытка обновить карту у пользователя, которого нет в данной комнате: room_id:%v user_id:%v.", Room.ID, userID)
		status = false
		message = "Пользователь с таким id не существует!"
		return status, message
	}

	return status, message
}