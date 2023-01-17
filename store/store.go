package store

import (
	"database/sql"
	"fmt"
	"github.com/xuri/excelize/v2"
	"strconv"
	"strings"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type OtherScores struct {
	UnknownPlayer     int
	LockScore         int
	UnknownPlusLocked int
}

type Storage interface {
	CommonGateway
	GetScores() map[string]int
	SetScores(scores map[string]int)
	GetPlaceScores() []int
	SetPlaceScores(placeScores []int)
	GetUnknownNames() map[string]bool
	SetUnknownNames(names map[string]bool)
	GetOtherScores() OtherScores
	SetOtherScores(oscores OtherScores)
	parseIntFromCell(f *excelize.File, sheet, axis string) (int, error)
	UpdateScores() []string
	IsKnownPlayer(name string) bool
	GetScore(name string, blocked bool, stake float64) int
	addLineToUnknownNames(name string)

	QueryRow(string, ...interface{}) *sql.Row
	Query(string, ...interface{}) (*sql.Rows, error)
	Exec(string, ...interface{}) (sql.Result, error)
}
type Store struct {
	*goqu.Database
}

func (d *Store) GetScores() map[string]int {
	type Player struct {
		name  string `db:"name"`
		score int    `db:"score"`
	}
	players := make([]Player, 0)
	d.GetElements("players", &players, 0, 0)
	scores := make(map[string]int)
	for _, p := range players {
		scores[p.name] = p.score
	}
	return scores
}
func (d *Store) SetScores(scores map[string]int) {
	tx, err := d.Begin()
	if err != nil {
		return
	}
	//d.Delete("*").From("players").Executor().Exec()
	//insert all new scores
	tx.Commit()
}
func (d *Store) GetPlaceScores() []int {
	return nil
}
func (d *Store) SetPlaceScores(scores []int) {
}
func (d *Store) GetUnknownNames() map[string]bool {
	return nil
}
func (d *Store) SetUnknownNames(names map[string]bool) {
}

func (d *Store) GetOtherScores() OtherScores {
	return OtherScores{}
}
func (d *Store) SetOtherScores(oscores OtherScores) {
}
func NewStore(conn, migration string) (Storage, error) {
	var err error

	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, fmt.Errorf("open db error: %v", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("connect db error: %v", err)
	}

	// накатываем миграции
	mg, err := migrate.New(migration, conn)
	if err != nil {
		return nil, fmt.Errorf("open migration error: %v", err)
	}

	err = mg.Up()
	switch err {
	case nil:
		break
	case migrate.ErrNoChange:
		err = nil
	default:
		return nil, fmt.Errorf("apply migration error: %v", err)
	}
	_, _ = mg.Close()

	d := goqu.New("postgres", db)
	return &Store{d}, err
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

	nscores := map[string]int{}
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
				_, ok := nscores[lowervalue]
				if ok {
					repeatingNames = append(repeatingNames, value)
				}
				nscores[lowervalue] = score
			}
		}
	}
	s.SetScores(nscores)
	/*s.LockScore, err = s.parseIntFromCell(f, lists[0], "Q2")
	if err != nil {
		return repeatingNames
	}*/
	oscores := OtherScores{LockScore: 0}
	oscores.UnknownPlusLocked, err = s.parseIntFromCell(f, lists[0], "P2")
	if err != nil {
		return repeatingNames
	}
	oscores.UnknownPlayer, err = s.parseIntFromCell(f, lists[0], "R2")
	if err != nil {
		return repeatingNames
	}
	pscores := make([]int, 0)
	for i, letter := 0, 'U'; letter <= 'Z'; i, letter = i+1, letter+1 {
		pscores[i], err = s.parseIntFromCell(f, lists[0], string([]byte{byte(letter), byte('2')}))
	}
	s.SetPlaceScores(pscores)
	s.SetOtherScores(oscores)
	return repeatingNames

}
func (s *Store) IsKnownPlayer(name string) bool {
	lowername := strings.ToLower(name)
	_, ok := s.GetScores()[lowername]
	return ok
}
func (s *Store) GetScore(name string, blocked bool, stake float64) int {
	lowername := strings.ToLower(name)
	sc, ok := s.GetScores()[lowername]
	if ok {
		return sc
	} else {
		if stake >= 80.0 {
			s.addLineToUnknownNames(name)
		}

		if blocked {
			return s.GetOtherScores().UnknownPlusLocked
		} else {
			return s.GetOtherScores().UnknownPlayer
		}
	}
}
func (s *Store) addLineToUnknownNames(name string) {
	_, contains := s.GetUnknownNames()[name]
	if contains {
		return
	}
	f, err := excelize.OpenFile("out.xlsx")
	if err != nil {
		f = excelize.NewFile()
		f.NewSheet("Sheet1")
	}

	unknownNames := map[string]bool{}
	for number := 1; ; number++ {
		axis := "A" + strconv.Itoa(number)
		val, err := f.GetCellValue("Sheet1", axis)
		if err != nil || val == "" {
			f.SetCellValue("Sheet1", axis, name)
			unknownNames[name] = true
			break
		} else {
			unknownNames[val] = true
		}
	}
	s.SetUnknownNames(unknownNames)
	err = f.SaveAs("out.xlsx")
	if err != nil {
		f.Save()
	}
}
