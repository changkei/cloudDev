package lru

import "testing"

func TestReadWrite(t *testing.T) {
	testlru := NewCache(3)
	testkey := "key"
	testval := "val"
	if err := testlru.Put(testkey, testval); err != nil {
		t.Error("Write test failed", err)
	}
	if val, err := testlru.Get(testkey); err != nil {
		t.Error("Read test failed", err)
	} else if val.(string) != "val" {
		t.Error("Read failed with incorrect return value", val)
	}
}

func TestWriteWithEviction(t *testing.T) {
	testlru := NewCache(3)
	inputs := make(map[string]string)
	inputs["key1"] = "val1"
	inputs["key2"] = "val2"
	inputs["key3"] = "val3"
	inputs["key4"] = "val4"
	for k, v := range inputs {
		testlru.Put(k, v)
	}
	if _, err := testlru.Get("key1"); err == nil {
		t.Error("LRU replacement policy test failed")
	}
	// testlru.Put("key1", "val1")
	// testlru.Put("key2", "val2")
	// testlru.Put("key3", "val3")
	// testlru.Put("key4", "val4")
	// testlru.Put("key3", "val3")
	// testlru.Put("key1", "val1")
}
