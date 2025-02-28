package excel_template_engines

import (
	"super-supply-chain/models"
	"super-supply-chain/utils"
)

func CreateCostCalculation(data *models.ExcelData, fileName string) (string, error) {

	utils.LogJson(data)

	return "CostCalculation", nil
}
