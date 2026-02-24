# OrganApi

All URIs are relative to *http://localhost*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**organIdMemberUserIdGet**](#organidmemberuseridget) | **GET** /organ/{id}/member/{userId} | Get settings for a user within an organ|
|[**organIdMemberUserIdPatch**](#organidmemberuseridpatch) | **PATCH** /organ/{id}/member/{userId} | Update settings for a user within an organ|

# **organIdMemberUserIdGet**
> UserOrgan organIdMemberUserIdGet()

Get organ-specific settings like nickname/username

### Example

```typescript
import {
    OrganApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new OrganApi(configuration);

let id: number; //Organ ID (default to undefined)
let userId: number; //User ID (default to undefined)

const { status, data } = await apiInstance.organIdMemberUserIdGet(
    id,
    userId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **id** | [**number**] | Organ ID | defaults to undefined|
| **userId** | [**number**] | User ID | defaults to undefined|


### Return type

**UserOrgan**

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

# **organIdMemberUserIdPatch**
> UserOrgan organIdMemberUserIdPatch(updateParams)

Update organ-specific settings like nickname/username

### Example

```typescript
import {
    OrganApi,
    Configuration,
    UpdateMemberSettingsParams
} from './api';

const configuration = new Configuration();
const apiInstance = new OrganApi(configuration);

let id: number; //Organ ID (default to undefined)
let userId: number; //User ID (default to undefined)
let updateParams: UpdateMemberSettingsParams; //Settings input

const { status, data } = await apiInstance.organIdMemberUserIdPatch(
    id,
    userId,
    updateParams
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **updateParams** | **UpdateMemberSettingsParams**| Settings input | |
| **id** | [**number**] | Organ ID | defaults to undefined|
| **userId** | [**number**] | User ID | defaults to undefined|


### Return type

**UserOrgan**

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

