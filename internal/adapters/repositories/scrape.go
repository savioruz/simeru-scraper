package repositories

import (
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/savioruz/simeru-scraper/internal/cores/entities"
	"github.com/savioruz/simeru-scraper/pkg/constant"
	"log"
	"strings"
	"time"
)

func (s *DB) ScrapeStudyPrograms(ctx context.Context, opts ...chromedp.ExecAllocatorOption) error {
	allocCtx, cancel := chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()

	chromeCtx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	var faculty []entities.Faculty
	var studyPrograms = make(map[string][]entities.StudyPrograms)
	var allStudyPrograms []entities.StudyPrograms

	err := chromedp.Run(chromeCtx, chromedp.Tasks{
		chromedp.Navigate("https://simeru.uad.ac.id/?mod=auth&sub=auth"),

		chromedp.SendKeys(`input[name="user"]`, "mhs"),
		chromedp.SendKeys(`input[name="pass"]`, "mhs"),
		chromedp.Click(`input[type="submit"]`),
		chromedp.Sleep(2 * time.Second),

		chromedp.Navigate("https://simeru.uad.ac.id/?mod=laporan_baru&sub=jadwal_prodi&do=daftar"),

		chromedp.Evaluate(`Array.from(document.querySelectorAll('select[name="fakultas"] option')).map(option => ({ value: option.value, name: option.text }))`, &faculty),

		chromedp.ActionFunc(func(ctx context.Context) error {
			for _, f := range faculty {
				if f.Value == "" {
					continue
				}

				// Set the faculty value
				err := chromedp.Run(ctx, chromedp.SetValue(`select[name="fakultas"]`, f.Value))
				if err != nil {
					log.Printf("Error setting Faculty value: %s", f.Value)
				}

				// Trigger JavaScript to update Prodi options
				err = chromedp.Run(ctx, chromedp.Evaluate(fmt.Sprintf(`popUpFak(document.form_cari.prodi, %q)`, f.Value), nil))
				if err != nil {
					log.Printf("Error triggering popUpFak for Faculty: %s", f.Value)
				}

				// Get the Prodi options for the selected Faculty
				var valueStudyPrograms []entities.StudyPrograms
				err = chromedp.Run(ctx, chromedp.Evaluate(`Array.from(document.querySelectorAll('select[name="prodi"] option')).filter(option => option.value !== "").map(option => ({ value: option.value, name: option.text }))`, &valueStudyPrograms))
				if err != nil {
					log.Printf("Error getting Prodi options for Faculty: %s", f.Value)
				}

				// Modify the `id` to include faculty name
				for i := range valueStudyPrograms {
					valueStudyPrograms[i].Faculty = strings.ToLower(f.Name)
				}

				// Store the study programs for this faculty
				studyPrograms[f.Value] = valueStudyPrograms

				// Append the study programs to the allStudyPrograms slice
				allStudyPrograms = append(allStudyPrograms, valueStudyPrograms...)
			}
			return nil
		}),
	})

	if err != nil {
		return err
	}

	// Store the scraped study programs into Redis
	for facultyID, prodiList := range studyPrograms {
		var facultyName string
		for _, f := range faculty {
			if f.Value == facultyID {
				facultyName = strings.ToLower(f.Name)
				break
			}
		}

		// Cache key for study programs by faculty
		key := fmt.Sprintf("studyPrograms:faculty:%s", strings.ReplaceAll(facultyName, " ", "_"))
		err = s.cache.Set(key, prodiList, constant.DefaultExpiration)
		if err != nil {
			log.Printf("Error setting cache for key: %s", key)
			continue
		}
		log.Printf("Study programs cached for Faculty: %s", facultyName)
	}
	err = s.cache.Set("studyPrograms:all", allStudyPrograms, constant.DefaultExpiration)
	if err != nil {
		log.Printf("Error setting cache for key: studyPrograms:all")
		return err
	}

	log.Printf("All study programs cached under key: studyPrograms:all")
	return nil
}

