package dahuarpc

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type AuthParam struct {
	Encryption string `json:"encryption"`
	Random     string `json:"random"`
	Realm      string `json:"realm"`
}

// HashPassword runs the hashing algorithm for the password.
func (a AuthParam) HashPassword(username, password string) string {
	switch a.Encryption {
	case "Basic":
		return base64.StdEncoding.EncodeToString(fmt.Appendf(nil, "%s:%s", username, password))
	case "Default":
		return strings.ToUpper(fmt.Sprintf("%x",
			md5.Sum(fmt.Appendf(nil,
				"%s:%s:%s",
				username,
				a.Random,
				strings.ToUpper(fmt.Sprintf(
					"%x",
					md5.Sum(fmt.Appendf(nil, "%s:%s:%s", username, a.Realm, password))))))))
	default:
		return password
	}
}

type Timestamp string

func NewTimestamp(date time.Time, deviceLocation *time.Location) Timestamp {
	if date.IsZero() {
		return ""
	}
	return Timestamp(date.In(deviceLocation).Format("2006-01-02 15:04:05"))
}

func (t Timestamp) Parse(deviceLocation *time.Location) (time.Time, error) {
	if t == "" {
		return time.Time{}.UTC(), nil
	}

	format := "2006-01-02 15:04:05"
	if strings.HasSuffix(string(t), "PM") || strings.HasSuffix(string(t), "AM") {
		format = "2006-01-02 03:04:05 PM"
	}

	date, err := time.ParseInLocation(format, string(t), deviceLocation)
	if err != nil {
		return date, err
	}

	return date.UTC(), nil
}

// ExtractFilePathTags extracts tags from the file path.
// Tags are strings surrounded by brackets
func ExtractFilePathTags(filePath string) []string {
	search := filePath
	idx := strings.LastIndex(filePath, "/")
	if idx != -1 {
		search = filePath[idx:]
	}

	var tags []string
	tokens := strings.Split(search, "[")
	for i := 1; i < len(tokens); i++ {
		if end := strings.Index(tokens[i], "]"); end != -1 {
			tags = append(tags, tokens[i][:end])
		}
	}

	return tags
}

// Integer is for types that are supposed to be integers but for some reason the device returns a float.
type Integer int64

func (s *Integer) UnmarshalJSON(data []byte) error {
	var number float64
	if err := json.Unmarshal(data, &number); err != nil {
		return err
	}

	*s = Integer(number)

	return nil
}

func (s Integer) Integer() int64 {
	return int64(s)
}

// URL is the HTTP RPC API URL.
func URL(u *url.URL) string {
	return fmt.Sprintf("%s://%s/RPC2", u.Scheme, u.Hostname())
}

// LoginURL is the HTTP RPC API URL for login.
func LoginURL(u *url.URL) string {
	return fmt.Sprintf("%s://%s/RPC2_Login", u.Scheme, u.Hostname())
}

// LoadFileURL is the HTTP URL for accessing files.
// The file path must be absolute.
func LoadFileURL(u *url.URL, filePath string) string {
	return fmt.Sprintf("%s://%s/RPC_Loadfile%s", u.Scheme, u.Hostname(), filePath)
}

// Cookie creates a session cookie.
func Cookie(session string) string {
	return fmt.Sprintf("WebClientSessionID=%s; DWebClientSessionID=%s; DhWebClientSessionID=%s", session, session, session)
}

func DefaultTimeSection() TimeSection {
	return TimeSection{
		Number: 0,
		Start:  0,
		End:    24 * time.Hour,
	}
}

// NewTimeSectionFromString (e.g. "1 08:01:45-16:16:22").
func NewTimeSectionFromString(s string) (TimeSection, error) {
	splitBySpace := strings.Split(s, " ")
	if len(splitBySpace) != 2 {
		return TimeSection{}, fmt.Errorf("invalid number of spaces: %d", len(splitBySpace))
	}

	splitByDash := strings.Split(splitBySpace[1], "-")
	if len(splitByDash) != 2 {
		return TimeSection{}, fmt.Errorf("invalid number of dashes: %d", len(splitByDash))
	}

	start, err := durationFromTimeString(splitByDash[0])
	if err != nil {
		return TimeSection{}, err
	}

	end, err := durationFromTimeString(splitByDash[1])
	if err != nil {
		return TimeSection{}, err
	}

	number, err := strconv.Atoi(splitBySpace[0])
	if err != nil {
		return TimeSection{}, err
	}

	return TimeSection{
		Number: number,
		Start:  start,
		End:    end,
	}, nil
}

