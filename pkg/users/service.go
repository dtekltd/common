package users

import (
	"fmt"
	"reflect"

	"github.com/dtekltd/common/database"
)

func BuildReferrers(acc *Account) (*AccountReferrers, error) {
	ref := AccountReferrers{ID: acc.ID}
	rv := reflect.ValueOf(&ref).Elem()
	setRef(acc, &rv, acc.ReferrerID.Uint64, 1)
	if err := database.DB.Create(ref).Error; err != nil {
		return nil, err
	}
	return &ref, nil
}

func setRef(acc *Account, rv *reflect.Value, refID uint64, idx int) {
	if refID > 0 && idx <= 10 {
		fv := rv.FieldByName(fmt.Sprintf("R%d", idx))
		if fv.IsValid() && fv.CanSet() {
			fv.SetUint(refID)
		}
		var nextRefIDs []database.NullableUint64
		if err := database.DB.Model(acc).
			Select("referrer_id").Where("id=?", refID).
			Find(&nextRefIDs).Error; err == nil {
			idx += 1
			setRef(acc, rv, nextRefIDs[0].Uint64, idx)
		}
	}
}
