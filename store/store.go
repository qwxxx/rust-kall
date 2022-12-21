package store

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"strconv"
	"strings"
)

type Store struct {
	scores            map[string]int
	PlaceScores       []int
	UnknownPlayer     int
	LockScore         int
	UnknownPlusLocked int
	UnknownNames      map[string]bool
}

func NewStore() *Store {
	s := Store{}
	s.scores = map[string]int{}
	s.PlaceScores = make([]int, 6)
	s.UnknownNames = map[string]bool{}
	return &s
}
func (s *Store) parseIntFromCell(f *excelize.File, sheet, axis string) (int, error) {
	value, err := f.GetCellValue(sheet, axis)
	if err != nil {
		return 0, err
	}

	ret, err := strconv.Atoi(value)

	return ret, err
}
func (s *Store) UpdateScores() []string {

	repeatingNames := make([]string, 0)
	f, err := excelize.OpenFile("in.xlsx")
	if err != nil {
		return repeatingNames
	}
	lists := f.GetSheetList()
	if len(lists) == 0 {
		return repeatingNames
	}

	s.scores = map[string]int{}
	for letter, score := 'A', -2; letter <= 'O'; letter, score = letter+1, score+1 {
		emptyCount := 0
		for num := 2; ; num = num + 1 {
			numstr := strconv.Itoa(num)
			a := append([]byte{byte(letter)}, []byte(numstr)...)
			value, err := f.GetCellValue(lists[0], string(a))

			if value == "" || err != nil {
				emptyCount++
				if err != nil {
					fmt.Println(err)
				}
				if emptyCount >= 10 {
					break
				}
			} else {
				emptyCount = 0
				lowervalue := strings.ToLower(value)
				_, ok := s.scores[lowervalue]
				if ok {
					repeatingNames = append(repeatingNames, value)
				}
				s.scores[lowervalue] = score
			}
		}
	}

	/*s.LockScore, err = s.parseIntFromCell(f, lists[0], "Q2")
	if err != nil {
		return repeatingNames
	}*/
	s.LockScore = 0
	s.UnknownPlusLocked, err = s.parseIntFromCell(f, lists[0], "P2")
	if err != nil {
		return repeatingNames
	}
	s.UnknownPlayer, err = s.parseIntFromCell(f, lists[0], "R2")
	if err != nil {
		return repeatingNames
	}

	for i, letter := 0, 'U'; letter <= 'Z'; i, letter = i+1, letter+1 {
		s.PlaceScores[i], err = s.parseIntFromCell(f, lists[0], string([]byte{byte(letter), byte('2')}))
	}

	return repeatingNames

}
func (s *Store) IsKnownPlayer(name string) bool {
	lowername := strings.ToLower(name)
	_, ok := s.scores[lowername]
	return ok
}
func (s *Store) GetScore(name string, blocked bool, stake float64) int {
	lowername := strings.ToLower(name)
	sc, ok := s.scores[lowername]
	if ok {
		return sc
	} else {
		if stake >= 50.0 {
			s.addLineToUnknownNames(name)
		}

		if blocked {
			return s.UnknownPlusLocked
		} else {
			return s.UnknownPlayer
		}
	}
}
func (s *Store) addLineToUnknownNames(name string) {
	_, contains := s.UnknownNames[name]
	if contains {
		return
	}
	f, err := excelize.OpenFile("out.xlsx")
	if err != nil {
		f = excelize.NewFile()
		f.NewSheet("Sheet1")
	}
	s.UnknownNames = map[string]bool{}
	for number := 1; ; number++ {
		axis := "A" + strconv.Itoa(number)
		val, err := f.GetCellValue("Sheet1", axis)
		if err != nil || val == "" {
			f.SetCellValue("Sheet1", axis, name)
			s.UnknownNames[name] = true
			break
		} else {
			s.UnknownNames[val] = true
		}
	}
	err = f.SaveAs("out.xlsx")
	if err != nil {
		f.Save()
	}
}
