package excel_template_engines

import (
	"testing"
)

func TestGetProductName(t *testing.T) {
	if str := GetProductName("进口全脂奶粉(≥26%)"); str != "全脂奶粉" {
		t.Errorf("GetSpeci error, but %s got", str)
	}
	if str1 := GetProductName("进口全脂乳粉(≥26%)"); str1 != "全脂奶粉" {
		t.Errorf("GetSpeci error, but %s got", str1)
	}
	if str2 := GetProductName("塑料粒子XX"); str2 != "" {
		t.Errorf("GetSpeci error, but %s got", str2)
	}
}
