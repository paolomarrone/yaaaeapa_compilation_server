# yaaaeapa compilation server

A web server listening for .c and .h files to be compiled as a shared library representing a yaaaeapa audio plugin.

Page `https://address:10002/uploadfiles` accepts requests whose body contains the .c and .h files in stringified JSON form like:

```json
[
	{
		"name": "filename",
		"str": "filecontent"
	},
	{
		"name": "filename",
		"str": "filecontent"
	}
]
```
The header's req must contain the "Target-Arch" info. Supported formats are "x86_64" and "arm64".

See test.js for an example


### Requirements: 
- Must be run on a x86_64 linux machine (for now)
- gcc
- aarch64-linux-gnu-gcc

### Execution
```bash
go run main.go
```

### Behaviour

- By default, it listens on 10002 port
- It saves the files in a temporary directory and compiles them as a shared library
- It Sends the shared library as http response
- It deletes the files