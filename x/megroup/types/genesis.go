package types

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Groups:               []GroupInfo{},
		GroupMembers:         []GroupMember{},
		MemberJoinedList:     []MemberJoined{},
		GroupMemberCountList: []GroupMemberCount{},
		// this line is used by starport scaffolding # genesis/types/default
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// Check for duplicated ID in group
	//todo:稍后在做
	return nil
	/*
		groupIdMap := make(map[uint64]bool)
		groupCount := gs.GetGroupCount()
		for _, elem := range gs.GroupList {
			if _, ok := groupIdMap[elem.Id]; ok {
				return fmt.Errorf("duplicated id for group")
			}
			if elem.Id >= groupCount {
				return fmt.Errorf("group id should be lower or equal than the last id")
			}
			groupIdMap[elem.Id] = true
		}
		// Check for duplicated ID in groupMember
		groupMemberIdMap := make(map[uint64]bool)
		var groupMemberCount uint64
		for _, elem := range gs.GroupMemberCountList {
			groupMemberCount += elem.Num
		}

		for _, elem := range gs.GroupMemberList {
			if _, ok := groupMemberIdMap[elem.Id]; ok {
				return fmt.Errorf("duplicated id for groupMember")
			}
			if elem.Id >= groupMemberCount {
				return fmt.Errorf("groupMember id should be lower or equal than the last id")
			}
			groupMemberIdMap[elem.Id] = true
		}
		// Check for duplicated index in memberJoined
		memberJoinedIndexMap := make(map[string]struct{})

		for _, elem := range gs.MemberJoinedList {
			index := string(MemberJoinedKey(elem.Address))
			if _, ok := memberJoinedIndexMap[index]; ok {
				return fmt.Errorf("duplicated index for memberJoined")
			}
			memberJoinedIndexMap[index] = struct{}{}
		}
		// Check for duplicated index in groupMemberCount
		groupMemberCountIndexMap := make(map[string]struct{})

		for _, elem := range gs.GroupMemberCountList {
			index := string(GroupMemberCountKey(elem.GroupId))
			if _, ok := groupMemberCountIndexMap[index]; ok {
				return fmt.Errorf("duplicated index for groupMemberCount")
			}
			groupMemberCountIndexMap[index] = struct{}{}
		}
		// this line is used by starport scaffolding # genesis/types/validate

		return gs.Params.Validate()

	*/
}
