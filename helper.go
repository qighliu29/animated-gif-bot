package main

// func readChanBreak(c <-chan interface{}, cb func(interface{}) bool) error {
// 	for t := range c {
// 		switch t.(type) {
// 		case error:
// 			return t.(error)
// 		default:
// 			if !cb(t) {
// 				return nil
// 			}
// 		}
// 	}
// 	return nil
// }

func readChanUntilClose(c <-chan interface{}, cb func(interface{})) error {
	for t := range c {
		switch t.(type) {
		case error:
			return t.(error)
		default:
			cb(t)
		}
	}
	return nil
}
