package resolver

import (
	"strconv"

	"gitlab.lan/Rightnao-site/microservices/grpc-proto/companyRPC"

	"gitlab.lan/Rightnao-site/microservices/grpc-proto/userRPC"
)

func NullableString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func NullableStringClause(s string, isNull bool) *string {
	if isNull {
		return nil
	}
	return &s
}

func NullToString(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

func NullBoolToBool(s *bool) bool {
	if s != nil {
		return *s
	}
	return false
}

func NullToInt32(i *int32) int32 {
	if i != nil {
		return int32(*i)
	}
	return int32(0)
}

func PermissionToRPC(g string) userRPC.PermissionType {
	switch g {
	case "me":
		return userRPC.PermissionType_ME
	case "members":
		return userRPC.PermissionType_MEMBERS
	case "my_connections":
		return userRPC.PermissionType_MY_CONNECTIONS
	}

	return userRPC.PermissionType_NONE
}

func NullPermissionInputToRPC(p *PermissionInput) *userRPC.Permission {
	if p == nil {
		return &userRPC.Permission{
			Type: userRPC.PermissionType_NONE,
		}
	}

	return &userRPC.Permission{
		Type: PermissionToRPC(p.Type),
	}
}

func PermissionInputToRPC(p PermissionInput) *userRPC.Permission {
	return &userRPC.Permission{
		Type: PermissionToRPC(p.Type),
	}
}

func NullIDToInt32(g *string) int32 {
	if g != nil {
		n, err := strconv.Atoi(*g)
		if err != nil {
			return int32(0)
		}
		return int32(n)
	}
	return int32(0)
}

func NullIDToUint32(g *string) uint32 {
	if g != nil {
		n, err := strconv.Atoi(*g)
		if err != nil {
			return uint32(0)
		}
		return uint32(n)
	}
	return uint32(0)
}

func IDToInt32(g string) int32 {
	n, err := strconv.Atoi(g)
	if err != nil {
		return int32(0)
	}
	return int32(n)

}

func Int32ToID(g int32) string {
	return strconv.Itoa(int(g))
}

func NullStringArrayToStringArray(a *[]string) []string {
	if a == nil {
		return []string{}
	}
	return *a
}

func RPCSizeArrayToNullStringArray(s []companyRPC.Size) []string {
	strArray := make([]string, 0, len(s))

	for i := range strArray {
		strArray = append(strArray, string((s)[i]))
	}

	return strArray
}

func Nullint32ToUint32(g *int32) uint32 {
	if g == nil {
		return 0
	}
	return uint32(*g)
}

func Uint32Toint32(g uint32) int32 {
	return int32(g)
}

func NullInt32ArrayToUint32Array(s *[]int32) []uint32 {
	if s == nil {
		return []uint32{}
	}

	ui := make([]uint32, 0, len(*s))

	for i := range ui {
		ui = append(ui, uint32((*s)[i]))
	}

	return ui
}

func NullUInt32ArrayToInt32Array(s []uint32) []int32 {

	i32 := make([]int32, 0, len(s))

	for i := range i32 {
		i32 = append(i32, int32((s)[i]))
	}

	return i32
}

func float32ArrayToFloat64Array(s []float32) []float64 {
	n := make([]float64, 0, len(s))

	for i := range s {
		n = append(n, float64(s[i]))
	}

	return n
}
