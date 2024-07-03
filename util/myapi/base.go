package myapi

type (
	ApiInfo[TYPE any] struct {
		Method string
		Path   string
		Module string
		Type   TYPE
	}

	ApiManager[TYPE any] interface {
		SetGroup(group string, op TYPE) ApiManager[TYPE]
		AddGET(path string, module string, op TYPE) ApiManager[TYPE]
		AddPOST(path string, module string, op TYPE) ApiManager[TYPE]
		AddPUT(path string, module string, op TYPE) ApiManager[TYPE]
		AddDELETE(path string, module string, op TYPE) ApiManager[TYPE]
		ValueBy(method string, path string) (ApiInfo[TYPE], bool)
		List() []ApiInfo[TYPE]
	}
)
