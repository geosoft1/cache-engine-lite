cache-engine-lite</sub>
====
Data sharing cloud platform.

- [Features](#features)
- [Parameters](#parameters)
- [Generate self-signed certificates](#generate-self-signed-certificates)
- [REST API](#rest-api)
  - [Create a new key](#create-a-new-key)
  - [Delete a key](#delete-a-key)
  - [Show entire cache](#show-entire-cache)
  - [Save the cache to file](#save-the-cache-to-file)
  - [Get a key](#get-a-key)
  - [Get a group of keys](#get-a-group-of-keys)
  - [Update a key](#update-a-key)
  - [Update a key (query connector)](#update-a-key-query-connector)
  - [Get version](#get-version)
  - [File server](#file-server)
- [Known issues](#known-issues)
- [FAQ](#faq)
- [Resources](#resources)
- [LICENSE](LICENSE)

<div style="page-break-after: always;"></div>

## Features

- Share data structures over HTTP protocol.
- Query connector for [web of things](https://en.wikipedia.org/wiki/Web_of_Things).
- Remote administration through [RESTful](https://en.wikipedia.org/wiki/Representational_state_transfer) API.
- File server for various resources.

## Parameters

Name|Description
---|---
`-conf`|Cache file name, default value `cache.json`.
`-http`|Listening address and port for HTTP server, default value `8080`.
`-https`|Listening address and port for HTTPS server, default value `8090`.

>Authorization token must be passed as environment variable.

	XAUTHTOKEN=35A6E ./cache-engine-lite

>To generate the tokens see [How can I generate strong keys or tokens?](#how-can-i-generate-strong-keys-or-tokens).

>Using under a reverse proxy needs an instance id passed as environment variable. The instance id can be a random key or a simple number and it's important to be preceded by the character `/` otherwise will not work. This id should be used further as route prefix for all endpoints.

	XAUTHTOKEN=35A6E INSTANCEID=/0 ./cache-engine-lite

## Generate self-signed certificates

	openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout server.key -out server.crt

>Examples from this document are made with `curl`, do not forget to use it with `-k` option for HTTPS connections with self-signed certificates.
 
<div style="page-break-after: always;"></div>

## REST API

**Standards**

- JSON [RFC4627](https://tools.ietf.org/html/rfc4627).
- [UTF-8](https://en.wikipedia.org/wiki/UTF-8) encoding.
- Item timestamp [RFC3339](https://tools.ietf.org/html/rfc3339) <sup>1</sup>.
- HTTP query [RFC3986](https://tools.ietf.org/html/rfc3986).
- HTTP responses time [RFC7231#section-7.1.1.1](https://tools.ietf.org/html/rfc7231#section-7.1.1.1) (see also [TimeFormat](https://golang.org/src/net/http/server.go?s=27987:28038#L903)).
- HTTP return codes [RFC7231#page-49](https://tools.ietf.org/html/rfc7231#page-49) and [RFC4918#section-11.5](https://tools.ietf.org/html/rfc4918#section-11.5).
- [CORS](https://www.w3.org/TR/cors/) <sup>2</sup>.

<sup>1</sup> Depending how is compilated can use [RFC3339](https://golang.org/src/time/format.go?s=3825:3867#L62) or [RFC3339Nano](https://golang.org/src/time/format.go?s=3868:3919#L62) format.

<sup>2</sup> Open APIs need cross-origin resource sharing.

**Return codes**

Code|Description
---|---
`200 OK`                    |Request has succeeded.
`201 Created`               |Key was created.
`204 No Content`            |No content or other party deleted the data.
`400 Bad Request`           |Wrong input JSON.
`401 Unauthorized`          |Wrong input `X-Auth-Token`.
`404 Not Found`             |Wrong input URL.
`409 Conflict`              |Key already exists.
`500 Internal Server Error` |Something has gone wrong on the server.
`507 Insufficient Storage`  |Server disk is out of space.

<div style="page-break-after: always;"></div>

### Create a new key

Resources|Values 
---|---
**Endpoint**|`POST /admin/keys`
**Request headers**|`X-Auth-Token: [TOKEN]`
**Request parameters**|`key=[KEY]`
**Response headers**|`HTTP/1.1 [code]`<br> `Access-Control-Allow-Origin: *`<br> `Location: /keys/{key}`<br> `Date: [timeformat]`<br> `Content-Length: 0`
**Response payload**|-
**Return codes**|`201 Created`<br> `401 Unauthorized`<br> `409 Conflict`

**Sample request**

	curl -X POST -i "localhost:8080/admin/keys" -H "X-Auth-Token: 35A6E" -d "key=37D4B"

**Sample response**

	HTTP/1.1 201 Created
	Access-Control-Allow-Origin: *
	Location: /keys/37D4B
	Date: Mon, 14 Oct 2019 15:44:55 GMT
	Content-Length: 0


>To generate the keys see [How can I generate strong keys or tokens?](#how-can-i-generate-strong-keys-or-tokens).

>Remember to save structure after this operation otherwise will be lost at restart (see [Save the cache to file](#save-the-cache-to-file)).

<div style="page-break-after: always;"></div>

### Delete a key
Resources|Values
---|---
**Endpoint**|`DELETE /admin/keys/{key}`
**Request headers**|`X-Auth-Token: [TOKEN]`
**Request parameters**|-
**Response headers**|`HTTP/1.1 [code]`<br> `Access-Control-Allow-Origin: *`<br> `Date: [timeformat]`<br> `Content-Length: 0`
**Response payload**|-
**Return codes**|`200 OK`<br> `401 Unauthorized`<br> `404 Not Found`

**Sample request:**

	curl -X DELETE -i "localhost:8080/admin/keys/37D4B" -H "X-Auth-Token: 35A6E"

**Sample response:**

	HTTP/1.1 200 OK
	Access-Control-Allow-Origin: *
	Date: Mon, 14 Oct 2019 15:44:55 GMT
	Content-Length: 0

>Remember to save structure after this operation otherwise will be lost at restart (see [Save the cache to file](#save-the-cache-to-file)).

<div style="page-break-after: always;"></div>

### Show entire cache

Resources|Values
---|---
**Endpoint**|`GET /admin/keys`
**Request headers**|`X-Auth-Token: [TOKEN]`
**Request parameters**|-
**Response headers**|`HTTP/1.1 [code]`<br> `Access-Control-Allow-Origin: *`<br> `Content-Size: [size]`<br> `Content-Type: application/json`<br> `Date: [timeformat]`<br> `Content-Length: [length]`
**Response payload**|`JSON RFC4627`
**Return codes**|`200 OK`<br> `401 Unauthorized`<br> `500 Internal Server Error`

**Sample request:**

	curl -X GET -i "localhost:8080/admin/keys" -H "X-Auth-Token: 35A6E"

**Sample response:**

	HTTP/1.1 200 OK
	Access-Control-Allow-Origin: *
	Content-Size: 1
	Content-Type: application/json
	Date: Mon, 14 Oct 2019 15:44:55 GMT
	Content-Length: 85
	
	{
		"37D4B": {
			"Tem": [
				"1",
				"2"
			]
		}
	}

`Content-Size` is a custom header meaning the number of items in the cache.

<div style="page-break-after: always;"></div>

### Save the cache to file

Resources|Values
---|---
**Endpoint**|`PUT /admin/keys`
**Request headers**|`X-Auth-Token: [TOKEN]`
**Request parameters**|-
**Response headers**|`HTTP/1.1 [code]`<br> `Access-Control-Allow-Origin: *`<br> `Date: [timeformat]`<br> `Content-Length: 0`
**Response payload**|-
**Return codes**|`200 OK`<br> `401 Unauthorized`<br> `500 Internal Server Error`<br> `507 Insufficient Storage`

**Sample request:**

	curl -X PUT -i "localhost:8080/admin/keys" -H "X-Auth-Token: 35A6E"

**Sample response:**

	HTTP/1.1 200 OK
	Access-Control-Allow-Origin: *
	Date: Mon, 14 Oct 2019 15:44:55 GMT
	Content-Length: 0

>Default cache file is `cache.json`. It's no need to take care about this file because it will be created if not exists. At start, the engine will load the cache from this file.

>If the disk is out of space `507 Insufficient Storage` is returned.

<div style="page-break-after: always;"></div>

### Get a key

Resources|Values
---|---
**Endpoint**|`GET /keys/{key}`
**Request headers**|`Cache-Control: no-store`
**Request parameters**|-
**Response headers**|`HTTP/1.1 [code]`<br> `Access-Control-Allow-Origin: *`<br> `Content-Type: application/json`<br> `Date: [timeformat]`<br> `Content-Length: [length]`
**Response payload**|`JSON RFC4627`
**Return codes**|`200 OK`<br> `204 No Content`<br> `404 Not Found`<br> `500 Internal Server Error`

**Sample request:**

	curl -X GET -i "localhost:8080/keys/37D4B" -H "Cache-Control: no-store"

**Sample response:**

	HTTP/1.1 200 OK
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Date: Mon, 14 Oct 2019 15:44:55 GMT
	Content-Length: 100
	
	{
		"Tem": [
			"1",
			"2"
		],
		"_time_": "2019-10-14T15:44:00Z"
	}

>Timestamp is available only for items updated with [query connector](#update-a-key-query-connector) as `_time_` in [RFC3339](https://golang.org/src/time/format.go?s=3825:3867#L62) or [RFC3339Nano](https://golang.org/src/time/format.go?s=3868:3919#L62) format.

>Using `Cache-Control: no-store` will retrieve and delete the item.

<div style="page-break-after: always;"></div>

### Get a group of keys

Resources|Values
---|---
**Endpoint**|`GET /keys?key=[KEY]`
**Request headers**|-
**Request parameters**|`Query RFC3986`
**Response headers**|`HTTP/1.1 [code]`<br> `Access-Control-Allow-Origin: *`<br> `Content-Size: [size]`<br> `Content-Type: application/json`<br> `Date: [timeformat]`<br> `Content-Length: [length]`
**Response payload**|`JSON RFC4627`
**Return codes**|`200 OK`<br> `204 No Content`<br> `500 Internal Server Error`

**Sample request:**

	curl -X GET -i "localhost:8080/keys?key=37D4B&key=37D4C&key=NONEXISTENTKEY"

**Sample response:**

	HTTP/1.1 200 OK
	Access-Control-Allow-Origin: *
	Content-Size: 3
	Content-Type: application/json
	Date: Mon, 14 Oct 2019 15:44:55 GMT
	Content-Length: 159
	
	[
		{
			"Tem": [
				"1",
				"2"
			]
		},
		{
			"Tem": [
				"1",
				"3"
			]
		},
		null
	]

`Content-Size` is a custom header meaning the number of items in the cache.

>Nonexistent or empty keys appears as `null` in the corresponding position in array.

<div style="page-break-after: always;"></div>

### Update a key

Resources|Values
---|---
**Endpoint**|`PUT /keys/{key}`
**Request headers**|-
**Request parameters**|`JSON RFC4627`
**Response headers**|`HTTP/1.1 [code]`<br> `Access-Control-Allow-Origin: *`<br> `Date: [timeformat]`<br> `Content-Length: 0`
**Response payload**|-
**Return codes**|`200 OK`<br> `400 Bad Request`<br> `404 Not Found`

**Sample request:**

	curl -X PUT -i "localhost:8080/keys/37D4B"  -d '{"Tem": [ "1", "2" ]}'

**Sample response:**

	HTTP/1.1 200 OK
	Access-Control-Allow-Origin: *
	Date: Mon, 14 Oct 2019 15:44:55 GMT
	Content-Length: 0

>Other party can't update a key until an administrator create one for him (see [Create a new key](#create-a-new-key)).

<div style="page-break-after: always;"></div>

### Update a key (query connector)

Resources|Values
---|---
**Endpoint**|`GET /update?key=[KEY]&name=value`
**Request headers**|-
**Request parameters**|`Query RFC3986`
**Response headers**|`HTTP/1.1 [code]`<br> `Access-Control-Allow-Origin: *`<br> `Date: [timeformat]`<br> `Content-Length: 0`
**Response payload**|-
**Return codes**|`200 OK`<br> `404 Not Found`

**Sample request:**

    curl -X GET -i "localhost:8080/update?key=37D4B&Tem=1&Tem=2"

**Sample response:**

	HTTP/1.1 200 OK
	Access-Control-Allow-Origin: *
	Date: Mon, 14 Oct 2019 15:44:55 GMT
	Content-Length: 0

>Other party can't update a key until an administrator create one for him (see [Create a new key](#create-a-new-key)).

>UTC timestamp is automatically added to the query as `_time_` in [RFC3339](https://golang.org/src/time/format.go?s=3825:3867#L62) or [RFC3339Nano](https://golang.org/src/time/format.go?s=3868:3919#L62) format.

<div style="page-break-after: always;"></div>

### Get version

Resources|Values
---|---
**Endpoint**|`GET /version`
**Request headers**|-
**Request parameters**|-
**Response headers**|`HTTP/1.1 200 OK`<br> `Access-Control-Allow-Origin: *`<br> `Content-Type: text/plain`<br> `Server: go[version] (GOOS)`<br> `Date: [timeformat]`<br> `Content-Length: [length]`
**Response payload**|`text`
**Return codes**|`200 OK`

**Sample request:**

	curl -X GET -i "localhost:8080/version"

**Sample response:**

	HTTP/1.1 200 OK
	Access-Control-Allow-Origin: *
	Content-Type: text/plain
	Server: go1.13 (linux)
	Date: Mon, 14 Oct 2019 15:44:55 GMT
	Content-Length: 15
	
	1.0.0+20191014

>Version string respect semantic versioning (see [semver.org](https://semver.org)).

<div style="page-break-after: always;"></div>

### File server

Resources|Values
---|---
**Endpoint**|`GET /static/`
**Request headers**|-
**Request parameters**|-
**Response headers**|`HTTP/1.1 200 OK`<br> `Access-Control-Allow-Origin: *`<br> `Content-Type: text/html; charset=utf-8`<br> `Last-Modified: [timeformat]`<br> `Date: [timeformat]`<br> `Content-Length: [length]`
**Response payload**|`text/html`
**Return codes**|`200 OK`<br> `404 Not Found`

**Sample request:**

	curl -X GET -i "localhost:8080/static/"

**Sample response:**

	HTTP/1.1 200 OK
	Access-Control-Allow-Origin: *
	Content-Type: text/html; charset=utf-8
	Last-Modified: Mon, 14 Oct 2019 15:44:55 GMT
	Date: Mon, 14 Oct 2019 15:44:55 GMT
	Content-Length: 13
	
	<pre>
	</pre>

>By default the file server expose the `static` folder content to the world, don't forget to create an `index.html` file inside to secure the folder.

<div style="page-break-after: always;"></div>

## Known issues

`XAUTHTOKEN` is passed as environment variable for now.

## Faq

### How can I generate strong keys or tokens?

Some free online keys generators are available like [randomkeygen.com](https://randomkeygen.com).

### How can I be sure I push a right JSON in the cache?

Simply check with a free JSON formatter like [jsonformatter.curiousconcept.com](https://jsonformatter.curiousconcept.com/).

### How can I run from automatically?

Recommended way is to do thid from `crontab` at reboot.

	crontab -e
	@reboot XAUTHTOKEN=35A6E $HOME/src/cache-engine-lite/cache-engine-lite > cache-engine-lite.log 2>&1
	
### Can I save the cache periodically?

Yes. The easiest way is from `crontab` (eg. at every 5 minutes).

	crontab -e
	*/5 * * * * curl -X PUT -i "localhost:8080/admin/keys" -H "X-Auth-Token: 35A6E"

### How can I replicate the cache on other machine?

	curl -X GET "cache-engine-ip:port/admin/keys" -H "X-Auth-Token: 35A6E" > $HOME/cache.json
	
>You can do this manually or periodically from crontab.

### How can I get the number of items in the cache?

	curl -X GET -sI "localhost:8080/admin/keys" -H "x-auth-token: 35A6E" | grep Content-Size

### Can I store metadata that describes other keys?

Yes, simply create a key with data about other keys. Your particular application will read this key to locate and read informations about the other keys.

	"37D4A": {
		"37D4B": "Servers room",
		"37D4C": "Alice's restaurant"
	}

### How can I keep my data secret?

A simple method is encoding your data with a password.

	curl -X PUT -i "localhost:8080/keys/37D4B" -d '{"enc": "'`echo secret=thing | 
	openssl enc -aes-128-cbc -a -salt -pass pass:1234`'"}'

Now, unauthorized eyes will see only the encrypted string.

	curl -X GET "localhost:8080/keys/37D4B"
	{
		"enc": "U2FsdGVkX1++SyexS7zE3yV6fiMd2Fsr/Ao7qdBWhzE="
	}

The other party must decode the data with the same password.

	curl -s -X GET "localhost:8080/keys/37D4B" | 
	grep enc | sed 's/"//g' | awk '{print $2}' | openssl enc -aes-128-cbc -a -d -salt -pass pass:1234
	secret=thing

## Resources

[RESTful API Design. Best Practices in a Nutshell.](https://phauer.com/2015/restful-api-design-best-practices/)<br>
[List of HTTP header fields](https://en.wikipedia.org/wiki/List_of_HTTP_header_fields#Standard_request_fields)
