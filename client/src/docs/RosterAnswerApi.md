# RosterAnswerApi

All URIs are relative to *http://localhost*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**createRosterAnswer**](#createrosteranswer) | **POST** /roster/answer | Create a new roster shift answer|
|[**updateRosterAnswer**](#updaterosteranswer) | **PATCH** /roster/answer/{id} | Updates a roster answer with the new value|

# **createRosterAnswer**
> RosterAnswer createRosterAnswer(createParams)


### Example

```typescript
import {
    RosterAnswerApi,
    Configuration,
    RosterAnswerCreateRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new RosterAnswerApi(configuration);

let createParams: RosterAnswerCreateRequest; //Roster answer input

const { status, data } = await apiInstance.createRosterAnswer(
    createParams
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **createParams** | **RosterAnswerCreateRequest**| Roster answer input | |


### Return type

**RosterAnswer**

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

# **updateRosterAnswer**
> RosterAnswer updateRosterAnswer(updateParams)


### Example

```typescript
import {
    RosterAnswerApi,
    Configuration,
    RosterAnswerUpdateRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new RosterAnswerApi(configuration);

let id: number; //Roster Answer ID (default to undefined)
let updateParams: RosterAnswerUpdateRequest; //New answer value

const { status, data } = await apiInstance.updateRosterAnswer(
    id,
    updateParams
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **updateParams** | **RosterAnswerUpdateRequest**| New answer value | |
| **id** | [**number**] | Roster Answer ID | defaults to undefined|


### Return type

**RosterAnswer**

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

