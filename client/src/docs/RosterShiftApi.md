# RosterShiftApi

All URIs are relative to *http://localhost*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**createRosterShift**](#createrostershift) | **POST** /roster/shift | Create a new roster shift|
|[**deleteRosterShift**](#deleterostershift) | **DELETE** /roster/shift/{id} | Deletes a roster shift|
|[**updateRosterShift**](#updaterostershift) | **PATCH** /roster/shift/{id} | Update a roster shift|

# **createRosterShift**
> RosterShift createRosterShift(createParams)


### Example

```typescript
import {
    RosterShiftApi,
    Configuration,
    ShiftCreateRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new RosterShiftApi(configuration);

let createParams: ShiftCreateRequest; //Roster shift input

const { status, data } = await apiInstance.createRosterShift(
    createParams
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **createParams** | **ShiftCreateRequest**| Roster shift input | |


### Return type

**RosterShift**

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

# **deleteRosterShift**
> string deleteRosterShift()


### Example

```typescript
import {
    RosterShiftApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new RosterShiftApi(configuration);

let id: number; //Roster Answer ID (default to undefined)

const { status, data } = await apiInstance.deleteRosterShift(
    id
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **id** | [**number**] | Roster Answer ID | defaults to undefined|


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

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **updateRosterShift**
> RosterShift updateRosterShift(updateParams)


### Example

```typescript
import {
    RosterShiftApi,
    Configuration,
    ShiftUpdateRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new RosterShiftApi(configuration);

let id: number; //Roster Shift ID (default to undefined)
let updateParams: ShiftUpdateRequest; //Update input

const { status, data } = await apiInstance.updateRosterShift(
    id,
    updateParams
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **updateParams** | **ShiftUpdateRequest**| Update input | |
| **id** | [**number**] | Roster Shift ID | defaults to undefined|


### Return type

**RosterShift**

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
|**404** | Not Found |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

