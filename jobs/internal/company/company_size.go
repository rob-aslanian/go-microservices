package company

// Size containg size of company. Amount of eployees.
type Size string

const (
	// SizeUnknown uknown amount
	SizeUnknown Size = "Unknown"
	// Size1To10 from 1 to 10 amount
	Size1To10 Size = "SIZE_1_10_EMPLOYEES"
	// Size11To50 from 11 to 50 amount
	Size11To50 Size = "SIZE_11_50_EMPLOYEES"
	// Size51To200 from 51 to 200 amount
	Size51To200 Size = "SIZE_51_200_EMPLOYEES"
	// Size201To500 from 201 to 500 amount
	Size201To500 Size = "SIZE_201_500_EMPLOYEES"
	// Size501To1000 from 501 to 1000 amount
	Size501To1000 Size = "SIZE_501_1000_EMPLOYEES"
	// Size1001To5000 from 1001 to 5000 amount
	Size1001To5000 Size = "SIZE_1001_5000_EMPLOYEES"
	// Size5001To10000 from 5001 to 10000 amount
	Size5001To10000 Size = "SIZE_5001_10000_EMPLOYEES"
	// Size10001Plus from 10001
	Size10001Plus Size = "SIZE_10001_PLUS_EMPLOYEES"
)
