# ShiftGroupApi

All URIs are relative to *http://localhost*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**createShiftGroup**](#createshiftgroup) | **POST** /roster/shift-groups | Create a new shift group|
|[**getShiftGroup**](#getshiftgroup) | **GET** /roster/shift-groups/{id} | Get a specific shift group by ID|
|[**getShiftGroupPriorities**](#getshiftgrouppriorities) | **GET** /roster/shift-groups/{id}/priority | Get a shift group priorities for a shift group|
|[**getShiftGroups**](#getshiftgroups) | **GET** /roster/shift-groups | Get all shift groups for an organ|
|[**updateShiftGroupPriority**](#updateshiftgrouppriority) | **PUT** /roster/shift-groups/{id}/priority | Update a shift group priority|

# **createShiftGroup**
> ShiftGroup createShiftGroup(params)


### Example

```typescript
import {
    ShiftGroupApi,
    Configuration,
    ShiftGroupCreateRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new ShiftGroupApi(configuration);

let params: ShiftGroupCreateRequest; //Shift Group Details

const { status, data } = await apiInstance.createShiftGroup(
    params
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **params** | **ShiftGroupCreateRequest**| Shift Group Details | |


### Return type

**ShiftGroup**

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**201** | Created |  -  |
|**400** | Bad Request |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getShiftGroup**
> ShiftGroup getShiftGroup()


### Example

```typescript
import {
    ShiftGroupApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new ShiftGroupApi(configuration);

let id: number; //Shift Group ID (default to undefined)

const { status, data } = await apiInstance.getShiftGroup(
    id
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **id** | [**number**] | Shift Group ID | defaults to undefined|


### Return type

**ShiftGroup**

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |
|**404** | Not Found |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getShiftGroupPriorities**
> Array<ShiftGroupPriority> getShiftGroupPriorities()


### Example

```typescript
import {
    ShiftGroupApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new ShiftGroupApi(configuration);

let id: number; //ShiftGroup ID (default to undefined)

const { status, data } = await apiInstance.getShiftGroupPriorities(
    id
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **id** | [**number**] | ShiftGroup ID | defaults to undefined|


### Return type

**Array<ShiftGroupPriority>**

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |
|**400** | Invalid request |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getShiftGroups**
> Array<ShiftGroup> getShiftGroups()


### Example

```typescript
import {
    ShiftGroupApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new ShiftGroupApi(configuration);

let organId: number; //Organ ID (default to undefined)

const { status, data } = await apiInstance.getShiftGroups(
    organId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **organId** | [**number**] | Organ ID | defaults to undefined|


### Return type

**Array<ShiftGroup>**

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

# **updateShiftGroupPriority**
> ShiftGroupPriority updateShiftGroupPriority(updateParams)


### Example

```typescript
import {
    ShiftGroupApi,
    Configuration,
    GroupPriorityUpdateParam
} from './api';

const configuration = new Configuration();
const apiInstance = new ShiftGroupApi(configuration);

let id: number; //ShiftGroup ID (default to undefined)
let updateParams: GroupPriorityUpdateParam; //Update parameters

const { status, data } = await apiInstance.updateShiftGroupPriority(
    id,
    updateParams
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **updateParams** | **GroupPriorityUpdateParam**| Update parameters | |
| **id** | [**number**] | ShiftGroup ID | defaults to undefined|


### Return type

**ShiftGroupPriority**

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |
|**400** | Invalid request |  -  |
|**404** | SavedShift not found |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

