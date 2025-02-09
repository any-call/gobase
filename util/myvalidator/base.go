package myvalidator

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"regexp"
	"strings"
)

// email verify
func ValidEmail(email string) bool {
	emailRegexp := regexp.MustCompile(`(?m)^(((((((((\s? +)?(\(((\s? +)?(([!-'*-[\]-~]*)|(\\([ -~]|\s))))*(\s? +)?\)))(\s? +)?)|(\s? +))?([A-Za-z0-9!#-'*+\/=?^_\x60{|}~-])+((((\s? +)?(\(((\s? +)?(([!-'*-[\]-~]*)|(\\([ -~]|\s))))*(\s? +)?\)))(\s? +)?)|(\s? +))?)|(((((\s? +)?(\(((\s? +)?(([!-'*-[\]-~]*)|(\\([ -~]|\s))))*(\s? +)?\)))(\s? +)?)|(\s? +))?"((\s? +)?(([!#-[\]-~])|(\\([ -~]|\s))))*(\s? +)?"))?)?(((((\s? +)?(\(((\s? +)?(([!-'*-[\]-~]*)|(\\([ -~]|\s))))*(\s? +)?\)))(\s? +)?)|(\s? +))?<(((((((\s? +)?(\(((\s? +)?(([!-'*-[\]-~]*)|(\\([ -~]|\s))))*(\s? +)?\)))(\s? +)?)|(\s? +))?(([A-Za-z0-9!#-'*+\/=?^_\x60{|}~-])+(\.([A-Za-z0-9!#-'*+\/=?^_\x60{|}~-])+)*)((((\s? +)?(\(((\s? +)?(([!-'*-[\]-~]*)|(\\([ -~]|\s))))*(\s? +)?\)))(\s? +)?)|(\s? +))?)|(((((\s? +)?(\(((\s? +)?(([!-'*-[\]-~]*)|(\\([ -~]|\s))))*(\s? +)?\)))(\s? +)?)|(\s? +))?"((\s? +)?(([!#-[\]-~])|(\\([ -~]|\s))))*(\s? +)?"))@((((((\s? +)?(\(((\s? +)?(([!-'*-[\]-~]*)|(\\([ -~]|\s))))*(\s? +)?\)))(\s? +)?)|(\s? +))?(([A-Za-z0-9!#-'*+\/=?^_\x60{|}~-])+(\.([A-Za-z0-9!#-'*+\/=?^_\x60{|}~-])+)*)((((\s? +)?(\(((\s? +)?(([!-'*-[\]-~]*)|(\\([ -~]|\s))))*(\s? +)?\)))(\s? +)?)|(\s? +))?)|(((((\s? +)?(\(((\s? +)?(([!-'*-[\]-~]*)|(\\([ -~]|\s))))*(\s? +)?\)))(\s? +)?)|(\s? +))?\[((\s? +)?([!-Z^-~]))*(\s? +)?\]((((\s? +)?(\(((\s? +)?(([!-'*-[\]-~]*)|(\\([ -~]|\s))))*(\s? +)?\)))(\s? +)?)|(\s? +))?)))>((((\s? +)?(\(((\s? +)?(([!-'*-[\]-~]*)|(\\([ -~]|\s))))*(\s? +)?\)))(\s? +)?)|(\s? +))?))|(((((((\s? +)?(\(((\s? +)?(([!-'*-[\]-~]*)|(\\([ -~]|\s))))*(\s? +)?\)))(\s? +)?)|(\s? +))?(([A-Za-z0-9!#-'*+\/=?^_\x60{|}~-])+(\.([A-Za-z0-9!#-'*+\/=?^_\x60{|}~-])+)*)((((\s? +)?(\(((\s? +)?(([!-'*-[\]-~]*)|(\\([ -~]|\s))))*(\s? +)?\)))(\s? +)?)|(\s? +))?)|(((((\s? +)?(\(((\s? +)?(([!-'*-[\]-~]*)|(\\([ -~]|\s))))*(\s? +)?\)))(\s? +)?)|(\s? +))?"((\s? +)?(([!#-[\]-~])|(\\([ -~]|\s))))*(\s? +)?"))@((((((\s? +)?(\(((\s? +)?(([!-'*-[\]-~]*)|(\\([ -~]|\s))))*(\s? +)?\)))(\s? +)?)|(\s? +))?(([A-Za-z0-9!#-'*+\/=?^_\x60{|}~-])+(\.([A-Za-z0-9!#-'*+\/=?^_\x60{|}~-])+)*)((((\s? +)?(\(((\s? +)?(([!-'*-[\]-~]*)|(\\([ -~]|\s))))*(\s? +)?\)))(\s? +)?)|(\s? +))?)|(((((\s? +)?(\(((\s? +)?(([!-'*-[\]-~]*)|(\\([ -~]|\s))))*(\s? +)?\)))(\s? +)?)|(\s? +))?\[((\s? +)?([!-Z^-~]))*(\s? +)?\]((((\s? +)?(\(((\s? +)?(([!-'*-[\]-~]*)|(\\([ -~]|\s))))*(\s? +)?\)))(\s? +)?)|(\s? +))?))))$`)
	if !emailRegexp.MatchString(strings.ToLower(email)) {
		return false
	}
	return true
}

