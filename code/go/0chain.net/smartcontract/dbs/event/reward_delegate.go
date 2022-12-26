package event

import (
	"0chain.net/smartcontract/common"
	"0chain.net/smartcontract/dbs"
	"0chain.net/smartcontract/stakepool/spenum"
	"github.com/0chain/common/core/currency"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RewardDelegate struct {
	gorm.Model
	Amount      currency.Coin `json:"amount"`
	BlockNumber int64         `json:"block_number" gorm:"index:idx_block,priority:1"`
	PoolID      string        `json:"pool_id" gorm:"index:idx_pool,priority:2"`
	RewardType  spenum.Reward `json:"reward_type" gorm:"index:idx_reward_type,priority:3"`
}

func (edb *EventDb) insertDelegateReward(updates []dbs.StakePoolReward, round int64) error {
	var drs []RewardDelegate
	for _, sp := range updates {
		for poolId, amount := range sp.DelegateRewards {
			dr := RewardDelegate{
				Amount:      amount,
				BlockNumber: round,
				PoolID:      poolId,
				RewardType:  sp.RewardType,
			}
			drs = append(drs, dr)
		}
		for poolId, amount := range sp.DelegatePenalties {
			dp := RewardDelegate{
				Amount:      amount,
				BlockNumber: round,
				PoolID:      poolId,
				RewardType:  sp.RewardType,
			}
			drs = append(drs, dp)
		}

	}
	if len(drs) == 0 {
		return nil
	}
	return edb.Get().Create(&drs).Error
}

func (edb *EventDb) GetDelegateRewards(limit common.Pagination, PoolId string, start, end int64) ([]RewardDelegate, error) {
	var rds []RewardDelegate
	query := edb.Get().Model(&RewardDelegate{})
	if PoolId == "" {
		if start == end {
			query = query.Where("block_number = ?", start)
		} else {
			query = query.Where("block_number >= ? AND block_number < ?", start, end)
		}
	} else {
		if start == end {
			query = query.Where("pool_id = ? AND block_number = ?", PoolId, start)
		} else {
			query = query.Where("pool_id = ? AND block_number >= ? AND block_number < ?", PoolId, start, end)
		}
	}
	return rds, query.Offset(limit.Offset).
		Limit(limit.Limit).
		Order(clause.OrderByColumn{
			Column: clause.Column{Name: "block_number"},
			Desc:   limit.IsDescending,
		}).Scan(&rds).Error
}