# AuthApi

All URIs are relative to *http://localhost*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**authCallbackGet**](#authcallbackget) | **GET** /auth/callback | Handle OAuth2 Callback|
|[**authRedirectGet**](#authredirectget) | **GET** /auth/redirect | Redirect to OIDC provider|

# **authCallbackGet**
> { [key: string]: string; } authCallbackGet()

Validates state, exchanges code for token, and returns user info

### Example

```typescript
import {
    AuthApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new AuthApi(configuration);

let state: string; //State returned from provider (default to undefined)
let code: string; //Authorization code from provider (default to undefined)

const { status, data } = await apiInstance.authCallbackGet(
    state,
    code
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **state** | [**string**] | State returned from provider | defaults to undefined|
| **code** | [**string**] | Authorization code from provider | defaults to undefined|


### Return type

**{ [key: string]: string; }**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: */*


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | User info and token |  -  |
|**400** | Bad request: missing or invalid state |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **authRedirectGet**
> string authRedirectGet()

Generates state, sets a cookie, and redirects to Google OIDC

### Example

```typescript
import {
    AuthApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new AuthApi(configuration);

let state: string; //State returned from provider (default to undefined)

const { status, data } = await apiInstance.authRedirectGet(
    state
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **state** | [**string**] | State returned from provider | defaults to undefined|


### Return type

**string**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: */*


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | redirect |  -  |
|**500** | pkg server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