// mobile verify
func ValidPhone(mobileNum string) bool {
	regular := "^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199|(147))\\d{8}$"

	reg := regexp.MustCompile(regular)
	return reg.MatchString(mobileNum)
}

func IsFixedLengthDigits(s string, len int) bool {
	re := regexp.MustCompile(fmt.Sprintf(`^\d{%d}$`, len)) // ^\d{6}$ 匹配 字符串是指定长度的的数据
	return re.MatchString(s)
}

// 基于字符串名称调用方法，并传递参数
func CallMethod(obj any, methodName string, args ...any) ([]any, error) {
	// 获取对象的Value
	value := reflect.ValueOf(obj)

	// 查找方法
	method := value.MethodByName(methodName)
	if !method.IsValid() {
		return nil, fmt.Errorf("方法 %s 不存在", methodName)
	}

	// 准备参数
	methodType := method.Type()
	if len(args) != methodType.NumIn() {
		return nil, fmt.Errorf("方法 %s 需要 %d 个参数, 但提供了 %d 个", methodName, methodType.NumIn(), len(args))
	}

	in := make([]reflect.Value, len(args))
	for i, arg := range args {
		// 检查参数类型是否匹配
		if reflect.TypeOf(arg) != methodType.In(i) {
			return nil, fmt.Errorf("参数 %d 类型不匹配：预期 %s, 得到 %s", i, methodType.In(i), reflect.TypeOf(arg))
		}
		in[i] = reflect.ValueOf(arg)
	}

	// 调用方法
	resultValues := method.Call(in)

	// 将结果转换为[]any类型
	results := make([]any, len(resultValues))
	for i, result := range resultValues {
		results[i] = result.Interface()
	}

	return results, nil
}

func ExportedMethodNames(obj any, isUpper bool) []string {
	// 获取对象的类型
	objType := reflect.TypeOf(obj)

	// 列表存储方法名
	var methodNames []string

	// 遍历对象的所有方法
	for i := 0; i < objType.NumMethod(); i++ {
		method := objType.Method(i)
		// 只添加导出的方法
		if isUpper {
			if strings.ToUpper(method.Name[:1]) == method.Name[:1] {
				methodNames = append(methodNames, method.Name)
			}
		} else {
			methodNames = append(methodNames, method.Name)
		}

	}

	return methodNames
}

// 提取特定结构体的导出方法及其注释
func ExtractMethodComments(filePath, structName string) (map[string]string, error) {
	// 创建文件集
	fset := token.NewFileSet()

	// 解析文件
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	// 存储方法名及其注释
	methods := make(map[string]string)

	// 遍历 AST
	ast.Inspect(node, func(n ast.Node) bool {
		// 只处理函数声明
		if fn, ok := n.(*ast.FuncDecl); ok {
			// 检查接收者是否为目标结构体
			if fn.Recv != nil {
				recvType := fmt.Sprintf("%s", fn.Recv.List[0].Type)
				recvType = strings.Trim(recvType, "*") // 移除指针符号
				if recvType == structName && ast.IsExported(fn.Name.Name) {
					// 提取注释
					comment := ""
					if fn.Doc != nil {
						comment = strings.TrimSpace(fn.Doc.Text())
					}
					methods[fn.Name.Name] = comment
				}
			}
		}
		return true
	})

	return methods, nil
}
