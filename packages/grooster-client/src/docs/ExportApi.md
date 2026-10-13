# ExportApi

All URIs are relative to *http://localhost*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**exportRosterIdGet**](#exportrosteridget) | **GET** /export/roster/{id} | Export roster assignments as PNG|

# **exportRosterIdGet**
> File exportRosterIdGet()

Generates and downloads a PNG image containing the shift assignments for a specific roster.

### Example

```typescript
import {
    ExportApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new ExportApi(configuration);

let id: number; //Roster ID (default to undefined)

const { status, data } = await apiInstance.exportRosterIdGet(
    id
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **id** | [**number**] | Roster ID | defaults to undefined|


### Return type

**File**

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: image/png


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |
|**400** | Invalid ID format |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

