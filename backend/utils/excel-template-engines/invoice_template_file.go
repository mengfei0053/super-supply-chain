package excel_template_engines

import "github.com/xuri/excelize/v2"

type workbookWriter func(*excelize.File) error

func generateFileFromTemplate(templatePath string, outputPath string, writeWorkbook workbookWriter) error {
	workbook, err := excelize.OpenFile(templatePath)
	if err != nil {
		return err
	}
	defer workbook.Close()

	if err := writeWorkbook(workbook); err != nil {
		return err
	}

	return workbook.SaveAs(outputPath)
}

func generateFileFromInvoiceTemplate(outputPath string, writeWorkbook workbookWriter) error {
	return generateFileFromTemplate(GetInvoiceTmpPath(), outputPath, writeWorkbook)
}

func generateFileFromChangjiuTemplate(outputPath string, writeWorkbook workbookWriter) error {
	return generateFileFromTemplate(GetChangejiuTmpPath(), outputPath, writeWorkbook)
}
