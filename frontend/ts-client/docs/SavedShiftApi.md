# SavedShiftApi

All URIs are relative to *http://localhost*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**getSavedRoster**](#getsavedroster) | **GET** /roster/saved-shift/{id} | Get all saved shifts for a specific roster|
|[**rosterSave**](#rostersave) | **POST** /roster/{id}/save | Save a specific roster|
|[**updateSavedShift**](#updatesavedshift) | **PATCH** /roster/saved-shift/{id} | Update a specific saved shift|

# **getSavedRoster**
> Array<SavedShift> getSavedRoster()


### Example

```typescript
import {
    SavedShiftApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new SavedShiftApi(configuration);

let id: number; //Roster ID (default to undefined)

const { status, data } = await apiInstance.getSavedRoster(
    id
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **id** | [**number**] | Roster ID | defaults to undefined|


### Return type

**Array<SavedShift>**

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Saved Shifts |  -  |
|**400** | Invalid request |  -  |
|**404** | SavedShift not found |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **rosterSave**
> string rosterSave()


### Example

```typescript
import {
    SavedShiftApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new SavedShiftApi(configuration);

let id: number; //Roster ID (default to undefined)

const { status, data } = await apiInstance.rosterSave(
    id
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **id** | [**number**] | Roster ID | defaults to undefined|


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

# **updateSavedShift**
> SavedShift updateSavedShift(updateParams)


### Example

```typescript
import {
    SavedShiftApi,
    Configuration,
    SavedShiftUpdateRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new SavedShiftApi(configuration);

let id: number; //SavedShift ID (default to undefined)
let updateParams: SavedShiftUpdateRequest; //Update data

const { status, data } = await apiInstance.updateSavedShift(
    id,
    updateParams
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **updateParams** | **SavedShiftUpdateRequest**| Update data | |
| **id** | [**number**] | SavedShift ID | defaults to undefined|


### Return type

**SavedShift**

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

