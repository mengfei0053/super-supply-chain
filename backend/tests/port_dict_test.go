package tests

import (
	"testing"

	"super-supply-chain/models"
	excel_template_engines "super-supply-chain/utils/excel-template-engines"
)

func TestGetPortInfoFromDict(t *testing.T) {
	setupTestDB(t, &models.BaseDict{})
	excel_template_engines.InvalidatePortDictCache()
	models.DB.Create(&models.BaseDict{
		Key: "上海洋山保税", Value: "上海", Type: "港口字典",
		PortName: "上海口岸", ExtraPay: "洋山补差",
	})
	models.DB.Create(&models.BaseDict{
		Key: "广州黄埔", Value: "广州", Type: "港口字典",
		PortName: "广州口岸", ExtraPay: "补差",
	})
	excel_template_engines.InvalidatePortDictCache()

	info := excel_template_engines.GetPortInfo("上海洋山保税")
	if info.Addr != "上海" || info.PortName != "上海口岸" || info.ExtraPay != "洋山补差" {
		t.Fatalf("unexpected port info: %+v", info)
	}
	info2 := excel_template_engines.GetPortInfo("广州黄埔")
	if info2.Addr != "广州" || info2.PortName != "广州口岸" {
		t.Fatalf("unexpected port info: %+v", info2)
	}
	if excel_template_engines.GetArrivalPort("上海洋山保税") != "上海" {
		t.Fatalf("GetArrivalPort mismatch")
	}
}
