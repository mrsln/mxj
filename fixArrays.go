package mxj

import "strings"

// Fixing artifacts after xml parsing (
// sometimes we can't tell whether it's an array or an object
// because there is only one child).
func (mv Map) FixArrays(pluralKey string) {
	// TODO (mrsln): show an example
	m := map[string]interface{}(mv)
	fixArrays(m, pluralKey)
}

func fixArrays(m interface{}, pluralKey string) {
	// TODO: add extra check for subKeys or ancestors
	originalPluralKey := pluralKey
	sep := "."
	subKey := ""
	isArrayCanBeInObject := false
	if strings.Contains(originalPluralKey, sep) {
		tmpAry := strings.Split(originalPluralKey, sep)
		pluralKey = tmpAry[0]
		subKey = tmpAry[1]
		isArrayCanBeInObject = true
	}
	switch m.(type) {
	case map[string]interface{}:
		mm := m.(map[string]interface{})
		for key, val := range mm {
			if key == pluralKey {
				// scenarios:
				// 1) key: object with a single key: object -> key: array
				// 2) key: object with a single key: array  -> key: array
				// 3) key: object with multiple keys: object -> key -> the key with an array -> array
				switch val.(type) {
				case map[string]interface{}:
					vval := val.(map[string]interface{})
					if len(vval) == 1 && !isArrayCanBeInObject { // scenario 1 or 2
						var newVal []interface{}
						for _, vvval := range vval {
							switch vvval.(type) {
							case map[string]interface{}: // 1
								newVal = append(newVal, vvval)
							case []interface{}: // 2
								newVal = vvval.([]interface{})
							}
						}
						mm[key] = newVal
					} else if isArrayCanBeInObject {
						if vvvval, ok := vval[subKey]; ok {
							switch vvvval.(type) {
							case map[string]interface{}: // 3
								mmm := mm[key].(map[string]interface{})
								mmm[subKey] = []interface{}{vvvval}
							}
						}

					}
				}
			} else {
				fixArrays(val, originalPluralKey)
			}
		}
	case []interface{}:
		mm := m.([]interface{})
		for _, val := range mm {
			fixArrays(val, originalPluralKey)
		}
	}
}
