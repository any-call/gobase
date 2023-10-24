package myxdb

import "fmt"

type XDBMaker struct {
	*Maker
}

func NewXDBMakerByVector(srcFile string, dstFile string) (ret *XDBMaker, err error) {
	maker, err := NewMaker(VectorIndexPolicy, srcFile, dstFile)
	if err != nil {
		return nil, err
	}

	return &XDBMaker{maker}, nil
}

func NewXDBMakerByBtree(srcFile string, dstFile string) (ret *XDBMaker, err error) {
	maker, err := NewMaker(BTreeIndexPolicy, srcFile, dstFile)
	if err != nil {
		return nil, err
	}

	return &XDBMaker{maker}, nil
}

func (self *XDBMaker) GenXDBFile() error {
	err := self.Init()
	if err != nil {
		return fmt.Errorf("failed Init: %s\n", err)
	}

	err = self.Start()
	if err != nil {
		return fmt.Errorf("failed Start: %s\n", err)
	}

	err = self.End()
	if err != nil {
		return fmt.Errorf("failed End: %s\n", err)
	}

	return nil
}
