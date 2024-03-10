package lib

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func IsElementExistsInSlice(slice []string, elem string) bool {
	for _, e := range slice {
		if e == elem {
			return true
		}
	}
	return false
}

func ReadFile(file_path string) string {
	content, err := os.ReadFile(file_path)
	if err != nil {
		fmt.Printf("[ERROR] ReadFile [%v] err: %v\n", file_path, err)
		return ""
	}
	return string(content)
}

func UpdateFile(str string, file_path string) bool {
	file, err := os.OpenFile(file_path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Printf("[ERROR] OpenFile file [%v] err: %v\n", file_path, err)
		return false
	}
	if _, err := file.WriteString(str); err != nil {
		fmt.Printf("[ERROR] WriteString [%v] err: %v\n", str, err)
		return false
	}
	defer file.Close()
	return true
}

func FormatPrint(header []string, body [][]string) bool {
	if len(header) == 0 || len(body) == 0 {
		fmt.Printf("[ERROR] header is: %v, body is: %v\n", header, body)
		return false
	}
	if len(header) != len(body[0]) {
		fmt.Printf("[ERROR] header is: %v len:%v, body 0 is: %v len:%v\n", header, len(header), body[0], len(body[0]))
		return false
	}
	for _, i := range body[1:] {
		if len(i) != len(body[0]) {
			fmt.Printf("[ERROR] %v length is: %v, %v length is: %v\n", body[0], len(body[0]), i, len(i))
			return false
		}
	}

	columnLengthList := []int{}
	for _, h := range header {
		columnLengthList = append(columnLengthList, len(h))
	}
	for _, b := range body {
		for j := 0; j < len(b); j++ {
			if columnLengthList[j] < len(b[j]) {
				columnLengthList[j] = len(b[j])
			}
		}
	}

	for _, row := range append([][]string{header}, body...) {
		var rowStr string
		for m := 0; m < len(columnLengthList); m++ {
			colLength := columnLengthList[m] + 2
			rowStr += row[m] + strings.Repeat(" ", colLength-len(row[m]))
		}
		fmt.Println(rowStr)
	}
	return true
}

func ReduceDuplicateSpace(str string) string {
	newStr := regexp.MustCompile(`\s+`).ReplaceAllString(str, " ")
	return strings.TrimSpace(newStr)
}

func CheckPodName(podName string, podNameList []string) bool {
	if len(podNameList) == 0 {
		return true
	} else {
		for _, pn := range podNameList {
			if strings.Contains(podName, pn) {
				return true
			}
		}
	}
	return false
}

func MatchResName(resName string, resNameList []string) bool {
	if len(resNameList) == 0 {
		return true
	} else {
		for _, name := range resNameList {
			if strings.Contains(resName, name) {
				return true
			}
		}
	}
	return false
}

func AppendResName(listNeedToBeExtend []string, resName string, resNameList []string) []string {
	if MatchResName(resName, resNameList) {
		listNeedToBeExtend = append(listNeedToBeExtend, resName)
	}
	return listNeedToBeExtend
}

func SortMapKeys(inputMap any) []string {
	var keys []string
	switch m := inputMap.(type) {
	case map[string]string:
		for k := range m {
			keys = append(keys, k)
		}
	case map[string][]string:
		for k := range m {
			keys = append(keys, k)
		}
	case map[string][][]string:
		for k := range m {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	return keys
}

func SortDataListWithIndex0(data [][]string) {
	sort.Slice(data, func(i, j int) bool { return data[i][0] < data[j][0] })
}

func HoursConverter(timeStr string) string {
	if !strings.Contains(timeStr, "h") {
		return timeStr
	}
	var newTimeStr string
	totalHourStr := strings.Split(timeStr, "h")[0]
	totalHour, err := strconv.Atoi(totalHourStr)
	if err != nil {
		fmt.Println("[ERROR] HoursConverter err:", err)
		return timeStr
	}
	switch {
	case totalHour > 8760:
		year := totalHour / 8760
		day := (totalHour % 8760) / 24
		newTimeStr = fmt.Sprintf("%vy%vd", year, day)
	case totalHour > 720:
		month := totalHour / 720
		day := (totalHour % 720) / 24
		newTimeStr = fmt.Sprintf("%vm%vd", month, day)
	case totalHour > 48:
		day := totalHour / 24
		hour := totalHour % 24
		newTimeStr = fmt.Sprintf("%vd%vh", day, hour)
	default:
		newTimeStr = timeStr
	}
	return newTimeStr
}
