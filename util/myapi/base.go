package myapi

type (
	ApiInfo[TYPE any] struct {
		Method string `json:"method"`
		Path   string `json:"path"`
		Module TYPE   `json:"module"`
	}

	ApiManager[TYPE any] interface {
		SetGroup(group string) ApiManager[TYPE]
		AddGET(path string, module TYPE) ApiManager[TYPE]
		AddPOST(path string, module TYPE) ApiManager[TYPE]
		AddPUT(path string, module TYPE) ApiManager[TYPE]
		AddDELETE(path string, module TYPE) ApiManager[TYPE]
		ValueBy(method string, path string) (ApiInfo[TYPE], bool)
		List() []ApiInfo[TYPE]
	}
)
