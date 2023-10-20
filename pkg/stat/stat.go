package stat

import (
	"fmt"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
	"github.com/nakabonne/tstorage"
	"github.com/pkg/errors"
	tele "gopkg.in/telebot.v3"
	"io"
	"strings"
	"time"
)

type Statistic struct {
	tst tstorage.Storage
}

func NewStatistic(path string) (*Statistic, error) {
	tst, err := tstorage.NewStorage(
		tstorage.WithTimestampPrecision(tstorage.Seconds),
		tstorage.WithDataPath(path),
	)
	if err != nil {
		return nil, errors.Wrap(err, "NewStatistic")
	}

	return &Statistic{tst: tst}, nil
}

func (s *Statistic) Stop() {
	s.tst.Close()
}

func (s *Statistic) Count(name string) {
	s.tst.InsertRows([]tstorage.Row{
		{
			Metric:    name,
			DataPoint: tstorage.DataPoint{Timestamp: time.Now().Unix(), Value: 1},
		},
	})
}

func (s *Statistic) Add(name string, value time.Duration) {
	s.tst.InsertRows([]tstorage.Row{
		{
			Metric:    name,
			DataPoint: tstorage.DataPoint{Timestamp: time.Now().Unix(), Value: value.Seconds()},
		},
	})
}

// HandlerStat /stat - статистика.
func (s *Statistic) HandlerStat(c tele.Context) error {
	xy, err := s.GetXY("", "", "h")
	if err != nil {
		return err
	}

	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("Время на сервере: %v\n", time.Now().String()))
	buf.WriteString("Час: Количество Время\n")

	for i := 0; i < len(xy.X); i++ {
		buf.WriteString(fmt.Sprintf("%v: %v %v\n", xy.X[i], xy.Ycount[i], xy.Ydur[i]))
	}

	return c.Send(buf.String())
}

func (s *Statistic) Chart(start string, stop string, delim string, w io.Writer) error {
	xy, err := s.GetXY(start, stop, delim)
	if err != nil {
		return err
	}

	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{
			Theme: types.ThemeInfographic,
		}),
		charts.WithTitleOpts(opts.Title{
			Title:    "Line chart",
			Subtitle: fmt.Sprintf("%s - %s", start, stop),
		}),
	)

	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{
			Theme: types.ThemeInfographic,
		}),
		charts.WithTitleOpts(opts.Title{
			Title:    "Line chart in Go",
			Subtitle: fmt.Sprintf("%s - %s", start, stop),
		}),
		charts.WithToolboxOpts(opts.Toolbox{
			Show: true,
			Feature: &opts.ToolBoxFeature{
				SaveAsImage: &opts.ToolBoxFeatureSaveAsImage{
					Show:  true,
					Type:  "png",
					Title: fmt.Sprintf("%s - %s", start, stop),
				},
			}},
		),
	)
	line.SetXAxis(xy.X).
		AddSeries("Count", convertInt(xy.Ycount)).
		AddSeries("Time", convertFloat(xy.Ydur)).
		SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: true}))

	return line.Render(w)
}

//func convert[t []int | []float64](sl t) []opts.LineData {
//	ld := make([]opts.LineData, 0, len(sl))
//	for _, v := range sl {
//		ld = append(ld, opts.LineData{Value: v})
//	}
//	return ld
//}

func convertFloat(sl []float64) []opts.LineData {
	ld := make([]opts.LineData, 0, len(sl))
	for _, v := range sl {
		ld = append(ld, opts.LineData{Value: v})
	}
	return ld
}

func convertInt(sl []int) []opts.LineData {
	ld := make([]opts.LineData, 0, len(sl))
	for _, v := range sl {
		ld = append(ld, opts.LineData{Value: v})
	}
	return ld
}

