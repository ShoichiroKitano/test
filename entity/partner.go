package entity

// TODO: 取引先企業名は暗号化する
type Partner struct {
	ID          uint64
	CorporateID uint64
}

type Partners []*Partner

func (partners Partners) FindByPartnerID(id uint64) *Partner {
	for _, p := range partners {
		if p.ID == id {
			return p
		}
	}
	return nil
}
