package utils


func RemoveElement(slice []int64, element int64) []int64 {
	for i, v := range slice {
		if v == element {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}