func (s *DB) ScrapeSchedule(ctx context.Context, opts ...chromedp.ExecAllocatorOption) error {
	allocCtx, cancel := chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()

	chromeCtx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	var faculty []entities.Faculty
	var studyPrograms = make(map[string][]entities.StudyPrograms)
	err := chromedp.Run(chromeCtx, chromedp.Tasks{
		chromedp.Navigate("https://simeru.uad.ac.id/?mod=auth&sub=auth"),

		chromedp.SendKeys(`input[name="user"]`, "mhs"),
		chromedp.SendKeys(`input[name="pass"]`, "mhs"),
		chromedp.Click(`input[type="submit"]`),
		chromedp.Sleep(2 * time.Second),

		chromedp.Navigate("https://simeru.uad.ac.id/?mod=laporan_baru&sub=jadwal_prodi&do=daftar"),

		chromedp.Evaluate(`Array.from(document.querySelectorAll('select[name="fakultas"] option')).map(option => ({ value: option.value, name: option.text }))`, &faculty),

		// Loop through each Faculty value and scrape StudyPrograms options
		chromedp.ActionFunc(func(ctx context.Context) error {
			for _, f := range faculty {
				if f.Value == "" {
					continue
				}

				// Set the Faculty value
				err := chromedp.Run(ctx, chromedp.SetValue(`select[name="fakultas"]`, f.Value))
				if err != nil {
					log.Printf("Error setting Faculty value: %s", f.Value)
				}

				// Trigger the popUpFak JavaScript function for the selected Faculty
				err = chromedp.Run(ctx, chromedp.Evaluate(fmt.Sprintf(`popUpFak(document.form_cari.prodi, %q)`, f.Value), nil))
				if err != nil {
					log.Printf("Error triggering popUpFak for Faculty: %s", f.Value)
				}

				// Get the Prodi options for the selected Faculty
				var valueStudyPrograms []entities.StudyPrograms
				err = chromedp.Run(ctx, chromedp.Evaluate(`Array.from(document.querySelectorAll('select[name="prodi"] option')).filter(option => option.value !== "").map(option => ({ value: option.value, name: option.text }))`, &valueStudyPrograms))
				if err != nil {
					log.Printf("Error getting Prodi options for Faculty: %s", f.Value)
				}

				// Store the Faculty and corresponding Prodi options in the map
				studyPrograms[f.Value] = valueStudyPrograms
			}
			return nil
		}),

		chromedp.ActionFunc(func(ctx context.Context) error {
			for facultyID, studyProgramLists := range studyPrograms {
				for _, prodi := range studyProgramLists {
					prodi.Name = strings.ToLower(prodi.Name)
					log.Printf("Scraping data for Faculty: %s, Study Program: %s-[%s]", facultyID, prodi.Name, prodi.Value)
					tableData, err := s.scrapeRowData(ctx, facultyID, prodi.Value)
					if err != nil {
						log.Printf("Error scraping table data for Faculty: %s, Study Program: %s", facultyID, prodi.Value)
						continue
					}

					rowsByHari := make(map[string][]entities.RowData)

					for i := range tableData {
						row := tableData[i]
						rowsByHari[row.Hari] = append(rowsByHari[row.Hari], row)
					}

					for day, rows := range rowsByHari {
						if day == "tidak ada data" {
							continue
						}

						key := fmt.Sprintf("schedule:studyPrograms:%s:day:%s", strings.ReplaceAll(prodi.Name, " ", "_"), strings.ToLower(day))
						err = s.cache.Set(key, rows, constant.DefaultExpiration)
						if err != nil {
							log.Printf("Error setting cache for key: %s", key)
						}
						log.Printf("Schedule cache set for key: %s", key)
					}
				}
			}
			return nil
		}),
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *DB) scrapeRowData(ctx context.Context, facultyID, studyProgramsID string) ([]entities.RowData, error) {
	var tableData []entities.RowData

	err := chromedp.Run(ctx,
		chromedp.Navigate("https://simeru.uad.ac.id/?mod=laporan_baru&sub=jadwal_prodi&do=daftar"),

		chromedp.SetValue(`select[name=fakultas]`, facultyID, chromedp.NodeVisible),
		chromedp.SetValue(`select[name=prodi]`, studyProgramsID, chromedp.NodeVisible),
		chromedp.Click(`input[name=submit]`, chromedp.NodeVisible),
		chromedp.WaitVisible(`table.table-list`),

		chromedp.Evaluate(`
        [...document.querySelectorAll("table.table-list tr")].slice(1).map(tr => {
            const cells = tr.children;
            return {
                Hari: cells[0] ? cells[0].textContent.trim() : "",
                Kode: cells[1] ? cells[1].textContent.trim() : "",
                Matkul: cells[2] ? cells[2].textContent.trim() : "",
                Kelas: cells[3] ? cells[3].textContent.trim() : "",
                Sks: cells[4] ? cells[4].textContent.trim() : "",
                Jam: cells[5] ? cells[5].textContent.trim() : "",
                Semester: cells[6] ? cells[6].textContent.trim() : "",
                Dosen: cells[7] ? cells[7].textContent.trim() : "",
                Ruang: cells[8] ? cells[8].textContent.trim() : ""
            }
        })`, &tableData))

	if err != nil {
		return nil, fmt.Errorf("error evaluating page: %w", err)
	}

	var lastDay string
	for i := range tableData {
		if tableData[i].Hari != "" {
			lastDay = tableData[i].Hari
		} else {
			tableData[i].Hari = lastDay
		}
	}

	return tableData, nil
}