func (s *Statistic) GetXY(start string, stop string, delim string) (*XY, error) {
	tStart, err := time.Parse(time.DateTime, start)
	if err != nil {
		tStart = getStartOfToday()
	}

	tStop, err := time.Parse(time.DateTime, stop)
	if err != nil {
		tStop = time.Now()
	}

	if delim == "m" {
		if !isDateWithinLastDays(tStart, -2) {
			return nil, errors.New("invalid date")
		}
		return s.GetXYmin(tStart, tStop), nil
	}

	if !isDateWithinLastDays(tStart, -14) {
		return nil, errors.New("invalid date")
	}

	return s.GetXYhour(tStart, tStop), nil
}

func (s *Statistic) GetXYmin(start, stop time.Time) *XY {
	mins := splitByMinute(start, stop)
	xy := NewXY()
	for i, v := range mins {
		points, _ := s.tst.Select("req_durations", nil, v, v+59)
		s := sr(points)
		xy.Add(fmt.Sprintf("%v", i+1), len(points), s)
	}

	return xy
}

func (s *Statistic) GetXYhour(start, stop time.Time) *XY {
	hours := splitByHour(start, stop)
	xy := NewXY()
	for i, v := range hours {
		points, _ := s.tst.Select("req_durations", nil, v, v+(60*60-1))
		s := sr(points)
		xy.Add(fmt.Sprintf("%v", i+1), len(points), s)
	}

	return xy
}

func isDateWithinLastDays(date time.Time, days int) bool {
	now := time.Now().UTC()
	lastWeek := now.AddDate(0, 0, days)
	return !date.Before(lastWeek)
}

func sr(points []*tstorage.DataPoint) float64 {
	if points == nil {
		return 0
	}
	var sum float64
	for _, v := range points {
		sum = sum + v.Value
	}

	return sum / float64(len(points))
}

type XY struct {
	X      []string  // date
	Ycount []int     // count
	Ydur   []float64 // float
}

func NewXY() *XY {
	return &XY{X: make([]string, 0), Ycount: make([]int, 0), Ydur: make([]float64, 0)}
}

func (xy *XY) Add(x string, yCount int, yDur float64) {
	xy.X = append(xy.X, x)
	xy.Ycount = append(xy.Ycount, yCount)
	xy.Ydur = append(xy.Ydur, yDur)
}

func (s *Statistic) BotMiddleware() tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			startedAt := time.Now()
			next(c)
			elapsed := time.Since(startedAt)
			s.Add("req_durations", elapsed)
			return nil
		}
	}
}

func getStartOfToday() time.Time {
	now := time.Now().UTC()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	return startOfDay
}

func splitByHour(timeStart time.Time, timeEnd time.Time) []int64 {
	var result []int64
	duration := timeEnd.Sub(timeStart)
	hours := int(duration.Hours())
	for i := 0; i < hours; i++ {
		tss := timeStart.Add(time.Duration(i+1) * time.Hour).Unix()
		result = append(result, tss)
	}
	return result
}

func splitByMinute(timeStart time.Time, timeEnd time.Time) []int64 {
	var result []int64
	duration := timeEnd.Sub(timeStart)
	minutes := int(duration.Minutes())
	for i := 0; i < minutes; i++ {
		tss := timeStart.Add(time.Duration(i+1) * time.Minute).Unix()
		result = append(result, tss)
	}
	return result
}

func splitTimeByHour(timeStart time.Time, timeEnd time.Time) []time.Time {
	var result []time.Time
	duration := timeEnd.Sub(timeStart)
	hours := int(duration.Hours())
	for i := 0; i < hours; i++ {
		result = append(result, timeStart.Add(time.Duration(i)*time.Hour))
	}
	return result
}

func splitTimeByMinute(timeStart time.Time, timeEnd time.Time) []time.Time {
	var result []time.Time
	duration := timeEnd.Sub(timeStart)
	minutes := int(duration.Minutes())
	for i := 0; i < minutes; i++ {
		result = append(result, timeStart.Add(time.Duration(i)*time.Minute))
	}
	return result
}
