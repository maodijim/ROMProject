package main

import (
	Cmd "ROMProject/Cmds"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
	"time"
)

const (
	defaultSheetName = "Sheet1"
	defaultFileName  = "tradeRecords.xlsx"
)

var (
	cellTitle = map[string]string{
		"日期":  "A1",
		"道具":  "B1",
		"数量":  "C1",
		"价格":  "D1",
		"单价":  "E1",
		"买/卖": "F1",
		"购买人": "G1",
	}
)

type TradeExcel struct {
	SheetName string
	ExcelPath string
	file      *excelize.File
}

func (t *TradeExcel) CreateExcel() {
	f := excelize.NewFile()
	for n, c := range cellTitle {
		_ = f.SetCellValue(t.SheetName, c, n)
	}
	f.SetActiveSheet(f.GetSheetIndex(t.SheetName))
	t.file = f
}

func (t *TradeExcel) WriteExcel() {
	log.Infof("saving %s", t.ExcelPath)
	if err := t.file.SaveAs(t.ExcelPath); err != nil {
		log.Errorf("failed to save excel file %s: %s", t.ExcelPath, err)
	}
}

func (t *TradeExcel) ReadExcel() {
	f, err := excelize.OpenFile(t.ExcelPath)
	if err != nil {
		log.Errorf("Failed to open %s", t.ExcelPath)
	} else {
		f.SetActiveSheet(f.GetSheetIndex(t.SheetName))
		t.file = f
	}
}

func (t *TradeExcel) AddRecord(info *Cmd.LogItemInfo, itemName string) {
	if t.file == nil {
		log.Warn("excel file is nil")
	}
	rows, err := t.file.GetRows(t.SheetName)
	if err != nil {
		log.Errorf("failed to get row in sheet %s", t.SheetName)
		return
	}
	maxRow := strconv.Itoa(len(rows) + 1)
	todayDate := time.Now().Format("01/02")
	_ = t.file.SetCellValue(t.SheetName, "A"+maxRow, todayDate)
	_ = t.file.SetCellValue(t.SheetName, "B"+maxRow, itemName)
	_ = t.file.SetCellValue(t.SheetName, "C"+maxRow, *info.Count)
	_ = t.file.SetCellValue(t.SheetName, "D"+maxRow, *info.Price*uint64(*info.Count))
	_ = t.file.SetCellValue(t.SheetName, "E"+maxRow, *info.Price)
	if info.GetLogtype() == Cmd.EOperType_EOperType_NormalSell ||
		info.GetLogtype() == Cmd.EOperType_EOperType_PublicitySellSuccess {
		_ = t.file.SetCellValue(t.SheetName, "F"+maxRow, "卖")
	} else {
		_ = t.file.SetCellValue(t.SheetName, "F"+maxRow, "买")
	}
	_ = t.file.SetCellValue(t.SheetName, "G"+maxRow, *info.NameInfo.Name)

}

func NewTradeExcel(excelPath, sheetName string) *TradeExcel {
	if excelPath == "" {
		excelPath = defaultFileName
	}
	if sheetName == "" {
		sheetName = defaultSheetName
	}
	tradeExcel := &TradeExcel{
		SheetName: sheetName,
		ExcelPath: excelPath,
	}
	if _, err := os.Stat(excelPath); os.IsNotExist(err) {
		tradeExcel.CreateExcel()
	} else {
		tradeExcel.ReadExcel()
	}
	return tradeExcel
}
