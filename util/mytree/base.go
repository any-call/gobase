package mytree

import (
	"strconv"
	"strings"
)

func FindNodeByID[T any](listData []T, uid string, getChildrenFn func(T) []T) (T, bool) {
	listIDs, err := ParseTreeNodePath(uid)
	if err != nil {
		var zeroValue T
		return zeroValue, false
	}

	return FindNodeByPath(listData, listIDs, getChildrenFn)
}

func FindNodeByPath[T any](listData []T, listIDs []int, getChildrenFn func(T) []T) (T, bool) {
	var zeroValue T
	if listIDs == nil || len(listIDs) == 0 {
		return zeroValue, false
	}
	if len(listData) == 0 {
		return zeroValue, false
	}

	for i := range listIDs {
		if listData != nil && listIDs[i] < len(listData) {
			if i+1 == len(listIDs) { // 如果是最后一个节点
				return listData[listIDs[i]], true
			} else {
				if getChildrenFn == nil {
					return zeroValue, false
				}

				listData = getChildrenFn(listData[listIDs[i]])
			}
		}
	}
	return zeroValue, false
}

func ParseTreeNodePath(uid string) ([]int, error) {
	parts := strings.Split(uid, "-")
	uids := make([]int, len(parts))
	var err error
	for i, part := range parts {
		uids[i], err = strconv.Atoi(part)
		if err != nil {
			return nil, err
		}
	}
	return uids, nil
}
