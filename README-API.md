# AuthGo
Microservicio de Autentificación.

## Version: 1.0

**Contact information:**  
Nestor Marsollier  
nmarsollier@gmail.com  

---
### /rabbit/logout

#### PUT
##### Summary

Mensage Rabbit

##### Description

SendLogout envía un broadcast a rabbit con logout. Esto no es Rest es RabbitMQ.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| body | body | Token deshabilitado | Yes | [rabbit.message](#rabbitmessage) |

##### Responses

| Code | Description |
| ---- | ----------- |

---
### /v1/user

#### POST
##### Summary

Registrar Usuario

##### Description

Registra un nuevo usuario en el sistema.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| body | body | Informacion de ususario | Yes | [user.SignUpRequest](#usersignuprequest) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | User Token | [rest.tokenResponse](#resttokenresponse) |
| 400 | Bad Request | [errs.ValidationErr](#errsvalidationerr) |
| 401 | Unauthorized | [engine.ErrorData](#engineerrordata) |
| 404 | Not Found | [engine.ErrorData](#engineerrordata) |
| 500 | Internal Server Error | [engine.ErrorData](#engineerrordata) |

### /v1/user/password

#### POST
##### Summary

Cambiar Password

##### Description

Cambia la contraseña del usuario actual.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| body | body | Passwords | Yes | [rest.changePasswordBody](#restchangepasswordbody) |
| Authorization | header | bearer {token} | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | No Content |  |
| 400 | Bad Request | [errs.ValidationErr](#errsvalidationerr) |
| 401 | Unauthorized | [engine.ErrorData](#engineerrordata) |
| 404 | Not Found | [engine.ErrorData](#engineerrordata) |
| 500 | Internal Server Error | [engine.ErrorData](#engineerrordata) |

### /v1/user/signin

#### POST
##### Summary

Login

##### Description

Loguea un usuario en el sistema.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| body | body | Sign in information | Yes | [user.SignInRequest](#usersigninrequest) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | User Token | [rest.tokenResponse](#resttokenresponse) |
| 400 | Bad Request | [errs.ValidationErr](#errsvalidationerr) |
| 401 | Unauthorized | [engine.ErrorData](#engineerrordata) |
| 404 | Not Found | [engine.ErrorData](#engineerrordata) |
| 500 | Internal Server Error | [engine.ErrorData](#engineerrordata) |

### /v1/user/signout

#### GET
##### Summary

Logout

##### Description

Desloguea un usuario en el sistema, invalida el token.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| Authorization | header | bearer {token} | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | No Content |  |
| 500 | Error response | [engine.ErrorData](#engineerrordata) |

### /v1/users

#### GET
##### Summary

Listar Usuarios

##### Description

Obtiene información de todos los usuarios. El usuario logueado debe tener permisos "admin".

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| Authorization | header | bearer {token} | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Users | [ [rest.UserDataResponse](#restuserdataresponse) ] |
| 400 | Bad Request | [errs.ValidationErr](#errsvalidationerr) |
| 401 | Unauthorized | [engine.ErrorData](#engineerrordata) |
| 404 | Not Found | [engine.ErrorData](#engineerrordata) |
| 500 | Internal Server Error | [engine.ErrorData](#engineerrordata) |

### /v1/users/:userID/grant

#### POST
##### Summary

Haiblitar permisos

##### Description

Otorga permisos al usuario indicado, el usuario logueado tiene que tener permiso "admin".

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| userId | path | ID del usuario a habilitar permiso | Yes | string |
| Authorization | header | bearer {token} | Yes | string |
| body | body | Permisos a Habilitar | Yes | [rest.grantPermissionBody](#restgrantpermissionbody) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | No Content |  |
| 400 | Bad Request | [errs.ValidationErr](#errsvalidationerr) |
| 401 | Unauthorized | [engine.ErrorData](#engineerrordata) |
| 404 | Not Found | [engine.ErrorData](#engineerrordata) |
| 500 | Internal Server Error | [engine.ErrorData](#engineerrordata) |

### /v1/users/:userID/revoke

#### POST
##### Summary

Quitar permisos

##### Description

Quita permisos al usuario indicado, el usuario logueado tiene que tener permiso "admin".

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| userId | path | ID del usuario a quitar permiso | Yes | string |
| Authorization | header | bearer {token} | Yes | string |
| body | body | Permisos a Qutiar | Yes | [rest.grantPermissionBody](#restgrantpermissionbody) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | No Content |  |
| 400 | Bad Request | [errs.ValidationErr](#errsvalidationerr) |
| 401 | Unauthorized | [engine.ErrorData](#engineerrordata) |
| 404 | Not Found | [engine.ErrorData](#engineerrordata) |
| 500 | Internal Server Error | [engine.ErrorData](#engineerrordata) |

### /v1/users/:userId/disable

#### POST
##### Summary

Deshabilitar Usuario

##### Description

Deshabilita un usuario en el sistema. El usuario logueado debe tener permisos "admin".

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| userId | path | ID del usuario a deshabilitar | Yes | string |
| Authorization | header | bearer {token} | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | No Content |  |
| 400 | Bad Request | [errs.ValidationErr](#errsvalidationerr) |
| 401 | Unauthorized | [engine.ErrorData](#engineerrordata) |
| 404 | Not Found | [engine.ErrorData](#engineerrordata) |
| 500 | Internal Server Error | [engine.ErrorData](#engineerrordata) |

### /v1/users/:userId/enable

#### POST
##### Summary

Enable User

##### Description

Habilita un usuario en el sistema. El usuario logueado debe tener permisos "admin".

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| userId | path | ID del usuario a habilitar | Yes | string |
| Authorization | header | bearer {token} | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | No Content |  |
| 400 | Bad Request | [errs.ValidationErr](#errsvalidationerr) |
| 401 | Unauthorized | [engine.ErrorData](#engineerrordata) |
| 404 | Not Found | [engine.ErrorData](#engineerrordata) |
| 500 | Internal Server Error | [engine.ErrorData](#engineerrordata) |

### /v1/users/current

#### GET
##### Summary

Usuario Actual

##### Description

Obtiene información del usuario actual.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| Authorization | header | bearer {token} | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | User data | [rest.UserResponse](#restuserresponse) |
| 400 | Bad Request | [errs.ValidationErr](#errsvalidationerr) |
| 401 | Unauthorized | [engine.ErrorData](#engineerrordata) |
| 404 | Not Found | [engine.ErrorData](#engineerrordata) |
| 500 | Internal Server Error | [engine.ErrorData](#engineerrordata) |

---
### Models

#### engine.ErrorData

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| error | string |  | No |

#### errs.ValidationErr

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| messages | [ [errs.errField](#errserrfield) ] |  | No |

#### errs.errField

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| message | string |  | No |
| path | string |  | No |

#### rabbit.message

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| message | string | *Example:* `"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbklEIjoiNjZiNjBlYzhlMGYzYzY4OTUzMzJlOWNmIiwidXNlcklEIjoiNjZhZmQ3ZWU4YTBhYjRjZjQ0YTQ3NDcyIn0.who7upBctOpmlVmTvOgH1qFKOHKXmuQCkEjMV3qeySg"` | No |
| type | string | *Example:* `"logout"` | No |

#### rest.UserDataResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| enabled | boolean |  | No |
| id | string |  | No |
| login | string |  | No |
| name | string |  | No |
| permissions | [ string ] |  | No |

#### rest.UserResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| id | string |  | No |
| login | string |  | No |
| name | string |  | No |
| permissions | [ string ] |  | No |

#### rest.changePasswordBody

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| currentPassword | string |  | Yes |
| newPassword | string |  | Yes |

#### rest.grantPermissionBody

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| permissions | [ string ] |  | Yes |

#### rest.tokenResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| token | string |  | No |

#### user.SignInRequest

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| login | string |  | Yes |
| password | string |  | Yes |

#### user.SignUpRequest

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| login | string |  | Yes |
| name | string |  | Yes |
| password | string |  | Yes |
