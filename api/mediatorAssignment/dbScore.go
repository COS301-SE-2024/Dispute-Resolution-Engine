package mediatorassignment

import "time"

type DBScoreInput interface {
	GetScoreInput() (float64, error)
}

type DBScoreInputBase struct {
	DB     DBModel
	Column string
}

type DBScoreLastAssignmentstruct struct {
	DBScoreInputBase
}

func (d *DBScoreLastAssignmentstruct) GetScoreInput() (float64, error) {
	expertSummary, err := d.DB.GetExpertSummaryViewByColumn(d.Column)
	if err != nil {
		return 0, err
	}

	//calculate score current date - last assigned date
	score := float64(time.Until(expertSummary.LastAssignedDate).Hours())

	return score, nil
}

type DBScoreAssignedDisputes struct {
	DBScoreInputBase
}

func (d *DBScoreAssignedDisputes) GetScoreInput() (float64, error) {
	expertSummary, err := d.DB.GetExpertSummaryViewByColumn(d.Column)
	if err != nil {
		return 0, err
	}

	//calculate score assigned disputes
	score := float64(expertSummary.ActiveDisputeCount)

	return score, nil
}

type DBScoreRejectionCount struct {
	DBScoreInputBase
}

func (d *DBScoreRejectionCount) GetScoreInput() (float64, error) {
	expertSummary, err := d.DB.GetExpertSummaryViewByColumn(d.Column)
	if err != nil {
		return 0, err
	}

	//calculate score rejection count
	score := float64(expertSummary.RejectionPercentage)

	return score, nil
}