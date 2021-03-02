//Package filter 过滤用户不需要的数据， filter结构体存放用户需要的数据
package filter

import (
	"encoding/json"
	"sync"
)

var filter Filter

const metricFlied = "metric"

func init() {
	filter.Map = make(map[string][]string)
}

//SetFilter 初始化过滤数据
func SetFilter(bytes []byte) error {

	err := json.Unmarshal(bytes, &filter.Map)
	if err != nil {
		return err
	}

	//初始化metric
	filter.Map[metricFlied] = []string{}
	for key, val := range filter.Map {
		if key != metricFlied {
			filter.Map[metricFlied] = append(filter.Map[metricFlied], val...)
		}
	}
	return nil
}

//GetFilter 返回filter数据
func GetFilter() *Filter {
	return &filter
}

//Filter 用于过滤用户是否需要该metric的收集
type Filter struct {
	sync.Mutex
	Map map[string][]string
}

//Type 是否需要该类型， 需要返回true
func (f *Filter) Type(typ string) bool {
	_, ok := f.Map[typ]
	return ok
}

//SubMetric 是否需要子类型， 需要返回true
func (f *Filter) SubMetric(subMetric string) bool {
	return find(f.Map[metricFlied], subMetric)
}

//in 查看数组中是否有该数据
func find(arr []string, val string) bool {
	for _, item := range arr {
		if item == val {
			return true
		}
	}
	return false
}
