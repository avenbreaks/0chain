package services

import (
	"0chain.net/conductor/stores"
	"0chain.net/conductor/types"
	"0chain.net/conductor/utils"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var allocationStore = stores.GetAllocationStore()

type AllocationService struct {
	allocation *types.Allocation
	baseUrl    string
}

func NewAllocationService(baseUrl string) *AllocationService {
	return &AllocationService{
		allocation: &types.Allocation{},
		baseUrl:    baseUrl,
	}
}

func (s *AllocationService) CompareAllocationsValue() (bool, error) {
	// Read allocationID from file
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return false, err
	}

	filePath := filepath.Join(homeDir, ".zcn", "allocation.txt")
	allocationIDBytes, err := os.ReadFile(filePath)
	if err != nil {
		return false, err
	}

	allocationID := string(allocationIDBytes)

	remoteAllocation, err := s.getRemoteAllocation(allocationID)
	if err != nil {
		return false, err
	}

	localAllocation, err := allocationStore.GetLatest()
	if err != nil {
		return false, err
	}

	movedToChallengeDiffInFloat64 := float64(remoteAllocation.MovedToChallenge - localAllocation.MovedToChallenge)
	movedBackDiffInFloat64 := float64(remoteAllocation.MovedBack - localAllocation.MovedBack)

	if movedToChallengeDiffInFloat64 <= 1.05*movedBackDiffInFloat64 &&
		movedToChallengeDiffInFloat64 >= 0.95*movedBackDiffInFloat64 {
		return true, nil
	}

	return false, nil
}
func (s *AllocationService) getRemoteAllocation(allocationID string) (*types.Allocation, error) {
	url := fmt.Sprintf("%v/allocation?allocation_id=%s", s.baseUrl, allocationID)

	log.Printf("Getting allocation from %v\n", url)

	resp, err := utils.HttpGet(url, map[string]string{})
	if err != nil {
		return nil, err
	}

	alloc := &types.Allocation{}
	err = json.Unmarshal(resp, alloc)
	if err != nil {
		return nil, err
	}

	return alloc, nil
}
