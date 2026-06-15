# To-Do API — curl responses

Live HTTP responses from a local run of this service (commit this file to the repo as API documentation evidence).

To reproduce: from the repo root run `PORT=18080 go run ./cmd/server`, then execute the `curl` commands below (or re-run your own capture script). Default dev port is `8080`; this capture used **18080** to avoid conflicts.

Base URL: `http://127.0.0.1:18080`
Captured (UTC): 2026-06-15 18:06:15

---

## 1. POST /todos — create first item

```http
$ curl -i -sS -X POST http://127.0.0.1:18080/todos \
  -H 'Content-Type: application/json' \
  -d '{"task":"First task","due_date":"2026-06-20T18:00:00Z"}'
```

```
HTTP/1.1 201 Created
Content-Type: application/json
Date: Mon, 15 Jun 2026 18:06:15 GMT
Content-Length: 81

{"id":1,"task":"First task","due_date":"2026-06-20T18:00:00Z","completed":false}
```

## 2. POST /todos — create second item (earlier due date)

```http
$ curl -i -sS -X POST http://127.0.0.1:18080/todos \
  -H 'Content-Type: application/json' \
  -d '{"task":"Second task","due_date":"2026-06-18T10:00:00Z"}'
```

```
HTTP/1.1 201 Created
Content-Type: application/json
Date: Mon, 15 Jun 2026 18:06:15 GMT
Content-Length: 82

{"id":2,"task":"Second task","due_date":"2026-06-18T10:00:00Z","completed":false}
```

## 3. GET /todos/{id} — get by id (id=1)

```http
$ curl -i -sS http://127.0.0.1:18080/todos/1
```

```
HTTP/1.1 200 OK
Content-Type: application/json
Date: Mon, 15 Jun 2026 18:06:15 GMT
Content-Length: 81

{"id":1,"task":"First task","due_date":"2026-06-20T18:00:00Z","completed":false}
```

## 4. GET /todos — list (default: incomplete only, sorted by due_date)

```http
$ curl -i -sS http://127.0.0.1:18080/todos
```

```
HTTP/1.1 200 OK
Content-Type: application/json
Date: Mon, 15 Jun 2026 18:06:15 GMT
Content-Length: 165

[{"id":2,"task":"Second task","due_date":"2026-06-18T10:00:00Z","completed":false},{"id":1,"task":"First task","due_date":"2026-06-20T18:00:00Z","completed":false}]
```

## 5. PUT /todos/{id} — update id=1 (mark completed)

```http
$ curl -i -sS -X PUT http://127.0.0.1:18080/todos/1 \
  -H 'Content-Type: application/json' \
  -d '{"task":"First task (done)","due_date":"2026-06-20T18:00:00Z","completed":true}'
```

```
HTTP/1.1 200 OK
Content-Type: application/json
Date: Mon, 15 Jun 2026 18:06:15 GMT
Content-Length: 87

{"id":1,"task":"First task (done)","due_date":"2026-06-20T18:00:00Z","completed":true}
```

## 6. GET /todos — after completion (completed hidden by default)

```http
$ curl -i -sS http://127.0.0.1:18080/todos
```

```
HTTP/1.1 200 OK
Content-Type: application/json
Date: Mon, 15 Jun 2026 18:06:15 GMT
Content-Length: 84

[{"id":2,"task":"Second task","due_date":"2026-06-18T10:00:00Z","completed":false}]
```

## 7. GET /todos?include_completed=true — include completed

```http
$ curl -i -sS "http://127.0.0.1:18080/todos?include_completed=true"
```

```
HTTP/1.1 200 OK
Content-Type: application/json
Date: Mon, 15 Jun 2026 18:06:15 GMT
Content-Length: 171

[{"id":2,"task":"Second task","due_date":"2026-06-18T10:00:00Z","completed":false},{"id":1,"task":"First task (done)","due_date":"2026-06-20T18:00:00Z","completed":true}]
```

## 8. DELETE /todos/{id} — delete id=2

```http
$ curl -i -sS -X DELETE http://127.0.0.1:18080/todos/2
```

```
HTTP/1.1 204 No Content
Date: Mon, 15 Jun 2026 18:06:15 GMT

```

## 9. DELETE /todos/{id} — delete id=1

```http
$ curl -i -sS -X DELETE http://127.0.0.1:18080/todos/1
```

```
HTTP/1.1 204 No Content
Date: Mon, 15 Jun 2026 18:06:15 GMT

```

## 10. GET /todos/{id} — not found after delete

```http
$ curl -i -sS http://127.0.0.1:18080/todos/1
```

```
HTTP/1.1 404 Not Found
Content-Type: application/json
Date: Mon, 15 Jun 2026 18:06:15 GMT
Content-Length: 27

{"error":"todo not found"}
```

## 11. GET /todos — empty list

```http
$ curl -i -sS http://127.0.0.1:18080/todos
```

```
HTTP/1.1 200 OK
Content-Type: application/json
Date: Mon, 15 Jun 2026 18:06:15 GMT
Content-Length: 3

[]
```

## 12. POST /todos — validation error (empty task)

```http
$ curl -i -sS -X POST http://127.0.0.1:18080/todos \
  -H 'Content-Type: application/json' \
  -d '{"task":"","due_date":"2026-06-20T18:00:00Z"}'
```

```
HTTP/1.1 400 Bad Request
Content-Type: application/json
Date: Mon, 15 Jun 2026 18:06:15 GMT
Content-Length: 33

{"error":"task cannot be empty"}
```
