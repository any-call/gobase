# gobase  基于go ,封装一些常的框架 及功能库， 会持续保持更新

## util  功能库
## frame 基础框架
### mysignal 事件框架
### mybing 数据驱动框架

### 部分已封装库介绍
### mylist : 用于封装 切片的常用功能
### mysignal :基于QT 的信号 与 槽 的思路来封装，对象间通信号，可以定义多个信号,对象间的通信则通过信号与槽建立联接来实现通信。
### myvalidator :用于常用数据格式的验证
#### 支持格式如下：
#### enum :针对数组，字符串
#### valid: 针对 结构体，map, 切片 的递归遍历标识
#### 数值类：
#### min;max;range
#### 字符串类：
#### min_len/min_length ;max_len/max_length;range_len/range_length
#### 数组类：
#### arr_minlen/arr_minlength ;arr_maxlen/arr_maxlength;arr_rangelen/arr_rangelength
#### map类：
#### map_minlen/map_minlength ;map_maxlen/map_maxlength;map_rangelen/map_rangelength



### validator Demo 

#### type MyReq struct {
#### ID    int      `json:"id" validate:"min(10,不正确的ID) max(100, 不正确的ID值)"`
#### Name  string   `json:"name" validate:"min_len(1,用户名不能为空)"`
#### Sex   string   `json:"sex" validate:"enum(男|女,错误的性别)"`
#### MyArr []MyDept `validate:"arr_minlen(1,入参数组不能为空) valid(T)"`
#### MyMap map[string]MyDept
#### }

#### type MyDept struct {
#### DeptID int `json:"dept_id" validate:"range(1,10,错误的部门ID)"`
#### Name   string   `validate:"rangelen(6,10,名称长度必须是6-10)"`
#### }

#### func Test_validator(t *testing.T) {
#### req := MyReq{
#### ID:   100,
#### Name: "this is test",
#### Sex:  "男",
#### MyArr: []MyDept{{
#### DeptID: 5,
#### Name:   "lu889i",
#### }, {
#### DeptID: 1,
#### Name:   "jinguihua",
#### }},
#### MyMap: map[string]MyDept{
#### "12": {
#### DeptID: 1001,
#### Name:   "good",
#### },
#### },
#### }

#### 	if err := Validate(req); err != nil {
#### 	t.Error("err:", err)
#### 		return
#### 	}

#### 	t.Log("validate: ok ")
#### }


### 数据驱动 mybind

#### type Dept struct {
#### Name  string
#### Total int
#### }

#### func (d Dept) DataChanged(val any) {
#### fmt.Println("data change to :", val)
#### }

#### func Test_bind(t *testing.T) {
#### type myStruct struct {
#### Name string
#### Sex  int
#### }

#### var monitorObj map[string]int = make(map[string]int, 0)

#### listener := &Dept{
####		Name:  "jin",
####		Total: 10,
####	}

####	if err := AddListener(listener, &monitorObj); err != nil {
####		t.Error(err)
####		return
####	}

####	go func() {
####		var i int
####		for {
####			i += 5
####			SetData(func() {
####				monitorObj["name"] = i
####			})

####			time.Sleep(time.Millisecond)
####		}
####	}()

####	time.Sleep(time.Second * 10)
####	t.Log("run ok ")
#### }