# resgate - A RES gateway

*Web APIs with real-time resources in a simple way.*

A [Go](http://golang.org) project implementing a Real-time API gateway for the [RES protocol](docs/res-protocol.md) with [NATS server](https://nats.io/about/) as messaging system.

Used for building *scaleable*, *resilient*, *extensible*, and *secure* client web APIs based on *simple*, *stateless* micro-services serving *realtime* resources to web applications.

Simple, stateless, and scalable like REST, but updated in real-time.

## Documentation

For more in depth information on the protocol:

* [RES protocol](docs/res-protocol.md) - Entry point to get all the details
* [RES-Service protocol](docs/res-service-protocol.md) - How to write services
* [RES-Client protocol](docs/res-client-protocol.md) - How to write client libraries, if [ResClient](https://github.com/jirenius/resclient) doesn't fit your needs

## Quickstart

If you just want to start using resgate, and you have:
* [installed Go](https://golang.org/doc/install) and [set your `$GOPATH`](https://golang.org/cmd/go/#hdr-GOPATH_environment_variable)
* Added `$GOPATH/bin` (where your binaries ends up) to your `PATH`
* [installed NATS server](https://nats.io/download/nats-io/gnatsd/) and have it running.

Install and run resgate:

```
go get github.com/jirenius/resgate
resgate
```

## Writing a service
Because of the simplicity of the [RES-Service protocol](docs/res-service-protocol.md), a single-threaded service can be created without the need of a library, as with the [Hello world example](#hello-world-example) below. NATS supports [many other languages](https://nats.io/download/) as well.

For writing multi-threaded services in Go:

* Go - https://github.com/jirenius/go-res

*Have you written a library? Send a link to have it added to the README.*


## Hello world example

A simple example of a service exposing a single resource, `exampleService.myModel`, and client application that accesses the resource.  
For a more extensive example, see the [Resgate Test App](https://github.com/jirenius/resgate-test-app).

### Service (Node.js)

Also available as a [Go service](docs/go-service/main.go).

Create an empty folder and install the *nats* client:

```
npm install nats
```

Create file `service.js` :

```javascript
const nats = require('nats').connect("nats://localhost:4222");

let myModel = { message: "Hello world" };

// Access listener. Everyone gets read access and access to call the set-method
nats.subscribe('access.exampleService.myModel', (request, replyTo, subject) => {
	nats.publish(replyTo, JSON.stringify({ result: { get: true, call: "set" }}));
});

// Get listener. Reply with the json encoded model
nats.subscribe('get.exampleService.myModel', (request, replyTo, subject) => {
	nats.publish(replyTo, JSON.stringify({ result: { model: myModel }}));
});

// Set listener for updating the myModel.message property
nats.subscribe('call.exampleService.myModel.set', (request, replyTo, subject) => {
	let req = JSON.parse(request);
	let p = req.params || {};
	// Check if the message property was changed
	if (typeof p.message === 'string' && p.message !== myModel.message) {
		myModel.message = p.message;
		// The model is updated. Send a change event.
		nats.publish('event.exampleService.myModel.change', JSON.stringify({ message: p.message }));
	}
	// Reply success by sending an empty result
	nats.publish(replyTo, JSON.stringify({result: null}));
});

// System resets tells resgate that the service has been (re)started.
// Resgate will then update any cached resource from exampleService
nats.publish('system.reset', JSON.stringify({ resources: [ 'exampleService.>' ]}));

```

Start the service:

```
node service.js
```

### Client

**Using Chrome**  
Copy the javascript below to [esnextb.in](https://esnextb.in/) and try it out from there (make sure you have the [latest resclient version](https://www.npmjs.org/package/resclient) under *Package*).  
Or just try it out in this [CodePen](https://codepen.io/sjirenius/pen/vraZPZ).  

**Using some other browser**  
Some browsers won't allow accessing a non-encrypted websocket from an encrypted page. You can get around that by running the script locally using a webpack server, or some other similar tool.

Try running it in two separate tabs!

```javascript
import ResClient from 'resclient';

let client = new ResClient('ws://localhost:8080');

// Get the model from the service.
client.get('exampleService.myModel').then(model => {
	// Create an input element
	let input = document.createElement('input');
	input.value = model.message;
	document.body.appendChild(input);

	// Call set to update the remote model
	input.addEventListener('input', () => {
		model.set({ message: input.value });
	});

	// Listen for model change events.
	// The model will be unsubscribed after calling model.off
	model.on('change', () => {
		input.value = model.message;
	});
});
```

### Web  resource

Resources can be retrieved using ordinary HTTP GET requests:

**GET**  
```
http://localhost:8080/api/exampleService/myModel
```

Methods can be called using HTTP POST requests:

**POST**  
```
http://localhost:8080/api/exampleService/myModel/set
```
*Body*  
```
{ "message": "Updated through HTTP" }
```

## Usage
```
resgate [options]
```
| Option | Description |
|---|---|
| `-n, --nats <url>` | NATS Server URL |
| `-p, --port <port>` | Use port for clients |
| `-w, --wspath <path>` | Path to websocket |
| `-a, --apipath <path>` | Path to webresources |
| `-r, --reqtimeout <seconds>` | Timeout duration for NATS requests |
| `-u, --headauth <method>` | Resource method for header authentication |
| `    --tls` | Enable TLS |
| `    --tlscert <file>` | Server certificate file |
| `    --tlskey <file>` | Private key for server certificate |
| `-c, --config <file>` | Configuration file |
| `-h, --help` | Show usage message |


## Configuration
Configuration is a JSON encoded file. If no config file is found at the given path, a new file will be created with default values as follows.

### Properties

```javascript
{
	// URL to the NATS server
	"natsUrl": "nats://127.0.0.1:4222",
	// Timeout in seconds for NATS requests
	"requestTimeout": 5,
	// Port for the http server to listen on.
	// If the port value is missing or 0, standard http(s) port is used.
	"port": 8080,
	// Path for accessing the RES API websocket
	"wsPath": "/",
	// Path for accessing web resources
	"apiPath": "/api",
	// Header authentication resource method for web resources.
	// Prior to accessing the resource, this resource method will be
	// called, allowing an auth service to set a token using
	// information such as the request headers.
	// Missing value or null will disable header authentication.
	// Eg. "authService.headerLogin"
	"headerAuth": null,
	// Flag telling if tls encryption is enabled
	"tls": false,
	// Certificate file path for tls encryption
	"tlsCert": "",
	// Key file path for tls encryption
	"tlsKey": ""
}
```

## Running resgate

By design, resgate will exit if it fails to connect to the NATS server, or if it loses the connection.
This is to allow clients to try to reconnect to another resgate instance and resume from there, and to give resgate a fresh new start if something went wrong.

A simple bash script can keep it running:

```bash
#!/bin/bash
until ./resgate; do
    echo "resgate exited with code $?.  Restarting.." >&2
    sleep 2
done
```

## Contributing

The RES protocol and resgate is still under development, and is currently at a state where the protocol has settled, but the gateway has yet to be properly tested.

While it may be used in non-critical environments, it is not yet recommended to use the gateway for any critical systems. Any feedback on the protocol and its implementation is highly appreciated!

If you find any issues with the protocol or the gateway, feel free to [report them](https://github.com/jirenius/resgate/issues/new) as an Issue.
