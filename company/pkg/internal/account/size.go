package account

// Size ...
type Size string

// SIZE_UNKNOWN = 0;
// SIZE_SELF_EMPLOYED = 1;
// SIZE_1_10_EMPLOYEES = 2;
// SIZE_11_50_EMPLOYEES = 3;
// SIZE_51_200_EMPLOYEES = 4;
// SIZE_201_500_EMPLOYEES = 5;
// SIZE_501_1000_EMPLOYEES = 6;
// SIZE_1001_5000_EMPLOYEES = 7;
// SIZE_5001_10000_EMPLOYEES = 8;
// SIZE_10001_PLUS_EMPLOYEES = 9;

const (
	// SizeUnknown ...
	SizeUnknown Size = "SIZE_UNKNOWN"
	// SizeSelfEmployed ...
	SizeSelfEmployed Size = "SIZE_SELF_EMPLOYED"
	// SizeFrom1Till10Employees ...
	SizeFrom1Till10Employees Size = "SIZE_1_10_EMPLOYEES"
	// SizeFrom11Till50Employees ...
	SizeFrom11Till50Employees Size = "SIZE_11_50_EMPLOYEES"
	// SizeFrom51Till200Employees ...
	SizeFrom51Till200Employees Size = "SIZE_51_200_EMPLOYEES"
	// SizeFrom201Till500Employees ...
	SizeFrom201Till500Employees Size = "SIZE_201_500_EMPLOYEES"
	// SizeFrom501Till1000Employees ...
	SizeFrom501Till1000Employees Size = "SIZE_501_1000_EMPLOYEES"
	// SizeFrom1001Till5000Employees ...
	SizeFrom1001Till5000Employees Size = "SIZE_1001_5000_EMPLOYEES"
	// SizeFrom5001Till10000Employees ...
	SizeFrom5001Till10000Employees Size = "SIZE_5001_10000_EMPLOYEES"
	// SizeFrom10001AndMoreEmployees ...
	SizeFrom10001AndMoreEmployees Size = "SIZE_10001_PLUS_EMPLOYEES"
)
