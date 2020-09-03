package chart

import "sort"

type HighChartLineSeries struct {
	Name string      `json:"name"`
	Data [][]float64 `json:"data"`
}

func (s *HighChartLineSeries) AddDataPoint(t float64, v float64) {
	s.Data = append(s.Data, []float64{t, v})
}

func NewHighChartLineSeries(name string) *HighChartLineSeries {
	return &HighChartLineSeries{Name: name, Data: make([][]float64, 0)}
}

type HighChartLineSeriesGroup struct {
	GroupData map[string]*HighChartLineSeries `json:"group_data"`
}

func NewHighChartLineSeriesGroup() *HighChartLineSeriesGroup {
	return &HighChartLineSeriesGroup{GroupData: make(map[string]*HighChartLineSeries, 0)}
}

func (grp *HighChartLineSeriesGroup) GetOrCreateSeries(name string) *HighChartLineSeries {
	s, ok := grp.GroupData[name]
	if !ok {
		s = NewHighChartLineSeries(name)
		grp.GroupData[name] = s
	}
	return s
}

func (grp *HighChartLineSeriesGroup) GetData() []HighChartLineSeries {
	result := make([]HighChartLineSeries, 0)
	for _, series := range grp.GroupData {
		result = append(result, *series)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})
	return result
}
