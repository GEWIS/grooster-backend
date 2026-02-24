# ShiftGroupApi

All URIs are relative to *http://localhost*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**createShiftGroup**](#createshiftgroup) | **POST** /roster/shift-groups | Create a new shift group|
|[**getShiftGroup**](#getshiftgroup) | **GET** /roster/shift-groups/{id} | Get a specific shift group by ID|
|[**getShiftGroups**](#getshiftgroups) | **GET** /roster/shift-groups | Get all shift groups for an organ|

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

