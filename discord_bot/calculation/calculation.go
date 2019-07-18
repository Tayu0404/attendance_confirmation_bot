package calculation

import (
	"fmt"
//	"log"
	"encoding/json"
	"io/ioutil"
	"time"
	"math"
	"github.com/jmoiron/sqlx"
	_"github.com/go-sql-driver/mysql"
	"github.com/Tayu0404/attendance_rec/discord_bot/modules"
)

type Person struct {
	Month int `json:"month"`
	Days  int `json:"days"`
}

func parseJson() []Person {
	bytes, err := ioutil.ReadFile("schedule.json")
	if err != nil {
		fmt.Println("ReadFile : ", err)
	}

	var attendanceDays []Person
	if err := json.Unmarshal(bytes, &attendanceDays); err != nil {
		fmt.Println("Unmarshal : ", err)
    }
	return attendanceDays
}

func AttendanceRate (db *sqlx.DB, u string) (int, int, float64) {
	sud := module.SelectUserDB(db, u)
	ad := parseJson()

	var days int
	
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	t := fmt.Sprintln(time.Now().In(jst).Format("200601"))

	for _, v := range ad {
		jt := fmt.Sprintln(v.Month)
		if t == jt {
			days = v.Days
		}
	}
	attendanceDays := days-len(sud)
	attendanceRate:= float64(attendanceDays)/float64(days)
	return days, len(sud), round(attendanceRate,3)
}

func round(val float64, place int) float64 {
	shift := math.Pow(10, float64(place))
	return math.Floor(val * shift + .5) / shift
}
