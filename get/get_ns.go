package get

import (
	"fmt"
	"k8s_tools/lib"
	"time"

	v1 "k8s.io/api/core/v1"
)

func ShowNamespacesList(nsList *v1.NamespaceList, nsNameList []string) bool {
	fmt.Printf("***************** namespace list *****************\n")
	nsHeader := []string{"NAME", "STATUS", "AGE"}
	var data [][]string
	for _, ns := range nsList.Items {
		name := ns.GetName()
		if !lib.MatchResName(name, nsNameList) {
			continue
		}
		nsStatus := ns.Status
		status := fmt.Sprintf("%v", nsStatus.Phase)
		// Calculate the age
		createTime := ns.GetCreationTimestamp()
		age := time.Since(createTime.Time).Round(time.Second)
		ageS := lib.HoursConverter(age.String())
		data = append(data, []string{name, status, ageS})
	}
	if len(data) == 0 {
		fmt.Println("[ERROR] not found ns:", nsNameList)
		return false
	}
	return lib.FormatPrint(nsHeader, data)
}
