# OrganApi

All URIs are relative to *http://localhost*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**organIdMemberUserIdPatch**](#organidmemberuseridpatch) | **PATCH** /organ/{id}/member/{userId} | Update settings for a user within an organ|

# **organIdMemberUserIdPatch**
> ModelsUserOrgan organIdMemberUserIdPatch(updateParams)

Update organ-specific settings like nickname/username

### Example

```typescript
import {
    OrganApi,
    Configuration,
    ModelsUpdateMemberSettingsParams
} from './api';

const configuration = new Configuration();
const apiInstance = new OrganApi(configuration);

let id: number; //Organ ID (default to undefined)
let userId: number; //User ID (default to undefined)
let updateParams: ModelsUpdateMemberSettingsParams; //Settings input

const { status, data } = await apiInstance.organIdMemberUserIdPatch(
    id,
    userId,
    updateParams
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **updateParams** | **ModelsUpdateMemberSettingsParams**| Settings input | |
| **id** | [**number**] | Organ ID | defaults to undefined|
| **userId** | [**number**] | User ID | defaults to undefined|


### Return type

**ModelsUserOrgan**

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
|**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

