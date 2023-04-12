package mycond

func IfV[T any](b bool, trueV, falseV T) T {
	if b {
		return trueV
	}
	return falseV
}

func IfFun(conf bool, trueFun, falseFun func()) {
	if conf {
		if trueFun != nil {
			trueFun()
		}
	} else {
		if falseFun != nil {
			falseFun()
		}
	}
}
