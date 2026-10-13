# RosterCommentApi

All URIs are relative to *http://localhost*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**createRosterComment**](#createrostercomment) | **POST** /roster/comment | Create a new roster comment|
|[**getRosterComments**](#getrostercomments) | **GET** /roster/comment | Get all comments for a roster|

# **createRosterComment**
> RosterComment createRosterComment(createParams)


### Example

```typescript
import {
    RosterCommentApi,
    Configuration,
    CommentCreateRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new RosterCommentApi(configuration);

let createParams: CommentCreateRequest; //Roster comment input

const { status, data } = await apiInstance.createRosterComment(
    createParams
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **createParams** | **CommentCreateRequest**| Roster comment input | |


### Return type

**RosterComment**

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
|**403** | Forbidden |  -  |
|**404** | Not Found |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getRosterComments**
> Array<RosterComment> getRosterComments()


### Example

```typescript
import {
    RosterCommentApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new RosterCommentApi(configuration);

let rosterId: number; //Roster ID (default to undefined)

const { status, data } = await apiInstance.getRosterComments(
    rosterId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **rosterId** | [**number**] | Roster ID | defaults to undefined|


### Return type

**Array<RosterComment>**

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

