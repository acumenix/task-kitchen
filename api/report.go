package api

import (
	"fmt"
	"time"

	"github.com/guregu/dynamo"
	"github.com/pkg/errors"
)

type ReportStatus string

const (
	// ReportEditing means the report is under editing.
	ReportEditing ReportStatus = "edit"
	// ReportWorking means the report is fixed.
	ReportWorking = "work"
	// ReportDone means works of the day has been finished, not tasks are completed.
	ReportDone = "done"
)

type Report struct {
	PKey      string       `dynamo:"pk" json:"-"`
	SKey      string       `dynamo:"sk" json:"-"`
	UserID    string       `dynamo:"user_id" json:"user_id"`
	CreatedAt time.Time    `dynamo:"created_at" json:"created_at"`
	Status    ReportStatus `dynamo:"status"`

	table dynamo.Table
}

func toReportKey(userID string, date time.Time) (string, string) {
	pk := fmt.Sprintf("%s/report", userID)
	sk := date.Format("20060102")
	return pk, sk
}

func (x KitchenManager) GetReport(userID string, date time.Time) (*Report, error) {
	var report Report
	pk, sk := toReportKey(userID, date)

	if err := x.table.Get("pk", pk).Range("sk", dynamo.Equal, sk).One(&report); err != nil {
		if err.Error() == "dynamo: no item found" {
			return nil, nil
		}
		return nil, errors.Wrapf(err, "Fail to get report: %s %s", pk, sk)
	}

	report.table = x.table
	return &report, nil
}

func (x KitchenManager) FetchReport(userID string, begin, end time.Time) ([]Report, error) {
	var reports []Report
	pk, sk1 := toReportKey(userID, begin)
	_, sk2 := toReportKey(userID, end)

	err := x.table.Get("pk", pk).
		Range("sk", dynamo.Between, sk1, sk2).
		All(&reports)

	if err != nil {
		if err.Error() == "dynamo: no item found" {
			return nil, nil
		}
		return nil, err
	}

	for i := range reports {
		reports[i].table = x.table
	}
	return reports, nil
}

func (x KitchenManager) NewReport(userID string, date time.Time) (*Report, error) {
	var report Report
	pk, sk := toReportKey(userID, date)

	report = Report{
		PKey:      pk,
		SKey:      sk,
		UserID:    userID,
		CreatedAt: date,
		Status:    ReportEditing,
		table:     x.table,
	}

	if err := report.Save(); err != nil {
		return nil, err
	}

	return &report, nil
}

func (x *Report) Save() error {
	if x.Status != ReportEditing && x.Status != ReportWorking && x.Status != ReportDone {
		return newUserError(400, "Invalid report status: '%s'", x.Status)
	}

	if err := x.table.Put(x).Run(); err != nil {
		return errors.Wrapf(err, "Fail to save report: %s", x.PKey)
	}

	return nil
}

func (x *Report) Delete() error {
	if err := x.table.Delete("pk", x.PKey).Range("sk", x.SKey).Run(); err != nil {
		return errors.Wrapf(err, "Fail to delete report: %s", x.PKey)
	}

	return nil
}
