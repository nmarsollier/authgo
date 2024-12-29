# AuthGo
Microservicio de Autentificación.

## Version: 1.0

**Contact information:**  
Nestor Marsollier  
nmarsollier@gmail.com  

---
### /users/:userID/grant

#### POST
##### Summary

Haiblitar permisos

##### Description

Otorga permisos al usuario indicado, el usuario logueado tiene que tener permiso "admin".

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| userId | path | ID del usuario a habilitar permiso | Yes | string |
| Authorization | header | Bearer {token} | Yes | string |
| body | body | Permisos a Habilitar | Yes | [rest.grantPermissionBody](#restgrantpermissionbody) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | No Content |  |
| 400 | Bad Request | [errs.ValidationErr](#errsvalidationerr) |
| 401 | Unauthorized | [rst.ErrorData](#rsterrordata) |
| 404 | Not Found | [rst.ErrorData](#rsterrordata) |
| 500 | Internal Server Error | [rst.ErrorData](#rsterrordata) |

### /users/:userID/revoke

#### POST
##### Summary

Quitar permisos

##### Description

Quita permisos al usuario indicado, el usuario logueado tiene que tener permiso "admin".

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| userId | path | ID del usuario a quitar permiso | Yes | string |
| Authorization | header | Bearer {token} | Yes | string |
| body | body | Permisos a Qutiar | Yes | [rest.grantPermissionBody](#restgrantpermissionbody) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | No Content |  |
| 400 | Bad Request | [errs.ValidationErr](#errsvalidationerr) |
| 401 | Unauthorized | [rst.ErrorData](#rsterrordata) |
| 404 | Not Found | [rst.ErrorData](#rsterrordata) |
| 500 | Internal Server Error | [rst.ErrorData](#rsterrordata) |

### /users/:userId/disable

#### POST
##### Summary

Deshabilitar Usuario

##### Description

Deshabilita un usuario en el sistema. El usuario logueado debe tener permisos "admin".

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| userId | path | ID del usuario a deshabilitar | Yes | string |
| Authorization | header | Bearer {token} | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | No Content |  |
| 400 | Bad Request | [errs.ValidationErr](#errsvalidationerr) |
| 401 | Unauthorized | [rst.ErrorData](#rsterrordata) |
| 404 | Not Found | [rst.ErrorData](#rsterrordata) |
| 500 | Internal Server Error | [rst.ErrorData](#rsterrordata) |

### /users/:userId/enable

#### POST
##### Summary

Enable User

##### Description

Habilita un usuario en el sistema. El usuario logueado debe tener permisos "admin".

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| userId | path | ID del usuario a habilitar | Yes | string |
| Authorization | header | Bearer {token} | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | No Content |  |
| 400 | Bad Request | [errs.ValidationErr](#errsvalidationerr) |
| 401 | Unauthorized | [rst.ErrorData](#rsterrordata) |
| 404 | Not Found | [rst.ErrorData](#rsterrordata) |
| 500 | Internal Server Error | [rst.ErrorData](#rsterrordata) |

### /users/all

#### GET
##### Summary

Listar Usuarios

##### Description

Obtiene información de todos los usuarios. El usuario logueado debe tener permisos "admin".

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| Authorization | header | Bearer {token} | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Users | [ [user.UserData](#useruserdata) ] |
| 400 | Bad Request | [errs.ValidationErr](#errsvalidationerr) |
| 401 | Unauthorized | [rst.ErrorData](#rsterrordata) |
| 404 | Not Found | [rst.ErrorData](#rsterrordata) |
| 500 | Internal Server Error | [rst.ErrorData](#rsterrordata) |

### /users/current

#### GET
##### Summary

Usuario Actual

##### Description

Obtiene información del usuario actual.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| Authorization | header | Bearer {token} | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | User data | [user.UserData](#useruserdata) |
| 400 | Bad Request | [errs.ValidationErr](#errsvalidationerr) |
| 401 | Unauthorized | [rst.ErrorData](#rsterrordata) |
| 404 | Not Found | [rst.ErrorData](#rsterrordata) |
| 500 | Internal Server Error | [rst.ErrorData](#rsterrordata) |

### /users/password

#### POST
##### Summary

Cambiar Password

##### Description

Cambia la contraseña del usuario actual.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| body | body | Passwords | Yes | [rest.changePasswordBody](#restchangepasswordbody) |
| Authorization | header | Bearer {token} | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | No Content |  |
| 400 | Bad Request | [errs.ValidationErr](#errsvalidationerr) |
| 401 | Unauthorized | [rst.ErrorData](#rsterrordata) |
| 404 | Not Found | [rst.ErrorData](#rsterrordata) |
| 500 | Internal Server Error | [rst.ErrorData](#rsterrordata) |

### /users/signin

#### POST
##### Summary

Login

##### Description

Loguea un usuario en el sistema.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| body | body | Sign in information | Yes | [usecases.SignInRequest](#usecasessigninrequest) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | User Token | [usecases.TokenResponse](#usecasestokenresponse) |
| 400 | Bad Request | [errs.ValidationErr](#errsvalidationerr) |
| 401 | Unauthorized | [rst.ErrorData](#rsterrordata) |
| 404 | Not Found | [rst.ErrorData](#rsterrordata) |
| 500 | Internal Server Error | [rst.ErrorData](#rsterrordata) |

### /users/signout

#### GET
##### Summary

Logout

##### Description

Desloguea un usuario en el sistema, invalida el token.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| Authorization | header | Bearer {token} | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | No Content |  |
| 500 | Error response | [rst.ErrorData](#rsterrordata) |

### /users/signup

#### POST
##### Summary

Registrar Usuario

##### Description

Registra un nuevo usuario en el sistema.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| body | body | Informacion de ususario | Yes | [usecases.SignUpRequest](#usecasessignuprequest) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | User Token | [usecases.TokenResponse](#usecasestokenresponse) |
| 400 | Bad Request | [errs.ValidationErr](#errsvalidationerr) |
| 401 | Unauthorized | [rst.ErrorData](#rsterrordata) |
| 404 | Not Found | [rst.ErrorData](#rsterrordata) |
| 500 | Internal Server Error | [rst.ErrorData](#rsterrordata) |

---
### Models

#### errs.ValidationErr

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| messages | [ [errs.errField](#errserrfield) ] |  | No |

#### errs.errField

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| message | string |  | No |
| path | string |  | No |

#### rest.changePasswordBody

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| currentPassword | string |  | Yes |
| newPassword | string |  | Yes |

#### rest.grantPermissionBody

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| permissions | [ string ] |  | Yes |

#### rst.ErrorData

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| error | string |  | No |

#### usecases.SignInRequest

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| login | string |  | Yes |
| password | string |  | Yes |

#### usecases.SignUpRequest

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| login | string |  | Yes |
| name | string |  | Yes |
| password | string |  | Yes |

#### usecases.TokenResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| token | string |  | No |

#### user.UserData

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| enabled | boolean |  | No |
| id | string |  | No |
| login | string |  | No |
| name | string |  | No |
| permissions | [ string ] |  | No |
