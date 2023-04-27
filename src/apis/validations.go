package apis

import (
	//"net/http"

	"github.com/ini8labs/lsdb"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func validateUserID(str string) (primitive.ObjectID, error) {
	ObjectID, err := strToPrimitiveObjID(str)
	return ObjectID, err
}

func validateEventID(str string) (primitive.ObjectID, error) {
	ObjectID, err := strToPrimitiveObjID(str)
	return ObjectID, err
}

func validateBetUID(str string) (primitive.ObjectID, error) {
	ObjectID, err := strToPrimitiveObjID(str)
	return ObjectID, err
}

// func validateBetnumbers(str string) ([]int, error) {
// 	split := strings.Split(str, ",")
// 	strToInt := []int{}

// 	if len(split) < 1 {
// 		return nil, errMinNumbers
// 	}

// 	if len(split) > 5 {
// 		return nil, errMaxNumbers
// 	}

// 	for _, i := range split {
// 		j, err := strconv.Atoi(i)
// 		if err != nil || j < 1 || j > 90 {
// 			return nil, errNumberNotAllowed
// 		}

// 		strToInt = append(strToInt, j)
// 	}
// 	return strToInt, nil
// }

func validateBetnumbers(nums []int) ([]int, error) {
	if len(nums) < 1 {
		return nil, errMinNumbers
	}

	if len(nums) > 5 {
		return nil, errMaxNumbers
	}

	err := hasDuplicates(nums)
	if err != nil {
		return nil, errDuplicatedNumbers
	}

	for _, num := range nums {
		if num < 1 || num > 90 {
			return nil, errNumberNotAllowed
		}
	}

	return nums, nil
}

func hasDuplicates(nums []int) error {
	seen := make(map[int]bool)
	for _, num := range nums {
		if seen[num] {
			return errDuplicatedNumbers
		}
		seen[num] = true
	}
	return nil
}

func validateAmount(amount int) (int, error) {
	if amount < 1 {
		return 0, errMinAmount
	}
	return amount, nil
}

func requiredInfo(resp []lsdb.EventParticipantInfo) []UserBetsInfo {
	var respSlice []UserBetsInfo
	for i := 0; i < len(resp); i++ {
		var tempResp UserBetsInfo
		tempResp.BetUID = resp[i].BetUID
		tempResp.Amount = resp[i].Amount
		tempResp.BetNumbers = resp[i].BetNumbers

		respSlice = append(respSlice, tempResp)
	}
	return respSlice
}

func strToPrimitiveObjID(str string) (primitive.ObjectID, error) {
	eventUIDConv, err := primitive.ObjectIDFromHex(str)
	if err != nil {
		return primitive.NilObjectID, errInvalidEventID
	}

	return eventUIDConv, nil
}
