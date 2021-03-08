package models

type Result struct {
	Id     int64  `json:"id"`
	Status Status `json:"status"`
}

func NewResult(id int64, status Status) Result {
	return Result{id, status}
}

func (p *Result) GetId() int64 {
	return p.Id
}

func (p *Result) GetStatus() Status {
	return p.Status
}

func (p *Result) SetId(id int64) {
	p.Id = id
}

func (p *Result) SetStatus(status Status) {
	p.Status = status
}
