package functions

import (
	"fmt"
	"time"

	"github.com/gonext-tech/internal/queries"
)

func GetLink(name, table string, queries queries.InvoiceQueryParams) string {
	order := "desc"
	if queries.SortBy == name && queries.OrderBy == "desc" {
		order = "asc"
	}
	return fmt.Sprintf("/%s?sortBy=%s&orderBy=%s", table, name, order)
}

func Ternary(condition bool, a, b string) string {
	if condition {
		return a
	}
	return b
}

func GetStatusClass(status string) string {
	var class = "inline text-sm text-white text-center  font-normal capitalize px-2 py-1 rounded-lg "
	if status == "ACTIVE" || status == "UP" || status == "PAYIN" {
		class += "bg-blue-600"
	} else if status == "P_PAID" {
		class += "bg-purple-600"
	} else if status == "PAID" {
		class += "bg-green-600"
	} else if status == "NOT_ACTIVE" || status == "DOWN" || status == "PAYOUT" {
		class += "bg-red-600"
	} else if status == "TOPAY" || status == "PENDING" {
		class += "bg-yellow-500"
	}
	return class
}

func GetToday() string {
	now := time.Now()
	return now.Format("2006-01-02")
}

func GetDateValue(date *time.Time) string {
	if date != nil {
		return date.Format("2006-01-02")
	}
	return ""
}

func PrintDate(date *time.Time) string {
	if date == nil {
		return "-"
	}
	return date.Format("Jan, 02 2006")
}

func PrintName(name string) string {
	if name == "" {
		return "-"
	}
	return name
}