func NewTimeSection(number int, start, end time.Duration) TimeSection {
	return TimeSection{
		Number: number,
		Start:  start,
		End:    end,
	}
}

// durationFromTimeString (e.g. "08:01:45")
func durationFromTimeString(s string) (time.Duration, error) {
	arr := strings.Split(s, ":")
	if len(arr) != 3 {
		return 0, fmt.Errorf("invalid number of colons: %d", len(arr))
	}

	var numbers [3]int
	for i := range arr {
		var err error
		numbers[i], err = strconv.Atoi(arr[i])
		if err != nil {
			return 0, err
		}
	}

	return time.Duration(numbers[0])*time.Hour + time.Duration(numbers[1])*time.Minute + time.Duration(numbers[2])*time.Second, nil
}

type TimeSection struct {
	Number int
	Start  time.Duration
	End    time.Duration
}

func (s *TimeSection) UnmarshalJSON(data []byte) error {
	var str string
	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}

	res, err := NewTimeSectionFromString(str)
	if err != nil {
		return err
	}

	s.Number = res.Number
	s.Start = res.Start
	s.End = res.End

	return nil
}

func (s TimeSection) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s TimeSection) String() string {
	return fmt.Sprintf(
		"%d %02d:%02d:%02d-%02d:%02d:%02d",
		s.Number,
		int(s.Start.Hours()),
		int(s.Start.Minutes())%60,
		int(s.Start.Seconds())%60,
		int(s.End.Hours()),
		int(s.End.Minutes())%60,
		int(s.End.Seconds())%60,
	)
}

// NewTimeSection2FromString (e.g. "01:15:00-05:00:00 Night").
func NewTimeSection2FromString(s string) (TimeSection2, error) {
	splitBySpace := strings.Split(s, " ")
	if len(splitBySpace) != 2 {
		return TimeSection2{}, fmt.Errorf("invalid number of spaces: %d", len(splitBySpace))
	}

	splitByDash := strings.Split(splitBySpace[0], "-")
	if len(splitByDash) != 2 {
		return TimeSection2{}, fmt.Errorf("invalid number of dashes: %d", len(splitByDash))
	}

	start, err := durationFromTimeString(splitByDash[0])
	if err != nil {
		return TimeSection2{}, err
	}

	end, err := durationFromTimeString(splitByDash[1])
	if err != nil {
		return TimeSection2{}, err
	}

	return TimeSection2{
		Start:   start,
		End:     end,
		Profile: splitBySpace[1],
	}, nil
}

func TimeSectionDuration(t time.Time) time.Duration {
	return time.Duration(t.Hour())*time.Hour + time.Duration(t.Minute())*time.Minute + time.Duration(t.Second())*time.Second
}

func NewTimeSection2(start, end time.Duration, profile string) TimeSection2 {
	return TimeSection2{
		Start:   start,
		End:     end,
		Profile: profile,
	}
}

type TimeSection2 struct {
	Start   time.Duration
	End     time.Duration
	Profile string
}

func (s *TimeSection2) UnmarshalJSON(data []byte) error {
	var str string
	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}

	res, err := NewTimeSection2FromString(str)
	if err != nil {
		return err
	}

	s.Start = res.Start
	s.End = res.End
	s.Profile = res.Profile

	return nil
}

func (s TimeSection2) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s TimeSection2) String() string {
	return fmt.Sprintf(
		"%02d:%02d:%02d-%02d:%02d:%02d %s",
		int(s.Start.Hours()),
		int(s.Start.Minutes())%60,
		int(s.Start.Seconds())%60,
		int(s.End.Hours()),
		int(s.End.Minutes())%60,
		int(s.End.Seconds())%60,
		s.Profile,
	)
}
