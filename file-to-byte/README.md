# File To Byte Converter

Reads a given file and converts it to the corresponding byte array that can be used directly in Go code.


# Usage

```
file-to-byte -f /somewhere/over/the/rainbow/file.out
```

Copy the generated output and use it directly in your Go code

# Example output
```
var filedata = []byte{45, 45, 45, 45, 45, 66, 69, 71, 73, 78, 32, 67, 69, 82, 84, 73, 70, 73, 67, 65, 84, 69, 45, 45, 45, 45, 45, 10, 77, 73, 73, 68, 113, 84, 67, 67, 65, 90, 112, 68, 13, 10, 77, 102, 65, 87, 82, 43, 53, 79, 101, 81, 105, 78, 65, 112, 47, 98, 71, 52, 102, 106, 74, 111, 84, 100, 111, 113, 107, 117, 108, 53, 49, 43, 50, 98, 72, 72, 86, 114, 85, 61, 10, 45, 45, 45, 45, 45, 69, 78, 68, 32, 67, 69, 82, 84, 73, 70, 73, 67, 65, 84, 69, 45, 45, 45, 45, 45, 10}
```

# Possible uses
1. Storing resource file content within Go binary for use during execution time.

# Credits
https://stackoverflow.com/a/28071360/3296350