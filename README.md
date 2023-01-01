cache-engine-lite
====
Data sharing cloud platform.

- [Features](#features)
- [Parameters](#parameters)
- [Generate self-signed certificates](#generate-self-signed-certificates)
- [Deploying with Docker](#deploying-with-docker)
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

## Getting started

	XAUTHTOKEN=35A6E ./cache-engine-lite

>The authorization token must be passed as environment variable. To generate the tokens see [How can I generate strong keys or tokens?](#how-can-i-generate-strong-keys-or-tokens).

>Using under a reverse proxy needs an instance id passed as environment variable. The instance id can be a random key or a simple number and it's important to be preceded by the character `/` otherwise will not work. This id should be used further as route prefix for all endpoints.

	XAUTHTOKEN=35A6E INSTANCEID=/0 ./cache-engine-lite

**Parameters**

Name|Description
:---|:---
`-conf`|Cache file name, default value `cache.json`.
`-http`|Listening address and port for HTTP server, default value `8080`.
`-https`|Listening address and port for HTTPS server, default value `8090`.

## Generate self-signed certificates

	openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout server.key -out server.crt

>Examples from this document are made with `curl`, do not forget to use it with `-k` option for HTTPS connections with self-signed certificates.
 
## Deploying with Docker

	docker build -t cache-engine-lite .
	docker run -dit --restart unless-stopped --name cache-engine-lite -p 8080:8080 cache-engine-lite
	docker logs cache-engine-lite

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

<sup>2</sup> Open APIs need cross-origin resource sharing. To disable CORS on individual requests (eg. keep the traffic low for IoT devices) send `Cors-Control: no-cors` header.

**Return codes**

Code|Description
:---|:---
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

**Endpoint**|`POST /admin/keys`
:---|:---
**Request headers**|`X-Auth-Token: [TOKEN]`
**Request parameters**|`key=[KEY]`
**Response headers**|`HTTP/1.1 [code]`<br> `Access-Control-Allow-Headers: X-Auth-Token, Cache-Control`<br> `Access-Control-Allow-Methods: GET, POST, PATCH, PUT, DELETE, OPTIONS`<br> `Access-Control-Allow-Origin: *`<br> `Access-Control-Expose-Headers: *`<br> `Location: /keys/{key}`<br> `Date: [timeformat]`<br> `Content-Length: 0`
**Response payload**|-
**Return codes**|`201 Created`<br> `401 Unauthorized`<br> `409 Conflict`

**Sample request**

	curl -X POST -i "localhost:8080/admin/keys" -H "X-Auth-Token: 35A6E" -d "key=37D4B"

**Sample response**

	HTTP/1.1 201 Created
	Access-Control-Allow-Headers: X-Auth-Token, Cache-Control
	Access-Control-Allow-Methods: GET, POST, PATCH, PUT, DELETE, OPTIONS
	Access-Control-Allow-Origin: *
	Access-Control-Expose-Headers: *
	Location: /keys/37D4B
	Date: Fri, 24 Jan 2020 09:22:18 GMT
	Content-Length: 0

>To generate the keys see [How can I generate strong keys or tokens?](#how-can-i-generate-strong-keys-or-tokens).

>Remember to save structure after this operation otherwise will be lost at restart (see [Save the cache to file](#save-the-cache-to-file)).

<div style="page-break-after: always;"></div>

### Delete a key

**Endpoint**|`DELETE /admin/keys/{key}`
:---|:---
**Request headers**|`X-Auth-Token: [TOKEN]`
**Request parameters**|-
**Response headers**|`HTTP/1.1 [code]`<br> `Access-Control-Allow-Headers: X-Auth-Token, Cache-Control`<br> `Access-Control-Allow-Methods: GET, POST, PATCH, PUT, DELETE, OPTIONS`<br> `Access-Control-Allow-Origin: *`<br> `Access-Control-Expose-Headers: *`<br> `Date: [timeformat]`<br> `Content-Length: 0`
**Response payload**|-
**Return codes**|`200 OK`<br> `401 Unauthorized`<br> `404 Not Found`

**Sample request:**

	curl -X DELETE -i "localhost:8080/admin/keys/37D4B" -H "X-Auth-Token: 35A6E"

**Sample response:**

	HTTP/1.1 200 OK
	Access-Control-Allow-Headers: X-Auth-Token, Cache-Control
	Access-Control-Allow-Methods: GET, POST, PATCH, PUT, DELETE, OPTIONS
	Access-Control-Allow-Origin: *
	Access-Control-Expose-Headers: *
	Date: Fri, 24 Jan 2020 09:23:57 GMT
	Content-Length: 0

>Remember to save structure after this operation otherwise will be lost at restart (see [Save the cache to file](#save-the-cache-to-file)).

<div style="page-break-after: always;"></div>

### Show entire cache

**Endpoint**|`GET /admin/keys`
:---|:---
**Request headers**|`X-Auth-Token: [TOKEN]`
**Request parameters**|-
**Response headers**|`HTTP/1.1 [code]`<br> `Access-Control-Allow-Headers: X-Auth-Token, Cache-Control`<br> `Access-Control-Allow-Methods: GET, POST, PATCH, PUT, DELETE, OPTIONS`<br> `Access-Control-Allow-Origin: *`<br> `Access-Control-Expose-Headers: *`<br> `Content-Size: [size]`<br> `Content-Type: application/json`<br> `Date: [timeformat]`<br> `Content-Length: [length]`
**Response payload**|`JSON RFC4627`
**Return codes**|`200 OK`<br> `401 Unauthorized`<br> `500 Internal Server Error`

**Sample request:**

	curl -X GET -i "localhost:8080/admin/keys" -H "X-Auth-Token: 35A6E"

**Sample response:**

	HTTP/1.1 200 OK
	Access-Control-Allow-Headers: X-Auth-Token, Cache-Control
	Access-Control-Allow-Methods: GET, POST, PATCH, PUT, DELETE, OPTIONS
	Access-Control-Allow-Origin: *
	Access-Control-Expose-Headers: *
	Content-Size: 1
	Content-Type: application/json
	Date: Fri, 24 Jan 2020 09:25:39 GMT
	Content-Length: 85
	
	{
	    "37D4B": {
	        "Tem": [
	            1,
	            2
	        ]
	    }
	}

`Content-Size` is a custom header meaning the number of items in the cache.

<div style="page-break-after: always;"></div>

### Save the cache to file

**Endpoint**|`PUT /admin/keys`
:---|:---
**Request headers**|`X-Auth-Token: [TOKEN]`
**Request parameters**|-
**Response headers**|`HTTP/1.1 [code]`<br> `Access-Control-Allow-Headers: X-Auth-Token, Cache-Control`<br> `Access-Control-Allow-Methods: GET, POST, PATCH, PUT, DELETE, OPTIONS`<br> `Access-Control-Allow-Origin: *`<br> `Access-Control-Expose-Headers: *`<br> `Date: [timeformat]`<br> `Content-Length: 0`
**Response payload**|-
**Return codes**|`200 OK`<br> `401 Unauthorized`<br> `500 Internal Server Error`<br> `507 Insufficient Storage`

**Sample request:**

	curl -X PUT -i "localhost:8080/admin/keys" -H "X-Auth-Token: 35A6E"

**Sample response:**

	HTTP/1.1 200 OK
	Access-Control-Allow-Headers: X-Auth-Token, Cache-Control
	Access-Control-Allow-Methods: GET, POST, PATCH, PUT, DELETE, OPTIONS
	Access-Control-Allow-Origin: *
	Access-Control-Expose-Headers: *
	Date: Fri, 24 Jan 2020 09:27:16 GMT
	Content-Length: 0

>Default cache file is `cache.json`. It's no need to take care about this file because it will be created if not exists. At start, the engine will load the cache from this file.

>If the disk is out of space `507 Insufficient Storage` is returned.

<div style="page-break-after: always;"></div>

### Get a key

**Endpoint**|`GET /keys/{key}`
:---|:---
**Request headers**|`Cache-Control: no-store`
**Request parameters**|-
**Response headers**|`HTTP/1.1 [code]`<br> `Access-Control-Allow-Headers: X-Auth-Token, Cache-Control`<br> `Access-Control-Allow-Methods: GET, POST, PATCH, PUT, DELETE, OPTIONS`<br> `Access-Control-Allow-Origin: *`<br> `Access-Control-Expose-Headers: *`<br> `Content-Type: application/json`<br> `Date: [timeformat]`<br> `Content-Length: [length]`
**Response payload**|`JSON RFC4627`
**Return codes**|`200 OK`<br> `204 No Content`<br> `404 Not Found`<br> `500 Internal Server Error`

**Sample request:**

	curl -X GET -i "localhost:8080/keys/37D4B" -H "Cache-Control: no-store"

**Sample response:**

	HTTP/1.1 200 OK
	Access-Control-Allow-Headers: X-Auth-Token, Cache-Control
	Access-Control-Allow-Methods: GET, POST, PATCH, PUT, DELETE, OPTIONS
	Access-Control-Allow-Origin: *
	Access-Control-Expose-Headers: *
	Content-Type: application/json
	Date: Fri, 24 Jan 2020 09:29:16 GMT
	Content-Length: 48
	
	{
	    "Tem": [
	        1,
	        2
	    ]
	}

>Timestamp is available only for items updated with [query connector](#update-a-key-query-connector) as `_time_` in [RFC3339](https://golang.org/src/time/format.go?s=3825:3867#L62) or [RFC3339Nano](https://golang.org/src/time/format.go?s=3868:3919#L62) format.

>Using `Cache-Control: no-store` will retrieve and delete the item.

<div style="page-break-after: always;"></div>

### Get a group of keys

**Endpoint**|`GET /keys?key=[KEY]`
:---|:---
**Request headers**|-
**Request parameters**|`Query RFC3986`
**Response headers**|`HTTP/1.1 [code]`<br> `Access-Control-Allow-Headers: X-Auth-Token, Cache-Control`<br> `Access-Control-Allow-Methods: GET, POST, PATCH, PUT, DELETE, OPTIONS`<br> `Access-Control-Allow-Origin: *`<br> `Access-Control-Expose-Headers: *`<br> `Content-Size: [size]`<br> `Content-Type: application/json`<br> `Date: [timeformat]`<br> `Content-Length: [length]`
**Response payload**|`JSON RFC4627`
**Return codes**|`200 OK`<br> `204 No Content`<br> `500 Internal Server Error`

**Sample request:**

	curl -X GET -i "localhost:8080/keys?key=37D4B&key=37D4C&key=NONEXISTENTKEY"

**Sample response:**

	HTTP/1.1 200 OK
	Access-Control-Allow-Headers: X-Auth-Token, Cache-Control
	Access-Control-Allow-Methods: GET, POST, PATCH, PUT, DELETE, OPTIONS
	Access-Control-Allow-Origin: *
	Access-Control-Expose-Headers: *
	Content-Size: 3
	Content-Type: application/json
	Date: Fri, 24 Jan 2020 09:32:39 GMT
	Content-Length: 159
	
	[
	    {
	        "Tem": [
	            1,
	            2
	        ]
	    },
	    {
	        "Tem": [
	            1,
	            3
	        ]
	    },
	    null
	]

`Content-Size` is a custom header meaning the number of items in the cache.

>Nonexistent or empty keys appears as `null` in the corresponding position in array.

<div style="page-break-after: always;"></div>

### Update a key

**Endpoint**|`PUT\|POST /keys/{key}`
:---|:---
**Request headers**|-
**Request parameters**|`JSON RFC4627`
**Response headers**|`HTTP/1.1 [code]`<br> `Access-Control-Allow-Headers: X-Auth-Token, Cache-Control`<br> `Access-Control-Allow-Methods: GET, POST, PATCH, PUT, DELETE, OPTIONS`<br> `Access-Control-Allow-Origin: *`<br> `Access-Control-Expose-Headers: *`<br> `Date: [timeformat]`<br> `Content-Length: 0`
**Response payload**|-
**Return codes**|`200 OK`<br> `400 Bad Request`<br> `404 Not Found`

**Sample request:**

	curl -X PUT -i "localhost:8080/keys/37D4B"  -d '{"Tem": [ 1, 2 ]}'

**Sample response:**

	HTTP/1.1 200 OK
	Access-Control-Allow-Headers: X-Auth-Token, Cache-Control
	Access-Control-Allow-Methods: GET, POST, PATCH, PUT, DELETE, OPTIONS
	Access-Control-Allow-Origin: *
	Access-Control-Expose-Headers: *
	Date: Fri, 24 Jan 2020 09:25:27 GMT
	Content-Length: 0

>Other party can't update a key until an administrator create one for him (see [Create a new key](#create-a-new-key)).

>To create persistent connections with the engine use `POST` method instead `PUT`. If you do not want a persistent connection use `Connection: close` header in your request. However, this header is mandatary with `PUT` method to avoid receiving `400 Bad Request` error code.

<div style="page-break-after: always;"></div>

### Update a key (query connector)

**Endpoint**|`GET /update?key=[KEY]&name=value`
:---|:---
**Request headers**|-
**Request parameters**|`Query RFC3986`
**Response headers**|`HTTP/1.1 [code]`<br> `Access-Control-Allow-Headers: X-Auth-Token, Cache-Control`<br> `Access-Control-Allow-Methods: GET, POST, PATCH, PUT, DELETE, OPTIONS`<br> `Access-Control-Allow-Origin: *`<br> `Access-Control-Expose-Headers: *`<br> `Date: [timeformat]`<br> `Content-Length: 0`
**Response payload**|-
**Return codes**|`200 OK`<br> `404 Not Found`

**Sample request:**

    curl -X GET -i "localhost:8080/update?key=37D4B&Tem=1&Tem=2"

**Sample response:**

	HTTP/1.1 200 OK
	Access-Control-Allow-Headers: X-Auth-Token, Cache-Control
	Access-Control-Allow-Methods: GET, POST, PATCH, PUT, DELETE, OPTIONS
	Access-Control-Allow-Origin: *
	Access-Control-Expose-Headers: *
	Date: Fri, 24 Jan 2020 09:35:16 GMT
	Content-Length: 0

>Other party can't update a key until an administrator create one for him (see [Create a new key](#create-a-new-key)).

>UTC timestamp is automatically added to the query as `_time_` in [RFC3339](https://golang.org/src/time/format.go?s=3825:3867#L62) or [RFC3339Nano](https://golang.org/src/time/format.go?s=3868:3919#L62) format.

>Query string fields are converted to values according to the data type of each field.

<div style="page-break-after: always;"></div>

### Get version

**Endpoint**|`GET /version`
:---|:---
**Request headers**|-
**Request parameters**|-
**Response headers**|`HTTP/1.1 200 OK`<br> `Access-Control-Allow-Headers: X-Auth-Token, Cache-Control`<br> `Access-Control-Allow-Methods: GET, POST, PATCH, PUT, DELETE, OPTIONS`<br> `Access-Control-Allow-Origin: *`<br> `Access-Control-Expose-Headers: *`<br> `Content-Type: text/plain`<br> `Server: go[version] (GOOS)`<br> `Date: [timeformat]`<br> `Content-Length: [length]`
**Response payload**|`text`
**Return codes**|`200 OK`

**Sample request:**

	curl -X GET -i "localhost:8080/version"

**Sample response:**

	HTTP/1.1 200 OK
	Access-Control-Allow-Headers: X-Auth-Token, Cache-Control
	Access-Control-Allow-Methods: GET, POST, PATCH, PUT, DELETE, OPTIONS
	Access-Control-Allow-Origin: *
	Access-Control-Expose-Headers: *
	Content-Type: text/plain
	Server: go1.13 (linux)
	Date: Fri, 24 Jan 2020 09:13:11 GMT
	Content-Length: 14
	
	1.4.3-20200124

>Version string respect semantic versioning (see [semver.org](https://semver.org)).

<div style="page-break-after: always;"></div>

### File server

**Endpoint**|`GET /static/`
:---|:---
**Request headers**|-
**Request parameters**|-
**Response headers**|`HTTP/1.1 200 OK`<br> `Access-Control-Allow-Headers: X-Auth-Token, Cache-Control`<br> `Access-Control-Allow-Methods: GET, POST, PATCH, PUT, DELETE, OPTIONS`<br> `Access-Control-Allow-Origin: *`<br> `Access-Control-Expose-Headers: *`<br> `Content-Type: text/html; charset=utf-8`<br> `Last-Modified: [timeformat]`<br> `Date: [timeformat]`<br> `Content-Length: [length]`
**Response payload**|`text/html`
**Return codes**|`200 OK`<br> `404 Not Found`

**Sample request:**

	curl -X GET -i "localhost:8080/static/"

**Sample response:**

	HTTP/1.1 200 OK
	Access-Control-Allow-Headers: X-Auth-Token, Cache-Control
	Access-Control-Allow-Methods: GET, POST, PATCH, PUT, DELETE, OPTIONS
	Access-Control-Allow-Origin: *
	Access-Control-Expose-Headers: *
	Content-Type: text/html; charset=utf-8
	Last-Modified: Fri, 24 Jan 2020 09:18:42 GMT
	Date: Fri, 24 Jan 2020 09:18:47 GMT
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

Recommended way is to do this from `crontab` at reboot.

	crontab -e
	@reboot XAUTHTOKEN=35A6E $HOME/cache-engine-lite/cache-engine-lite > cache-engine-lite.log 2>&1
	
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

### Can I generate datetime into the cache?

	nano $HOME/cache-engine:datetime
	#!/bin/bash	
	while true; do
	  curl -X PUT "localhost:8080/keys/datetime" -d "{ \"unixtime\": `date '+%s'` }"
	  sleep 1;
	done
	chmod +x cache-engine:datetime
	crontab -e	
	@reboot $HOME/cache-engine:datetime

## Resources

[RESTful API Design. Best Practices in a Nutshell.](https://phauer.com/2015/restful-api-design-best-practices/)<br>
[List of HTTP header fields](https://en.wikipedia.org/wiki/List_of_HTTP_header_fields#Standard_request_fields)
