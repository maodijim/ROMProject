package gameConnection

import (
	"sync"
	"testing"

	Cmd "ROMProject/Cmds"
	"ROMProject/utils"
	"github.com/mohae/deepcopy"
)

func TestGameConnection_EnchantPreviewContains(t *testing.T) {
	type fields struct {
		Role            *utils.RoleInfo
		Mutex           *sync.RWMutex
		BuffItems       map[uint32]utils.BuffItem
		BuffItemsByName map[string]utils.BuffItemByName
	}
	type args struct {
		equipGuid string
		preview   *EnchantCompare
	}
	atkType := Cmd.EAttrType_EATTRTYPE_ATK
	atkVal := uint32(100)
	atkValLow := uint32(1)
	guid := "0"
	// 尖锐3
	buffId := uint32(500043)
	notBuffId := uint32(500045)
	items := utils.NewItemsLoader("", "", "")
	role := utils.NewRole()
	role.PackItems = make(map[Cmd.EPackType]map[string]*Cmd.ItemData)
	role.PackItems[Cmd.EPackType_EPACKTYPE_EQUIP] = map[string]*Cmd.ItemData{
		"0": &Cmd.ItemData{
			Base: &Cmd.ItemInfo{
				Guid: &guid,
			},
			Enchant: &Cmd.EnchantData{
				Attrs: []*Cmd.EnchantAttr{
					{
						Type:  &atkType,
						Value: &atkValLow,
					},
				},
			},
			Previewenchant: &Cmd.EnchantData{
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
	}
	atkValHigh := uint32(150)
	role2 := deepcopy.Copy(role).(*utils.RoleInfo)
	role2.PackItems[Cmd.EPackType_EPACKTYPE_EQUIP]["0"].Enchant.Attrs[0].Value = &atkValHigh
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
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
			want: true,
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
			want: false,
		},
		{
			name: "TestGameConnection_EnchantPreviewContains",
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
					Attrs: []*EnchantAttrCompare{
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
			want: true,
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
					Attrs: []*EnchantAttrCompare{
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
			want: false,
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
			if got := g.EnchantPreviewContains(tt.args.equipGuid, tt.args.preview); got != tt.want {
				t.Errorf("EnchantPreviewContains() = %v, want %v", got, tt.want)
			}
		})
	}
}
