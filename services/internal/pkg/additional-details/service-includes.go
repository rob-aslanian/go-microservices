package additionaldetails

// ServiceIncludes is for additional details
type ServiceIncludes string

const (
	// ServiceIncludesSourceFile ...
	ServiceIncludesSourceFile ServiceIncludes = "source_file"
	// ServiceIncludesPrintReady ...
	ServiceIncludesPrintReady ServiceIncludes = "print_ready"
	// ServiceIncludesPhotoEditing ...
	ServiceIncludesPhotoEditing ServiceIncludes = "photo_editing"
	// ServiceIncludesCustomGraphics ...
	ServiceIncludesCustomGraphics ServiceIncludes = "custom_graphics"
	// ServiceIncludesStockPhotos ...
	ServiceIncludesStockPhotos ServiceIncludes = "stock_photos"
)
