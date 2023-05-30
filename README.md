# yaaaeapa compilation server

A web server listening for .c and .h files to be compiled as a shared library representing a yaaaeapa audio plugin.

Files must be passed as a stringified JSON object in the form:
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

- It listens on 10002 port
- It saves the files in a temporary directory and compiles them as a shared library
- It Sends the shared library as http response