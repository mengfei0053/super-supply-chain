package batch_print

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"os/exec"
)

func main() {
	excelFile := "/Users/menghongfei/Desktop/yifan/2025-06-03/结算单/25SYF151 结算单.xlsx"

	// 1. 用 excelize 打开 Excel
	f, err := excelize.OpenFile(excelFile)
	if err != nil {
		fmt.Printf("无法打开 Excel 文件: %v\n", err)
		return
	}
	defer f.Close()

	// 2. 可以在这里修改 Excel 内容（可选）
	// f.SetCellValue("Sheet1", "A1", "Hello, Printer!")

	// 3. 保存修改（可选）
	// if err := f.Save(); err != nil {
	// 	fmt.Printf("保存失败: %v\n", err)
	// 	return
	// }

	// 4. 调用系统打印命令
	cmd := exec.Command("lp", excelFile) // macOS/Linux
	// cmd := exec.Command("start", "/min", "excel", "/p", excelFile) // Windows

	err = cmd.Run()
	if err != nil {
		fmt.Printf("打印失败: %v\n", err)
		return
	}
	fmt.Println("打印任务已发送！")
}
