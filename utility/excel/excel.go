// Package excel 导入导出封装
//
// 基于 xuri/excelize/v2 实现，支持 xls 与 xlsx
package excel

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

// ============================== Excel 导入导出 ==============================

// FileOpener 可打开的上传文件接口
//
// ghttp.UploadFile 与 multipart.FileHeader 均满足该接口
type FileOpener interface {
	Open() (multipart.File, error)
}

// ReadSheet 读取上传文件的首个工作表
//
// 返回表头与数据行列表，首行视为表头，所有单元格值统一转为字符串；
// 空行会被跳过
func ReadSheet(file FileOpener) (headers []string, rows [][]string, err error) {
	f, err := file.Open()
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()
	xl, err := excelize.OpenReader(f)
	if err != nil {
		return nil, nil, fmt.Errorf("解析Excel失败: %v", err)
	}
	defer xl.Close()
	sheet := xl.GetSheetName(0)
	all, err := xl.GetRows(sheet)
	if err != nil {
		return nil, nil, err
	}
	if len(all) == 0 {
		return nil, nil, fmt.Errorf("Excel内容为空")
	}
	for _, h := range all[0] {
		headers = append(headers, strings.TrimSpace(h))
	}
	for _, row := range all[1:] {
		if isBlankRow(row) {
			continue
		}
		cells := make([]string, len(headers))
		for i := range headers {
			if i < len(row) {
				cells[i] = strings.TrimSpace(row[i])
			}
		}
		rows = append(rows, cells)
	}
	return headers, rows, nil
}

// WriteSheet 写出 xlsx 文件
//
// 文件写入指定目录，文件名由基础名与时间戳拼接而成；
// 返回文件名与完整路径
func WriteSheet(dir, baseName string, headers []string, rows [][]any) (fileName, filePath string, err error) {
	if err = os.MkdirAll(dir, 0o755); err != nil {
		return "", "", err
	}
	xl := excelize.NewFile()
	defer xl.Close()
	sheet := xl.GetSheetName(0)
	_ = xl.SetSheetName(sheet, "data")
	sheet = "data"
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		if err = xl.SetCellValue(sheet, cell, h); err != nil {
			return "", "", err
		}
	}
	for r, row := range rows {
		for c, v := range row {
			cell, _ := excelize.CoordinatesToCellName(c+1, r+2)
			if err = xl.SetCellValue(sheet, cell, v); err != nil {
				return "", "", err
			}
		}
	}
	fileName = fmt.Sprintf("%s_%s.xlsx", baseName, time.Now().Format("20060102150405"))
	filePath = filepath.Join(dir, fileName)
	if err = xl.SaveAs(filePath); err != nil {
		return "", "", err
	}
	return fileName, filePath, nil
}

// isBlankRow 判断是否为空行
//
// 所有单元格均为空时返回 true
func isBlankRow(row []string) bool {
	for _, c := range row {
		if strings.TrimSpace(c) != "" {
			return false
		}
	}
	return true
}
