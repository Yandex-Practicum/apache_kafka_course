package domain

import (
	"encoding/json"
	"errors"
)

type MType = string
type MName = string

var NotExistsType = errors.New("type does not exist")

var (
	Gauge   MType = "gauge"
	Counter MType = "counter"
)

var (
	Alloc       MName = "Alloc"
	FreeMemory  MName = "FreeMemory"
	PollCount   MName = "PollCount"
	TotalMemory MName = "TotalMemory"
)

type Metric struct {
	Type        MType  `json:"Type"`
	Name        MName  `json:"Name"`
	Description string `json:"Description"`
	Value       int64  `json:"Value"`
}

func (m Metric) GetType() (MType, error) {
	switch m.Name {
	case Alloc:
	case FreeMemory:
	case TotalMemory:
		return Gauge, nil
	case PollCount:
		return Counter, nil
	}
	return "", NotExistsType
}

func GetMetricNames() []MName {
	return []MName{
		Alloc,
		FreeMemory,
		PollCount,
		TotalMemory,
	}
}

func (m Metric) Bytes() ([]byte, error) {
	return json.Marshal(m)
}

func NewMetric(data []byte) (Metric, error) {
	var m Metric
	err := json.Unmarshal(data, &m)
	return m, err
}

type List struct {
	Alloc       Metric `json:"Alloc"`
	FreeMemory  Metric `json:"FreeMemory"`
	PollCount   Metric `json:"PollCount"`
	TotalMemory Metric `json:"TotalMemory"`
}

func (l List) Array() []Metric {
	return []Metric{
		l.Alloc,
		l.FreeMemory,
		l.PollCount,
		l.TotalMemory,
	}
}
