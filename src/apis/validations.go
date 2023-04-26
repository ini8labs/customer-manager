package apis

import (
	"strconv"
	"strings"
)

func (s Server) Convert(str string) ([]int, error) {
	split := strings.Split(str, ",")
	strToInt := []int{}

	if len(split) < 1 {
		s.Logger.Error(errMinNumbers)
		return nil, errMinNumbers
	}

	if len(split) > 5 {
		s.Logger.Error(errMaxNumbers)
		return nil, errMaxNumbers
	}

	for _, i := range split {
		j, err := strconv.Atoi(i)
		if err != nil || j < 1 || j > 90 {
			s.Logger.Error(errNumberNotAllowed)
			return nil, errNumberNotAllowed
		}

		strToInt = append(strToInt, j)
	}
	return strToInt, nil
}
