package minersc

import (
	cstate "0chain.net/chaincore/chain/state"
	"0chain.net/chaincore/node"
	"0chain.net/chaincore/state"
	"0chain.net/smartcontract/dbs"
	"0chain.net/smartcontract/dbs/event"
	"encoding/json"
	"fmt"
)

func sharderTableToSharderNode(edbSharder *event.Sharder) *MinerNode {

	var status = node.NodeStatusInactive
	if edbSharder.Active {
		status = node.NodeStatusActive
	}
	msn := SimpleNode{
		ID:                edbSharder.SharderID,
		N2NHost:           edbSharder.N2NHost,
		Host:              edbSharder.Host,
		Port:              edbSharder.Port,
		Path:              edbSharder.Path,
		PublicKey:         edbSharder.PublicKey,
		ShortName:         edbSharder.ShortName,
		BuildTag:          edbSharder.BuildTag,
		TotalStaked:       int64(edbSharder.TotalStaked),
		Delete:            edbSharder.Delete,
		DelegateWallet:    edbSharder.DelegateWallet,
		ServiceCharge:     edbSharder.ServiceCharge,
		NumberOfDelegates: edbSharder.NumberOfDelegates,
		MinStake:          edbSharder.MinStake,
		MaxStake:          edbSharder.MaxStake,
		Stat: Stat{
			GeneratorRewards: edbSharder.Rewards,
			GeneratorFees:    edbSharder.Fees,
		},
		LastHealthCheck: edbSharder.LastHealthCheck,
		Status:          status,
	}

	return &MinerNode{
		SimpleNode: &msn,
	}

}

func sharderNodeToSharderTable(sn *MinerNode) event.Sharder {

	return event.Sharder{
		SharderID:         sn.ID,
		N2NHost:           sn.N2NHost,
		Host:              sn.Host,
		Port:              sn.Port,
		Path:              sn.Path,
		PublicKey:         sn.PublicKey,
		ShortName:         sn.ShortName,
		BuildTag:          sn.BuildTag,
		TotalStaked:       state.Balance(sn.TotalStaked),
		Delete:            sn.Delete,
		DelegateWallet:    sn.DelegateWallet,
		ServiceCharge:     sn.ServiceCharge,
		NumberOfDelegates: sn.NumberOfDelegates,
		MinStake:          sn.MinStake,
		MaxStake:          sn.MaxStake,
		LastHealthCheck:   sn.LastHealthCheck,
		Rewards:           sn.Stat.GeneratorRewards,
		Fees:              sn.Stat.GeneratorFees,
		Active:            sn.Status == node.NodeStatusActive,
		Longitude:         0,
		Latitude:          0,
	}
}

func emitAddSharder(sn *MinerNode, balances cstate.StateContextI) error {

	data, err := json.Marshal(sharderNodeToSharderTable(sn))
	if err != nil {
		return fmt.Errorf("marshalling sharder: %v", err)
	}

	balances.EmitEvent(event.TypeStats, event.TagAddSharder, sn.ID, string(data))

	return nil
}

func emitAddOrOverwriteSharder(sn *MinerNode, balances cstate.StateContextI) error {

	data, err := json.Marshal(sharderNodeToSharderTable(sn))
	if err != nil {
		return fmt.Errorf("marshalling sharder: %v", err)
	}

	balances.EmitEvent(event.TypeStats, event.TagAddOrOverwriteSharder, sn.ID, string(data))

	return nil
}

func emitUpdateSharder(sn *MinerNode, balances cstate.StateContextI, updateStatus bool) error {

	dbUpdates := dbs.DbUpdates{
		Id: sn.ID,
		Updates: map[string]interface{}{
			"n2n_host":            sn.N2NHost,
			"host":                sn.Host,
			"port":                sn.Port,
			"path":                sn.Path,
			"public_key":          sn.PublicKey,
			"short_name":          sn.ShortName,
			"build_tag":           sn.BuildTag,
			"total_staked":        sn.TotalStaked,
			"delete":              sn.Delete,
			"delegate_wallet":     sn.DelegateWallet,
			"service_charge":      sn.ServiceCharge,
			"number_of_delegates": sn.NumberOfDelegates,
			"min_stake":           sn.MinStake,
			"max_stake":           sn.MaxStake,
			"last_health_check":   sn.LastHealthCheck,
			"rewards":             sn.SimpleNode.Stat.GeneratorRewards,
			"fees":                sn.SimpleNode.Stat.GeneratorFees,
			"longitude":           0,
			"latitude":            0,
		},
	}

	if updateStatus {
		dbUpdates.Updates["active"] = sn.Status == node.NodeStatusActive
	}

	data, err := json.Marshal(dbUpdates)
	if err != nil {
		return fmt.Errorf("marshalling update: %v", err)
	}
	balances.EmitEvent(event.TypeStats, event.TagUpdateSharder, sn.ID, string(data))
	return nil
}

func emitDeleteSharder(id string, balances cstate.StateContextI) error {

	balances.EmitEvent(event.TypeStats, event.TagDeleteSharder, id, id)
	return nil
}