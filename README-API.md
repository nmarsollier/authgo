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
| 400 | Bad Request | [app_errors.ErrValidation](#app_errorserrvalidation) |
| 401 | Unauthorized | [app_errors.OtherErrors](#app_errorsothererrors) |
| 404 | Not Found | [app_errors.OtherErrors](#app_errorsothererrors) |
| 500 | Internal Server Error | [app_errors.OtherErrors](#app_errorsothererrors) |

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
| 400 | Bad Request | [app_errors.ErrValidation](#app_errorserrvalidation) |
| 401 | Unauthorized | [app_errors.OtherErrors](#app_errorsothererrors) |
| 404 | Not Found | [app_errors.OtherErrors](#app_errorsothererrors) |
| 500 | Internal Server Error | [app_errors.OtherErrors](#app_errorsothererrors) |

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
| 400 | Bad Request | [app_errors.ErrValidation](#app_errorserrvalidation) |
| 401 | Unauthorized | [app_errors.OtherErrors](#app_errorsothererrors) |
| 404 | Not Found | [app_errors.OtherErrors](#app_errorsothererrors) |
| 500 | Internal Server Error | [app_errors.OtherErrors](#app_errorsothererrors) |

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
| 500 | Error response | [app_errors.OtherErrors](#app_errorsothererrors) |

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
| 400 | Bad Request | [app_errors.ErrValidation](#app_errorserrvalidation) |
| 401 | Unauthorized | [app_errors.OtherErrors](#app_errorsothererrors) |
| 404 | Not Found | [app_errors.OtherErrors](#app_errorsothererrors) |
| 500 | Internal Server Error | [app_errors.OtherErrors](#app_errorsothererrors) |

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
| 400 | Bad Request | [app_errors.ErrValidation](#app_errorserrvalidation) |
| 401 | Unauthorized | [app_errors.OtherErrors](#app_errorsothererrors) |
| 404 | Not Found | [app_errors.OtherErrors](#app_errorsothererrors) |
| 500 | Internal Server Error | [app_errors.OtherErrors](#app_errorsothererrors) |

---
### /v1/users/:userID/grant

#### POST
##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | User Token | [rest.tokenResponse](#resttokenresponse) |
| 400 | Bad Request | [app_errors.ErrValidation](#app_errorserrvalidation) |
| 401 | Unauthorized | [app_errors.OtherErrors](#app_errorsothererrors) |
| 404 | Not Found | [app_errors.OtherErrors](#app_errorsothererrors) |
| 500 | Internal Server Error | [app_errors.OtherErrors](#app_errorsothererrors) |

### /v1/users/:userID/revoke

#### POST
##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | User Token | [rest.revokePermissionBody](#restrevokepermissionbody) |
| 400 | Bad Request | [app_errors.ErrValidation](#app_errorserrvalidation) |
| 401 | Unauthorized | [app_errors.OtherErrors](#app_errorsothererrors) |
| 404 | Not Found | [app_errors.OtherErrors](#app_errorsothererrors) |
| 500 | Internal Server Error | [app_errors.OtherErrors](#app_errorsothererrors) |

### /v1/users/:userId/disable

#### POST
##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | User Token | [rest.tokenResponse](#resttokenresponse) |
| 400 | Bad Request | [app_errors.ErrValidation](#app_errorserrvalidation) |
| 401 | Unauthorized | [app_errors.OtherErrors](#app_errorsothererrors) |
| 404 | Not Found | [app_errors.OtherErrors](#app_errorsothererrors) |
| 500 | Internal Server Error | [app_errors.OtherErrors](#app_errorsothererrors) |

### /v1/users/:userId/enable

#### POST
##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | User Token | [rest.tokenResponse](#resttokenresponse) |
| 400 | Bad Request | [app_errors.ErrValidation](#app_errorserrvalidation) |
| 401 | Unauthorized | [app_errors.OtherErrors](#app_errorsothererrors) |
| 404 | Not Found | [app_errors.OtherErrors](#app_errorsothererrors) |
| 500 | Internal Server Error | [app_errors.OtherErrors](#app_errorsothererrors) |

---
### Models

#### app_errors.ErrField

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| message | string |  | No |
| path | string |  | No |

#### app_errors.ErrValidation

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| messages | [ [app_errors.ErrField](#app_errorserrfield) ] |  | No |

#### app_errors.OtherErrors

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| error | string |  | No |

#### rabbit.message

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| message | string |  | No |
| type | string |  | No |

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

#### rest.revokePermissionBody

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| permissions | [ string ] |  | Yes |
| user | string |  | Yes |

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
