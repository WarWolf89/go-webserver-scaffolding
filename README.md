## General notes
* code has comments where needed
* use the `make` commands to quickly spin up resources and start testing against the server with curl(see api docs for example commands)
* github actions pipeline runs redis instance in a service container so that e2e tests pass, I recommend for development purposes installing the redis server locally (or run it as a container and then expose the port on the host)

## API docs 

<details>
 <summary><code>GET</code> <code><b>/api/v1/fruits</b></code> <code>(gets all fruits from the redis root hash key 'basket')</code></summary>

##### Parameters

> None

##### Responses

> | http code     | content-type                      | response                                                            |
> |---------------|-----------------------------------|---------------------------------------------------------------------|
> | `200`         | `text/plain;charset=UTF-8`        | json array string                                                         |

##### Example cURL

> ```javascript
>  curl -X GET -H "Content-Type: application/json" http://localhost:8080/api/v1/fruits
> ```

</details>

<details>
 <summary><code>POST</code> <code><b>/api/v1/fruits</b></code> <code>(creates a fruit)</code></summary>

##### Parameters

> | name      |  type     | data type               | description                                                           |
> |-----------|-----------|-------------------------|-----------------------------------------------------------------------|
> | None      |  required | object (JSON or YAML)   | N/A  |

##### Responses

Header has uri in the form of `Location: v1/api/fruits:e89493f1-2645-48a7-9a42-b2073a69027e`

> | http code     | content-type                      | response                                                            |
> |---------------|-----------------------------------|---------------------------------------------------------------------|
> | `201`         | `text/plain;charset=UTF-8`        |    OK                                                        |
> | `400`         | `application/json`                | `{"code":"400","message":"Request body has incorrect data:"}`                            |`                            |

##### Example cURL

> ```javascript
>    curl http://localhost:8080/api/v1/fruits
>    --include \    
>    --header "Content-Type: application/json" \
>    --request "POST" \
>    --data '{"name": "orange","color": "orange"}'
> ```

</details>

<details>
  <summary><code>GET</code> <code><b>/api/v1/fruits:{d}</b></code> <code>(gets the fruit by it's id)</code></summary>

##### Parameters

> | name              |  type     | data type      | description                         |
> |-------------------|-----------|----------------|-------------------------------------|
> | `id` |  required | uuid (v4)   | The specific fruit id        |

##### Responses

> | http code     | content-type                      | response                                                            |
> |---------------|-----------------------------------|---------------------------------------------------------------------|
> | `200`         | `text/plain;charset=UTF-8`        |    `{id: uuid, name: string, color:string}`                                                          |
> | `400`         | `application/json`                | `{"code":"400","message":"no ID provided as a pathparam"}`                            |
> | `404`         | `application/json`                | `{"code":"404","message":"No fruit with id:{id}"}` 

##### Example cURL

> ```javascript
>  curl -X GET -H "Content-Type: application/json" http://localhost:8080/v1/api/fruits:e89493f1-2645-48a7-9a42-b2073a69027e
> ```
</details>



