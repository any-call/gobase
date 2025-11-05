package mylist

func Prepend[E any](slice []E, element E) []E {
	if slice == nil {
		return []E{element}
	}
	return append([]E{element}, slice...)
}

func Range[T int | int8 | int16 | int32 | int64 | float32 | float64](start, end T) []T {
	if end < start {
		return []T{}
	}
	n := int(end - start + 1)
	out := make([]T, n)
	for i := 0; i < n; i++ {
		out[i] = start + T(i)
	}
	return out
}

// 求并集
func Union[E comparable](slice1, slice2 []E) []E {
	m := make(map[E]int)
	for _, v := range slice1 {
		m[v]++
	}

	for _, v := range slice2 {
		times, _ := m[v]
		if times == 0 {
			slice1 = append(slice1, v)
		}
	}
	return slice1
}

// 求交集
func Intersect[E comparable](slice1, slice2 []E) []E {
	m := make(map[E]int)
	nn := make([]E, 0)
	for _, v := range slice1 {
		m[v]++
	}

	for _, v := range slice2 {
		times, _ := m[v]
		if times == 1 {
			nn = append(nn, v)
		}
	}
	return nn
}

// 求差集 slice1-并集 :取Slice的 差集部分
func Difference[E comparable](slice1, slice2 []E) []E {
	m := make(map[E]int)
	nn := make([]E, 0)
	inter := Intersect[E](slice1, slice2)
	for _, v := range inter {
		m[v]++
	}

	for _, value := range slice1 {
		times, _ := m[value]
		if times == 0 {
			nn = append(nn, value)
		}
	}
	return nn
}

// 删除掉数组中指定的元素
func DelItem[E comparable](list []E, delE E) []E {
	for i := 0; i < len(list); i++ {
		if list[i] == delE {
			return append(list[:i], list[i+1:]...)
		}
	}

	return list
}

func DelItemAtIndex[E comparable](list []E, index int) []E {
	if index < 0 || index >= len(list) {
		return list // 下标无效，返回原切片
	}
	return append(list[:index], list[index+1:]...)
}

// 检测数组中是否存在该元素
func IsExistItem[E comparable](list []E, ele E) bool {
	for i := 0; i < len(list); i++ {
		if list[i] == ele {
			return true
		}
	}

	return false
}

// 数组去重
func RemoveDuplicItem[E comparable](list []E) []E {
	set := make(map[E]int, len(list))
	ret := make([]E, 0)
	for _, v := range list {
		_, ok := set[v]
		if ok {
			continue
		}
		set[v] = 1
		ret = append(ret, v)

	}

	return ret
}

func RemoveDupKey[T any, K comparable](items []T, keyFn func(item T) K) []T {
	if keyFn == nil {
		return items
	}

	result := make([]T, 0, len(items))
	seen := make(map[K]struct{})

	for _, item := range items {
		k := keyFn(item)
		if _, ok := seen[k]; !ok {
			seen[k] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

// a list to b list
func ToList[F, T any](org *List[F], f func(o F) T) *List[T] {
	ret := NewList[T]()
	if org != nil && f != nil {
		org.Range(func(index int, v F) {
			ret.Append(f(v))
		})
	}

	return ret
}
