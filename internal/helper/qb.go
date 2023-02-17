// Package helper
// Author Daud Valentino
package helper

import (
	"fmt"

	"GolangBookingApp/pkg/util"

	"github.com/spf13/cast"
)

// StructQueryWhere create query builder from struct
func StructQueryWhere(i interface{}, hideDeleted bool, tag string) (q string, vals []interface{}, limit, page uint64, err error) {

	var cols []string
	var startDate, endDate, sortOrder, groupBy string

	if i == nil {
		return q, vals, limit, page, nil
	}

	data, err := util.StructToMap(i, tag)
	if err != nil {
		return q, vals, limit, page, err
	}

	if len(data) == 0 {
		return q, vals, limit, page, err
	}

	for k, x := range data {
		if k == "page" {
			page = cast.ToUint64(x)
			continue
		}

		if k == "limit" {
			limit = cast.ToUint64(x)
			continue
		}

		if k == "start_date" {
			startDate = cast.ToString(x)
			continue
		}

		if k == "end_date" {
			endDate = cast.ToString(x)
			continue
		}

		if k == "sort_order" {
			sortOrder = cast.ToString(x)
			continue
		}

		if k == "group_by" {
			groupBy = cast.ToString(x)
			continue
		}

		vals = append(vals, x)
		cols = append(cols, k)
	}

	if len(cols) > 0 && !hideDeleted {
		q = fmt.Sprintf(`WHERE %s`, util.StringJoin(cols, "=? AND ", "=?"))
	}

	if len(cols) > 0 && hideDeleted {
		q = fmt.Sprintf(`WHERE %s AND deleted_at = '0000-00-00 00:00:00'`, util.StringJoin(cols, "=? AND ", "=?"))
	}

	if len(cols) < 1 && hideDeleted {
		q = fmt.Sprint(`WHERE deleted_at = '0000-00-00 00:00:00'`)
	}

	if startDate != "" && endDate != "" {
		q = fmt.Sprintf(`%s AND ( DATE(created_at) >= ?  AND DATE(created_at) <= ? )`, q)
		if hideDeleted {
			q = fmt.Sprintf(`%s AND ( DATE(created_at) >= ?  AND DATE(created_at) <= ? ) AND deleted_at = '0000-00-00 00:00:00'`, q)
		}

		if len(cols) == 0 && !hideDeleted {
			q = fmt.Sprint(`WHERE ( DATE(created_at) >= ?  AND DATE(created_at) <= ? )`)
		}

		if len(cols) == 0 && hideDeleted {
			q = fmt.Sprint(`WHERE ( DATE(created_at) >= ?  AND DATE(created_at) <= ? ) AND deleted_at = '0000-00-00 00:00:00'`)
		}

		vals = append(vals, startDate, endDate)
	}

	if groupBy != "" {
		q = fmt.Sprintf("%s GROUP BY %s", q, groupBy)
	}

	if sortOrder != "" {
		q = fmt.Sprintf("%s ORDER BY created_at %s", q, sortOrder)
	}

	return q, vals, limit, page, err
}
