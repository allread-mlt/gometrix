package gometrix

import (
	"fmt"
)

func joinTagsToName(name string, tagMap []MetricTag) string {
	if len(tagMap) == 0 {
		return name
	}

	for _, v := range tagMap {
		name = fmt.Sprintf("%v_%v_%v", name, v.name, v.value)
	}
	return name
}
