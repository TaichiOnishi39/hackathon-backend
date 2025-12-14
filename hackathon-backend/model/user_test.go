package model

import "testing"

func TestUserReqForHTTPPost_Validate(t *testing.T) {
	type fields struct {
		Name string
		Age  int
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "正常系: 正しいデータ",
			fields:  fields{Name: "Taro", Age: 25},
			wantErr: false, // エラーは出ないはず
		},
		{
			name:    "異常系: 名前が空",
			fields:  fields{Name: "", Age: 25},
			wantErr: true, // エラーが出るはず
		},
		{
			name:    "異常系: 名前が長すぎる(51文字)",
			fields:  fields{Name: "123456789012345678901234567890123456789012345678901", Age: 25},
			wantErr: true,
		},
		{
			name:    "異常系: 年齢が若すぎる",
			fields:  fields{Name: "Jiro", Age: 19},
			wantErr: true,
		},
		{
			name:    "異常系: 年齢が高すぎる",
			fields:  fields{Name: "Saburo", Age: 81},
			wantErr: true,
		},
		{
			name:    "境界値: 年齢が20歳(OK)",
			fields:  fields{Name: "Shiro", Age: 20},
			wantErr: false,
		},
		{
			name:    "境界値: 年齢が80歳(OK)",
			fields:  fields{Name: "Goro", Age: 80},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserReqForHTTPPost{
				Name: tt.fields.Name,
				Age:  tt.fields.Age,
			}
			if err := u.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
