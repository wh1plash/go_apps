package main

import (
	"archive/zip"
	"encoding/csv"
	"fmt"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/xuri/excelize/v2"
)

func main() {
	// Создаем папку temp
	tempDir := "temp"
	os.MkdirAll(tempDir, os.ModePerm)

	// Разархивируем файлы из zip
	err := unzip("archive.zip", tempDir)
	if err != nil {
		fmt.Println("Ошибка разархивирования:", err)
		return
	}

	// Изменяем формат csv на xlsx
	err = convertCSVToXLSX(tempDir)
	if err != nil {
		fmt.Println("Ошибка конвертации:", err)
		return
	}

	// Заархивируем в файл с текущей датой
	err = zipFiles("archive_"+time.Now().Format("20060102")+".zip", tempDir)
	if err != nil {
		fmt.Println("Ошибка архивирования:", err)
		return
	}

	fmt.Println("Процесс завершен успешно!")
}

func unzip(zipFile string, dest string) error {
	r, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		dstFile, err := os.OpenFile(fpath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		defer dstFile.Close()

		srcFile, err := f.Open()
		if err != nil {
			return err
		}
		defer srcFile.Close()

		if _, err := io.Copy(dstFile, srcFile); err != nil {
			return err
		}
	}

	return nil
}

func convertCSVToXLSX(dir string) error {
	files, err := filepath.Glob(filepath.Join(dir, "*.csv"))
	if err != nil {
		return err
	}

	totalRows := 0 // Переменная для подсчета общего количества строк

	for _, file := range files {
		xlsxFile := file[:len(file)-len(filepath.Ext(file))] + ".xlsx"
		f := excelize.NewFile()

		csvFile, err := os.Open(file)
		if err != nil {
			return err
		}
		defer csvFile.Close()

		// Декодируем файл из Windows-1251 в UTF-8
		reader := transform.NewReader(csvFile, charmap.Windows1251.NewDecoder())

		// Создаем новый CSV Reader с разделителем ';'
		csvReader := csv.NewReader(reader)
		csvReader.Comma = ';' // Устанавливаем разделитель на ';'

		// Читаем данные из CSV файла
		records, err := csvReader.ReadAll()
		if err != nil {
			return err
		}

		totalRows += len(records) // Увеличиваем счетчик на количество строк

		for i, record := range records {
			for j, value := range record {
				cell, _ := excelize.CoordinatesToCellName(j+1, i+1)
				f.SetCellStr("Sheet1", cell, value)
			}
		}

		if err := f.SaveAs(xlsxFile); err != nil {
			return err
		}

		// Удаляем оригинальный csv файл
		os.Remove(file)
	}

	return nil // Возвращаем общее количество строк
}

func zipFiles(zipFile string, dir string) error {
	newZipFile, err := os.Create(zipFile)
	if err != nil {
		return err
	}
	defer newZipFile.Close()

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	return filepath.Walk(dir, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if fi.IsDir() {
			return nil
		}

		f, err := os.Open(file)
		if err != nil {
			return err
		}
		defer f.Close()

		w, err := zipWriter.Create(filepath.Base(file))
		if err != nil {
			return err
		}
		_, err = io.Copy(w, f)
		return err
	})
}
