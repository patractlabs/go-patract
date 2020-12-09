package main

import (
	"fmt"

	gsrpc "github.com/centrifuge/go-substrate-rpc-client"
	"github.com/centrifuge/go-substrate-rpc-client/config"
	"github.com/centrifuge/go-substrate-rpc-client/types"
)

func printEventRecords(events *types.EventRecords) {
	// Show what we are busy with
	for _, e := range events.Balances_Endowed {
		fmt.Printf("\tBalances:Endowed:: (phase=%#v)\n", e.Phase)
		fmt.Printf("\t\t%#x, %v\n", e.Who, e.Balance)
	}
	for _, e := range events.Balances_DustLost {
		fmt.Printf("\tBalances:DustLost:: (phase=%#v)\n", e.Phase)
		fmt.Printf("\t\t%#x, %v\n", e.Who, e.Balance)
	}
	for _, e := range events.Balances_Transfer {
		fmt.Printf("\tBalances:Transfer:: (phase=%#v)\n", e.Phase)
		fmt.Printf("\t\t%v, %v, %v\n", e.From, e.To, e.Value)
	}
	for _, e := range events.Balances_BalanceSet {
		fmt.Printf("\tBalances:BalanceSet:: (phase=%#v)\n", e.Phase)
		fmt.Printf("\t\t%v, %v, %v\n", e.Who, e.Free, e.Reserved)
	}
	for _, e := range events.Balances_Deposit {
		fmt.Printf("\tBalances:Deposit:: (phase=%#v)\n", e.Phase)
		fmt.Printf("\t\t%v, %v\n", e.Who, e.Balance)
	}
	for _, e := range events.Grandpa_NewAuthorities {
		fmt.Printf("\tGrandpa:NewAuthorities:: (phase=%#v)\n", e.Phase)
		fmt.Printf("\t\t%v\n", e.NewAuthorities)
	}
	for _, e := range events.Grandpa_Paused {
		fmt.Printf("\tGrandpa:Paused:: (phase=%#v)\n", e.Phase)
	}
	for _, e := range events.Grandpa_Resumed {
		fmt.Printf("\tGrandpa:Resumed:: (phase=%#v)\n", e.Phase)
	}
	for _, e := range events.ImOnline_HeartbeatReceived {
		fmt.Printf("\tImOnline:HeartbeatReceived:: (phase=%#v)\n", e.Phase)
		fmt.Printf("\t\t%#x\n", e.AuthorityID)
	}
	for _, e := range events.ImOnline_AllGood {
		fmt.Printf("\tImOnline:AllGood:: (phase=%#v)\n", e.Phase)
	}
	for _, e := range events.ImOnline_SomeOffline {
		fmt.Printf("\tImOnline:SomeOffline:: (phase=%#v)\n", e.Phase)
		fmt.Printf("\t\t%v\n", e.IdentificationTuples)
	}
	for _, e := range events.Indices_IndexAssigned {
		fmt.Printf("\tIndices:IndexAssigned:: (phase=%#v)\n", e.Phase)
		fmt.Printf("\t\t%#x%v\n", e.AccountID, e.AccountIndex)
	}
	for _, e := range events.Indices_IndexFreed {
		fmt.Printf("\tIndices:IndexFreed:: (phase=%#v)\n", e.Phase)
		fmt.Printf("\t\t%v\n", e.AccountIndex)
	}
	for _, e := range events.Offences_Offence {
		fmt.Printf("\tOffences:Offence:: (phase=%#v)\n", e.Phase)
		fmt.Printf("\t\t%v%v\n", e.Kind, e.OpaqueTimeSlot)
	}
	for _, e := range events.Session_NewSession {
		fmt.Printf("\tSession:NewSession:: (phase=%#v)\n", e.Phase)
		fmt.Printf("\t\t%v\n", e.SessionIndex)
	}
	for _, e := range events.Staking_Reward {
		fmt.Printf("\tStaking:Reward:: (phase=%#v)\n", e.Phase)
		fmt.Printf("\t\t%v\n", e.Amount)
	}
	for _, e := range events.Staking_Slash {
		fmt.Printf("\tStaking:Slash:: (phase=%#v)\n", e.Phase)
		fmt.Printf("\t\t%#x%v\n", e.AccountID, e.Balance)
	}
	for _, e := range events.Staking_OldSlashingReportDiscarded {
		fmt.Printf("\tStaking:OldSlashingReportDiscarded:: (phase=%#v)\n", e.Phase)
		fmt.Printf("\t\t%v\n", e.SessionIndex)
	}
	for _, e := range events.System_ExtrinsicSuccess {
		fmt.Printf("\tSystem:ExtrinsicSuccess:: (phase=%#v)\n", e.Phase)
	}
	for _, e := range events.System_ExtrinsicFailed {
		fmt.Printf("\tSystem:ExtrinsicFailed:: (phase=%#v)\n", e.Phase)
		fmt.Printf("\t\t%v\n", e.DispatchError)
	}
	for _, e := range events.System_CodeUpdated {
		fmt.Printf("\tSystem:CodeUpdated:: (phase=%#v)\n", e.Phase)
	}
	for _, e := range events.System_NewAccount {
		fmt.Printf("\tSystem:NewAccount:: (phase=%#v)\n", e.Phase)
		fmt.Printf("\t\t%#x\n", e.Who)
	}
	for _, e := range events.System_KilledAccount {
		fmt.Printf("\tSystem:KilledAccount:: (phase=%#v)\n", e.Phase)
		fmt.Printf("\t\t%#X\n", e.Who)
	}
}

// queryEvent Query the system events and extract information from them. This example runs until exited via Ctrl-C
func queryEvent() {
	// Create our API with a default connection to the local node
	api, err := gsrpc.NewSubstrateAPI(config.Default().RPCURL)
	if err != nil {
		panic(err)
	}

	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		panic(err)
	}

	// Subscribe to system events via storage
	key, err := types.CreateStorageKey(meta, "System", "Events", nil, nil)
	if err != nil {
		panic(err)
	}

	sub, err := api.RPC.State.SubscribeStorageRaw([]types.StorageKey{key})
	if err != nil {
		panic(err)
	}
	defer sub.Unsubscribe()

	// outer for loop for subscription notifications
	for {
		set := <-sub.Chan()
		// inner loop for the changes within one of those notifications
		for _, chng := range set.Changes {
			if !types.Eq(chng.StorageKey, key) || !chng.HasStorageData {
				// skip, we are only interested in events with content
				continue
			}

			// Decode the event records
			events := types.EventRecords{}
			err = types.EventRecordsRaw(chng.StorageData).DecodeEventRecords(meta, &events)
			if err != nil {
				panic(err)
			}

			printEventRecords(&events)
		}
	}
}

// queryEventByBlock Query the system events and extract information from them. This example runs until exited via Ctrl-C
func queryEventByBlock(hash types.Hash) {
	// Create our API with a default connection to the local node
	api, err := gsrpc.NewSubstrateAPI(config.Default().RPCURL)
	if err != nil {
		panic(err)
	}

	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		panic(err)
	}

	// Subscribe to system events via storage
	key, err := types.CreateStorageKey(meta, "System", "Events", nil, nil)
	if err != nil {
		panic(err)
	}

	set, err := api.RPC.State.GetStorageRaw(key, hash)

	// Decode the event records
	events := types.EventRecords{}
	err = types.EventRecordsRaw(*set).DecodeEventRecords(meta, &events)
	if err != nil {
		panic(err)
	}

	printEventRecords(&events)

}
