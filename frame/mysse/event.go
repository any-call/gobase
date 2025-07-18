package mysse

type Event struct {
	ID    string
	Event string
	Data  string
}

func (e *Event) Encode() []byte { //sse 数据封装格式
	var buf []byte
	if e.ID != "" {
		buf = append(buf, []byte("id: "+e.ID+"\n")...)
	}
	if e.Event != "" {
		buf = append(buf, []byte("event: "+e.Event+"\n")...)
	}
	for _, line := range splitLines(e.Data) {
		buf = append(buf, []byte("data: "+line+"\n")...)
	}
	buf = append(buf, []byte("\n")...)
	return buf
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}
