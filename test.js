#!/bin/node

(async function() {

	const https = require('https')
	
	const yaaaeapaToSharedLibrary = async function (files, compilationServerUrl, compilationServerPort, arch, onSuccessCb, onFailureCb) {

		const postData = JSON.stringify(files);
		
		const req = https.request(
			{
				hostname: compilationServerUrl,
				port: compilationServerPort,
				path: '/uploadfiles',
				method: 'POST',
				
				headers: {
					'Content-Type': 'application/json',
					'Content-Length': Buffer.byteLength(postData),
					'Target-arch': arch
				}
			}, (res) => {
				var bufs = [];
				res.on('data', (chunk) => {
					bufs.push(chunk);
				});
				res.on('end', () => {
					if (res.statusCode == 500) {
						onFailureCb("Compilation Server Error:\n" + res.statusCode + "\nin response: \n" + bufs.join(''));
						return;
					}
					if (res.headers["compilation-result"] != "ok") {
						onFailureCb("Compilation Server Error:\n" + res.statusCode + "\n" + res.headers["Compilation-log"]);
						return;
					}
					var buf = Buffer.concat(bufs);
					onSuccessCb(buf);
				});
			}
		);

		req.on('error', (e) => {
			onFailureCb("Compilation Server Error:\nProblem with request: " + e.message);
		});

		req.write(postData);
		req.end(); 
	}


	// Test
	
	const compilationServerUrl = "localhost";
	const compilationServerPort = 10002;

	const files = [
		{
			name: "test.h",
			str:  "\
				void 	yaaaeapa_init (void); \n \
				void  	yaaaeapa_fini (void); \n \
			"
		},
		{
			name: "test.c",
			str:  "\
				#include \"test.h\" \n \
				void 	yaaaeapa_init (void) {} \n \
				void  	yaaaeapa_fini (void) {} \n \
			"
		}
	];

	const arch = "arm64" // arm64, x86_64

	const onSuccessCb = (m) => { console.log("Success", m) };
	const onFailureCb = (m) => { console.log("Failure", m) };

	try {
		yaaaeapaToSharedLibrary(files, compilationServerUrl, compilationServerPort, arch, onSuccessCb, onFailureCb);
	} catch (e) {
		console.log("Error", e)
	}

}());
