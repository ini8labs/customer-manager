package apis

import (
	//"net/http"

	"regexp"
	"strconv"
	"time"

	"github.com/ini8labs/lsdb"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func validateID(str string) (primitive.ObjectID, error) {
	ObjectID, err := strToPrimitiveObjID(str)
	return ObjectID, err
}

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

func requiredInfoUserBets(resp []lsdb.EventParticipantInfo) []UserBetsInfo {
	var respSlice []UserBetsInfo
	for i := 0; i < len(resp); i++ {
		var tempResp UserBetsInfo
		tempResp.BetUID = resp[i].BetUID.Hex()
		tempResp.Amount = resp[i].ParticipantInfo.Amount
		tempResp.BetNumbers = resp[i].ParticipantInfo.BetNumbers

		respSlice = append(respSlice, tempResp)
	}
	return respSlice
}

func requiredUserInfo(resp *lsdb.UserInfo) UserInformation {
	var respSlice UserInformation
	respSlice.Name = resp.Name
	respSlice.Phone = resp.Phone
	respSlice.EMail = resp.EMail
	return respSlice
}

func requiredEventInfo(resp []lsdb.LotteryEventInfo) []EventsInfo {
	var resultSlice []EventsInfo
	for i := 0; i < len(resp); i++ {
		var tempResp EventsInfo
		tempResp.EventUID = primitiveToString(resp[i].EventUID)
		tempResp.EventDate = convertPrimitiveToTime(resp[i].EventDate)
		tempResp.EventName = resp[i].Name
		tempResp.EventType = resp[i].EventType

		resultSlice = append(resultSlice, tempResp)
	}
	return resultSlice
}

func strToPrimitiveObjID(str string) (primitive.ObjectID, error) {
	eventUIDConv, err := primitive.ObjectIDFromHex(str)
	if err != nil {
		return primitive.NilObjectID, errInvalidID
	}

	return eventUIDConv, nil
}

func validateEventType(s string) error {
	for i := 0; i < len(eventsList); i++ {
		if eventsList[i] == s {
			return nil
		}
	}
	return errIncorrectEventType
}

// func getParticipantsInfoByEventIDLoop(resp []lsdb.LotteryEventInfo) {
// 	for i, _ := range resp {

// 		EventUID := resp[i].EventUID
// 		resp2, err := server.Client.GetParticipantsInfoByEventID(EventUID)
// 		if err != nil {
// 			Logrus.Error(err.Error())
// 			http.JSON(http.StatusBadRequest, err.Error())
// 			return
// 		}

// 		requiredBetsByEventType(resp2)

// 	}
// }

func requiredBetsByEventType(resp []lsdb.EventParticipantInfo) {
	for i, _ := range resp {
		var temp UserBetsInfoByEvent
		temp.Amount = resp[i].Amount
		temp.BetNumbers = resp[i].BetNumbers
		respSlice = append(respSlice, temp)
	}
}

func stringToPrimitive(s string) primitive.ObjectID {
	var a, _ = primitive.ObjectIDFromHex(s)
	return a
}

func convertTimeToPrimitive(date Date) primitive.DateTime {

	d := time.Date(date.Year, time.Month(date.Month), date.Day, 0, 0, 0, 0, time.Local)
	dNew := primitive.NewDateTimeFromTime(d)
	return dNew
}

func convertPrimitiveToTime(date primitive.DateTime) Date {
	t := date.Time()

	return Date{
		Day:   t.Day(),
		Month: int(t.Month()),
		Year:  t.Year(),
	}
}

func primitiveToString(p primitive.ObjectID) string {
	return p.Hex()
}

func initializeEventInfo(resp []lsdb.LotteryEventInfo) []EventsInfo {
	var arr []EventsInfo

	for i := 0; i < len(resp); i++ {
		eventinfo := EventsInfo{
			EventUID:  primitiveToString(resp[i].EventUID),
			EventDate: convertPrimitiveToTime(resp[i].EventDate),
			EventName: resp[i].Name,
			EventType: resp[i].EventType,
		}

		arr = append(arr, eventinfo)
	}
	return arr
}

func validateEventId(str string, resp []EventsInfo) bool {
	eventIdExist := true
	for i := 0; i < len(resp); i++ {
		if resp[i].EventUID == str {
			eventIdExist = true
			break
		}
		if resp[i].EventUID != str {
			eventIdExist = false
		}
	}
	return eventIdExist
}

func validatePhoneNumberString(phone string) error {

	pattern := `^[0-9]{10}$`

	regex := regexp.MustCompile(pattern)

	isValid := regex.MatchString(phone)
	if !isValid {
		return errIncorrectPhoneNo
	}

	return nil
}

func validatePhoneNumberInt(phone int64) error {
	phoneString := strconv.FormatInt(phone, 10)

	if err := validatePhoneNumberString(phoneString); err != nil {
		return err
	}
	return nil
}
