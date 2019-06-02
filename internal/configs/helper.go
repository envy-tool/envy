package configs

// IsReservedHelperType returns true if the given helper type would create
// a helper that cannot be referenced, due to its type name instead indicating
// that a different kind of object is being referenced.
func IsReservedHelperType(proposed string) bool {
	_, exists := reservedHelperTypeNames[proposed]
	return exists
}
