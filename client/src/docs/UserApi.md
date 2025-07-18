# UserApi

All URIs are relative to *http://localhost*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**userCreatePost**](#usercreatepost) | **POST** /user/create | CreateRoster a new user|
|[**userGet**](#userget) | **GET** /user/ | Get all users with optional filtering|
|[**userIdDelete**](#useriddelete) | **DELETE** /user/{id} | DeleteRoster a user|
|[**userIdGet**](#useridget) | **GET** /user/{id} | Get user by ID|

# **userCreatePost**
> User userCreatePost(createParams)

create user

### Example

```typescript
import {
    UserApi,
    Configuration,
    UserCreateRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new UserApi(configuration);

let createParams: UserCreateRequest; //User input

const { status, data } = await apiInstance.userCreatePost(
    createParams
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **createParams** | **UserCreateRequest**| User input | |


### Return type

**User**

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |
|**400** | Bad Request |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **userGet**
> Array<User> userGet()

Retrieve a list of users with optional query parameter filtering

### Example

```typescript
import {
    UserApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new UserApi(configuration);

let organId: number; //Organ ID (optional) (default to undefined)
let gewisId: number; //GEWIS ID (optional) (default to undefined)

const { status, data } = await apiInstance.userGet(
    organId,
    gewisId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **organId** | [**number**] | Organ ID | (optional) defaults to undefined|
| **gewisId** | [**number**] | GEWIS ID | (optional) defaults to undefined|


### Return type

**Array<User>**

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |
|**400** | Bad Request |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **userIdDelete**
> string userIdDelete()


### Example

```typescript
import {
    UserApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new UserApi(configuration);

let id: number; //User ID (default to undefined)

const { status, data } = await apiInstance.userIdDelete(
    id
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **id** | [**number**] | User ID | defaults to undefined|


### Return type

**string**

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |
|**400** | Bad Request |  -  |
|**404** | Not Found |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **userIdGet**
> User userIdGet()

Retrieve a specific user by their unique ID

### Example

```typescript
import {
    UserApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new UserApi(configuration);

let id: number; //User ID (default to undefined)

const { status, data } = await apiInstance.userIdGet(
    id
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **id** | [**number**] | User ID | defaults to undefined|


### Return type

**User**

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |
|**400** | Bad Request |  -  |
|**404** | Not Found |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

