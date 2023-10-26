package main

func ParseQuery(b []byte) (version []byte, command []byte, value []byte) {
	version = b[0:4]
	command = b[4:12]
	value = b[12:]

	return version, command, value
}

func ParseSet(b []byte) (key, value []byte) {
	return key, value
}
