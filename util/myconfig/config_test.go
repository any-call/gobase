package myconfig

import (
	"reflect"
	"testing"
)

type testAConfig struct {
	Host     string
	SiteName string
	SiteDesc string
}

type testBConfig struct {
	Host     string
	Enabled  bool
	Port     int64
	Ratio    float64
	Base     testAConfig
	ListA    []testAConfig
	Labels   map[string]string
	Empty    []string
	Optional *testAConfig
}

func TestEncodeStruct(t *testing.T) {
	cfg := testBConfig{
		Host:    "example.com",
		Enabled: true,
		Port:    443,
		Ratio:   1.25,
		Base: testAConfig{
			Host:     "base.example.com",
			SiteName: "Base",
		},
		ListA: []testAConfig{
			{Host: "a.example.com", SiteName: "A"},
			{Host: "b.example.com", SiteName: "B"},
		},
		Labels: map[string]string{"env": "prod"},
		Empty:  []string{},
	}

	items, err := Encode(cfg)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("list items:", items)

	got := itemMap(items)
	assertItem(t, got, "Host", "example.com", TypeString)
	assertItem(t, got, "Enabled", "true", TypeBool)
	assertItem(t, got, "Port", "443", TypeInt64)
	assertItem(t, got, "Ratio", "1.25", TypeFloat64)
	assertItem(t, got, "Base", `{"Host":"base.example.com","SiteName":"Base","SiteDesc":""}`, TypeJSON)
	assertItem(t, got, "ListA", `[{"Host":"a.example.com","SiteName":"A","SiteDesc":""},{"Host":"b.example.com","SiteName":"B","SiteDesc":""}]`, TypeJSON)
	assertItem(t, got, "Labels", `{"env":"prod"}`, TypeJSON)
	assertItem(t, got, "Empty", `[]`, TypeJSON)
	assertItem(t, got, "Optional", "", TypeNil)
}

func TestEncodeDecodeStruct(t *testing.T) {
	optional := &testAConfig{Host: "optional.example.com", SiteName: "Optional"}
	want := testBConfig{
		Host:     "example.com",
		Enabled:  true,
		Port:     443,
		Ratio:    1.25,
		Base:     testAConfig{Host: "base.example.com", SiteName: "Base"},
		ListA:    []testAConfig{{Host: "a.example.com", SiteName: "A"}},
		Labels:   map[string]string{"env": "prod"},
		Empty:    []string{},
		Optional: optional,
	}

	items, err := Encode(&want)
	if err != nil {
		t.Fatal(err)
	}

	var got testBConfig
	if err = Decode(items, &got); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("decode mismatch:\nwant: %#v\ngot:  %#v", want, got)
	}
}

func TestEncodeDecodeRootValue(t *testing.T) {
	items, err := Encode("example.com")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 {
		t.Fatalf("want 1 item, got %d", len(items))
	}
	assertItem(t, itemMap(items), rootKey, "example.com", TypeString)

	var got string
	if err = Decode(items, &got); err != nil {
		t.Fatal(err)
	}
	if got != "example.com" {
		t.Fatalf("want example.com, got %q", got)
	}
}

func TestDecodeKeepsMissingFieldAndIgnoresUnknownKey(t *testing.T) {
	cfg := testAConfig{SiteDesc: "default"}
	items := []Item{
		{Key: "Host", Value: "example.com", Type: TypeString},
		{Key: "Unknown", Value: "ignored", Type: TypeString},
	}

	if err := Decode(items, &cfg); err != nil {
		t.Fatal(err)
	}
	if cfg.Host != "example.com" {
		t.Fatalf("want example.com, got %q", cfg.Host)
	}
	if cfg.SiteDesc != "default" {
		t.Fatalf("want default value preserved, got %q", cfg.SiteDesc)
	}
}

func TestDecodeInvalidValue(t *testing.T) {
	var cfg struct {
		Enabled bool
	}
	err := Decode([]Item{{Key: "Enabled", Value: "invalid", Type: TypeBool}}, &cfg)
	if err == nil {
		t.Fatal("want parse error")
	}
}

func TestDecodeDuplicateKey(t *testing.T) {
	var cfg testAConfig
	err := Decode([]Item{
		{Key: "Host", Value: "a", Type: TypeString},
		{Key: "Host", Value: "b", Type: TypeString},
	}, &cfg)
	if err == nil {
		t.Fatal("want duplicate key error")
	}
}

func itemMap(items []Item) map[string]Item {
	ret := make(map[string]Item, len(items))
	for _, item := range items {
		ret[item.Key] = item
	}
	return ret
}

func assertItem(t *testing.T, items map[string]Item, key, value, valueType string) {
	t.Helper()

	item, ok := items[key]
	if !ok {
		t.Fatalf("item %q not found", key)
	}
	if item.Value != value || item.Type != valueType {
		t.Fatalf("item %q: want value=%q type=%q, got value=%q type=%q",
			key, value, valueType, item.Value, item.Type)
	}
}
