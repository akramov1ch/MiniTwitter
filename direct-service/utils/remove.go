package utils


func RemoveElement(slice []int32, element int32) []int32 {
	for i, v := range slice {
		if v == element {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}