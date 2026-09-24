package adapter

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw3-logs/internal/domain"
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw3-logs/internal/domain/models"
)

func WriteMarkdown(w io.Writer, config domain.Config, dto models.DTO) {
	bw := bufio.NewWriter(w)

	files := make([]string, 0)
	for _, f := range dto.Files {
		files = append(files, fmt.Sprintf("%s ", f))
	}
	fileStr := strings.Join(files, ",")

	fmt.Fprintln(bw, "#### Общая информация")
	fmt.Fprintln(bw)
	fmt.Fprintln(bw, "|        Метрика        |                 Значение             |")
	fmt.Fprintln(bw, "|:---------------------:|-------------------------------------:|")
	fmt.Fprintf(bw, "|       Файл(-ы)        | %s |\n", fileStr)
	fmt.Fprintf(bw, "|    Начальная дата     |   %s |\n", config.From)
	fmt.Fprintf(bw, "|     Конечная дата     |   %s |\n", config.To)
	fmt.Fprintf(bw, "|  Количество запросов  |       %d |\n", dto.TotalRequestsCount)
	fmt.Fprintf(bw, "| Средний размер ответа |       %f |\n", dto.ResponseSizeInBytes.Average)
	fmt.Fprintf(bw, "|  95p размера ответа   |       %f |\n", dto.ResponseSizeInBytes.P95)
	fmt.Fprintln(bw)

	fmt.Fprintln(bw, "#### Запрашиваемые ресурсы")
	fmt.Fprintln(bw)
	fmt.Fprintln(bw, "|     Ресурс      | Количество |")
	fmt.Fprintln(bw, "|:---------------:|-----------:|")
	for _, r := range dto.Resources {
		fmt.Fprintf(bw, "|  `%s`  |      %d |\n", r.Resource, r.TotalRequestsCount)
	}
	fmt.Fprintln(bw)

	fmt.Fprintln(bw, "#### Коды ответа")
	fmt.Fprintln(bw)
	fmt.Fprintln(bw, "| Код |          Имя          | Количество |")
	fmt.Fprintln(bw, "|:---:|:---------------------:|-----------:|")
	for _, rc := range dto.ResponseCodesInfo {
		name := http.StatusText(rc.Code)
		fmt.Fprintf(bw, "| %d | %s |       %d |\n", rc.Code, name, rc.TotalResponsesCount)
	}
	bw.Flush()
}

func WriteJson(w io.Writer, dto models.DTO) error {
	wr := bufio.NewWriter(w)
	b, err := json.MarshalIndent(dto, "", "	")
	if err != nil {
		return err
	}

	_, err = wr.Write(b)
	if err != nil {
		return err
	}
	wr.Flush()
	return nil
}

func WriteAdoc(w io.Writer, config domain.Config, dto models.DTO) {
	bw := bufio.NewWriter(w)

	files := make([]string, 0)
	for _, f := range dto.Files {
		files = append(files, fmt.Sprintf("%s ", f))
	}
	filesCell := strings.Join(files, ", ")

	fmt.Fprintln(bw, "==== Общая информация")
	fmt.Fprintln(bw)
	fmt.Fprintln(bw, "[options=\"header\"]")
	fmt.Fprintln(bw, "|===")
	fmt.Fprintln(bw, "| Метрика | Значение")
	fmt.Fprintf(bw, "| Файл(-ы) | %s\n", filesCell)
	fmt.Fprintf(bw, "| Начальная дата | %s\n", config.From)
	fmt.Fprintf(bw, "| Конечная дата | %s\n", config.To)
	fmt.Fprintf(bw, "| Количество запросов | %d\n", dto.TotalRequestsCount)
	fmt.Fprintf(bw, "| Средний размер ответа | %f\n", dto.ResponseSizeInBytes.Average)
	fmt.Fprintf(bw, "| 95p размера ответа | %f\n", dto.ResponseSizeInBytes.P95)
	fmt.Fprintln(bw, "|===")
	fmt.Fprintln(bw)

	fmt.Fprintln(bw, "==== Запрашиваемые ресурсы")
	fmt.Fprintln(bw)
	fmt.Fprintln(bw, "[options=\"header\"]")
	fmt.Fprintln(bw, "|===")
	fmt.Fprintln(bw, "| Ресурс | Количество")
	for _, r := range dto.Resources {
		fmt.Fprintf(bw, "| %s | %d\n", r.Resource, r.TotalRequestsCount)
	}
	fmt.Fprintln(bw, "|===")
	fmt.Fprintln(bw)

	fmt.Fprintln(bw, "==== Коды ответа")
	fmt.Fprintln(bw)
	fmt.Fprintln(bw, "[options=\"header\"]")
	fmt.Fprintln(bw, "|===")
	fmt.Fprintln(bw, "| Код | Имя | Количество")
	for _, rc := range dto.ResponseCodesInfo {
		name := http.StatusText(rc.Code)
		fmt.Fprintf(bw, "| %d | %s | %d\n", rc.Code, name, rc.TotalResponsesCount)
	}
	fmt.Fprintln(bw, "|===")
	bw.Flush()

}
