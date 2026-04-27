package gameConnection

import (
	"sync"
	"testing"

	Cmd "ROMProject/Cmds"
	"ROMProject/utils"
)

func TestGameConnection_EnchantPreviewContains(t *testing.T) {
	type fields struct {
		Role            *RoleInfo
		Mutex           *sync.RWMutex
		BuffItems       map[uint32]utils.BuffItem
		BuffItemsByName map[string]utils.BuffItemByName
	}
	type args struct {
		equipGuid     string
		preview       *EnchantCompare
		bothCondition bool
		allAttrMatch  bool
	}
	hpType := Cmd.EAttrType_EATTRTYPE_HP
	hpVal := uint32(1000)
	hpValLow := uint32(1)
	hpValHigh := uint32(2000)
	atkDefType := Cmd.EAttrType_EATTRTYPE_DEF
	atkDefVal := uint32(500)
	atkType := Cmd.EAttrType_EATTRTYPE_ATK
	atkVal := uint32(100)
	atkValLow := uint32(1)
	atkValHigh := uint32(150)
	mAtkType := Cmd.EAttrType_EATTRTYPE_MATK
	matkVal := uint32(50)
	matkValLow := uint32(1)
	matkValHigh := uint32(200)
	guid0 := "0"
	guid1 := "1"
	guid2 := "2"
	guid3 := "3"
	// 尖锐3
	buffId := uint32(500043)
	notBuffId := uint32(500045)
	items := utils.NewItemsLoader("", "", "")
	role := NewRole()
	role.PackItems = make(map[Cmd.EPackType]map[string]*Cmd.ItemData)
	role.PackItems[Cmd.EPackType_EPACKTYPE_EQUIP] = map[string]*Cmd.ItemData{
		"0": {
			Base: &Cmd.ItemInfo{
				Guid: &guid0,
			},
			Enchant: &Cmd.EnchantData{
				Attrs: []*Cmd.EnchantAttr{
					{
						Type:  &atkType,
						Value: &atkValLow,
					},
				},
			},
			Previewenchant: []*Cmd.EnchantData{
				{
					Extras: []*Cmd.EnchantExtra{
						{
							Buffid: &buffId,
						},
					},
					Attrs: []*Cmd.EnchantAttr{
						{
							Type:  &atkType,
							Value: &atkVal,
						},
					},
				},
			},
		},
	}
	role2 := NewRole()
	role2.PackItems = make(map[Cmd.EPackType]map[string]*Cmd.ItemData)
	role2.PackItems[Cmd.EPackType_EPACKTYPE_EQUIP] = map[string]*Cmd.ItemData{
		"0": {
			Base: &Cmd.ItemInfo{
				Guid: &guid0,
			},
			Enchant: &Cmd.EnchantData{
				Attrs: []*Cmd.EnchantAttr{
					{
						Type:  &atkType,
						Value: &atkValHigh,
					},
				},
			},
		},
		"1": {
			Base: &Cmd.ItemInfo{
				Guid: &guid1,
			},
			Enchant: &Cmd.EnchantData{
				Attrs: []*Cmd.EnchantAttr{
					{
						Type:  &mAtkType,
						Value: &matkValLow,
					},
					{
						Type:  &atkType,
						Value: &atkValLow,
					},
				},
			},
			Previewenchant: []*Cmd.EnchantData{
				{
					Extras: []*Cmd.EnchantExtra{
						{
							Buffid: &buffId,
						},
					},
					Attrs: []*Cmd.EnchantAttr{
						{
							Type:  &mAtkType,
							Value: &matkVal,
						},
						{
							Type:  &atkType,
							Value: &atkVal,
						},
					},
				},
			},
		},
		"2": {
			Base: &Cmd.ItemInfo{
				Guid: &guid2,
			},
			Enchant: &Cmd.EnchantData{
				Attrs: []*Cmd.EnchantAttr{
					{
						Type:  &hpType,
						Value: &atkValLow,
					},
					{
						Type:  &atkType,
						Value: &atkValLow,
					},
				},
			},
			Previewenchant: []*Cmd.EnchantData{
				{
					Extras: []*Cmd.EnchantExtra{
						{
							Buffid: &buffId,
						},
					},
					Attrs: []*Cmd.EnchantAttr{
						{
							Type:  &mAtkType,
							Value: &matkValHigh,
						},
						{
							Type:  &atkType,
							Value: &atkVal,
						},
					},
				},
			},
		},
		// item to test 3 attrs
		"3": {
			Base: &Cmd.ItemInfo{
				Guid: &guid3,
			},
			Enchant: &Cmd.EnchantData{
				Attrs: []*Cmd.EnchantAttr{
					{
						Type:  &atkDefType,
						Value: &atkDefVal,
					},
					{
						Type:  &hpType,
						Value: &hpValLow,
					},
					{
						Type:  &atkType,
						Value: &atkValLow,
					},
				},
			},
			Previewenchant: []*Cmd.EnchantData{
				{
					Extras: []*Cmd.EnchantExtra{
						{
							Buffid: &buffId,
						},
					},
					Attrs: []*Cmd.EnchantAttr{
						{
							Type:  &hpType,
							Value: &hpValLow,
						},
						{
							Type:  &mAtkType,
							Value: &matkValLow,
						},
						{
							Type:  &atkType,
							Value: &atkValLow,
						},
					},
				},
				{
					Extras: []*Cmd.EnchantExtra{
						{
							Buffid: &buffId,
						},
					},
					Attrs: []*Cmd.EnchantAttr{
						{
							Type:  &hpType,
							Value: &hpValHigh,
						},
						{
							Type:  &mAtkType,
							Value: &matkValHigh,
						},
						{
							Type:  &atkType,
							Value: &atkValHigh,
						},
					},
				},
			},
		},
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		want          bool
		wantTargetNum int
	}{
		{
			name: "TestGameConnection_EnchantPreviewContains_hasExtras",
			fields: fields{
				Role:            role,
				BuffItems:       items.BuffItems,
				BuffItemsByName: items.BuffItemsByName,
			},
			args: args{
				equipGuid: "0",
				preview: &EnchantCompare{
					EnchantData: Cmd.EnchantData{
						Extras: []*Cmd.EnchantExtra{
							{
								Buffid: &buffId,
							},
						},
					},
				},
			},
			want:          true,
			wantTargetNum: 0,
		},
		{
			name: "TestGameConnection_EnchantPreviewContains_NoMatchExtras",
			fields: fields{
				Role:            role,
				BuffItems:       items.BuffItems,
				BuffItemsByName: items.BuffItemsByName,
			},
			args: args{
				equipGuid: "0",
				preview: &EnchantCompare{
					EnchantData: Cmd.EnchantData{
						Extras: []*Cmd.EnchantExtra{
							{
								Buffid: &notBuffId,
							},
						},
					},
				},
			},
			want:          false,
			wantTargetNum: 0,
		},
		{
			name: "TestGameConnection_EnchantPreviewContainsHigherAttr",
			fields: fields{
				BuffItems:       items.BuffItems,
				BuffItemsByName: items.BuffItemsByName,
				Role:            role,
			},
			args: args{
				equipGuid: "0",
				preview: &EnchantCompare{
					EnchantData: Cmd.EnchantData{
						Extras: []*Cmd.EnchantExtra{},
					},
					Attrs: [][]*EnchantAttrCompare{
						{
							{
								EnchantAttr: Cmd.EnchantAttr{
									Type:  &atkType,
									Value: &atkVal,
								},
								Condition: ">=",
							},
						},
					},
				},
			},
			want:          true,
			wantTargetNum: 0,
		},
		{
			name: "TestGameConnection_EnchantPreviewContains_HasHighterAttr_Now",
			fields: fields{
				BuffItems:       items.BuffItems,
				BuffItemsByName: items.BuffItemsByName,
				Role:            role2,
			},
			args: args{
				equipGuid: "0",
				preview: &EnchantCompare{
					EnchantData: Cmd.EnchantData{
						Extras: []*Cmd.EnchantExtra{},
					},
					Attrs: [][]*EnchantAttrCompare{
						{
							{
								EnchantAttr: Cmd.EnchantAttr{
									Type:  &atkType,
									Value: &atkVal,
								},
								Condition: ">=",
							},
						},
					},
				},
			},
			want:          false,
			wantTargetNum: 0,
		},
		{
			name: "TestGameConnection_EnchantPreviewContains_AllAttrMatch_True",
			fields: fields{
				BuffItems:       items.BuffItems,
				BuffItemsByName: items.BuffItemsByName,
				Role:            role2,
			},
			args: args{
				equipGuid:     "1",
				bothCondition: false,
				allAttrMatch:  true,
				preview: &EnchantCompare{
					EnchantData: Cmd.EnchantData{
						Extras: []*Cmd.EnchantExtra{},
					},
					Attrs: [][]*EnchantAttrCompare{
						{
							{
								EnchantAttr: Cmd.EnchantAttr{
									Type:  &mAtkType,
									Value: &matkVal,
								},
								Condition: ">=",
							},
							{
								EnchantAttr: Cmd.EnchantAttr{
									Type:  &atkType,
									Value: &atkVal,
								},
								Condition: ">=",
							},
						},
					},
				},
			},
			want:          true,
			wantTargetNum: 0,
		},
		{
			name: "TestGameConnection_EnchantPreviewContains_AllAttrMatch_False",
			fields: fields{
				BuffItems:       items.BuffItems,
				BuffItemsByName: items.BuffItemsByName,
				Role:            role2,
			},
			args: args{
				equipGuid:     "1",
				bothCondition: false,
				allAttrMatch:  true,
				preview: &EnchantCompare{
					EnchantData: Cmd.EnchantData{
						Extras: []*Cmd.EnchantExtra{},
					},
					Attrs: [][]*EnchantAttrCompare{
						{
							{
								EnchantAttr: Cmd.EnchantAttr{
									Type:  &mAtkType,
									Value: &matkValLow,
								},
								Condition: ">=",
							},
							{
								EnchantAttr: Cmd.EnchantAttr{
									Type:  &atkType,
									Value: &atkValHigh,
								},
								Condition: ">",
							},
						},
					},
				},
			},
			want:          false,
			wantTargetNum: 0,
		},
		{
			name: "TestGameConnection_EnchantPreviewContains_BothCondition",
			fields: fields{
				BuffItems:       items.BuffItems,
				BuffItemsByName: items.BuffItemsByName,
				Role:            role2,
			},
			args: args{
				equipGuid:     "2",
				bothCondition: true,
				allAttrMatch:  false,
				preview: &EnchantCompare{
					EnchantData: Cmd.EnchantData{
						Extras: []*Cmd.EnchantExtra{
							{
								Buffid: &buffId,
							},
						},
					},
					Attrs: [][]*EnchantAttrCompare{
						{
							{
								EnchantAttr: Cmd.EnchantAttr{
									Type:  &mAtkType,
									Value: &matkValLow,
								},
								Condition: ">=",
							},
							{
								EnchantAttr: Cmd.EnchantAttr{
									Type:  &atkType,
									Value: &atkValLow,
								},
								Condition: ">=",
							},
						},
					},
				},
			},
			want:          true,
			wantTargetNum: 0,
		},
		{
			name: "TestGameConnection_EnchantPreviewContains_BothCondition_AllAttrNotMatch",
			fields: fields{
				BuffItems:       items.BuffItems,
				BuffItemsByName: items.BuffItemsByName,
				Role:            role2,
			},
			args: args{
				equipGuid:     "2",
				bothCondition: true,
				allAttrMatch:  true,
				preview: &EnchantCompare{
					EnchantData: Cmd.EnchantData{
						Extras: []*Cmd.EnchantExtra{
							{
								Buffid: &buffId,
							},
						},
					},
					Attrs: [][]*EnchantAttrCompare{
						{
							{
								EnchantAttr: Cmd.EnchantAttr{
									Type:  &mAtkType,
									Value: &matkValLow,
								},
								Condition: ">=",
							},
							{
								EnchantAttr: Cmd.EnchantAttr{
									Type:  &atkType,
									Value: &atkValHigh,
								},
								Condition: ">=",
							},
						},
					},
				},
			},
			want:          false,
			wantTargetNum: 0,
		},
		{
			name: "TestGameConnection_EnchantPreviewContains_ThreeAttrs_AllMatch",
			fields: fields{
				BuffItems:       items.BuffItems,
				BuffItemsByName: items.BuffItemsByName,
				Role:            role2,
			},
			args: args{
				equipGuid:     "3",
				bothCondition: false,
				allAttrMatch:  true,
				preview: &EnchantCompare{
					EnchantData: Cmd.EnchantData{
						Extras: []*Cmd.EnchantExtra{},
					},
					Attrs: [][]*EnchantAttrCompare{
						{
							{
								EnchantAttr: Cmd.EnchantAttr{
									Type:  &hpType,
									Value: &hpVal,
								},
								Condition: ">=",
							},
							{
								EnchantAttr: Cmd.EnchantAttr{
									Type:  &mAtkType,
									Value: &matkVal,
								},
								Condition: ">=",
							},
							{
								EnchantAttr: Cmd.EnchantAttr{
									Type:  &atkType,
									Value: &atkVal,
								},
								Condition: ">=",
							},
						},
					},
				},
			},
			want:          true,
			wantTargetNum: 1,
		},
		// test item 3 not all attrs match
		{
			name: "TestGameConnection_EnchantPreviewContains_ThreeAttrs_NotAllMatch",
			fields: fields{
				BuffItems:       items.BuffItems,
				BuffItemsByName: items.BuffItemsByName,
				Role:            role2,
			},
			args: args{
				equipGuid:     "3",
				bothCondition: false,
				allAttrMatch:  true,
				preview: &EnchantCompare{
					EnchantData: Cmd.EnchantData{
						Extras: []*Cmd.EnchantExtra{},
					},
					Attrs: [][]*EnchantAttrCompare{
						{
							{
								EnchantAttr: Cmd.EnchantAttr{
									Type:  &hpType,
									Value: &hpVal,
								},
								Condition: ">=",
							},
							{
								EnchantAttr: Cmd.EnchantAttr{
									Type:  &mAtkType,
									Value: &matkValHigh,
								},
								Condition: ">",
							},
							{
								EnchantAttr: Cmd.EnchantAttr{
									Type:  &atkType,
									Value: &atkVal,
								},
								Condition: ">=",
							},
						},
					},
				},
			},
			want:          false,
			wantTargetNum: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GameConnection{
				Role:            tt.fields.Role,
				Mutex:           tt.fields.Mutex,
				BuffItems:       tt.fields.BuffItems,
				BuffItemsByName: tt.fields.BuffItemsByName,
			}
			got, targetNum := g.EnchantPreviewContains(tt.args.equipGuid, tt.args.preview, tt.args.bothCondition, tt.args.allAttrMatch)
			if got != tt.want {
				t.Errorf("EnchantPreviewContains() = %v, want %v", got, tt.want)
			}
			if int(targetNum) != tt.wantTargetNum {
				t.Errorf("EnchantPreviewContains() targetNum = %v, want %v", targetNum, tt.wantTargetNum)
			}
		})
	}
}
