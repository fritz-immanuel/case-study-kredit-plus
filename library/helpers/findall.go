package helpers

import (
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"case-study-kredit-plus/library/types"

	"github.com/gin-gonic/gin"
	"github.com/leekchan/accounting"
)

type TableStatus struct {
	Name string
	ID   int
}

type ParentTableStatus struct {
	ID       int
	StatusID int
}

type StatusClient struct {
	Err    error
	Status int
	Jwt    string
	IP     [][]string
}

type LastID struct {
	ID uint64
}

type buffer struct {
	r         []byte
	runeBytes [utf8.UTFMax]byte
}

type FindAllConditionParam struct {
	FindAll      bool
	CountPage    int
	CountSize    int
	StartData    int
	ErrorStatus  int
	ErrorMessage string
	Error        error
}

const ()

// find all
func FilterFindAll(c *gin.Context) (string, string) {
	page := c.Query("Page")
	size := c.Query("Size")
	if c.Query("Page") == "" {
		page = "-1"
	}
	if c.Query("Size") == "" {
		size = "10"
	}

	return page, size
}

// find all multifunction
func FilterFindAllParam(c *gin.Context) types.FindAllParams {
	var statusID string
	var sort string

	sortName := Underscore(c.Query("SortName"))
	sortBy := strings.ToLower(c.Query("SortBy"))

	if c.Query("SortName") == "" {
		sortName = "id"
	}

	if c.Query("SortBy") == "" {
		sortBy = "DESC"
	}

	// if c.Query("StatusID") == "" {
	statusID = c.Query("StatusID")
	// }

	explodeStatus := strings.Split(statusID, ",")
	for _, vStatus := range explodeStatus {
		if vStatus != "-1" && vStatus != "" {
			JoinStringStatus := strings.Join(explodeStatus, "','")
			statusID = "status_id IN ('" + JoinStringStatus + "')"
			break
		} else {
			statusID = ""
			break
		}
	}

	if sortName != "" {
		sort = GetSortBy(sortName, sortBy)
	}

	dataFinder := DataFinder(c.Query("KeywordName"), c.Query("Keyword"))
	page, _ := strconv.Atoi(c.Query("Page"))
	size, _ := strconv.Atoi(c.Query("Size"))
	findallparams := types.FindAllParams{Page: page, Size: size, StatusID: statusID, DataFinder: dataFinder, SortName: sortName, SortBy: sort}

	return findallparams
}

func sanitize(text string) string {
	return strings.NewReplacer("'", "", `"`, "").Replace(text)
}

// keyword like full text search
func DataFinder(keywordname string, keyword string) string {
	str := "1=1"
	if keywordname != "" && keyword != "" {
		ExplodeParam := strings.Split(keywordname, ",")
		str += " AND ( "
		strTmp := ""
		for _, vParam := range ExplodeParam {
			date := strings.Contains(vParam, "date")
			if date {
				t, errDate := time.Parse("2006-01-02", keyword)
				if errDate == nil {
					keyword = t.Format("2006-01-02")
				}

				t, errDate = time.Parse("02-01-2006", keyword)
				if errDate == nil {
					keyword = t.Format("2006-01-02")
				}
			}

			if strTmp != "" {
				strTmp += " or "
			}

			strTmp += " " + sanitize(Underscore(vParam)) + " LIKE '%" + keyword + "%' "
		}
		str += strTmp
		str += " )"
	}

	return str
}

func GetSortBy(sortName string, sortBy string) string {
	var sort string
	var sortNameArr []string
	var sortByArr []string

	checkMultipleSortName := strings.Contains(sortName, ",")
	checkMultipleSortBy := strings.Contains(sortBy, ",")
	if checkMultipleSortName {
		explodeSortName := strings.Split(sortName, ",")
		sortNameArr = append(sortNameArr, explodeSortName...)
	} else {
		sortNameArr = append(sortNameArr, sortName)
	}

	if checkMultipleSortBy {
		explodeSortBy := strings.Split(sortBy, ",")
		sortByArr = append(sortByArr, explodeSortBy...)
	} else {
		sortByArr = append(sortByArr, sortBy)
	}

	for k, v := range sortNameArr {
		var str string
		lenSortBy := len(sortByArr)
		lenSortName := len(sortNameArr)
		if lenSortBy-1 >= k {
			str = v + " " + sortByArr[k]
		} else {
			str = v + " " + sortByArr[lenSortBy-1]
		}

		if lenSortName-1 != k {
			str = str + ","
		}

		sort = sort + str
	}

	return sort
}

func (b *buffer) write(r rune) {
	if r < utf8.RuneSelf {
		b.r = append(b.r, byte(r))
		return
	}
	n := utf8.EncodeRune(b.runeBytes[0:], r)
	b.r = append(b.r, b.runeBytes[0:n]...)
}

func (b *buffer) indent() {
	if len(b.r) > 0 {
		b.r = append(b.r, '_')
	}
}

func (b *buffer) indentSpace() {
	if len(b.r) > 0 {
		b.r = append(b.r, ' ')
	}
}

// set camelcase model name to table name with underscore
func Underscore(s string) string {
	b := buffer{
		r: make([]byte, 0, len(s)),
	}
	var m rune
	var w bool
	for _, ch := range s {
		if unicode.IsUpper(ch) {
			if m != 0 {
				if !w {
					b.indent()
					w = true
				}
				b.write(m)
			}
			m = unicode.ToLower(ch)
		} else if unicode.IsSpace(ch) {
			if m != 0 {
				b.indentSpace()
				m = 0
				w = false
			}
		} else {
			if m != 0 {
				b.indent()
				b.write(m)
				m = 0
				w = false
			}
			b.write(ch)
		}
	}
	if m != 0 {
		if !w {
			b.indent()
		}
		b.write(m)
	}

	// handle ID camel case
	strReplace := []byte(string(b.r))
	countID := strings.Count(string(strReplace), "i_d")
	if countID >= 1 {
		len := len(strReplace)
		for i := 0; i < len; i++ {
			if strReplace[i] == 'i' {
				if strReplace[i+1] == '_' {
					if strReplace[i+2] == 'd' {
						strReplace[i+1] = ' '
					}
				}
			}
		}
	}
	return strings.Replace(string(strReplace), " ", "", -1)
}

// format rupiah
func ConvertRupiah(value int, symbol bool) string {
	var strSymbol string
	if symbol {
		strSymbol = "Rp. "
	}
	ac := accounting.Accounting{Symbol: strSymbol, Precision: 2}

	Strings := ac.FormatMoney(value)

	return Strings
}
